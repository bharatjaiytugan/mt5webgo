//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// callbackReceive is the delegate type for data receive callbacks.
type callbackReceive func(data []byte, done chan struct{})

// queueData contains data for adding to the send queue.
type queueData struct {
	Callback  callbackReceive
	Command   string
	Params    map[string]string
	IsReceive bool
	ID        int
	Done      chan struct{} // replaces ManualResetEvent
}

// newQueueData creates a new queueData with default values.
func newQueueData() *queueData {
	return &queueData{
		IsReceive: true,
		Done:      make(chan struct{}, 1),
	}
}

// dataReceived holds data received from the MT server.
type dataReceived struct {
	ID         int
	Command    string
	Data       [][]byte
	DataLength int64
	DateCreate time.Time
	IsFinished bool
	Callback   callbackReceive
	Done       chan struct{}
}

// dataSend holds data to be sent to the MT server.
type dataSend struct {
	ID         int
	Command    string
	Data       []byte
	DateCreate time.Time
}

// Async connection constants.
const (
	waitLockerTimeout      = 500 * time.Millisecond
	periodPingSend         = 20 // seconds
	timeoutLiveQueue       = 120 // seconds
	maxPacketSize          = 65536 - 1
	waitThreadTimeout      = 100 * time.Millisecond
	waitThreadClearTimeout = 1000 * time.Millisecond
)

// mtAsyncConnect manages asynchronous communication with the MetaTrader 5 server.
type mtAsyncConnect struct {
	lockerQueueSend    sync.RWMutex
	lockerQueueReceive sync.RWMutex
	lockerSocket       sync.RWMutex

	queueSend    []*dataSend
	queueReceive map[int]*dataReceived

	stopSignal chan struct{}
	socket     *MTNetSocket
	crypt      *MTCrypt
	cryptMode  EnCryptModes
	timeout    int

	numberClientCommand int32
	lastDateSend        time.Time

	running bool
}

// newMTAsyncConnect creates a new async connection instance.
func newMTAsyncConnect() *mtAsyncConnect {
	return &mtAsyncConnect{
		queueReceive: make(map[int]*dataReceived),
		timeout:      5000,
	}
}

// Start initializes and starts the async connection goroutines.
func (ac *mtAsyncConnect) Start(socket *MTNetSocket, cryptMode EnCryptModes, timeout int, cryptRand []byte, password string) bool {
	if ac.running {
		ac.Stop()
	}
	if socket == nil {
		logWrite(LogTypeError, "socket info is nil")
		return false
	}
	ac.cryptMode = cryptMode
	ac.timeout = timeout
	ac.stopSignal = make(chan struct{})
	// set crypt information for AES
	ac.crypt = &MTCrypt{}
	ac.crypt.SetCryptAESRand(cryptRand, password)
	ac.socket = socket
	ac.running = true
	// start receive goroutine
	go ac.threadReceiveData()
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, "receive thread started")
	}
	// start send goroutine
	go ac.threadSendData()
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, "send thread started")
	}
	// start clear queues goroutine
	go ac.threadClearQueues()
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, "clear queue thread started")
	}
	ac.lastDateSend = time.Now()
	return true
}

// Add adds a command to the send/receive queues.
func (ac *mtAsyncConnect) Add(data *queueData) bool {
	if data == nil {
		logWrite(LogTypeError, "can not add empty data to queue")
		return false
	}
	number := int(atomic.AddInt32(&ac.numberClientCommand, 1))
	if number >= MAX_CLIENT_COMMAND {
		atomic.StoreInt32(&ac.numberClientCommand, 1)
	}
	data.ID = number
	// add to receive queue
	ac.lockerQueueReceive.Lock()
	if _, exists := ac.queueReceive[data.ID]; exists {
		ac.lockerQueueReceive.Unlock()
		logWrite(LogTypeError, fmt.Sprintf("adding data to received failed, #%d already exists in list", data.ID))
		return false
	}
	ac.queueReceive[data.ID] = &dataReceived{
		Callback:   data.Callback,
		Done:       data.Done,
		Command:    data.Command,
		ID:         data.ID,
		IsFinished: false,
		DateCreate: time.Now(),
		Data:       make([][]byte, 0),
		DataLength: 0,
	}
	ac.lockerQueueReceive.Unlock()
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("add data to receive: #%d, command: %s", data.ID, data.Command))
	}
	// get bytes for send
	bytesSend := getDataSend(data.Command, data.Params)
	// add to send queue
	ac.lockerQueueSend.Lock()
	ac.queueSend = append(ac.queueSend, &dataSend{
		ID:         data.ID,
		Command:    data.Command,
		Data:       bytesSend,
		DateCreate: time.Now(),
	})
	ac.lockerQueueSend.Unlock()
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("add data to send: #%d, command: %s", data.ID, data.Command))
	}
	return true
}

