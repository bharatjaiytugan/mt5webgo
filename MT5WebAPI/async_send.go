//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"time"
)

// timeoutGetData is the timeout for receiving data from the MT server (10 minutes).
const timeoutGetData = 10 * 60 * time.Second

// mtAsyncSend handles asynchronous sending and receiving of data.
type mtAsyncSend struct {
	done   chan struct{}
	answer []byte
}

// newMTAsyncSend creates a new async send instance.
func newMTAsyncSend() *mtAsyncSend {
	return &mtAsyncSend{}
}

// BeginSend initiates sending a command to the MT server via the async connection.
func (s *mtAsyncSend) BeginSend(connect *mtAsyncConnect, command string, params map[string]string, isNeedReceived bool) bool {
	if connect == nil {
		logWrite(LogTypeError, "connection is nil")
		return false
	}
	s.done = make(chan struct{}, 1)
	s.answer = nil
	// create queue data
	qd := newQueueData()
	qd.Command = command
	qd.IsReceive = isNeedReceived
	qd.Done = s.done
	qd.Callback = s.receiveData
	qd.Params = params

	if connect.Add(qd) {
		if !isNeedReceived {
			// no answer needed, signal done immediately
			select {
			case s.done <- struct{}{}:
			default:
			}
		}
		return true
	}
	// data not added to queue, signal done
	select {
	case s.done <- struct{}{}:
	default:
	}
	return false
}

// EndSend waits for the send/receive operation to complete and returns the answer.
func (s *mtAsyncSend) EndSend() []byte {
	select {
	case <-s.done:
	case <-time.After(timeoutGetData):
	}
	return s.answer
}

// receiveData is the callback invoked when data is received from the server.
func (s *mtAsyncSend) receiveData(data []byte, done chan struct{}) {
	s.answer = data
	if done != nil {
		select {
		case done <- struct{}{}:
		default:
		}
	} else {
		logWrite(LogTypeError, "event is nil")
	}
}
