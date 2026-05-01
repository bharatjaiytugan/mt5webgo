//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"sync"
)

// Log type constants.
const (
	LogTypeError = 1
	LogTypeDebug = 2
)

// CallbackLogWrite is a user callback function for writing logs.
type CallbackLogWrite func(logType int, message string)

var (
	logAgent    string
	logCallback CallbackLogWrite
	logMu       sync.RWMutex
	logStatus   int = LogTypeError
)

// IsWriteDebugLog returns true if debug logging is enabled.
func IsWriteDebugLog() bool {
	logMu.RLock()
	defer logMu.RUnlock()
	return (logStatus & LogTypeDebug) > 0
}

// SetWriteDebugLog enables or disables debug logging.
func SetWriteDebugLog(enabled bool) {
	logMu.Lock()
	defer logMu.Unlock()
	if enabled {
		logStatus |= LogTypeDebug
	} else {
		logStatus &^= LogTypeDebug
	}
}

// logInit initializes the logging subsystem with an agent name and callback.
func logInit(agent string, callback CallbackLogWrite) {
	logMu.Lock()
	defer logMu.Unlock()
	logAgent = agent
	logCallback = callback
}

// logWrite writes a log message if the callback is set and the log type is enabled.
func logWrite(logType int, message string) {
	logMu.Lock()
	defer logMu.Unlock()
	if logCallback == nil {
		return
	}
	if (logType & logStatus) > 0 {
		logCallback(logType, fmt.Sprintf("%s: %s", logAgent, message))
	}
}
