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

// orderHandler implements order-related protocol commands.
type orderHandler struct {
	apiBase
}

// OrderGet retrieves an order by ticket.
func (h *orderHandler) OrderGet(ticket uint64) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_ORDER_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send order get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_ORDER_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	order := parseOrderFromJSON(jsonStr)
	if order == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return order, MT_RET_OK
}

// OrderGetTotal retrieves the total number of orders for a user.
func (h *orderHandler) OrderGetTotal(login uint64) (uint32, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_ORDER_GET_TOTAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send orders total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_ORDER_GET_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.ParseUint(totalStr, 10, 32)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return uint32(total), MT_RET_OK
}

// OrderGetPage retrieves orders by page for a user.
func (h *orderHandler) OrderGetPage(login uint64, offset, total uint32) ([]*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_ORDER_GET_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send orders page get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_ORDER_GET_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse orders page get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	orders := parseOrdersFromJSON(jsonStr)
	if orders == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return orders, MT_RET_OK
}
