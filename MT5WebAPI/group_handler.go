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

// groupHandler implements group-related protocol commands.
type groupHandler struct {
	apiBase
}

// GroupTotal retrieves the total number of groups on the server.
func (h *groupHandler) GroupTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_GROUP_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send group total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_GROUP_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// GroupNext retrieves the group config at the given position.
func (h *groupHandler) GroupNext(pos uint32) (*MTConGroup, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatUint32(pos),
	}
	answer := h.send(WEB_CMD_GROUP_NEXT, params)
	if answer == nil {
		logWrite(LogTypeError, "send group next failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_GROUP_NEXT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group next failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	group := parseConGroupFromJSON(jsonStr)
	if group == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return group, MT_RET_OK
}

// GroupGet retrieves a group config by name.
func (h *groupHandler) GroupGet(name string) (*MTConGroup, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_GROUP: name,
	}
	answer := h.send(WEB_CMD_GROUP_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send group get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_GROUP_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	group := parseConGroupFromJSON(jsonStr)
	if group == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return group, MT_RET_OK
}

// GroupAdd adds a new group to the server.
func (h *groupHandler) GroupAdd(group *MTConGroup) (*MTConGroup, MTRetCode) {
	if group == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: groupToJSON(group),
	}
	answer := h.send(WEB_CMD_GROUP_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send group add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_GROUP_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	newGroup := parseConGroupFromJSON(jsonStr)
	if newGroup == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return newGroup, MT_RET_OK
}

// GroupDelete deletes a group by name.
func (h *groupHandler) GroupDelete(name string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_GROUP: name,
	}
	answer := h.send(WEB_CMD_GROUP_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send group delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_GROUP_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse group delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// parseConGroupFromJSON parses an MTConGroup from a JSON string.
// MT5 server sends all numeric values as strings, so we parse manually.
func parseConGroupFromJSON(jsonStr string) *MTConGroup {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing group from json failed: %v", err))
		return nil
	}
	g := &MTConGroup{
		Group:              jsonGetString(raw, "Group"),
		Server:             jsonGetUint64(raw, "Server"),
		PermissionsFlags:   EnGroupPermissionsFlags(jsonGetUint64(raw, "PermissionsFlags")),
		AuthMode:           EnAuthMode(jsonGetUint32(raw, "AuthMode")),
		AuthPasswordMin:    jsonGetUint32(raw, "AuthPasswordMin"),
		AuthOTPMode:        EnAuthOTPMode(jsonGetUint32(raw, "AuthOTPMode")),
		Company:            jsonGetString(raw, "Company"),
		CompanyPage:        jsonGetString(raw, "CompanyPage"),
		CompanyEmail:       jsonGetString(raw, "CompanyEmail"),
		CompanySupportPage:  jsonGetString(raw, "CompanySupportPage"),
		CompanySupportEmail: jsonGetString(raw, "CompanySupportEmail"),
		CompanyCatalog:     jsonGetString(raw, "CompanyCatalog"),
		Currency:           jsonGetString(raw, "Currency"),
		CurrencyDigits:     jsonGetUint32(raw, "CurrencyDigits"),
		ReportsMode:        EnReportsMode(jsonGetUint32(raw, "ReportsMode")),
		ReportsFlags:       EnReportsFlags(jsonGetUint32(raw, "ReportsFlags")),
		ReportsSMTP:        jsonGetString(raw, "ReportsSMTP"),
		ReportsSMTPLogin:   jsonGetString(raw, "ReportsSMTPLogin"),
		ReportsSMTPPass:    jsonGetString(raw, "ReportsSMTPPass"),
		NewsMode:           EnNewsMode(jsonGetUint32(raw, "NewsMode")),
		NewsCategory:       jsonGetString(raw, "NewsCategory"),
		MailMode:           EnMailMode(jsonGetUint32(raw, "MailMode")),
		TradeFlags:         EnGroupTradeFlags(jsonGetUint64(raw, "TradeFlags")),
		TradeTransferMode:  EnTransferMode(jsonGetUint32(raw, "TradeTransferMode")),
		TradeInterestrate:  jsonGetFloat64(raw, "TradeInterestrate"),
		TradeVirtualCredit: jsonGetFloat64(raw, "TradeVirtualCredit"),
		MarginMode:         EnGroupMarginMode(jsonGetUint32(raw, "MarginMode")),
		MarginFreeMode:     EnFreeMarginMode(jsonGetUint32(raw, "MarginFreeMode")),
		MarginSOMode:       EnStopOutMode(jsonGetUint32(raw, "MarginSOMode")),
		MarginCall:         jsonGetFloat64(raw, "MarginCall"),
		MarginStopOut:      jsonGetFloat64(raw, "MarginStopOut"),
		DemoLeverage:       jsonGetUint32(raw, "DemoLeverage"),
		DemoDeposit:        jsonGetFloat64(raw, "DemoDeposit"),
		LimitHistory:       EnHistoryLimit(jsonGetUint32(raw, "LimitHistory")),
		LimitOrders:        jsonGetUint32(raw, "LimitOrders"),
		LimitSymbols:       jsonGetUint32(raw, "LimitSymbols"),
		LimitPositions:     jsonGetUint32(raw, "LimitPositions"),
		MarginFreeProfitMode: EnMarginFreeProfitMode(jsonGetUint32(raw, "MarginFreeProfitMode")),
	}
	return g
}

// groupToJSON serializes MTConGroup to JSON using JSONWriter (matches C# MTGroupJson.ToJson).
func groupToJSON(g *MTConGroup) string {
	if g == nil {
		return "{}"
	}
	w := &JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeString("Group", g.Group)
	w.WriteAttributeUint64("Server", g.Server)
	w.WriteAttributeUint64("PermissionsFlags", uint64(g.PermissionsFlags))
	w.WriteAttributeUint32("AuthMode", uint32(g.AuthMode))
	w.WriteAttributeUint32("AuthPasswordMin", g.AuthPasswordMin)
	w.WriteAttributeUint32("AuthOTPMode", uint32(g.AuthOTPMode))
	w.WriteAttributeString("Company", g.Company)
	w.WriteAttributeString("CompanyPage", g.CompanyPage)
	w.WriteAttributeString("CompanyEmail", g.CompanyEmail)
	w.WriteAttributeString("CompanySupportPage", g.CompanySupportPage)
	w.WriteAttributeString("CompanySupportEmail", g.CompanySupportEmail)
	w.WriteAttributeString("CompanyCatalog", g.CompanyCatalog)
	w.WriteAttributeString("Currency", g.Currency)
	w.WriteAttributeUint32("CurrencyDigits", g.CurrencyDigits)
	w.WriteAttributeUint32("ReportsMode", uint32(g.ReportsMode))
	w.WriteAttributeUint32("ReportsFlags", uint32(g.ReportsFlags))
	w.WriteAttributeString("ReportsSMTP", g.ReportsSMTP)
	w.WriteAttributeString("ReportsSMTPLogin", g.ReportsSMTPLogin)
	w.WriteAttributeString("ReportsSMTPPass", g.ReportsSMTPPass)
	w.WriteAttributeUint32("NewsMode", uint32(g.NewsMode))
	w.WriteAttributeString("NewsCategory", g.NewsCategory)
	w.WriteAttributeUint32("MailMode", uint32(g.MailMode))
	w.WriteAttributeUint64("TradeFlags", uint64(g.TradeFlags))
	w.WriteAttributeUint32("TradeTransferMode", uint32(g.TradeTransferMode))
	w.WriteAttributeFloat64("TradeInterestrate", g.TradeInterestrate)
	w.WriteAttributeFloat64("TradeVirtualCredit", g.TradeVirtualCredit)
	w.WriteAttributeUint32("MarginMode", uint32(g.MarginMode))
	w.WriteAttributeUint32("MarginSOMode", uint32(g.MarginSOMode))
	w.WriteAttributeUint32("MarginFreeMode", uint32(g.MarginFreeMode))
	w.WriteAttributeFloat64("MarginCall", g.MarginCall)
	w.WriteAttributeFloat64("MarginStopOut", g.MarginStopOut)
	w.WriteAttributeUint32("MarginFreeProfitMode", uint32(g.MarginFreeProfitMode))
	w.WriteAttributeUint32("DemoLeverage", g.DemoLeverage)
	w.WriteAttributeFloat64("DemoDeposit", g.DemoDeposit)
	w.WriteAttributeUint32("LimitHistory", uint32(g.LimitHistory))
	w.WriteAttributeUint32("LimitOrders", g.LimitOrders)
	w.WriteAttributeUint32("LimitSymbols", g.LimitSymbols)
	w.WriteAttributeUint32("LimitPositions", g.LimitPositions)
	// commissions
	w.WriteAttributeRaw("Commissions", commissionsToJSON(g.Commissions))
	// symbols
	w.WriteAttributeRaw("Symbols", groupSymbolsToJSON(g.Symbols))
	w.WriteEndObject()
	return w.String()
}

func commissionsToJSON(commissions []MTConCommission) string {
	if commissions == nil || len(commissions) == 0 {
		return "[]"
	}
	result := "["
	for i, c := range commissions {
		if i > 0 {
			result += ","
		}
		result += commissionToJSON(&c)
	}
	return result + "]"
}

func commissionToJSON(c *MTConCommission) string {
	if c == nil {
		return "{}"
	}
	w := &JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeString("Name", c.Name)
	w.WriteAttributeString("Description", c.Description)
	w.WriteAttributeString("Path", c.Path)
	w.WriteAttributeUint32("Mode", uint32(c.Mode))
	w.WriteAttributeUint32("RangeMode", uint32(c.RangeMode))
	w.WriteAttributeUint32("ChargeMode", uint32(c.ChargeMode))
	w.WriteAttributeString("TurnoverCurrency", c.TurnoverCurrency)
	w.WriteAttributeUint32("EntryMode", uint32(c.EntryMode))
	w.WriteAttributeRaw("Tiers", tiersToJSON(c.Tiers))
	w.WriteEndObject()
	return w.String()
}

func tiersToJSON(tiers []MTConCommTier) string {
	if tiers == nil || len(tiers) == 0 {
		return "[]"
	}
	result := "["
	for i, t := range tiers {
		if i > 0 {
			result += ","
		}
		w := &JSONWriter{}
		w.WriteBeginObject()
		w.WriteAttributeUint32("Mode", uint32(t.Mode))
		w.WriteAttributeUint32("Type", uint32(t.Type))
		w.WriteAttributeFloat64("Value", t.Value)
		w.WriteAttributeFloat64("Minimal", t.Minimal)
		w.WriteAttributeFloat64("Maximal", t.Maximal)
		w.WriteAttributeFloat64("RangeFrom", t.RangeFrom)
		w.WriteAttributeFloat64("RangeTo", t.RangeTo)
		w.WriteAttributeString("Currency", t.Currency)
		w.WriteEndObject()
		result += w.String()
	}
	return result + "]"
}

func groupSymbolsToJSON(symbols []MTConGroupSymbol) string {
	if symbols == nil || len(symbols) == 0 {
		return "[]"
	}
	result := "["
	for i, s := range symbols {
		if i > 0 {
			result += ","
		}
		result += groupSymbolToJSON(&s)
	}
	return result + "]"
}

func groupSymbolToJSON(s *MTConGroupSymbol) string {
	if s == nil {
		return "{}"
	}
	w := &JSONWriter{}
	w.WriteBeginObject()
	w.WriteAttributeString("Path", s.Path)
	w.WriteAttributeUint32("TradeMode", uint32(s.TradeMode))
	w.WriteAttributeUint32("ExecMode", uint32(s.ExecMode))
	w.WriteAttributeUint32("FillFlags", uint32(s.FillFlags))
	w.WriteAttributeUint32("ExpirFlags", uint32(s.ExpirFlags))
	w.WriteAttributeInt32("SpreadDiff", int32(s.SpreadDiff))
	w.WriteAttributeInt32("SpreadDiffBalance", int32(s.SpreadDiffBalance))
	w.WriteAttributeInt32("StopsLevel", int32(s.StopsLevel))
	w.WriteAttributeInt32("FreezeLevel", int32(s.FreezeLevel))
	w.WriteEndObject()
	return w.String()
}
