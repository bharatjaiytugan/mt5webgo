//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// serverHandler implements server management protocol commands.
type serverHandler struct {
	apiBase
}

// ServerRestart sends a server restart command.
func (h *serverHandler) ServerRestart() MTRetCode {
	answer := h.send(WEB_CMD_SERVER_RESTART, nil)
	if answer == nil {
		logWrite(LogTypeError, "send server restart failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result server restart answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SERVER_RESTART, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server restart failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}
