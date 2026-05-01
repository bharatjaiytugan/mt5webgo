//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// WebSocket command constants for config CRUD entities.
const (
	// Leverage
	WEB_CMD_LEVERAGE_TOTAL  = "LEVERAGE_TOTAL"
	WEB_CMD_LEVERAGE_NEXT   = "LEVERAGE_NEXT"
	WEB_CMD_LEVERAGE_GET    = "LEVERAGE_GET"
	WEB_CMD_LEVERAGE_ADD    = "LEVERAGE_ADD"
	WEB_CMD_LEVERAGE_DELETE = "LEVERAGE_DELETE"
	WEB_CMD_LEVERAGE_SHIFT  = "LEVERAGE_SHIFT"
	// Firewall
	WEB_CMD_FIREWALL_TOTAL   = "FIREWALL_TOTAL"
	WEB_CMD_FIREWALL_NEXT    = "FIREWALL_NEXT"
	WEB_CMD_FIREWALL_GET_ALL = "FIREWALL_GET_ALL"
	WEB_CMD_FIREWALL_ADD     = "FIREWALL_ADD"
	WEB_CMD_FIREWALL_DELETE  = "FIREWALL_DELETE"
	WEB_CMD_FIREWALL_SHIFT   = "FIREWALL_SHIFT"
	// Holiday
	WEB_CMD_HOLIDAY_TOTAL   = "HOLIDAY_TOTAL"
	WEB_CMD_HOLIDAY_NEXT    = "HOLIDAY_NEXT"
	WEB_CMD_HOLIDAY_GET_ALL = "HOLIDAY_GET_ALL"
	WEB_CMD_HOLIDAY_ADD     = "HOLIDAY_ADD"
	WEB_CMD_HOLIDAY_DELETE  = "HOLIDAY_DELETE"
	WEB_CMD_HOLIDAY_SHIFT   = "HOLIDAY_SHIFT"
	// Manager
	WEB_CMD_MANAGER_TOTAL   = "MANAGER_TOTAL"
	WEB_CMD_MANAGER_NEXT    = "MANAGER_NEXT"
	WEB_CMD_MANAGER_GET     = "MANAGER_GET"
	WEB_CMD_MANAGER_GET_ALL = "MANAGER_GET_ALL"
	WEB_CMD_MANAGER_ADD     = "MANAGER_ADD"
	WEB_CMD_MANAGER_DELETE  = "MANAGER_DELETE"
	WEB_CMD_MANAGER_SHIFT   = "MANAGER_SHIFT"
	// Route
	WEB_CMD_ROUTE_TOTAL   = "ROUTE_TOTAL"
	WEB_CMD_ROUTE_NEXT    = "ROUTE_NEXT"
	WEB_CMD_ROUTE_GET     = "ROUTE_GET"
	WEB_CMD_ROUTE_GET_ALL = "ROUTE_GET_ALL"
	WEB_CMD_ROUTE_ADD     = "ROUTE_ADD"
	WEB_CMD_ROUTE_DELETE  = "ROUTE_DELETE"
	WEB_CMD_ROUTE_SHIFT   = "ROUTE_SHIFT"
	// Spread
	WEB_CMD_SPREAD_TOTAL   = "SPREAD_TOTAL"
	WEB_CMD_SPREAD_NEXT    = "SPREAD_NEXT"
	WEB_CMD_SPREAD_GET_ALL = "SPREAD_GET_ALL"
	WEB_CMD_SPREAD_ADD     = "SPREAD_ADD"
	WEB_CMD_SPREAD_DELETE  = "SPREAD_DELETE"
	WEB_CMD_SPREAD_SHIFT   = "SPREAD_SHIFT"
	// HistorySync
	WEB_CMD_HISTORY_SYNC_TOTAL   = "HISTORY_SYNC_TOTAL"
	WEB_CMD_HISTORY_SYNC_NEXT    = "HISTORY_SYNC_NEXT"
	WEB_CMD_HISTORY_SYNC_GET_ALL = "HISTORY_SYNC_GET_ALL"
	WEB_CMD_HISTORY_SYNC_ADD     = "HISTORY_SYNC_ADD"
	WEB_CMD_HISTORY_SYNC_DELETE  = "HISTORY_SYNC_DELETE"
	WEB_CMD_HISTORY_SYNC_SHIFT   = "HISTORY_SYNC_SHIFT"
	WEB_CMD_HISTORY_SYNC_START   = "HISTORY_SYNC_START"
	// SymbolGroup
	WEB_CMD_SYMBOL_GROUP_TOTAL  = "SYMBOL_GROUP_TOTAL"
	WEB_CMD_SYMBOL_GROUP_NEXT   = "SYMBOL_GROUP_NEXT"
	WEB_CMD_SYMBOL_GROUP_LIST   = "SYMBOL_GROUP_LIST"
	WEB_CMD_SYMBOL_GROUP_ADD    = "SYMBOL_GROUP_ADD"
	WEB_CMD_SYMBOL_GROUP_DELETE = "SYMBOL_GROUP_DELETE"
	WEB_CMD_SYMBOL_GROUP_SHIFT  = "SYMBOL_GROUP_SHIFT"
)

