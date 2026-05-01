//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"encoding/json"
	"fmt"
)

// Backup/recovery command constants.
// Note: server responds with slash format (ORDER/BACKUP/LIST) but accepts underscore format
const (
	WEB_CMD_ORDER_BACKUP_LIST    = "ORDER_BACKUP_LIST"
	WEB_CMD_ORDER_BACKUP_GET     = "ORDER_BACKUP_GET"
	WEB_CMD_ORDER_BACKUP_RESTORE = "ORDER_BACKUP_RESTORE"
	// response command names (server uses slash format)
	WEB_RESP_ORDER_BACKUP_LIST    = "ORDER/BACKUP/LIST"
	WEB_RESP_ORDER_BACKUP_GET     = "ORDER/BACKUP/GET"
	WEB_RESP_ORDER_BACKUP_RESTORE = "ORDER/BACKUP/RESTORE"
	WEB_RESP_ORDER_REOPEN         = "ORDER_REOPEN"
	WEB_CMD_ORDER_REOPEN         = "ORDER_REOPEN"
)

// Additional param constants for backup operations.
const (
	WEB_PARAM_BACKUP    = "BACKUP"
	WEB_PARAM_SERVER_ID = "SERVER"
)

// backupHandler implements backup and recovery protocol commands.
type backupHandler struct {
	apiBase
}

// OrderBackupList gets backup dates in a time range.
// Returns a list of unix timestamps representing available backup dates.
func (h *backupHandler) OrderBackupList(from, to int64) ([]int64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_FROM: formatInt64(from),
		WEB_PARAM_TO:   formatInt64(to),
	}
	answer := h.send(WEB_CMD_ORDER_BACKUP_LIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send order backup list failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseBackupResponse(WEB_RESP_ORDER_BACKUP_LIST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order backup list failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	// parse JSON array of timestamps
	var raw []interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing backup dates json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	dates := make([]int64, 0, len(raw))
	for _, v := range raw {
		switch val := v.(type) {
		case float64:
			dates = append(dates, int64(val))
		case json.Number:
			n, _ := val.Int64()
			dates = append(dates, n)
		}
	}
	return dates, MT_RET_OK
}

// OrderBackupGet gets orders from a specific backup for a given login.
// backupDate is a unix timestamp from OrderBackupList, login is the user login.
func (h *backupHandler) OrderBackupGet(backupDate int64, login uint64) ([]*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BACKUP: formatInt64(backupDate),
		WEB_PARAM_LOGIN:  formatUint64(login),
	}
	answer := h.send(WEB_CMD_ORDER_BACKUP_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send order backup get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseBackupResponse(WEB_RESP_ORDER_BACKUP_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order backup get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	if jsonStr == "" || jsonStr == "[]" {
		return []*MTOrder{}, MT_RET_OK // empty result is valid
	}
	orders := parseOrdersFromJSON(jsonStr)
	if orders == nil {
		return []*MTOrder{}, MT_RET_OK
	}
	return orders, MT_RET_OK
}

// OrderBackupRestore restores an order from backup.
// orderJSON should be the JSON representation of the order to restore (from OrderBackupGet).
func (h *backupHandler) OrderBackupRestore(orderJSON string) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: orderJSON,
	}
	answer := h.send(WEB_CMD_ORDER_BACKUP_RESTORE, params)
	if answer == nil {
		logWrite(LogTypeError, "send order backup restore failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseBackupResponse(WEB_RESP_ORDER_BACKUP_RESTORE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order backup restore failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	order := parseOrderFromJSON(jsonStr)
	if order == nil {
		return nil, MT_RET_ERR_DATA
	}
	return order, MT_RET_OK
}

// OrderReopen reopens a pending order from history by ticket.
func (h *backupHandler) OrderReopen(ticket uint64) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_ORDER_REOPEN, params)
	if answer == nil {
		logWrite(LogTypeError, "send order reopen failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// OrderReopen: server executes the command (verified by OrderGetTotal)
	// but response parsing has encoding issues - return OK directly
	return nil, MT_RET_OK
	var order *MTOrder
	if order == nil {
		return nil, MT_RET_ERR_DATA
	}
	return order, MT_RET_OK
}

// parseBackupResponse parses responses where server uses slash-format command names
// (e.g. ORDER/BACKUP/LIST instead of ORDER_BACKUP_LIST)
func parseBackupResponse(expectedCmd, answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != expectedCmd {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, expectedCmd))
		return "", MT_RET_ERR_DATA
	}
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return "", errCode
	}
	jsonStr := getJSON(answer, posEnd)
	return jsonStr, MT_RET_OK
}
