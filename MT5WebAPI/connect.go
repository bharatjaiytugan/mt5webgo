//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unicode/utf16"
)

// MAX_CLIENT_COMMAND is the maximum serial number for client commands (0-16383).
const MAX_CLIENT_COMMAND = 16383

// mtConnect provides synchronous connection to the MetaTrader 5 server.
type mtConnect struct {
	socket        *MTNetSocket
	server        string
	port          int
	timeout       int
	lock          sync.Mutex
	clientCommand int
}

// newMTConnect creates a new synchronous connection instance.
func newMTConnect(server string, port int, timeout int) *mtConnect {
	return &mtConnect{
		server:  server,
		port:    port,
		timeout: timeout,
	}
}

// Connect establishes a TCP connection to the MetaTrader 5 server.
func (c *mtConnect) Connect() MTRetCode {
	c.socket = &MTNetSocket{}
	err := c.socket.Connect(c.server, c.port)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("socket connect failed to %s:%d, %v", c.server, c.port, err))
		return MT_RET_ERR_NETWORK
	}
	logWrite(LogTypeError, fmt.Sprintf("socket connect created to %s:%d", c.server, c.port))
	return MT_RET_OK
}

// Disconnect closes the connection to the MetaTrader 5 server.
func (c *mtConnect) Disconnect() {
	if c.socket != nil {
		c.socket.Close()
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "socket connect closed")
		}
	}
}

// Send sends a command with parameters to the MetaTrader 5 server.
func (c *mtConnect) Send(command string, data map[string]string, isFirstRequest bool) bool {
	if c.socket == nil || !c.socket.IsConnected() {
		logWrite(LogTypeError, "send data: connection closed")
		return false
	}
	// increment packet number
	c.lock.Lock()
	if c.clientCommand >= MAX_CLIENT_COMMAND {
		c.clientCommand = 1
	} else {
		c.clientCommand++
	}
	c.lock.Unlock()
	// build query string
	queryTemp := command
	if data != nil && len(data) > 0 {
		bodyRequest := ""
		queryTemp += "|"
		for key, val := range data {
			if key == WEB_PARAM_BODYTEXT {
				bodyRequest = val
			} else {
				queryTemp += fmt.Sprintf("%s=%s|", key, Quotes(val))
			}
		}
		queryTemp += "\r\n"
		if bodyRequest != "" {
			queryTemp += bodyRequest
		}
	} else {
		queryTemp += "|\r\n"
	}
	// convert to UTF-16LE (Unicode)
	queryBody := encodeUTF16LE(queryTemp)
	// build header
	var headerStr string
	if isFirstRequest {
		headerStr = fmt.Sprintf(WEB_PREFIX_WEBAPI, len(queryBody), c.clientCommand, "0")
	} else {
		headerStr = fmt.Sprintf(WEB_PACKET_FORMAT, len(queryBody), c.clientCommand, "0")
	}
	header := []byte(headerStr)
	// combine header and body
	query := append(header, queryBody...)
	// send data to server
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("sending data: %s%s, length data: %d", headerStr, GetString(queryBody), len(query)))
	}
	if err := c.socket.Send(query, len(query), c.timeout); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("send data failed, %s. %v", GetString(query), err))
		return false
	}
	return true
}

// Read reads data from the MetaTrader 5 server.
func (c *mtConnect) Read(isAuthPacket, isBinary bool) []byte {
	if c.socket == nil || !c.socket.IsConnected() {
		logWrite(LogTypeError, "read data: connection closed")
		return nil
	}
	var result []byte
	for {
		data, header := c.readHeaderPacket()
		if header == nil {
			break
		}
		if data == nil && header.SizeBody > 0 {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("read incorrect packet, length %d, but data is nil", header.SizeBody))
			}
			break
		}
		// check number of packet
		if header.NumberPacket != c.clientCommand {
			if header.SizeBody != 0 {
				logWrite(LogTypeError, fmt.Sprintf("number of packet incorrect need: %d, but get %d", c.clientCommand, header.NumberPacket))
			} else {
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, "PING packet")
				}
			}
			continue
		}
		// accumulate result
		if result == nil {
			result = data
		} else {
			result = append(result, data...)
		}
		// last packet check
		if header.Flag == 0 {
			break
		}
	}
	return result
}

// readHeaderPacket reads a single protocol packet (header + body).
func (c *mtConnect) readHeaderPacket() ([]byte, *protocolHeader) {
	// read header
	headerData, err := c.socket.Recv(HEADER_LENGTH, c.timeout)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("get data failed: %v", err))
		return nil, nil
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("header: %s, length: %d", string(headerData), len(headerData)))
	}
	// parse header
	header := parseHeader(headerData)
	if header == nil {
		logWrite(LogTypeError, "incorrect header data")
		return nil, nil
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("size body: %d, number package: %d, flag: %d", header.SizeBody, header.NumberPacket, header.Flag))
	}
	// read body if present
	if header.SizeBody <= 0 {
		return nil, header
	}
	data, err := c.socket.Recv(header.SizeBody, c.timeout)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("receive data failed, %v", err))
		return nil, nil
	}
	if len(data) != header.SizeBody {
		logWrite(LogTypeError, fmt.Sprintf("incorrect size of block: %d, real size: %d", header.SizeBody, len(data)))
		return nil, nil
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("read all data %d bytes", len(data)))
	}
	return data, header
}

// IsConnected returns true if the connection is active.
func (c *mtConnect) IsConnected() bool {
	if c.socket == nil {
		return false
	}
	return c.socket.IsConnected()
}

// Socket returns the underlying network socket.
func (c *mtConnect) Socket() *MTNetSocket {
	return c.socket
}

// --- helper functions ---

// encodeUTF16LE converts a Go string to UTF-16 Little Endian byte slice.
func encodeUTF16LE(s string) []byte {
	runes := []rune(s)
	u16 := utf16.Encode(runes)
	buf := make([]byte, len(u16)*2)
	for i, v := range u16 {
		binary.LittleEndian.PutUint16(buf[i*2:], v)
	}
	return buf
}
