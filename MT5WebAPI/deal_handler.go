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

// dealHandler implements deal-related protocol commands.
type dealHandler struct {
	apiBase
}

// DealGet retrieves a deal by ticket.
func (h *dealHandler) DealGet(ticket uint64) (*MTDeal, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_DEAL_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_DEAL_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	deal := parseDealFromJSON(jsonStr)
	if deal == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return deal, MT_RET_OK
}

// DealGetTotal retrieves the total number of deals for a user in a time range.
func (h *dealHandler) DealGetTotal(login uint64, from, to int64) (uint32, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_FROM:  formatInt64(from),
		WEB_PARAM_TO:    formatInt64(to),
	}
	answer := h.send(WEB_CMD_DEAL_GET_TOTAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_DEAL_GET_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.ParseUint(totalStr, 10, 32)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return uint32(total), MT_RET_OK
}

// DealGetPage retrieves deals by page for a user in a time range.
func (h *dealHandler) DealGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTDeal, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_FROM:   formatInt64(from),
		WEB_PARAM_TO:     formatInt64(to),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_DEAL_GET_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send deals page get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_DEAL_GET_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deals page get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	deals := parseDealsFromJSON(jsonStr)
	if deals == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return deals, MT_RET_OK
}
