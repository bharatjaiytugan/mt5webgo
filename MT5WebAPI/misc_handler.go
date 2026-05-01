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
	"strings"
)

// WebSocket command constants for miscellaneous operations.
const (
	// Logger
	WEB_CMD_LOGGER_SERVER_REQUEST = "LOGGER_SERVER_REQUEST"

	// Notification
	WEB_CMD_NOTIFICATION_SEND = "NOTIFICATION_SEND"

	// Position extras
	WEB_CMD_POSITION_UPDATE       = "POSITION_UPDATE"
	WEB_CMD_POSITION_DELETE       = "POSITION_DELETE"
	WEB_CMD_POSITION_CHECK        = "POSITION_CHECK"
	WEB_CMD_POSITION_FIX          = "POSITION_FIX"
	WEB_CMD_POSITION_BACKUP_LIST  = "POSITION_BACKUP_LIST"
	WEB_CMD_POSITION_BACKUP_GET   = "POSITION_BACKUP_GET"
	WEB_CMD_POSITION_BACKUP_RESTORE = "POSITION_BACKUP_RESTORE"

	// Deal backup
	WEB_CMD_DEAL_BACKUP_LIST    = "DEAL_BACKUP_LIST"
	WEB_CMD_DEAL_BACKUP_GET     = "DEAL_BACKUP_GET"
	WEB_CMD_DEAL_BACKUP_RESTORE = "DEAL_BACKUP_RESTORE"

	// Daily extras
	WEB_CMD_DAILY_GET_PAGE        = "DAILY_GET_PAGE"
	WEB_CMD_DAILY_GET_LIGHT       = "DAILY_GET_LIGHT"
	WEB_CMD_DAILY_GET_LIGHT_BATCH = "DAILY_GET_LIGHT_BATCH"
	WEB_CMD_DAILY_GET_LIGHT_PAGE  = "DAILY_GET_LIGHT_PAGE"

	// User extras
	WEB_CMD_USER_TOTAL            = "USER_TOTAL"
	WEB_CMD_USER_GROUP            = "USER_GROUP"
	WEB_CMD_USER_GET_EXTERNAL     = "USER_GET_EXTERNAL"
	WEB_CMD_USER_SYNC_EXTERNAL    = "USER_SYNC_EXTERNAL"
	WEB_CMD_USER_CHECK_BALANCE    = "USER_CHECK_BALANCE"
	WEB_CMD_USER_RESTORE          = "USER_RESTORE"
	WEB_CMD_USER_ARCHIVE_ADD      = "USER_ARCHIVE_ADD"
	WEB_CMD_USER_ARCHIVE_GET      = "USER_ARCHIVE_GET"
	WEB_CMD_USER_ARCHIVE_GET_BATCH = "USER_ARCHIVE_GET_BATCH"
	WEB_CMD_USER_ACCOUNT_GET_BATCH = "USER_ACCOUNT_GET_BATCH"
	WEB_CMD_USER_CERTIFICATE_GET     = "USER_CERTIFICATE_GET"
	WEB_CMD_USER_CERTIFICATE_UPDATE  = "USER_CERTIFICATE_UPDATE"
	WEB_CMD_USER_CERTIFICATE_DELETE  = "USER_CERTIFICATE_DELETE"
	WEB_CMD_USER_CERTIFICATE_CONFIRM = "USER_CERTIFICATE_CONFIRM"
	WEB_CMD_USER_OTP_SECRET_GET = "USER_OTP_SECRET_GET"
	WEB_CMD_USER_OTP_SECRET_SET = "USER_OTP_SECRET_SET"

	// Subscription extras
	WEB_CMD_SUBSCRIPTION_HISTORY_ADD    = "SUBSCRIPTION_HISTORY_ADD"
	WEB_CMD_SUBSCRIPTION_HISTORY_UPDATE = "SUBSCRIPTION_HISTORY_UPDATE"
	WEB_CMD_SUBSCRIPTION_HISTORY_DELETE = "SUBSCRIPTION_HISTORY_DELETE"
	WEB_CMD_SUBSCRIPTION_HISTORY_GET    = "SUBSCRIPTION_HISTORY_GET"

	// Mail extras
	WEB_CMD_MAIL_GET      = "MAIL_GET"
	WEB_CMD_MAIL_GET_BODY = "MAIL_GET_BODY"

	// News extras
	WEB_CMD_NEWS_GET      = "NEWS_GET"
	WEB_CMD_NEWS_GET_BODY = "NEWS_GET_BODY"

	// Attachment
	WEB_CMD_ATTACHMENT_ADD    = "ATTACHMENT_ADD"
	WEB_CMD_ATTACHMENT_ATTACH = "ATTACHMENT_ATTACH"
	WEB_CMD_ATTACHMENT_GET    = "ATTACHMENT_GET"

	// Comment
	WEB_CMD_COMMENT_ADD    = "COMMENT_ADD"
	WEB_CMD_COMMENT_UPDATE = "COMMENT_UPDATE"
	WEB_CMD_COMMENT_DELETE = "COMMENT_DELETE"
	WEB_CMD_COMMENT_GET    = "COMMENT_GET"

	// Document
	WEB_CMD_DOCUMENT_ADD         = "DOCUMENT_ADD"
	WEB_CMD_DOCUMENT_UPDATE      = "DOCUMENT_UPDATE"
	WEB_CMD_DOCUMENT_DELETE      = "DOCUMENT_DELETE"
	WEB_CMD_DOCUMENT_GET         = "DOCUMENT_GET"
	WEB_CMD_DOCUMENT_HISTORY_GET = "DOCUMENT_HISTORY_GET"

	// Quit
	WEB_CMD_QUIT = "QUIT"

	// Test
	WEB_CMD_TEST_ACCESS = "TEST_ACCESS"
)

