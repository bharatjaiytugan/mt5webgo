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

// WebSocket command constants for server configuration.
const (
	WEB_CMD_COMMON_SET          = "COMMON_SET"
	WEB_CMD_SERVER_ADD          = "SERVER_ADD"
	WEB_CMD_SERVER_DELETE       = "SERVER_DELETE"
	WEB_CMD_SERVER_SHIFT        = "SERVER_SHIFT"
	WEB_CMD_SERVER_TOTAL        = "SERVER_TOTAL"
	WEB_CMD_SERVER_NEXT         = "SERVER_NEXT"
	WEB_CMD_SERVER_GET          = "SERVER_GET"
	WEB_CMD_TLS_CERTIFICATE_ADD    = "TLS_CERTIFICATE_ADD"
	WEB_CMD_TLS_CERTIFICATE_DELETE = "TLS_CERTIFICATE_DELETE"
	WEB_CMD_TLS_CERTIFICATE_SHIFT  = "TLS_CERTIFICATE_SHIFT"
	WEB_CMD_TLS_CERTIFICATE_TOTAL  = "TLS_CERTIFICATE_TOTAL"
	WEB_CMD_TLS_CERTIFICATE_NEXT   = "TLS_CERTIFICATE_NEXT"
	WEB_CMD_TIME_SET            = "TIME_SET"
)

// Additional param constants.
const (
	WEB_PARAM_SHIFT = "SHIFT"
)

// serverConfigHandler implements server configuration protocol commands.
type serverConfigHandler struct {
	apiBase
}

