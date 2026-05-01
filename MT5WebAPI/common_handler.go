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



// commonHandler implements the COMMON_GET protocol command.
type commonHandler struct {
	apiBase
}

// CommonGet retrieves the common server configuration.
func (h *commonHandler) CommonGet() (*MTConCommon, MTRetCode) {
	answer := h.send(WEB_CMD_COMMON_GET, nil)
	if answer == nil {
		logWrite(LogTypeError, "send common get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_COMMON_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse common get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	common := parseConCommonFromJSON(jsonStr)
	if common == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return common, MT_RET_OK
}

// parseConCommonFromJSON parses an MTConCommon from a JSON string.
// MT5 server sends all numeric values as strings, so we parse manually.
func parseConCommonFromJSON(jsonStr string) *MTConCommon {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing common config from json failed: %v", err))
		return nil
	}
	c := &MTConCommon{
		Name:               jsonGetString(raw, "Name"),
		Owner:              jsonGetString(raw, "Owner"),
		OwnerID:            jsonGetString(raw, "OwnerID"),
		OwnerHost:          jsonGetString(raw, "OwnerHost"),
		OwnerEmail:         jsonGetString(raw, "OwnerEmail"),
		Product:            jsonGetString(raw, "Product"),
		ExpirationLicense:  jsonGetInt64(raw, "ExpirationLicense"),
		ExpirationSupport:  jsonGetInt64(raw, "ExpirationSupport"),
		LimitTradeServers:  jsonGetUint32(raw, "LimitTradeServers"),
		LimitWebServers:    jsonGetUint32(raw, "LimitWebServers"),
		LimitAccounts:      jsonGetUint32(raw, "LimitAccounts"),
		LimitDeals:         jsonGetUint32(raw, "LimitDeals"),
		LimitSymbols:       jsonGetUint32(raw, "LimitSymbols"),
		LimitGroups:        jsonGetUint32(raw, "LimitGroups"),
		LiveUpdateMode:     EnUpdateMode(jsonGetUint32(raw, "LiveUpdateMode")),
		TotalUsers:         jsonGetUint32(raw, "TotalUsers"),
		TotalUsersReal:     jsonGetUint32(raw, "TotalUsersReal"),
		TotalDeals:         jsonGetUint32(raw, "TotalDeals"),
		TotalOrders:        jsonGetUint32(raw, "TotalOrders"),
		TotalOrdersHistory: jsonGetUint32(raw, "TotalOrdersHistory"),
		TotalPositions:     jsonGetUint32(raw, "TotalPositions"),
		AccountURL:         jsonGetString(raw, "AccountUrl"),
		AccountAuto:        jsonGetUint32(raw, "AccountAuto"),
	}
	return c
}