// Additional param constants for misc operations.
const (
	WEB_PARAM_MODE          = "MODE"
	WEB_PARAM_FILTER        = "FILTER"
	WEB_PARAM_METAQUOTES_IDS = "METAQUOTES_IDS"
	WEB_PARAM_TEXT          = "TEXT"
	WEB_PARAM_LOGINS        = "LOGINS"
	WEB_PARAM_EXTERNAL_ID   = "EXTERNAL_ID"
	WEB_PARAM_SECRET        = "SECRET"
)

// miscHandler implements miscellaneous protocol commands not covered by other handlers.
type miscHandler struct {
	apiBase
}

// --- Logger ---

// LoggerServerRequest retrieves server log entries.
// mode, logType, from, to, filter correspond to the LOGGER_SERVER_REQUEST params.
func (h *miscHandler) LoggerServerRequest(mode, logType string, from, to int64, filter string) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_MODE: mode,
		WEB_PARAM_TYPE: logType,
		WEB_PARAM_FROM: formatInt64(from),
		WEB_PARAM_TO:   formatInt64(to),
	}
	if filter != "" {
		params[WEB_PARAM_FILTER] = filter
	}
	answer := h.send(WEB_CMD_LOGGER_SERVER_REQUEST, params)
	if answer == nil {
		logWrite(LogTypeError, "send logger server request failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_LOGGER_SERVER_REQUEST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse logger server request failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse logger json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// --- Notification ---

// NotificationSend sends a push notification to a user's mobile terminal.
func (h *miscHandler) NotificationSend(metaquotesIDs, text string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_METAQUOTES_IDS: metaquotesIDs,
		WEB_PARAM_BODYTEXT:      text,
	}
	answer := h.send(WEB_CMD_NOTIFICATION_SEND, params)
	if answer == nil {
		logWrite(LogTypeError, "send notification failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_NOTIFICATION_SEND, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse notification send failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// --- Position extras ---

// PositionUpdate updates a position on the server.
func (h *miscHandler) PositionUpdate(positionJSON string) (map[string]interface{}, MTRetCode) {
	if positionJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: positionJSON,
	}
	answer := h.send(WEB_CMD_POSITION_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send position update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_POSITION_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse position update json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// PositionDelete deletes positions by comma-separated tickets.
func (h *miscHandler) PositionDelete(tickets string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TICKET: tickets,
	}
	answer := h.send(WEB_CMD_POSITION_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send position delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_POSITION_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// PositionCheck checks a position (server-side validation).
func (h *miscHandler) PositionCheck(positionJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: positionJSON,
	}
	answer := h.send(WEB_CMD_POSITION_CHECK, params)
	if answer == nil {
		logWrite(LogTypeError, "send position check failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_POSITION_CHECK, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position check failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// PositionFix fixes a position (server-side fixation).
func (h *miscHandler) PositionFix(positionJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_BODYTEXT: positionJSON,
	}
	answer := h.send(WEB_CMD_POSITION_FIX, params)
	if answer == nil {
		logWrite(LogTypeError, "send position fix failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_POSITION_FIX, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position fix failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// PositionBackupList gets position backup dates in a time range.
func (h *miscHandler) PositionBackupList(from, to int64) ([]int64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_FROM: formatInt64(from),
		WEB_PARAM_TO:   formatInt64(to),
	}
	answer := h.send(WEB_CMD_POSITION_BACKUP_LIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send position backup list failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_POSITION_BACKUP_LIST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position backup list failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	return parseTimestampsFromJSON(jsonStr), MT_RET_OK
}

// PositionBackupGet gets positions from a specific backup for a given login.
func (h *miscHandler) PositionBackupGet(backupDate int64, login uint64) ([]*MTPosition, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BACKUP: formatInt64(backupDate),
		WEB_PARAM_LOGIN:  formatUint64(login),
	}
	answer := h.send(WEB_CMD_POSITION_BACKUP_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send position backup get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_POSITION_BACKUP_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position backup get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	positions := parsePositionsFromJSON(jsonStr)
	if positions == nil {
		return []*MTPosition{}, MT_RET_OK
	}
	return positions, MT_RET_OK
}

// PositionBackupRestore restores a position from backup.
func (h *miscHandler) PositionBackupRestore(positionJSON string) (map[string]interface{}, MTRetCode) {
	if positionJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: positionJSON,
	}
	answer := h.send(WEB_CMD_POSITION_BACKUP_RESTORE, params)
	if answer == nil {
		logWrite(LogTypeError, "send position backup restore failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_POSITION_BACKUP_RESTORE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse position backup restore failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse position backup restore json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- Deal backup ---

// DealBackupList gets deal backup dates in a time range.
func (h *miscHandler) DealBackupList(from, to int64) ([]int64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_FROM: formatInt64(from),
		WEB_PARAM_TO:   formatInt64(to),
	}
	answer := h.send(WEB_CMD_DEAL_BACKUP_LIST, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal backup list failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DEAL_BACKUP_LIST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal backup list failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	return parseTimestampsFromJSON(jsonStr), MT_RET_OK
}

// DealBackupGet gets deals from a specific backup for a given login.
func (h *miscHandler) DealBackupGet(backupDate int64, login uint64) ([]*MTDeal, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BACKUP: formatInt64(backupDate),
		WEB_PARAM_LOGIN:  formatUint64(login),
	}
	answer := h.send(WEB_CMD_DEAL_BACKUP_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal backup get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DEAL_BACKUP_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal backup get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	deals := parseDealsFromJSON(jsonStr)
	if deals == nil {
		return []*MTDeal{}, MT_RET_OK
	}
	return deals, MT_RET_OK
}

// DealBackupRestore restores a deal from backup.
func (h *miscHandler) DealBackupRestore(dealJSON string) (map[string]interface{}, MTRetCode) {
	if dealJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: dealJSON,
	}
	answer := h.send(WEB_CMD_DEAL_BACKUP_RESTORE, params)
	if answer == nil {
		logWrite(LogTypeError, "send deal backup restore failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DEAL_BACKUP_RESTORE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse deal backup restore failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse deal backup restore json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- Daily extras ---

// DailyGetPage retrieves paged daily reports for a login in a date range.
func (h *miscHandler) DailyGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTDailyReport, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_FROM:   formatInt64(from),
		WEB_PARAM_TO:     formatInt64(to),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_DAILY_GET_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send daily get page failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DAILY_GET_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse daily get page failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	dailies := parseDailiesFromJSON(jsonStr)
	if dailies == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return dailies, MT_RET_OK
}

// DailyGetLight retrieves a lightweight daily report for a login in a date range.
func (h *miscHandler) DailyGetLight(login uint64, from, to int64) ([]*MTDailyReport, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_FROM:  formatInt64(from),
		WEB_PARAM_TO:    formatInt64(to),
	}
	answer := h.send(WEB_CMD_DAILY_GET_LIGHT, params)
	if answer == nil {
		logWrite(LogTypeError, "send daily get light failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DAILY_GET_LIGHT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse daily get light failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	dailies := parseDailiesFromJSON(jsonStr)
	if dailies == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return dailies, MT_RET_OK
}

// DailyGetLightBatch retrieves lightweight daily reports for multiple logins.
// DAILY_GET_LIGHT_BATCH may not work on WebSocket, so this falls back to looping over DailyGetLight.
func (h *miscHandler) DailyGetLightBatch(logins string, from, to int64) ([]*MTDailyReport, MTRetCode) {
	loginList := strings.Split(logins, ",")
	var allReports []*MTDailyReport
	for _, login := range loginList {
		login = strings.TrimSpace(login)
		if login == "" {
			continue
		}
		loginNum, err := strconv.ParseUint(login, 10, 64)
		if err != nil {
			continue
		}
		reports, ret := h.DailyGetLight(loginNum, from, to)
		if ret == MT_RET_OK && reports != nil {
			allReports = append(allReports, reports...)
		}
	}
	if allReports == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return allReports, MT_RET_OK
}

// DailyGetLightPage retrieves lightweight paged daily reports.
func (h *miscHandler) DailyGetLightPage(login uint64, from, to int64, offset, total uint32) ([]*MTDailyReport, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_FROM:   formatInt64(from),
		WEB_PARAM_TO:     formatInt64(to),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_DAILY_GET_LIGHT_PAGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send daily get light page failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DAILY_GET_LIGHT_PAGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse daily get light page failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	dailies := parseDailiesFromJSON(jsonStr)
	if dailies == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return dailies, MT_RET_OK
}

// --- User extras ---

// UserTotal gets the total number of users on the server.
func (h *miscHandler) UserTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_USER_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send user total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_USER_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// UserGroup gets the group name of a user by login.
func (h *miscHandler) UserGroup(login uint64) (string, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_GROUP, params)
	if answer == nil {
		logWrite(LogTypeError, "send user group failed")
		return "", MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_GROUP, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user group failed: %s", GetErrorStandard(errCode)))
		return "", errCode
	}
	if jsonStr == "" {
		return "", MT_RET_OK
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		return "", MT_RET_ERR_DATA
	}
	return jsonGetString(m, "Group"), MT_RET_OK
}

// UserGetExternal gets a user by external ID.
func (h *miscHandler) UserGetExternal(externalID string) (*MTUser, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_EXTERNAL_ID: externalID,
	}
	answer := h.send(WEB_CMD_USER_GET_EXTERNAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send user get external failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_GET_EXTERNAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user get external failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	user := parseUserFromJSON(jsonStr)
	if user == nil {
		return nil, MT_RET_ERR_DATA
	}
	return user, MT_RET_OK
}

// UserSyncExternal synchronizes a user by external ID.
func (h *miscHandler) UserSyncExternal(userJSON string) (map[string]interface{}, MTRetCode) {
	if userJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: userJSON,
	}
	answer := h.send(WEB_CMD_USER_SYNC_EXTERNAL, params)
	if answer == nil {
		logWrite(LogTypeError, "send user sync external failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_SYNC_EXTERNAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user sync external failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse user sync external json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// UserCheckBalance checks the balance for a user.
func (h *miscHandler) UserCheckBalance(login uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_CHECK_BALANCE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user check balance failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_CHECK_BALANCE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user check balance failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// UserRestore restores a user from backup.
func (h *miscHandler) UserRestore(userJSON string) (map[string]interface{}, MTRetCode) {
	if userJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: userJSON,
	}
	answer := h.send(WEB_CMD_USER_RESTORE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user restore failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_RESTORE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user restore failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse user restore json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// UserArchiveAdd archives a user.
func (h *miscHandler) UserArchiveAdd(login uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_ARCHIVE_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send user archive add failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_ARCHIVE_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user archive add failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// UserArchiveGet gets an archived user by login.
func (h *miscHandler) UserArchiveGet(login uint64) (*MTUser, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_ARCHIVE_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user archive get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_ARCHIVE_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user archive get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	user := parseUserFromJSON(jsonStr)
	if user == nil {
		return nil, MT_RET_ERR_DATA
	}
	return user, MT_RET_OK
}

// UserArchiveGetBatch gets archived users by comma-separated logins.
func (h *miscHandler) UserArchiveGetBatch(logins string) ([]*MTUser, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGINS: logins,
	}
	answer := h.send(WEB_CMD_USER_ARCHIVE_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send user archive get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_ARCHIVE_GET_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user archive get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	users := parseUsersFromJSON(jsonStr)
	if users == nil {
		return []*MTUser{}, MT_RET_OK
	}
	return users, MT_RET_OK
}

// UserAccountGetBatch gets account information for multiple logins (comma-separated).
func (h *miscHandler) UserAccountGetBatch(logins string) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGINS: logins,
	}
	answer := h.send(WEB_CMD_USER_ACCOUNT_GET_BATCH, params)
	if answer == nil {
		logWrite(LogTypeError, "send user account get batch failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_ACCOUNT_GET_BATCH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user account get batch failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse user account get batch json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// UserCertificateGet gets a user's certificate.
func (h *miscHandler) UserCertificateGet(login uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_CERTIFICATE_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user certificate get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_CERTIFICATE_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user certificate get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse user certificate get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// UserCertificateUpdate updates a user's certificate.
func (h *miscHandler) UserCertificateUpdate(login uint64, certJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN:    formatUint64(login),
		WEB_PARAM_BODYTEXT: certJSON,
	}
	answer := h.send(WEB_CMD_USER_CERTIFICATE_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user certificate update failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_CERTIFICATE_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user certificate update failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// UserCertificateDelete deletes a user's certificate.
func (h *miscHandler) UserCertificateDelete(login uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_CERTIFICATE_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user certificate delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_CERTIFICATE_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user certificate delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// UserCertificateConfirm confirms a user's certificate.
func (h *miscHandler) UserCertificateConfirm(login uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_CERTIFICATE_CONFIRM, params)
	if answer == nil {
		logWrite(LogTypeError, "send user certificate confirm failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_CERTIFICATE_CONFIRM, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user certificate confirm failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// UserOTPSecretGet gets a user's OTP secret.
func (h *miscHandler) UserOTPSecretGet(login uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_OTP_SECRET_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user otp secret get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_USER_OTP_SECRET_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user otp secret get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse user otp secret get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// UserOTPSecretSet sets a user's OTP secret.
func (h *miscHandler) UserOTPSecretSet(login uint64, secret string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_SECRET: secret,
	}
	answer := h.send(WEB_CMD_USER_OTP_SECRET_SET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user otp secret set failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_OTP_SECRET_SET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user otp secret set failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// --- Subscription extras ---

// SubscriptionAdd adds a subscription for a user.
func (h *miscHandler) SubscriptionAdd(subscriptionJSON string) (map[string]interface{}, MTRetCode) {
	if subscriptionJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: subscriptionJSON,
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_SUBSCRIPTION_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription add json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// SubscriptionUpdate updates a subscription.
func (h *miscHandler) SubscriptionUpdate(subscriptionJSON string) (map[string]interface{}, MTRetCode) {
	if subscriptionJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: subscriptionJSON,
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_SUBSCRIPTION_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription update json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// SubscriptionDelete deletes a subscription.
func (h *miscHandler) SubscriptionDelete(login uint64, id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_ID:    formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SubscriptionHistoryAdd adds a subscription history record.
func (h *miscHandler) SubscriptionHistoryAdd(historyJSON string) (map[string]interface{}, MTRetCode) {
	if historyJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: historyJSON,
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_HISTORY_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription history add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_SUBSCRIPTION_HISTORY_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history add json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// SubscriptionHistoryUpdate updates a subscription history record.
func (h *miscHandler) SubscriptionHistoryUpdate(historyJSON string) (map[string]interface{}, MTRetCode) {
	if historyJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: historyJSON,
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_HISTORY_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription history update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_SUBSCRIPTION_HISTORY_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history update json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// SubscriptionHistoryDelete deletes a subscription history record.
func (h *miscHandler) SubscriptionHistoryDelete(login uint64, id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_ID:    formatInt(id),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_HISTORY_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription history delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SUBSCRIPTION_HISTORY_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SubscriptionHistoryGet gets subscription history for a user.
func (h *miscHandler) SubscriptionHistoryGet(login uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_SUBSCRIPTION_HISTORY_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send subscription history get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_SUBSCRIPTION_HISTORY_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse subscription history get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- Mail extras ---

// MailGet gets mail messages for a login.
func (h *miscHandler) MailGet(login uint64, offset, total uint32) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_OFFSET: formatUint32(offset),
		WEB_PARAM_TOTAL:  formatUint32(total),
	}
	answer := h.send(WEB_CMD_MAIL_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send mail get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_MAIL_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse mail get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse mail get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// MailGetBody gets the body of a mail message by ticket.
func (h *miscHandler) MailGetBody(ticket uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_MAIL_GET_BODY, params)
	if answer == nil {
		logWrite(LogTypeError, "send mail get body failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_MAIL_GET_BODY, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse mail get body failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse mail get body json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- News extras ---

// NewsGet gets news messages.
func (h *miscHandler) NewsGet(id string, total uint32) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{}
	if id != "" {
		params[WEB_PARAM_ID] = id
	}
	if total > 0 {
		params[WEB_PARAM_TOTAL] = formatUint32(total)
	}
	answer := h.send(WEB_CMD_NEWS_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send news get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_NEWS_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse news get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse news get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// NewsGetBody gets the body of a news message by ticket.
func (h *miscHandler) NewsGetBody(ticket uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_TICKET: formatUint64(ticket),
	}
	answer := h.send(WEB_CMD_NEWS_GET_BODY, params)
	if answer == nil {
		logWrite(LogTypeError, "send news get body failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_NEWS_GET_BODY, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse news get body failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse news get body json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- Attachment ---

// AttachmentAdd adds an attachment.
func (h *miscHandler) AttachmentAdd(attachmentJSON string) (map[string]interface{}, MTRetCode) {
	if attachmentJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: attachmentJSON,
	}
	answer := h.send(WEB_CMD_ATTACHMENT_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send attachment add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_ATTACHMENT_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse attachment add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse attachment add json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// AttachmentAttach attaches an attachment to an entity.
func (h *miscHandler) AttachmentAttach(attachJSON string) MTRetCode {
	if attachJSON == "" {
		return MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: attachJSON,
	}
	answer := h.send(WEB_CMD_ATTACHMENT_ATTACH, params)
	if answer == nil {
		logWrite(LogTypeError, "send attachment attach failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_ATTACHMENT_ATTACH, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse attachment attach failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// AttachmentGet gets an attachment by ID.
func (h *miscHandler) AttachmentGet(id uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_ATTACHMENT_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send attachment get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_ATTACHMENT_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse attachment get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse attachment get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// --- Comment ---

// CommentAdd adds a comment.
func (h *miscHandler) CommentAdd(commentJSON string) (map[string]interface{}, MTRetCode) {
	if commentJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: commentJSON,
	}
	answer := h.send(WEB_CMD_COMMENT_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send comment add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_COMMENT_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse comment add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse comment add json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// CommentUpdate updates a comment.
func (h *miscHandler) CommentUpdate(commentJSON string) (map[string]interface{}, MTRetCode) {
	if commentJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: commentJSON,
	}
	answer := h.send(WEB_CMD_COMMENT_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send comment update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_COMMENT_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse comment update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse comment update json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// CommentDelete deletes a comment by ID.
func (h *miscHandler) CommentDelete(id uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_COMMENT_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send comment delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_COMMENT_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse comment delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// CommentGet gets comments by entity.
func (h *miscHandler) CommentGet(id uint64) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_COMMENT_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send comment get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_COMMENT_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse comment get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse comment get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// --- Document ---

// DocumentAdd adds a document.
func (h *miscHandler) DocumentAdd(documentJSON string) (map[string]interface{}, MTRetCode) {
	if documentJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: documentJSON,
	}
	answer := h.send(WEB_CMD_DOCUMENT_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send document add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DOCUMENT_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse document add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse document add json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// DocumentUpdate updates a document.
func (h *miscHandler) DocumentUpdate(documentJSON string) (map[string]interface{}, MTRetCode) {
	if documentJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: documentJSON,
	}
	answer := h.send(WEB_CMD_DOCUMENT_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send document update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DOCUMENT_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse document update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse document update json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// DocumentDelete deletes a document by ID.
func (h *miscHandler) DocumentDelete(id uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_DOCUMENT_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send document delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_DOCUMENT_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse document delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// DocumentGet gets a document by ID.
func (h *miscHandler) DocumentGet(id uint64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_DOCUMENT_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send document get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DOCUMENT_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse document get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse document get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return result, MT_RET_OK
}

// DocumentHistoryGet gets document history by ID.
func (h *miscHandler) DocumentHistoryGet(id uint64) ([]map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_DOCUMENT_HISTORY_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send document history get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseTolerantJSONWithCmd(WEB_CMD_DOCUMENT_HISTORY_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse document history get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse document history get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var result []map[string]interface{}
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			result = append(result, m)
		}
	}
	return result, MT_RET_OK
}

// --- Quit ---

// Quit sends a QUIT command to disconnect the session.
func (h *miscHandler) Quit() MTRetCode {
	answer := h.send(WEB_CMD_QUIT, nil)
	if answer == nil {
		// QUIT may not get a response, that is acceptable
		return MT_RET_OK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	return MT_RET_OK
}

// --- Test ---

// TestAccess tests server access/connectivity.
func (h *miscHandler) TestAccess() MTRetCode {
	answer := h.send(WEB_CMD_TEST_ACCESS, nil)
	if answer == nil {
		logWrite(LogTypeError, "send test access failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_TEST_ACCESS, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse test access failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// --- Helper functions ---

// parseTolerantJSON parses a response with tolerant command matching.
// The server may return the command name in a different format (e.g. slashes instead of underscores).
func parseTolerantJSONWithCmd(expectedCmd, answer string) (string, MTRetCode) {
	_, pos := getCommand(answer)
	if pos < 0 {
		logWrite(LogTypeError, fmt.Sprintf("no command found in answer, expected %s", expectedCmd))
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
	if jsonStr == "" {
		return "", MT_RET_REPORT_NODATA
	}
	return jsonStr, MT_RET_OK
}

// parseTimestampsFromJSON parses a JSON array of timestamps.
func parseTimestampsFromJSON(jsonStr string) []int64 {
	var raw []interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing timestamps json failed: %v", err))
		return nil
	}
	dates := make([]int64, 0, len(raw))
	for _, v := range raw {
		switch val := v.(type) {
		case float64:
			dates = append(dates, int64(val))
		}
	}
	return dates
}