// Disconnect closes the socket connection.
func (ac *mtAsyncConnect) Disconnect() {
	if ac.socket != nil {
		ac.lockerSocket.Lock()
		ac.socket.Close()
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "socket connect closed")
		}
		ac.lockerSocket.Unlock()
	}
}

// Stop stops all goroutines and closes the connection.
func (ac *mtAsyncConnect) Stop() {
	if !ac.running {
		return
	}
	// close connection
	if ac.socket != nil {
		ac.socket.Close()
	}
	// signal all goroutines to stop
	close(ac.stopSignal)
	ac.running = false
	// wait a bit for goroutines to finish
	time.Sleep(4 * waitThreadTimeout)
	// clear queues
	ac.sendRemoveAll()
	ac.receiveRemoveAll()
}

// IsConnected returns true if the async connection is active.
func (ac *mtAsyncConnect) IsConnected() bool {
	if ac.socket == nil {
		return false
	}
	return ac.socket.IsConnected()
}

// --- thread functions ---

// threadClearQueues periodically clears expired items from the queues.
func (ac *mtAsyncConnect) threadClearQueues() {
	for {
		select {
		case <-ac.stopSignal:
			ac.clearQueues()
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, "the clear queue thread is stopped")
			}
			return
		case <-time.After(waitThreadClearTimeout):
			ac.clearQueues()
		}
	}
}

// clearQueues removes expired items from both send and receive queues.
func (ac *mtAsyncConnect) clearQueues() {
	// clear send queue
	for {
		ds := ac.getFirstSend()
		if ds == nil {
			break
		}
		if ds.DateCreate.Add(time.Duration(timeoutLiveQueue) * time.Second).Before(time.Now()) {
			if ac.sendRemove(ds) {
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, fmt.Sprintf("remove element from send id: #%d, command: %s by timeout, date create: %s", ds.ID, ds.Command, ds.DateCreate.Format("2006.01.02 15:04")))
				}
			}
		} else {
			break
		}
	}
	// clear receive queue
	for {
		dr := ac.getFirstReceive()
		if dr == nil {
			break
		}
		if dr.DateCreate.Add(time.Duration(timeoutLiveQueue) * time.Second).Before(time.Now()) {
			if ac.receiveRemove(dr) {
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, fmt.Sprintf("remove element from receive id: #%d, command: %s by timeout, date create: %s", dr.ID, dr.Command, dr.DateCreate.Format("2006.01.02 15:04")))
				}
			}
		} else {
			break
		}
	}
}

// threadReceiveData is the goroutine that receives data from the MT server.
func (ac *mtAsyncConnect) threadReceiveData() {
	for {
		select {
		case <-ac.stopSignal:
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, "the receive thread is stopped")
			}
			return
		case <-time.After(waitThreadTimeout):
			if !ac.receiveFromMetaTraderServer() {
				return
			}
		}
	}
}

// receiveFromMetaTraderServer reads and processes packets from the MT server.
func (ac *mtAsyncConnect) receiveFromMetaTraderServer() bool {
	if ac.socket == nil || !ac.socket.IsConnected() {
		if IsWriteDebugLog() {
			logWrite(LogTypeError, "read data: connection closed")
		}
		ac.updateReceiveQueueConnectionError()
		return false
	}
	for {
		data, header := ac.readHeaderPacket()
		if header == nil {
			if ac.socket == nil || !ac.socket.IsConnected() {
				ac.updateReceiveQueueConnectionError()
				return false
			}
			break
		}
		if data == nil && header.SizeBody > 0 {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("read incorrect packet, length %d, but data is nil", header.SizeBody))
			}
			break
		}
		// decrypt if needed
		if ac.cryptMode != CRYPT_MODE_NONE {
			data = ac.crypt.DeCryptPacket(data)
		}
		// check ping/service packet
		if header.NumberPacket > MAX_CLIENT_COMMAND {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("get service packet, number: %d, length %d", header.NumberPacket, header.SizeBody))
			}
			continue
		}
		// update receive queue
		ac.updateReceiveQueue(header, data)
	}
	return true
}

// updateReceiveQueueConnectionError notifies all pending receivers of a connection error.
func (ac *mtAsyncConnect) updateReceiveQueueConnectionError() {
	ac.lockerQueueReceive.Lock()
	for _, received := range ac.queueReceive {
		received.IsFinished = true
		if received.Callback != nil {
			received.Callback(nil, received.Done)
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("run callback, #%d, command: %s, network error", received.ID, received.Command))
			}
		}
	}
	ac.lockerQueueReceive.Unlock()
	ac.receiveRemoveAll()
}

