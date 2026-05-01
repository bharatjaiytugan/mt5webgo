//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// customHandler implements the custom command protocol handler.
type customHandler struct {
	apiBase
}

// CustomSend sends a custom command to the MT server.
func (h *customHandler) CustomSend(command string, params map[string]string, body string) ([]byte, MTRetCode) {
	if params == nil {
		params = make(map[string]string)
	}
	if _, ok := params[WEB_PARAM_BODYTEXT]; !ok {
		params[WEB_PARAM_BODYTEXT] = body
	}
	answer := h.send(command, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", command))
		return nil, MT_RET_ERR_NETWORK
	}
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", GetString(answer)))
	}
	return answer, MT_RET_OK
}