// CommonSet updates the common server configuration.
// The server returns only RETCODE without JSON body, so we verify with CommonGet.
func (h *serverConfigHandler) CommonSet(config *MTConCommon) (*MTConCommon, MTRetCode) {
	if config == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: commonToJSON(config),
	}
	answer := h.send(WEB_CMD_COMMON_SET, params)
	if answer == nil {
		logWrite(LogTypeError, "send common set failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// server returns only RETCODE, no JSON body
	errCode := h.parseEmptyResult(WEB_CMD_COMMON_SET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse common set failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	// get updated config
	return h.commonGet()
}

// commonGet retrieves the common config (internal, reuses common_handler logic).
func (h *serverConfigHandler) commonGet() (*MTConCommon, MTRetCode) {
	answer := h.send(WEB_CMD_COMMON_GET, nil)
	if answer == nil {
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	jsonStr, errCode := parseSingleJSON(WEB_CMD_COMMON_GET, answerStr)
	if errCode != MT_RET_OK {
		return nil, errCode
	}
	common := parseConCommonFromJSON(jsonStr)
	if common == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return common, MT_RET_OK
}

// AccessServerTotal retrieves the total number of access servers.
func (h *serverConfigHandler) AccessServerTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_SERVER_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send server total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_SERVER_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// AccessServerNext retrieves the access server config at the given index.
func (h *serverConfigHandler) AccessServerNext(index int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(WEB_CMD_SERVER_NEXT, params)
	if answer == nil {
		logWrite(LogTypeError, "send server next failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SERVER_NEXT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server next failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse server next json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// AccessServerGet retrieves an access server config by ID.
func (h *serverConfigHandler) AccessServerGet(id int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatInt(id),
	}
	answer := h.send(WEB_CMD_SERVER_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send server get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SERVER_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse server get json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// AccessServerAdd adds a new access server configuration.
func (h *serverConfigHandler) AccessServerAdd(serverJSON string) (map[string]interface{}, MTRetCode) {
	if serverJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: serverJSON,
	}
	answer := h.send(WEB_CMD_SERVER_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send server add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SERVER_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse server add json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// AccessServerDelete deletes an access server by ID.
func (h *serverConfigHandler) AccessServerDelete(id int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_ID: formatInt(id),
	}
	answer := h.send(WEB_CMD_SERVER_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send server delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SERVER_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// AccessServerShift shifts an access server position.
func (h *serverConfigHandler) AccessServerShift(index, shift int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
		WEB_PARAM_SHIFT: formatInt(shift),
	}
	answer := h.send(WEB_CMD_SERVER_SHIFT, params)
	if answer == nil {
		logWrite(LogTypeError, "send server shift failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SERVER_SHIFT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse server shift failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// TLSCertTotal retrieves the total number of TLS certificates.
func (h *serverConfigHandler) TLSCertTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_TLS_CERTIFICATE_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send tls cert total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_TLS_CERTIFICATE_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// TLSCertNext retrieves the TLS certificate config at the given index.
func (h *serverConfigHandler) TLSCertNext(index int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(WEB_CMD_TLS_CERTIFICATE_NEXT, params)
	if answer == nil {
		logWrite(LogTypeError, "send tls cert next failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TLS_CERTIFICATE_NEXT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert next failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert next json failed: %v", err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// TLSCertAdd adds a new TLS certificate.
func (h *serverConfigHandler) TLSCertAdd(certJSON string) MTRetCode {
	if certJSON == "" {
		return MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: certJSON,
	}
	answer := h.send(WEB_CMD_TLS_CERTIFICATE_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send tls cert add failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_TLS_CERTIFICATE_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert add failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// TLSCertDelete deletes a TLS certificate by index.
func (h *serverConfigHandler) TLSCertDelete(index int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(WEB_CMD_TLS_CERTIFICATE_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send tls cert delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_TLS_CERTIFICATE_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// TLSCertShift shifts a TLS certificate position.
func (h *serverConfigHandler) TLSCertShift(index, shift int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
		WEB_PARAM_SHIFT: formatInt(shift),
	}
	answer := h.send(WEB_CMD_TLS_CERTIFICATE_SHIFT, params)
	if answer == nil {
		logWrite(LogTypeError, "send tls cert shift failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_TLS_CERTIFICATE_SHIFT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tls cert shift failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// TimeSet updates the server time configuration.
func (h *serverConfigHandler) TimeSet(config *MTConTime) (*MTConTime, MTRetCode) {
	if config == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: conTimeToJSON(config),
	}
	answer := h.send(WEB_CMD_TIME_SET, params)
	if answer == nil {
		logWrite(LogTypeError, "send time set failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TIME_SET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse time set failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	timeCon := parseConTimeFromJSON(jsonStr)
	if timeCon == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return timeCon, MT_RET_OK
}

// commonToJSON serializes MTConCommon to JSON using JSONWriter.
func commonToJSON(c *MTConCommon) string {
	if c == nil {
		return "{}"
	}
	w := &JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeString("Name", c.Name)
	w.WriteAttributeString("Owner", c.Owner)
	w.WriteAttributeString("OwnerID", c.OwnerID)
	w.WriteAttributeString("OwnerHost", c.OwnerHost)
	w.WriteAttributeString("OwnerEmail", c.OwnerEmail)
	w.WriteAttributeString("Product", c.Product)
	w.WriteAttributeInt64("ExpirationLicense", c.ExpirationLicense)
	w.WriteAttributeInt64("ExpirationSupport", c.ExpirationSupport)
	w.WriteAttributeUint32("LimitTradeServers", c.LimitTradeServers)
	w.WriteAttributeUint32("LimitWebServers", c.LimitWebServers)
	w.WriteAttributeUint32("LimitAccounts", c.LimitAccounts)
	w.WriteAttributeUint32("LimitDeals", c.LimitDeals)
	w.WriteAttributeUint32("LimitSymbols", c.LimitSymbols)
	w.WriteAttributeUint32("LimitGroups", c.LimitGroups)
	w.WriteAttributeUint32("LiveUpdateMode", uint32(c.LiveUpdateMode))
	w.WriteAttributeUint32("TotalUsers", c.TotalUsers)
	w.WriteAttributeUint32("TotalUsersReal", c.TotalUsersReal)
	w.WriteAttributeUint32("TotalDeals", c.TotalDeals)
	w.WriteAttributeUint32("TotalOrders", c.TotalOrders)
	w.WriteAttributeUint32("TotalOrdersHistory", c.TotalOrdersHistory)
	w.WriteAttributeUint32("TotalPositions", c.TotalPositions)
	w.WriteAttributeString("AccountUrl", c.AccountURL)
	w.WriteAttributeUint32("AccountAuto", c.AccountAuto)
	w.WriteEndObject()
	return w.String()
}

// conTimeToJSON serializes MTConTime to JSON using JSONWriter.
func conTimeToJSON(t *MTConTime) string {
	if t == nil {
		return "{}"
	}
	w := &JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeInt32("DaylightState", t.DaylightState)
	w.WriteAttributeInt32("Daylight", t.Daylight)
	w.WriteAttributeInt32("TimeZone", t.TimeZone)
	w.WriteAttributeString("TimeServer", t.TimeServer)
	// serialize Days array
	if t.Days != nil {
		daysJSON := daysToJSON(t.Days)
		w.WriteAttributeRaw("Days", daysJSON)
	}
	w.WriteEndObject()
	return w.String()
}

// daysToJSON serializes the Days 2D array to JSON.
func daysToJSON(days [][]EnTimeTableMode) string {
	if days == nil {
		return "[]"
	}
	b, err := json.Marshal(days)
	if err != nil {
		return "[]"
	}
	return string(b)
}
