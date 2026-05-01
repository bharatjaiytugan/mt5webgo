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

// Batch operation command constants.
const (
	WEB_CMD_GROUP_ADD_BATCH    = "GROUP_ADD_BATCH"
	WEB_CMD_GROUP_DELETE_BATCH = "GROUP_DELETE_BATCH"
	WEB_CMD_GROUP_LIST         = "GROUP_LIST"
	WEB_CMD_SYMBOL_ADD_BATCH   = "SYMBOL_ADD_BATCH"
	WEB_CMD_SYMBOL_DELETE_BATCH = "SYMBOL_DELETE_BATCH"
	WEB_CMD_ORDER_GET_BATCH    = "ORDER_GET_BATCH"
	WEB_CMD_ORDER_UPDATE       = "ORDER_UPDATE"
	WEB_CMD_ORDER_DELETE_DIRECT = "ORDER_DELETE"  // direct delete (not via dealer)
	WEB_CMD_ORDER_CANCEL       = "ORDER_CANCEL"   // move to history
	WEB_CMD_HISTORY_GET_BATCH  = "HISTORY_GET_BATCH"
	WEB_CMD_HISTORY_UPDATE     = "HISTORY_UPDATE"
	WEB_CMD_HISTORY_DELETE     = "HISTORY_DELETE"
	WEB_CMD_DEAL_GET_BATCH     = "DEAL_GET_BATCH"
	WEB_CMD_DEAL_UPDATE        = "DEAL_UPDATE"
	WEB_CMD_DEAL_DELETE        = "DEAL_DELETE"
	WEB_CMD_POSITION_GET_BATCH = "POSITION_GET_BATCH"
	WEB_CMD_USER_GET_BATCH     = "USER_GET_BATCH"
	WEB_CMD_USER_BACKUP_LIST   = "USER_BACKUP_LIST"
	WEB_CMD_USER_BACKUP_GET    = "USER_BACKUP_GET"
)

// batchHandler implements batch operation protocol commands.
type batchHandler struct {
	apiBase
}

// GroupList retrieves a list of group names matching the given mask.
func (h *batchHandler) GroupList(mask string) ([]string, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_GROUP: mask,
	}
	answer := h.send(WEB_CMD_GROUP_LIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send group list failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_GROUP_LIST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group list failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	var groups []string
	if err := json.Unmarshal([]byte(jsonStr), &groups); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing group list json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return groups, MT_RET_OK
}

