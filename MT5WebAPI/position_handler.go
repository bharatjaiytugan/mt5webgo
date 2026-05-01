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

// positionHandler implements position-related protocol commands.
type positionHandler struct {
	apiBase
}

// PositionGet retrieves a position by login and symbol.
func (h *positionHandler) PositionGet(login uint64, symbol string) (*MTPosition, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_SYMBOL: symbol,
	}
	answer := h.send(WEB_CMD_POSITION_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send position get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_POSITION_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	position := parsePositionFromJSON(jsonStr)
	if position == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return position, MT_RET_OK
}

// PositionGetTotal retrieves the total number of positions for a user.
func (h *positionHandler) PositionGetTotal(login uint64) (uint32, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_POSITION_GET_TOTAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send positions total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_POSITION_GET_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.ParseUint(totalStr, 10, 32)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return uint32(total), MT_RET_OK
}

// PositionGetPage retrieves positions by page for a user.
func (h *positionHandler) PositionGetPage(login uint64, offset, total uint32) ([]*MTPosition, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_POSITION_GET_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send positions page get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_POSITION_GET_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse positions page get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	positions := parsePositionsFromJSON(jsonStr)
	if positions == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return positions, MT_RET_OK
}