// updateReceiveQueue updates the receive queue with data from a received packet.
func (ac *mtAsyncConnect) updateReceiveQueue(header *protocolHeader, data []byte) {
	if header == nil {
		return
	}
	ac.lockerQueueReceive.Lock()
	defer ac.lockerQueueReceive.Unlock()

	received, exists := ac.queueReceive[header.NumberPacket]
	if !exists {
		return
	}
	// add data
	received.Data = append(received.Data, data)
	received.DataLength += int64(len(data))
	// check finished flag
	if header.Flag == 0 {
		received.IsFinished = true
	}
	if received.IsFinished {
		// combine all data into one buffer
		allData := CopyBufferList(received.Data, received.DataLength)
		// run callback
		if received.Callback != nil {
			received.Callback(allData, received.Done)
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("run callback, #%d, command: %s, total data: %d bytes", received.ID, received.Command, len(allData)))
			}
		}
		// remove from queue
		delete(ac.queueReceive, received.ID)
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, fmt.Sprintf("remove element from receive id: #%d, command: %s", received.ID, received.Command))
		}
	} else {
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, fmt.Sprintf("update receive id: #%d, command: %s, while %d bytes, %d packets", received.ID, received.Command, received.DataLength, len(received.Data)))
		}
	}
}

// readHeaderPacket reads a header+body packet from the socket.
func (ac *mtAsyncConnect) readHeaderPacket() ([]byte, *protocolHeader) {
	headerData, err := ac.socket.Recv(HEADER_LENGTH, 0)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("get data failed: %v", err))
		return nil, nil
	}
	if headerData == nil || len(headerData) < HEADER_LENGTH {
		logWrite(LogTypeError, fmt.Sprintf("incorrect header length, %d", len(headerData)))
		return nil, nil
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("header: %s, length: %d", string(headerData), len(headerData)))
	}
	header := parseHeader(headerData)
	if header == nil {
		logWrite(LogTypeError, "incorrect header data")
		return nil, nil
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("size body: %d, number package: %d, flag: %d", header.SizeBody, header.NumberPacket, header.Flag))
	}
	if header.SizeBody <= 0 {
		return nil, header
	}
	data, err := ac.socket.Recv(header.SizeBody, ac.timeout)
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

// threadSendData is the goroutine that sends data to the MT server.
func (ac *mtAsyncConnect) threadSendData() {
	for {
		select {
		case <-ac.stopSignal:
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, "the send thread is stopped")
			}
			return
		case <-time.After(waitThreadTimeout):
			ac.sendToMetaTraderServer()
		}
	}
}

// sendToMetaTraderServer processes the send queue.
func (ac *mtAsyncConnect) sendToMetaTraderServer() {
	for {
		ds := ac.getFirstSend()
		if ds == nil {
			break
		}
		ac.lastDateSend = time.Now()
		query := ac.getQuerySend(ds)
		if query == nil {
			continue
		}
		if ac.socket == nil || !ac.socket.IsConnected() {
			logWrite(LogTypeError, "send data: connection closed")
			if ac.sendRemove(ds) {
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, fmt.Sprintf("remove element from send id: #%d, command: %s", ds.ID, ds.Command))
				}
			}
			return
		}
		if err := ac.socket.Send(query, len(query), ac.timeout); err != nil {
			logWrite(LogTypeError, fmt.Sprintf("send data failed, %s. %v", GetString(query), err))
			if ac.sendRemove(ds) {
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, fmt.Sprintf("remove element from send id: #%d, command: %s", ds.ID, ds.Command))
				}
			}
			continue
		}
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, fmt.Sprintf("data sent, #%d, command: %s", ds.ID, ds.Command))
		}
		if ac.sendRemove(ds) {
			if IsWriteDebugLog() {
				logWrite(LogTypeDebug, fmt.Sprintf("remove element from send id: #%d, command: %s", ds.ID, ds.Command))
			}
		}
	}
	// check maybe send ping
	if ac.lastDateSend.Add(time.Duration(periodPingSend) * time.Second).Before(time.Now()) {
		ac.sendPingCommand()
	}
}

// getQuerySend builds the raw bytes for sending a packet.
func (ac *mtAsyncConnect) getQuerySend(ds *dataSend) []byte {
	queryBody := ds.Data
	// check large packet
	if queryBody != nil && len(queryBody) > maxPacketSize {
		return ac.getQuerySendLarge(ds)
	}
	// encrypt body if needed
	if ac.cryptMode == CRYPT_MODE_AES && ac.crypt != nil {
		queryBody = ac.crypt.CryptPacket(queryBody)
	}
	headerStr := fmt.Sprintf(WEB_PACKET_FORMAT, len(queryBody), ds.ID, "0")
	header := []byte(headerStr)
	query := append(header, queryBody...)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("sending data: %s%s, length data: %d", headerStr, GetString(queryBody), len(query)))
	}
	return query
}

