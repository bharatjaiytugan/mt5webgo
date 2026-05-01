//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// tradeHandler implements trade-related protocol commands.
type tradeHandler struct {
	apiBase
}

// TradeBalance changes a user's balance.
func (h *tradeHandler) TradeBalance(login uint64, dealType EnDealAction, balance float64, comment string) MTRetCode {
	_, ret := h.TradeBalanceWithTicket(login, dealType, balance, comment)
	return ret
}

// TradeBalanceWithTicket changes a user's balance and returns the resulting deal ticket.
func (h *tradeHandler) TradeBalanceWithTicket(login uint64, dealType EnDealAction, balance float64, comment string) (uint64, MTRetCode) {
	return h.TradeBalanceWithMarginCheck(login, dealType, balance, comment, true)
}

// TradeBalanceWithMarginCheck changes a user's balance with optional margin checking,
// and returns the resulting deal ticket.
func (h *tradeHandler) TradeBalanceWithMarginCheck(login uint64, dealType EnDealAction, balance float64, comment string, checkMargin bool) (uint64, MTRetCode) {
	checkMarginStr := "1"
	if !checkMargin {
		checkMarginStr = "0"
	}
	params := map[string]string{
		WEB_PARAM_LOGIN:        formatUint64(login),
		WEB_PARAM_TYPE:         formatUint32(uint32(dealType)),
		WEB_PARAM_BALANCE:      formatFloat64(balance),
		WEB_PARAM_COMMENT:      comment,
		WEB_PARAM_CHECK_MARGIN: checkMarginStr,
	}
	answer := h.send(WEB_CMD_TRADE_BALANCE, params)
	if answer == nil {
		logWrite(LogTypeError, "send trade balance failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer trade balance: %s", answerStr))
	}
	ticket, errCode := parseTicketAnswer(WEB_CMD_TRADE_BALANCE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse trade balance failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	return ticket, MT_RET_OK
}
