//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

// pingHandler implements the ping protocol command.
type pingHandler struct {
	apiBase
}

// PingSend sends a ping (keep-alive) to the server without expecting a response.
func (h *pingHandler) PingSend() MTRetCode {
	h.sendNoRecv("", nil)
	return MT_RET_OK
}
