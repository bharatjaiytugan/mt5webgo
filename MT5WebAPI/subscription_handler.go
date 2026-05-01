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

// WebSocket command constants for subscription operations.
const (
	// Subscription config
	WEB_CMD_SUBSCRIPTION_CFG_ADD    = "SUBSCRIPTION_CFG_ADD"
	WEB_CMD_SUBSCRIPTION_CFG_DELETE = "SUBSCRIPTION_CFG_DELETE"
	WEB_CMD_SUBSCRIPTION_CFG_SHIFT  = "SUBSCRIPTION_CFG_SHIFT"
	WEB_CMD_SUBSCRIPTION_CFG_TOTAL  = "SUBSCRIPTION_CFG_TOTAL"
	WEB_CMD_SUBSCRIPTION_CFG_NEXT   = "SUBSCRIPTION_CFG_NEXT"
	WEB_CMD_SUBSCRIPTION_CFG_GET    = "SUBSCRIPTION_CFG_GET"
	// Subscription management
	WEB_CMD_SUBSCRIPTION_JOIN   = "SUBSCRIPTION_JOIN"
	WEB_CMD_SUBSCRIPTION_CANCEL = "SUBSCRIPTION_CANCEL"
	WEB_CMD_SUBSCRIPTION_ADD    = "SUBSCRIPTION_ADD"
	WEB_CMD_SUBSCRIPTION_UPDATE = "SUBSCRIPTION_UPDATE"
	WEB_CMD_SUBSCRIPTION_DELETE = "SUBSCRIPTION_DELETE"
	WEB_CMD_SUBSCRIPTION_EXIST  = "SUBSCRIPTION_EXIST"
	WEB_CMD_SUBSCRIPTION_GET    = "SUBSCRIPTION_GET"
)

// subscriptionHandler implements subscription protocol commands.
type subscriptionHandler struct {
	apiBase
}

// --- 41. Subscription Config ---

// SubscriptionCfgTotal retrieves the total number of subscription configs.
func (h *subscriptionHandler) SubscriptionCfgTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_SUBSCRIPTION_CFG_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// SubscriptionCfgNext retrieves the subscription config at the given index.
func (h *subscriptionHandler) SubscriptionCfgNext(index int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_NEXT, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg next failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SUBSCRIPTION_CFG_NEXT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg next failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg next json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// SubscriptionCfgGet retrieves a subscription config by ID.
func (h *subscriptionHandler) SubscriptionCfgGet(id int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SUBSCRIPTION_CFG_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg get json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// SubscriptionCfgAdd adds a new subscription config from a JSON string.
func (h *subscriptionHandler) SubscriptionCfgAdd(cfgJSON string) (map[string]interface{}, MTRetCode) {
	if cfgJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: cfgJSON,
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SUBSCRIPTION_CFG_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg add json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// SubscriptionCfgDelete deletes a subscription config by ID.
func (h *subscriptionHandler) SubscriptionCfgDelete(id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_ID: formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_CFG_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SubscriptionCfgShift shifts a subscription config position.
func (h *subscriptionHandler) SubscriptionCfgShift(index, shift int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
		WEB_PARAM_SHIFT: formatInt(shift),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CFG_SHIFT, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cfg shift failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_CFG_SHIFT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cfg shift failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// --- 42. Subscription Management ---

// SubscriptionGet retrieves subscriptions for a given login.
func (h *subscriptionHandler) SubscriptionGet(login uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SUBSCRIPTION_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription get json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// SubscriptionExist checks whether a user is subscribed to a subscription config.
// Returns true if the subscription exists (retcode == 0).
func (h *subscriptionHandler) SubscriptionExist(login uint64, id int) (bool, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_ID:    formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_EXIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription exist failed")
		return false, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_EXIST, answerStr)
	if errCode == MT_RET_OK {
		return true, MT_RET_OK
	}
	// not found is not an error per se
	if errCode == MT_RET_ERR_NOTFOUND {
		return false, MT_RET_OK
	}
	return false, errCode
}

// SubscriptionJoin joins a user to a subscription config.
func (h *subscriptionHandler) SubscriptionJoin(login uint64, id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_ID:    formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_JOIN, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription join failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_JOIN, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription join failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SubscriptionCancel cancels a user's subscription to a subscription config.
func (h *subscriptionHandler) SubscriptionCancel(login uint64, id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_ID:    formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_CANCEL, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription cancel failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_CANCEL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription cancel failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}