// getQuerySendLarge builds multiple packets for large data.
func (ac *mtAsyncConnect) getQuerySendLarge(ds *dataSend) []byte {
	data := ds.Data
	var packets [][]byte
	for dataPos := 0; dataPos < len(data); {
		packetBodySize := len(data) - dataPos
		if packetBodySize > maxPacketSize {
			packetBodySize = maxPacketSize
		}
		packetBody := make([]byte, packetBodySize)
		copy(packetBody, data[dataPos:dataPos+packetBodySize])
		dataPos += packetBodySize
		// encrypt body if needed
		if ac.cryptMode == CRYPT_MODE_AES && ac.crypt != nil {
			packetBody = ac.crypt.CryptPacket(packetBody)
		}
		flag := "0"
		if dataPos < len(data) {
			flag = "1"
		}
		headerStr := fmt.Sprintf(WEB_PACKET_FORMAT, len(packetBody), ds.ID, flag)
		header := []byte(headerStr)
		packet := append(header, packetBody...)
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, fmt.Sprintf("sending data: %s%s, length data: %d", headerStr, GetString(packetBody), len(packet)))
		}
		packets = append(packets, packet)
	}
	if len(packets) == 1 {
		return packets[0]
	}
	// combine all packets
	totalLen := 0
	for _, p := range packets {
		totalLen += len(p)
	}
	query := make([]byte, 0, totalLen)
	for _, p := range packets {
		query = append(query, p...)
	}
	return query
}

// sendPingCommand sends a ping to the MT server to keep the connection alive.
func (ac *mtAsyncConnect) sendPingCommand() {
	ac.lastDateSend = time.Now()
	number := int(atomic.AddInt32(&ac.numberClientCommand, 1))
	headerStr := fmt.Sprintf(WEB_PACKET_FORMAT, 0, number, "0")
	header := []byte(headerStr)
	if ac.socket == nil || !ac.socket.IsConnected() {
		return
	}
	if err := ac.socket.Send(header, len(header), ac.timeout); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("send ping failed, %s. %v", headerStr, err))
		return
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("ping sent, %s", headerStr))
	}
}

// --- queue helper methods ---

func (ac *mtAsyncConnect) getFirstSend() *dataSend {
	ac.lockerQueueSend.RLock()
	defer ac.lockerQueueSend.RUnlock()
	if len(ac.queueSend) > 0 {
		return ac.queueSend[0]
	}
	return nil
}

func (ac *mtAsyncConnect) getFirstReceive() *dataReceived {
	ac.lockerQueueReceive.RLock()
	defer ac.lockerQueueReceive.RUnlock()
	var first *dataReceived
	for _, v := range ac.queueReceive {
		if first == nil || v.ID < first.ID {
			first = v
		}
	}
	return first
}

func (ac *mtAsyncConnect) sendRemove(ds *dataSend) bool {
	if ds == nil {
		return false
	}
	ac.lockerQueueSend.Lock()
	defer ac.lockerQueueSend.Unlock()
	for i, item := range ac.queueSend {
		if item == ds {
			ac.queueSend = append(ac.queueSend[:i], ac.queueSend[i+1:]...)
			return true
		}
	}
	return false
}

func (ac *mtAsyncConnect) sendRemoveAll() {
	ac.lockerQueueSend.Lock()
	defer ac.lockerQueueSend.Unlock()
	if len(ac.queueSend) > 0 {
		ac.queueSend = ac.queueSend[:0]
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "remove all element from send queue")
		}
	}
}

func (ac *mtAsyncConnect) receiveRemove(dr *dataReceived) bool {
	if dr == nil {
		return false
	}
	ac.lockerQueueReceive.Lock()
	defer ac.lockerQueueReceive.Unlock()
	if _, exists := ac.queueReceive[dr.ID]; exists {
		delete(ac.queueReceive, dr.ID)
		return true
	}
	return false
}

func (ac *mtAsyncConnect) receiveRemoveAll() {
	ac.lockerQueueReceive.Lock()
	defer ac.lockerQueueReceive.Unlock()
	if len(ac.queueReceive) > 0 {
		for _, dr := range ac.queueReceive {
			if dr.Callback != nil {
				dr.Callback(nil, dr.Done)
				if IsWriteDebugLog() {
					logWrite(LogTypeDebug, fmt.Sprintf("run callback, #%d, command: %s, total data: nil bytes", dr.ID, dr.Command))
				}
			}
		}
		ac.queueReceive = make(map[int]*dataReceived)
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "remove all element from receive queue")
		}
	}
}

// getDataSend builds the UTF-16LE encoded body for a command with parameters.
func getDataSend(command string, params map[string]string) []byte {
	queryTemp := command
	if params != nil && len(params) > 0 {
		bodyRequest := ""
		queryTemp += "|"
		for key, val := range params {
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
	return encodeUTF16LE(queryTemp)
}
