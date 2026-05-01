//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"strconv"
)

// historyHandler implements history order-related protocol commands.
type historyHandler struct {
	apiBase
}

// HistoryGet retrieves a history order by ticket.
func (h *historyHandler) HistoryGet(ticket uint64) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_HISTORY_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send history order get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_HISTORY_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	order := parseOrderFromJSON(jsonStr)
	if order == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return order, MT_RET_OK
}

// HistoryGetTotal retrieves the total number of history orders for a user in a time range.
func (h *historyHandler) HistoryGetTotal(login uint64, from, to int64) (uint32, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_FROM:  formatInt64(from),
		WEB_PARAM_TO:    formatInt64(to),
	}
	answer := h.send(WEB_CMD_HISTORY_GET_TOTAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send history order total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_HISTORY_GET_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse history order total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.ParseUint(totalStr, 10, 32)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return uint32(total), MT_RET_OK
}

// HistoryGetPage retrieves history orders by page for a user in a time range.
func (h *historyHandler) HistoryGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_FROM:   formatInt64(from),
		WEB_PARAM_TO:     formatInt64(to),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_HISTORY_GET_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send history orders page get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_HISTORY_GET_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse history orders page get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	orders := parseOrdersFromJSON(jsonStr)
	if orders == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return orders, MT_RET_OK
}