// configCrudHandler implements generic CRUD operations for config entities.
type configCrudHandler struct {
	apiBase
}

// --- Generic CRUD operations ---

// configTotal sends a TOTAL command and returns the count.
func (h *configCrudHandler) configTotal(cmdPrefix string) (int, MTRetCode) {
	cmd := cmdPrefix + "_TOTAL"
	answer := h.send(cmd, nil)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// configNext sends a NEXT command with an index and returns the parsed JSON object.
func (h *configCrudHandler) configNext(cmdPrefix string, index int) (map[string]interface{}, MTRetCode) {
	cmd := cmdPrefix + "_NEXT"
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", cmd, err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// configGet sends a GET command with the provided params and returns the parsed JSON object.
func (h *configCrudHandler) configGet(cmdPrefix string, params map[string]string) (map[string]interface{}, MTRetCode) {
	cmd := cmdPrefix + "_GET"
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", cmd, err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// configGetAll sends a GET_ALL command and returns the parsed JSON array.
func (h *configCrudHandler) configGetAll(cmdPrefix string) ([]map[string]interface{}, MTRetCode) {
	cmd := cmdPrefix + "_GET_ALL"
	answer := h.send(cmd, nil)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", cmd, err))
		return nil, MT_RET_REPORT_NODATA
	}
	var results []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			results = append(results, m)
		}
	}
	return results, MT_RET_OK
}

// configAdd sends an ADD command with the body JSON and returns the parsed JSON response.
func (h *configCrudHandler) configAdd(cmdPrefix string, bodyJSON string) (map[string]interface{}, MTRetCode) {
	if bodyJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	cmd := cmdPrefix + "_ADD"
	params := map[string]string{
		WEB_PARAM_BODYTEXT: bodyJSON,
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", cmd, err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// configDelete sends a DELETE command with the provided params.
func (h *configCrudHandler) configDelete(cmdPrefix string, params map[string]string) MTRetCode {
	cmd := cmdPrefix + "_DELETE"
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
	}
	return errCode
}

// configShift sends a SHIFT command with index and shift params.
func (h *configCrudHandler) configShift(cmdPrefix string, index, shift int) MTRetCode {
	cmd := cmdPrefix + "_SHIFT"
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
		WEB_PARAM_SHIFT: formatInt(shift),
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
	}
	return errCode
}

// --- Leverage (6 methods) ---

// LeverageTotal gets the total count of leverage configs.
func (h *configCrudHandler) LeverageTotal() (int, MTRetCode) {
	return h.configTotal("LEVERAGE")
}

// LeverageNext gets the leverage config at the given index.
func (h *configCrudHandler) LeverageNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("LEVERAGE", index)
}

// LeverageGet gets a leverage config by name.
func (h *configCrudHandler) LeverageGet(name string) (map[string]interface{}, MTRetCode) {
	return h.configGet("LEVERAGE", map[string]string{WEB_PARAM_NAME: name})
}

// LeverageAdd adds a new leverage config.
func (h *configCrudHandler) LeverageAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("LEVERAGE", bodyJSON)
}

// LeverageDelete deletes a leverage config by name.
func (h *configCrudHandler) LeverageDelete(name string) MTRetCode {
	return h.configDelete("LEVERAGE", map[string]string{WEB_PARAM_NAME: name})
}

// LeverageShift shifts a leverage config position.
func (h *configCrudHandler) LeverageShift(index, shift int) MTRetCode {
	return h.configShift("LEVERAGE", index, shift)
}

