//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"net"
	"time"
)

const maxTrySend = 3

// MTNetSocket wraps a TCP connection to a MetaTrader 5 server.
type MTNetSocket struct {
	conn   net.Conn
	server string
}

// Connect establishes a TCP connection to the given server and port.
func (s *MTNetSocket) Connect(server string, port int) error {
	s.server = server
	addr := fmt.Sprintf("%s:%d", server, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

// Send writes data to the connection with retry logic and a timeout.
func (s *MTNetSocket) Send(data []byte, dataLen, timeout int) error {
	if !s.IsConnected() {
		return fmt.Errorf("not connected")
	}
	dataPos := 0
	trySend := 0
	deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	if err := s.conn.SetWriteDeadline(deadline); err != nil {
		return err
	}
	for dataPos < dataLen {
		n, err := s.conn.Write(data[dataPos:dataLen])
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() && trySend < maxTrySend {
				time.Sleep(30 * time.Millisecond)
				trySend++
				continue
			}
			return err
		}
		dataPos += n
	}
	return nil
}

// Recv reads exactly dataLen bytes from the connection with a timeout.
func (s *MTNetSocket) Recv(dataLen, timeout int) ([]byte, error) {
	if !s.IsConnected() {
		return nil, fmt.Errorf("connection closed")
	}
	data := make([]byte, dataLen)
	dataPos := 0
	if timeout > 0 {
		deadline := time.Now().Add(time.Duration(timeout) * time.Millisecond)
		if err := s.conn.SetReadDeadline(deadline); err != nil {
			return nil, err
		}
	} else {
		// timeout 0 means no deadline (block indefinitely)
		if err := s.conn.SetReadDeadline(time.Time{}); err != nil {
			return nil, err
		}
	}
	for dataPos < dataLen {
		n, err := s.conn.Read(data[dataPos:])
		dataPos += n
		if n == 0 && err == nil {
			s.Close()
			return nil, fmt.Errorf("connection closed by remote server")
		}
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

// Close shuts down and closes the connection.
func (s *MTNetSocket) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
}

// IsConnected returns true if the socket is currently connected.
func (s *MTNetSocket) IsConnected() bool {
	return s.conn != nil
}

// Server returns the server address.
func (s *MTNetSocket) Server() string {
	return s.server
}