// GroupAddBatch adds multiple groups at once from a JSON array of group configs.
func (h *batchHandler) GroupAddBatch(groupsJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: groupsJSON,
	}
	answer := h.send(WEB_CMD_GROUP_ADD_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send group add batch failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_GROUP_ADD_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group add batch failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// GroupDeleteBatch deletes multiple groups by comma-separated names.
func (h *batchHandler) GroupDeleteBatch(groups string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_GROUP: groups,
	}
	answer := h.send(WEB_CMD_GROUP_DELETE_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send group delete batch failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_GROUP_DELETE_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group delete batch failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SymbolAddBatch adds multiple symbols at once from a JSON array of symbol configs.
func (h *batchHandler) SymbolAddBatch(symbolsJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: symbolsJSON,
	}
	answer := h.send(WEB_CMD_SYMBOL_ADD_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol add batch failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SYMBOL_ADD_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol add batch failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SymbolDeleteBatch deletes multiple symbols by comma-separated names.
func (h *batchHandler) SymbolDeleteBatch(symbols string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbols,
	}
	answer := h.send(WEB_CMD_SYMBOL_DELETE_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol delete batch failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SYMBOL_DELETE_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol delete batch failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// OrderGetBatch retrieves orders matching the given logins, groups, and tickets (comma-separated).
func (h *batchHandler) OrderGetBatch(logins string, groups string, tickets string) ([]*MTOrder, MTRetCode) {
	params := map[string]string{}
	if logins != "" {
		params[WEB_PARAM_LOGIN] = logins
	}
	if groups != "" {
		params[WEB_PARAM_GROUP] = groups
	}
	if tickets != "" {
		params[WEB_PARAM_TICKET] = tickets
	}
	answer := h.send(WEB_CMD_ORDER_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send order get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_ORDER_GET_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	orders := parseOrdersFromJSON(jsonStr)
	if orders == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return orders, MT_RET_OK
}

// OrderUpdateDirect updates an order directly (identified by ticket in JSON body).
func (h *batchHandler) OrderUpdateDirect(orderJSON string) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: orderJSON,
	}
	answer := h.send(WEB_CMD_ORDER_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send order update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_ORDER_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	order := parseOrderFromJSON(jsonStr)
	if order == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return order, MT_RET_OK
}

// OrderDeleteDirect deletes orders directly by comma-separated tickets.
func (h *batchHandler) OrderDeleteDirect(tickets string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TICKET: tickets,
	}
	answer := h.send(WEB_CMD_ORDER_DELETE_DIRECT, params)
	if answer == nil {
		logWrite(LogTypeError, "send order delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_ORDER_DELETE_DIRECT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// OrderCancel cancels orders by comma-separated tickets, moving them to history with CANCELED state.
func (h *batchHandler) OrderCancel(tickets string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TICKET: tickets,
	}
	answer := h.send(WEB_CMD_ORDER_CANCEL, params)
	if answer == nil {
		logWrite(LogTypeError, "send order cancel failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_ORDER_CANCEL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse order cancel failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// HistoryGetBatch retrieves historical orders matching the given filters.
func (h *batchHandler) HistoryGetBatch(logins string, groups string, tickets string, from, to int64) ([]*MTOrder, MTRetCode) {
	params := map[string]string{}
	if logins != "" {
		params[WEB_PARAM_LOGIN] = logins
	}
	if groups != "" {
		params[WEB_PARAM_GROUP] = groups
	}
	if tickets != "" {
		params[WEB_PARAM_TICKET] = tickets
	}
	params[WEB_PARAM_FROM] = formatInt64(from)
	params[WEB_PARAM_TO] = formatInt64(to)
	answer := h.send(WEB_CMD_HISTORY_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send history get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_HISTORY_GET_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse history get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	orders := parseOrdersFromJSON(jsonStr)
	if orders == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return orders, MT_RET_OK
}

// HistoryUpdateDirect updates a historical order directly (identified by ticket in JSON body).
func (h *batchHandler) HistoryUpdateDirect(orderJSON string) (*MTOrder, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: orderJSON,
	}
	answer := h.send(WEB_CMD_HISTORY_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send history update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_HISTORY_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse history update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	order := parseOrderFromJSON(jsonStr)
	if order == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return order, MT_RET_OK
}

// HistoryDeleteDirect deletes historical orders by comma-separated tickets.
func (h *batchHandler) HistoryDeleteDirect(tickets string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TICKET: tickets,
	}
	answer := h.send(WEB_CMD_HISTORY_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send history delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_HISTORY_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse history delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// DealGetBatch retrieves deals matching the given logins, groups, and tickets (comma-separated) within a time range.
func (h *batchHandler) DealGetBatch(logins string, groups string, tickets string, from, to int64) ([]*MTDeal, MTRetCode) {
	params := map[string]string{}
	if logins != "" {
		params[WEB_PARAM_LOGIN] = logins
	}
	if groups != "" {
		params[WEB_PARAM_GROUP] = groups
	}
	if tickets != "" {
		params[WEB_PARAM_TICKET] = tickets
	}
	params[WEB_PARAM_FROM] = formatInt64(from)
	params[WEB_PARAM_TO] = formatInt64(to)
	answer := h.send(WEB_CMD_DEAL_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseDealBatchResponse(WEB_CMD_DEAL_GET_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	deals := parseDealsFromJSON(jsonStr)
	if deals == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return deals, MT_RET_OK
}

// DealUpdateDirect updates a deal directly (identified by Deal ticket in JSON body).
func (h *batchHandler) DealUpdateDirect(dealJSON string) (*MTDeal, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: dealJSON,
	}
	answer := h.send(WEB_CMD_DEAL_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseDealBatchResponse(WEB_CMD_DEAL_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	deal := parseDealFromJSON(jsonStr)
	if deal == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return deal, MT_RET_OK
}

// DealDeleteDirect deletes deals by comma-separated tickets.
func (h *batchHandler) DealDeleteDirect(tickets string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TICKET: tickets,
	}
	answer := h.send(WEB_CMD_DEAL_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseDealEmptyResult(WEB_CMD_DEAL_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// parseDealBatchResponse parses deal batch responses, accepting both underscore and slash command formats.
// The server may respond with DEAL_GET_BATCH or DEAL/GET_BATCH format.
func parseDealBatchResponse(expectedCmd, answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != expectedCmd {
		// accept any command name the server returns (may use slash format)
		logWrite(LogTypeDebug, fmt.Sprintf("deal batch: server returned command '%s', expected '%s' - accepting", cmdReal, expectedCmd))
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
	if jsonStr == "" {
		return "", MT_RET_REPORT_NODATA
	}
	return jsonStr, MT_RET_OK
}

// parseDealEmptyResult parses a deal batch response that only contains RETCODE (no JSON body).
// Accepts both underscore and slash command formats from the server.
func (h *batchHandler) parseDealEmptyResult(expectedCmd, answer string) MTRetCode {
	cmdReal, pos := getCommand(answer)
	if cmdReal != expectedCmd {
		logWrite(LogTypeDebug, fmt.Sprintf("deal batch: server returned command '%s', expected '%s' - accepting", cmdReal, expectedCmd))
	}
	ba := newBaseAnswer()
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			ba.RetCode = param.Value
		}
	}
	return getRetCode(ba.RetCode)
}

// PositionGetBatch gets positions by logins, groups, or tickets.
func (h *batchHandler) PositionGetBatch(logins, groups, tickets string) ([]*MTPosition, MTRetCode) {
	params := map[string]string{}
	if logins != "" {
		params[WEB_PARAM_LOGIN] = logins
	}
	if groups != "" {
		params[WEB_PARAM_GROUP] = groups
	}
	if tickets != "" {
		params[WEB_PARAM_TICKET] = tickets
	}
	answer := h.send(WEB_CMD_POSITION_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send position get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" || jsonStr == "[]" {
		return []*MTPosition{}, MT_RET_OK
	}
	positions := parsePositionsFromJSON(jsonStr)
	if positions == nil {
		return []*MTPosition{}, MT_RET_OK
	}
	return positions, MT_RET_OK
}

// UserGetBatch gets users by logins or groups (comma-separated).
func (h *batchHandler) UserGetBatch(logins, groups string) ([]*MTUser, MTRetCode) {
	params := map[string]string{}
	if logins != "" {
		params[WEB_PARAM_LOGIN] = logins
	}
	if groups != "" {
		params[WEB_PARAM_GROUP] = groups
	}
	answer := h.send(WEB_CMD_USER_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send user get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format (server may use slash format)
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" || jsonStr == "[]" {
		return []*MTUser{}, MT_RET_OK
	}
	users := parseUsersFromJSON(jsonStr)
	if users == nil {
		return []*MTUser{}, MT_RET_OK
	}
	return users, MT_RET_OK
}

// UserBackupList gets user backup dates in a time range.
// Returns a list of unix timestamps representing available backup dates.
func (h *batchHandler) UserBackupList(from, to int64) ([]int64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_FROM: formatInt64(from),
		WEB_PARAM_TO:   formatInt64(to),
	}
	answer := h.send(WEB_CMD_USER_BACKUP_LIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send user backup list failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format (server may respond USER/BACKUP/LIST)
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user backup list failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" || jsonStr == "[]" {
		return []int64{}, MT_RET_OK
	}
	// parse JSON array of timestamps
	var raw []interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing user backup dates json failed: %v", err))
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

// UserBackupGet gets users from a specific backup for a given login.
// backupDate is a unix timestamp from UserBackupList, login is the user login.
func (h *batchHandler) UserBackupGet(backupDate int64, login uint64) ([]*MTUser, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BACKUP: formatInt64(backupDate),
		WEB_PARAM_LOGIN:  formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_BACKUP_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user backup get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format (server may respond USER/BACKUP/GET)
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user backup get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" || jsonStr == "[]" {
		return []*MTUser{}, MT_RET_OK
	}
	users := parseUsersFromJSON(jsonStr)
	if users == nil {
		return []*MTUser{}, MT_RET_OK
	}
	return users, MT_RET_OK
}