// --- Firewall (6 methods) ---

// FirewallTotal gets the total count of firewall configs.
func (h *configCrudHandler) FirewallTotal() (int, MTRetCode) {
	return h.configTotal("FIREWALL")
}

// FirewallNext gets the firewall config at the given index.
func (h *configCrudHandler) FirewallNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("FIREWALL", index)
}

// FirewallGetAll gets all firewall configs.
func (h *configCrudHandler) FirewallGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("FIREWALL")
}

// FirewallAdd adds a new firewall config.
func (h *configCrudHandler) FirewallAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("FIREWALL", bodyJSON)
}

// FirewallDelete deletes a firewall config by index.
func (h *configCrudHandler) FirewallDelete(index int) MTRetCode {
	return h.configDelete("FIREWALL", map[string]string{WEB_PARAM_INDEX: formatInt(index)})
}

// FirewallShift shifts a firewall config position.
func (h *configCrudHandler) FirewallShift(index, shift int) MTRetCode {
	return h.configShift("FIREWALL", index, shift)
}

// --- Holiday (6 methods) ---

// HolidayTotal gets the total count of holiday configs.
func (h *configCrudHandler) HolidayTotal() (int, MTRetCode) {
	return h.configTotal("HOLIDAY")
}

// HolidayNext gets the holiday config at the given index.
func (h *configCrudHandler) HolidayNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("HOLIDAY", index)
}

// HolidayGetAll gets all holiday configs.
func (h *configCrudHandler) HolidayGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("HOLIDAY")
}

// HolidayAdd adds a new holiday config.
func (h *configCrudHandler) HolidayAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("HOLIDAY", bodyJSON)
}

// HolidayDelete deletes a holiday config by index.
func (h *configCrudHandler) HolidayDelete(index int) MTRetCode {
	return h.configDelete("HOLIDAY", map[string]string{WEB_PARAM_INDEX: formatInt(index)})
}

// HolidayShift shifts a holiday config position.
func (h *configCrudHandler) HolidayShift(index, shift int) MTRetCode {
	return h.configShift("HOLIDAY", index, shift)
}

// --- Manager (7 methods) ---

// ManagerTotal gets the total count of manager configs.
func (h *configCrudHandler) ManagerTotal() (int, MTRetCode) {
	return h.configTotal("MANAGER")
}

// ManagerNext gets the manager config at the given index.
func (h *configCrudHandler) ManagerNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("MANAGER", index)
}

// ManagerGet gets a manager config by login.
func (h *configCrudHandler) ManagerGet(login uint64) (map[string]interface{}, MTRetCode) {
	return h.configGet("MANAGER", map[string]string{WEB_PARAM_LOGIN: formatUint64(login)})
}

// ManagerGetAll gets all manager configs.
func (h *configCrudHandler) ManagerGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("MANAGER")
}

// ManagerAdd adds a new manager config.
func (h *configCrudHandler) ManagerAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("MANAGER", bodyJSON)
}

// ManagerDelete deletes a manager config by login.
func (h *configCrudHandler) ManagerDelete(login uint64) MTRetCode {
	return h.configDelete("MANAGER", map[string]string{WEB_PARAM_LOGIN: formatUint64(login)})
}

// ManagerShift shifts a manager config position.
func (h *configCrudHandler) ManagerShift(index, shift int) MTRetCode {
	return h.configShift("MANAGER", index, shift)
}

// --- Route (7 methods) ---

// RouteTotal gets the total count of route configs.
func (h *configCrudHandler) RouteTotal() (int, MTRetCode) {
	return h.configTotal("ROUTE")
}

// RouteNext gets the route config at the given index.
func (h *configCrudHandler) RouteNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("ROUTE", index)
}

// RouteGet gets a route config by name.
func (h *configCrudHandler) RouteGet(name string) (map[string]interface{}, MTRetCode) {
	return h.configGet("ROUTE", map[string]string{WEB_PARAM_NAME: name})
}

// RouteGetAll gets all route configs.
func (h *configCrudHandler) RouteGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("ROUTE")
}

// RouteAdd adds a new route config.
func (h *configCrudHandler) RouteAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("ROUTE", bodyJSON)
}

// RouteDelete deletes a route config by name.
func (h *configCrudHandler) RouteDelete(name string) MTRetCode {
	return h.configDelete("ROUTE", map[string]string{WEB_PARAM_NAME: name})
}

// RouteShift shifts a route config position.
func (h *configCrudHandler) RouteShift(index, shift int) MTRetCode {
	return h.configShift("ROUTE", index, shift)
}

// --- Spread (6 methods) ---

// SpreadTotal gets the total count of spread configs.
func (h *configCrudHandler) SpreadTotal() (int, MTRetCode) {
	return h.configTotal("SPREAD")
}

// SpreadNext gets the spread config at the given index.
func (h *configCrudHandler) SpreadNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("SPREAD", index)
}

// SpreadGetAll gets all spread configs.
func (h *configCrudHandler) SpreadGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("SPREAD")
}

// SpreadAdd adds a new spread config.
func (h *configCrudHandler) SpreadAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("SPREAD", bodyJSON)
}

// SpreadDelete deletes a spread config by index.
func (h *configCrudHandler) SpreadDelete(index int) MTRetCode {
	return h.configDelete("SPREAD", map[string]string{WEB_PARAM_INDEX: formatInt(index)})
}

// SpreadShift shifts a spread config position.
func (h *configCrudHandler) SpreadShift(index, shift int) MTRetCode {
	return h.configShift("SPREAD", index, shift)
}

// --- HistorySync (7 methods) ---

// HistorySyncTotal gets the total count of history sync configs.
func (h *configCrudHandler) HistorySyncTotal() (int, MTRetCode) {
	return h.configTotal("HISTORY_SYNC")
}

// HistorySyncNext gets the history sync config at the given index.
func (h *configCrudHandler) HistorySyncNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("HISTORY_SYNC", index)
}

// HistorySyncGetAll gets all history sync configs.
func (h *configCrudHandler) HistorySyncGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.configGetAll("HISTORY_SYNC")
}

// HistorySyncAdd adds a new history sync config.
func (h *configCrudHandler) HistorySyncAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("HISTORY_SYNC", bodyJSON)
}

// HistorySyncDelete deletes a history sync config by index.
func (h *configCrudHandler) HistorySyncDelete(index int) MTRetCode {
	return h.configDelete("HISTORY_SYNC", map[string]string{WEB_PARAM_INDEX: formatInt(index)})
}

// HistorySyncShift shifts a history sync config position.
func (h *configCrudHandler) HistorySyncShift(index, shift int) MTRetCode {
	return h.configShift("HISTORY_SYNC", index, shift)
}

// HistorySyncStart triggers a history synchronization start.
func (h *configCrudHandler) HistorySyncStart() MTRetCode {
	cmd := WEB_CMD_HISTORY_SYNC_START
	answer := h.send(cmd, nil)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
	}
	return errCode
}

// --- SymbolGroup (6 methods) ---

// SymbolGroupTotal gets the total count of symbol group configs.
func (h *configCrudHandler) SymbolGroupTotal() (int, MTRetCode) {
	return h.configTotal("SYMBOL_GROUP")
}

// SymbolGroupNext gets the symbol group config at the given index.
func (h *configCrudHandler) SymbolGroupNext(index int) (map[string]interface{}, MTRetCode) {
	return h.configNext("SYMBOL_GROUP", index)
}

// SymbolGroupList retrieves a list of symbol group names.
func (h *configCrudHandler) SymbolGroupList() ([]string, MTRetCode) {
	cmd := WEB_CMD_SYMBOL_GROUP_LIST
	answer := h.send(cmd, nil)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", cmd))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", cmd, GetErrorStandard(errCode)))
		return nil, errCode
	}
	var names []string
	if err := json.Unmarshal([]byte(jsonStr), &names); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", cmd, err))
		return nil, MT_RET_REPORT_NODATA
	}
	return names, MT_RET_OK
}

// SymbolGroupAdd adds a new symbol group config.
func (h *configCrudHandler) SymbolGroupAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.configAdd("SYMBOL_GROUP", bodyJSON)
}

// SymbolGroupDelete deletes a symbol group config by name.
func (h *configCrudHandler) SymbolGroupDelete(name string) MTRetCode {
	return h.configDelete("SYMBOL_GROUP", map[string]string{WEB_PARAM_NAME: name})
}

// SymbolGroupShift shifts a symbol group config position.
func (h *configCrudHandler) SymbolGroupShift(index, shift int) MTRetCode {
	return h.configShift("SYMBOL_GROUP", index, shift)
}
