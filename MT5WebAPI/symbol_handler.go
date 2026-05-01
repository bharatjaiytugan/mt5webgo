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

// symbolHandler implements symbol-related protocol commands.
type symbolHandler struct {
	apiBase
}

// SymbolTotal retrieves the total number of symbols on the server.
func (h *symbolHandler) SymbolTotal() (int, MTRetCode) {
	answer := h.send(WEB_CMD_SYMBOL_TOTAL, nil)
	if answer == nil {
		logWrite(LogTypeError, "send symbol total failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(WEB_CMD_SYMBOL_TOTAL, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol total failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// SymbolNext retrieves the symbol config at the given position.
func (h *symbolHandler) SymbolNext(pos int) (*MTConSymbol, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(pos),
	}
	answer := h.send(WEB_CMD_SYMBOL_NEXT, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol next failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SYMBOL_NEXT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol next failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	sym := parseConSymbolFromJSON(jsonStr)
	if sym == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return sym, MT_RET_OK
}

// SymbolGet retrieves a symbol config by name.
func (h *symbolHandler) SymbolGet(name string) (*MTConSymbol, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: name,
	}
	answer := h.send(WEB_CMD_SYMBOL_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SYMBOL_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	sym := parseConSymbolFromJSON(jsonStr)
	if sym == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return sym, MT_RET_OK
}

// SymbolGetGroup retrieves a symbol config by name and group.
func (h *symbolHandler) SymbolGetGroup(name, group string) (*MTConSymbol, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: name,
		WEB_PARAM_GROUP:  group,
	}
	answer := h.send(WEB_CMD_SYMBOL_GET_GROUP, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol get group failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SYMBOL_GET_GROUP, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol get group failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	sym := parseConSymbolFromJSON(jsonStr)
	if sym == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return sym, MT_RET_OK
}

// SymbolAdd adds a new symbol to the server.
func (h *symbolHandler) SymbolAdd(symbol *MTConSymbol) (*MTConSymbol, MTRetCode) {
	if symbol == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	symbolJSON := SymbolToJSON(symbol)
	params := map[string]string{
		WEB_PARAM_BODYTEXT: symbolJSON,
	}
	answer := h.send(WEB_CMD_SYMBOL_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_SYMBOL_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	newSymbol := parseConSymbolFromJSON(jsonStr)
	if newSymbol == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return newSymbol, MT_RET_OK
}

// SymbolDelete deletes a symbol by name.
func (h *symbolHandler) SymbolDelete(name string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_SYMBOL: name,
	}
	answer := h.send(WEB_CMD_SYMBOL_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send symbol delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SYMBOL_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse symbol delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// parseConSymbolFromJSON parses an MTConSymbol from a JSON string using manual map-based parsing
// because MT5 sends all numbers as strings in JSON.
func parseConSymbolFromJSON(jsonStr string) *MTConSymbol {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing symbol from json failed: %v", err))
		return nil
	}
	sym := &MTConSymbol{}
	// String fields
	sym.Symbol = jsonGetString(m, "Symbol")
	sym.Path = jsonGetString(m, "Path")
	sym.ISIN = jsonGetString(m, "ISIN")
	sym.Description = jsonGetString(m, "Description")
	sym.International = jsonGetString(m, "International")
	sym.Basis = jsonGetString(m, "Basis")
	sym.Source = jsonGetString(m, "Source")
	sym.Page = jsonGetString(m, "Page")
	sym.CurrencyBase = jsonGetString(m, "CurrencyBase")
	sym.CurrencyProfit = jsonGetString(m, "CurrencyProfit")
	sym.CurrencyMargin = jsonGetString(m, "CurrencyMargin")
	sym.Category = jsonGetString(m, "Category")
	sym.Exchange = jsonGetString(m, "Exchange")
	sym.CFI = jsonGetString(m, "CFI")
	sym.Country = jsonGetString(m, "Country")
	// uint32 fields
	sym.CurrencyBaseDigits = jsonGetUint32(m, "CurrencyBaseDigits")
	sym.CurrencyProfitDigits = jsonGetUint32(m, "CurrencyProfitDigits")
	sym.CurrencyMarginDigits = jsonGetUint32(m, "CurrencyMarginDigits")
	sym.Color = jsonGetUint32(m, "Color")
	sym.ColorBackground = jsonGetUint32(m, "ColorBackground")
	sym.Digits = jsonGetUint32(m, "Digits")
	sym.TickBookDepth = jsonGetUint32(m, "TickBookDepth")
	sym.FilterSoft = jsonGetUint32(m, "FilterSoft")
	sym.FilterSoftTicks = jsonGetUint32(m, "FilterSoftTicks")
	sym.FilterHard = jsonGetUint32(m, "FilterHard")
	sym.FilterHardTicks = jsonGetUint32(m, "FilterHardTicks")
	sym.FilterDiscard = jsonGetUint32(m, "FilterDiscard")
	sym.FilterSpreadMax = jsonGetUint32(m, "FilterSpreadMax")
	sym.FilterSpreadMin = jsonGetUint32(m, "FilterSpreadMin")
	sym.FilterGap = jsonGetUint32(m, "FilterGap")
	sym.FilterGapTicks = jsonGetUint32(m, "FilterGapTicks")
	sym.QuotesTimeout = jsonGetUint32(m, "QuotesTimeout")
	sym.RETimeout = jsonGetUint32(m, "RETimeout")
	sym.IETimeout = jsonGetUint32(m, "IETimeout")
	sym.IESlipProfit = jsonGetUint32(m, "IESlipProfit")
	sym.IESlipLosing = jsonGetUint32(m, "IESlipLosing")
	sym.SpliceTimeDays = jsonGetUint32(m, "SpliceTimeDays")
	sym.SubscriptionsDelay = jsonGetUint32(m, "SubscriptionsDelay")
	// Enum/uint32 cast fields
	sym.TickFlags = EnTickFlags(jsonGetUint64(m, "TickFlags"))
	sym.ChartMode = EnChartMode(jsonGetUint32(m, "TickChartMode"))
	sym.TradeMode = EnTradeMode(jsonGetUint32(m, "TradeMode"))
	sym.CalcMode = EnCalcMode(jsonGetUint32(m, "CalcMode"))
	sym.ExecMode = EnExecutionMode(jsonGetUint32(m, "ExecMode"))
	sym.GTCMode = EnGTCMode(jsonGetUint32(m, "GTCMode"))
	sym.FillFlags = EnFillingFlags(jsonGetUint32(m, "FillFlags"))
	sym.ExpirFlags = EnExpirationFlags(jsonGetUint32(m, "ExpirFlags"))
	sym.OrderFlags = EnOrderFlags(jsonGetUint32(m, "OrderFlags"))
	sym.MarginCheckMode = EnMarginFlags(jsonGetUint32(m, "MarginCheckMode"))
	sym.MarginFlags = EnMarginFlags(jsonGetUint32(m, "MarginFlags"))
	sym.SwapMode = EnSwapMode(jsonGetUint32(m, "SwapMode"))
	sym.REFlags = EnRequestFlags(jsonGetUint32(m, "REFlags"))
	sym.IECheckMode = EnInstantMode(jsonGetUint32(m, "IECheckMode"))
	sym.OptionsMode = EnOptionMode(jsonGetUint32(m, "OptionMode"))
	sym.SpliceType = EnSpliceType(jsonGetUint32(m, "SpliceType"))
	sym.SpliceTimeType = EnSpliceTimeType(jsonGetUint32(m, "SpliceTimeType"))
	sym.IEFlags = EnTradeInstantFlags(jsonGetUint32(m, "IEFlags"))
	sym.Sector = EnSectors(jsonGetUint32(m, "Sector"))
	sym.Industry = EnIndustries(jsonGetUint32(m, "Industry"))
	// TradeFlags as uint64
	sym.TradeFlags = EnSymbolTradeFlags(jsonGetUint64(m, "TradeFlags"))
	// int32 fields
	sym.Spread = jsonGetInt32(m, "Spread")
	sym.SpreadBalance = jsonGetInt32(m, "SpreadBalance")
	sym.SpreadDiff = jsonGetInt32(m, "SpreadDiff")
	sym.SpreadDiffBalance = jsonGetInt32(m, "SpreadDiffBalance")
	sym.StopsLevel = jsonGetInt32(m, "StopsLevel")
	sym.FreezeLevel = jsonGetInt32(m, "FreezeLevel")
	sym.Swap3Day = jsonGetInt32(m, "Swap3Day")
	// float64 fields
	sym.Point = jsonGetFloat64(m, "Point")
	sym.Multiply = jsonGetFloat64(m, "Multiply")
	sym.TickValue = jsonGetFloat64(m, "TickValue")
	sym.TickSize = jsonGetFloat64(m, "TickSize")
	sym.ContractSize = jsonGetFloat64(m, "ContractSize")
	sym.MarginInitial = jsonGetFloat64(m, "MarginInitial")
	sym.MarginMaintenance = jsonGetFloat64(m, "MarginMaintenance")
	sym.MarginRateLiquidity = jsonGetFloat64(m, "MarginLiquidity")
	sym.MarginHedged = jsonGetFloat64(m, "MarginHedged")
	sym.MarginRateCurrency = jsonGetFloat64(m, "MarginCurrency")
	sym.MarginLong = jsonGetFloat64(m, "MarginLong")
	sym.MarginShort = jsonGetFloat64(m, "MarginShort")
	sym.MarginLimit = jsonGetFloat64(m, "MarginLimit")
	sym.MarginStop = jsonGetFloat64(m, "MarginStop")
	sym.MarginStopLimit = jsonGetFloat64(m, "MarginStopLimit")
	sym.SwapLong = jsonGetFloat64(m, "SwapLong")
	sym.SwapShort = jsonGetFloat64(m, "SwapShort")
	sym.PriceSettle = jsonGetFloat64(m, "PriceSettle")
	sym.PriceLimitMax = jsonGetFloat64(m, "PriceLimitMax")
	sym.PriceLimitMin = jsonGetFloat64(m, "PriceLimitMin")
	sym.PriceStrike = jsonGetFloat64(m, "PriceStrike")
	sym.FaceValue = jsonGetFloat64(m, "FaceValue")
	sym.AccruedInterest = jsonGetFloat64(m, "AccruedInterest")
	// int64 fields
	sym.TimeStart = jsonGetInt64(m, "TimeStart")
	sym.TimeExpiration = jsonGetInt64(m, "TimeExpiration")
	// uint64 volume fields
	sym.VolumeMinExt = jsonGetUint64(m, "VolumeMinExt")
	sym.VolumeMaxExt = jsonGetUint64(m, "VolumeMaxExt")
	sym.VolumeStepExt = jsonGetUint64(m, "VolumeStepExt")
	sym.VolumeLimitExt = jsonGetUint64(m, "VolumeLimitExt")
	sym.IEVolumeMaxExt = jsonGetUint64(m, "IEVolumeMaxExt")
	// MarginRateInitial map
	sym.MarginRateInitial = map[EnMarginRateTypes]float64{
		MarginRateBuy:           jsonGetFloat64(m, "MarginInitialBuy"),
		MarginRateSell:          jsonGetFloat64(m, "MarginInitialSell"),
		MarginRateBuyLimit:      jsonGetFloat64(m, "MarginInitialBuyLimit"),
		MarginRateSellLimit:     jsonGetFloat64(m, "MarginInitialSellLimit"),
		MarginRateBuyStop:       jsonGetFloat64(m, "MarginInitialBuyStop"),
		MarginRateSellStop:      jsonGetFloat64(m, "MarginInitialSellStop"),
		MarginRateBuyStopLimit:  jsonGetFloat64(m, "MarginInitialBuyStopLimit"),
		MarginRateSellStopLimit: jsonGetFloat64(m, "MarginInitialSellStopLimit"),
	}
	// MarginRateMaintenance map
	sym.MarginRateMaintenance = map[EnMarginRateTypes]float64{
		MarginRateBuy:           jsonGetFloat64(m, "MarginMaintenanceBuy"),
		MarginRateSell:          jsonGetFloat64(m, "MarginMaintenanceSell"),
		MarginRateBuyLimit:      jsonGetFloat64(m, "MarginMaintenanceBuyLimit"),
		MarginRateSellLimit:     jsonGetFloat64(m, "MarginMaintenanceSellLimit"),
		MarginRateBuyStop:       jsonGetFloat64(m, "MarginMaintenanceBuyStop"),
		MarginRateSellStop:      jsonGetFloat64(m, "MarginMaintenanceSellStop"),
		MarginRateBuyStopLimit:  jsonGetFloat64(m, "MarginMaintenanceBuyStopLimit"),
		MarginRateSellStopLimit: jsonGetFloat64(m, "MarginMaintenanceSellStopLimit"),
	}
	// Sessions
	sym.SessionsQuotes = parseSymbolSessions(m, "SessionsQuotes")
	sym.SessionsTrades = parseSymbolSessions(m, "SessionsTrades")
	return sym
}

// parseSymbolSessions parses session arrays from the JSON map.
// Sessions are nested arrays: [[{Open,OpenHours,Close,CloseHours},...],...]
func parseSymbolSessions(m map[string]interface{}, key string) [][]MTConSymbolSession {
	val, ok := m[key]
	if !ok || val == nil {
		return nil
	}
	outerArr, ok := val.([]interface{})
	if !ok {
		return nil
	}
	result := make([][]MTConSymbolSession, 0, len(outerArr))
	for _, dayVal := range outerArr {
		dayArr, ok := dayVal.([]interface{})
		if !ok {
			result = append(result, nil)
			continue
		}
		daySessions := make([]MTConSymbolSession, 0, len(dayArr))
		for _, sessVal := range dayArr {
			sessMap, ok := sessVal.(map[string]interface{})
			if !ok {
				continue
			}
			sess := MTConSymbolSession{
				Open:       jsonGetUint32(sessMap, "Open"),
				OpenHours:  jsonGetUint32(sessMap, "OpenHours"),
				Close:      jsonGetUint32(sessMap, "Close"),
				CloseHours: jsonGetUint32(sessMap, "CloseHours"),
			}
			daySessions = append(daySessions, sess)
		}
		result = append(result, daySessions)
	}
	return result
}

// SymbolToJSON serializes an MTConSymbol to JSON using JSONWriter, matching the C# MTSymbolJson.ToJson format.
func SymbolToJSON(symbol *MTConSymbol) string {
	w := NewJSONWriter()
	w.WriteBeginObject()
	w.WriteAttributeString("Symbol", symbol.Symbol)
	w.WriteAttributeString("Path", symbol.Path)
	w.WriteAttributeString("ISIN", symbol.ISIN)
	w.WriteAttributeString("Description", symbol.Description)
	w.WriteAttributeString("International", symbol.International)
	w.WriteAttributeString("Basis", symbol.Basis)
	w.WriteAttributeString("Source", symbol.Source)
	w.WriteAttributeString("Page", symbol.Page)
	w.WriteAttributeString("CurrencyBase", symbol.CurrencyBase)
	w.WriteAttributeUint32("CurrencyBaseDigits", symbol.CurrencyBaseDigits)
	w.WriteAttributeString("CurrencyProfit", symbol.CurrencyProfit)
	w.WriteAttributeUint32("CurrencyProfitDigits", symbol.CurrencyProfitDigits)
	w.WriteAttributeString("CurrencyMargin", symbol.CurrencyMargin)
	w.WriteAttributeUint32("CurrencyMarginDigits", symbol.CurrencyMarginDigits)
	w.WriteAttributeUint32("Color", symbol.Color)
	w.WriteAttributeUint32("ColorBackground", symbol.ColorBackground)
	w.WriteAttributeUint32("Digits", symbol.Digits)
	w.WriteAttributeFloat64("Point", symbol.Point)
	w.WriteAttributeFloat64("Multiply", symbol.Multiply)
	w.WriteAttributeUint64("TickFlags", uint64(symbol.TickFlags))
	w.WriteAttributeUint32("TickBookDepth", symbol.TickBookDepth)
	w.WriteAttributeUint32("TickChartMode", uint32(symbol.ChartMode))
	w.WriteAttributeUint32("FilterSoft", symbol.FilterSoft)
	w.WriteAttributeUint32("FilterSoftTicks", symbol.FilterSoftTicks)
	w.WriteAttributeUint32("FilterHard", symbol.FilterHard)
	w.WriteAttributeUint32("FilterHardTicks", symbol.FilterHardTicks)
	w.WriteAttributeUint32("FilterDiscard", symbol.FilterDiscard)
	w.WriteAttributeUint32("FilterSpreadMax", symbol.FilterSpreadMax)
	w.WriteAttributeUint32("FilterSpreadMin", symbol.FilterSpreadMin)
	w.WriteAttributeUint32("FilterGap", symbol.FilterGap)
	w.WriteAttributeUint32("FilterGapTicks", symbol.FilterGapTicks)
	w.WriteAttributeUint32("TradeMode", uint32(symbol.TradeMode))
	w.WriteAttributeUint64("TradeFlags", uint64(symbol.TradeFlags))
	w.WriteAttributeUint32("CalcMode", uint32(symbol.CalcMode))
	w.WriteAttributeUint32("ExecMode", uint32(symbol.ExecMode))
	w.WriteAttributeUint32("GTCMode", uint32(symbol.GTCMode))
	w.WriteAttributeUint32("FillFlags", uint32(symbol.FillFlags))
	w.WriteAttributeUint32("ExpirFlags", uint32(symbol.ExpirFlags))
	w.WriteAttributeUint32("OrderFlags", uint32(symbol.OrderFlags))
	w.WriteAttributeInt32("Spread", symbol.Spread)
	w.WriteAttributeInt32("SpreadBalance", symbol.SpreadBalance)
	w.WriteAttributeInt32("SpreadDiff", symbol.SpreadDiff)
	w.WriteAttributeInt32("SpreadDiffBalance", symbol.SpreadDiffBalance)
	w.WriteAttributeFloat64("TickValue", symbol.TickValue)
	w.WriteAttributeFloat64("TickSize", symbol.TickSize)
	w.WriteAttributeFloat64("ContractSize", symbol.ContractSize)
	w.WriteAttributeInt32("StopsLevel", symbol.StopsLevel)
	w.WriteAttributeInt32("FreezeLevel", symbol.FreezeLevel)
	w.WriteAttributeUint32("QuotesTimeout", symbol.QuotesTimeout)
	w.WriteAttributeUint64("VolumeMinExt", symbol.VolumeMinExt)
	w.WriteAttributeUint64("VolumeMaxExt", symbol.VolumeMaxExt)
	w.WriteAttributeUint64("VolumeStepExt", symbol.VolumeStepExt)
	w.WriteAttributeUint64("VolumeLimitExt", symbol.VolumeLimitExt)
	w.WriteAttributeUint32("MarginCheckMode", uint32(symbol.MarginCheckMode))
	w.WriteAttributeUint32("MarginFlags", uint32(symbol.MarginFlags))
	w.WriteAttributeFloat64("MarginInitial", symbol.MarginInitial)
	w.WriteAttributeFloat64("MarginMaintenance", symbol.MarginMaintenance)
	// MarginRateInitial
	if symbol.MarginRateInitial != nil {
		w.WriteAttributeFloat64("MarginInitialBuy", symbol.MarginRateInitial[MarginRateBuy])
		w.WriteAttributeFloat64("MarginInitialSell", symbol.MarginRateInitial[MarginRateSell])
		w.WriteAttributeFloat64("MarginInitialBuyLimit", symbol.MarginRateInitial[MarginRateBuyLimit])
		w.WriteAttributeFloat64("MarginInitialSellLimit", symbol.MarginRateInitial[MarginRateSellLimit])
		w.WriteAttributeFloat64("MarginInitialBuyStop", symbol.MarginRateInitial[MarginRateBuyStop])
		w.WriteAttributeFloat64("MarginInitialSellStop", symbol.MarginRateInitial[MarginRateSellStop])
		w.WriteAttributeFloat64("MarginInitialBuyStopLimit", symbol.MarginRateInitial[MarginRateBuyStopLimit])
		w.WriteAttributeFloat64("MarginInitialSellStopLimit", symbol.MarginRateInitial[MarginRateSellStopLimit])
	}
	// MarginRateMaintenance
	if symbol.MarginRateMaintenance != nil {
		w.WriteAttributeFloat64("MarginMaintenanceBuy", symbol.MarginRateMaintenance[MarginRateBuy])
		w.WriteAttributeFloat64("MarginMaintenanceSell", symbol.MarginRateMaintenance[MarginRateSell])
		w.WriteAttributeFloat64("MarginMaintenanceBuyLimit", symbol.MarginRateMaintenance[MarginRateBuyLimit])
		w.WriteAttributeFloat64("MarginMaintenanceSellLimit", symbol.MarginRateMaintenance[MarginRateSellLimit])
		w.WriteAttributeFloat64("MarginMaintenanceBuyStop", symbol.MarginRateMaintenance[MarginRateBuyStop])
		w.WriteAttributeFloat64("MarginMaintenanceSellStop", symbol.MarginRateMaintenance[MarginRateSellStop])
		w.WriteAttributeFloat64("MarginMaintenanceBuyStopLimit", symbol.MarginRateMaintenance[MarginRateBuyStopLimit])
		w.WriteAttributeFloat64("MarginMaintenanceSellStopLimit", symbol.MarginRateMaintenance[MarginRateSellStopLimit])
	}
	w.WriteAttributeFloat64("MarginLiquidity", symbol.MarginRateLiquidity)
	w.WriteAttributeFloat64("MarginHedged", symbol.MarginHedged)
	w.WriteAttributeFloat64("MarginCurrency", symbol.MarginRateCurrency)
	w.WriteAttributeFloat64("MarginLong", symbol.MarginLong)
	w.WriteAttributeFloat64("MarginShort", symbol.MarginShort)
	w.WriteAttributeFloat64("MarginLimit", symbol.MarginLimit)
	w.WriteAttributeFloat64("MarginStop", symbol.MarginStop)
	w.WriteAttributeFloat64("MarginStopLimit", symbol.MarginStopLimit)
	w.WriteAttributeUint32("SwapMode", uint32(symbol.SwapMode))
	w.WriteAttributeFloat64("SwapLong", symbol.SwapLong)
	w.WriteAttributeFloat64("SwapShort", symbol.SwapShort)
	w.WriteAttributeInt32("Swap3Day", symbol.Swap3Day)
	w.WriteAttributeInt64("TimeStart", symbol.TimeStart)
	w.WriteAttributeInt64("TimeExpiration", symbol.TimeExpiration)
	// Sessions
	w.WriteAttributeRaw("SessionsQuotes", sessionsToJSON(symbol.SessionsQuotes))
	w.WriteAttributeRaw("SessionsTrades", sessionsToJSON(symbol.SessionsTrades))
	w.WriteAttributeUint32("REFlags", uint32(symbol.REFlags))
	w.WriteAttributeUint32("RETimeout", symbol.RETimeout)
	w.WriteAttributeUint32("IECheckMode", uint32(symbol.IECheckMode))
	w.WriteAttributeUint32("IETimeout", symbol.IETimeout)
	w.WriteAttributeUint32("IESlipProfit", symbol.IESlipProfit)
	w.WriteAttributeUint32("IESlipLosing", symbol.IESlipLosing)
	w.WriteAttributeUint64("IEVolumeMaxExt", symbol.IEVolumeMaxExt)
	w.WriteAttributeFloat64("PriceSettle", symbol.PriceSettle)
	w.WriteAttributeFloat64("PriceLimitMax", symbol.PriceLimitMax)
	w.WriteAttributeFloat64("PriceLimitMin", symbol.PriceLimitMin)
	w.WriteAttributeFloat64("PriceStrike", symbol.PriceStrike)
	w.WriteAttributeUint32("OptionMode", uint32(symbol.OptionsMode))
	w.WriteAttributeFloat64("FaceValue", symbol.FaceValue)
	w.WriteAttributeFloat64("AccruedInterest", symbol.AccruedInterest)
	w.WriteAttributeUint32("SpliceType", uint32(symbol.SpliceType))
	w.WriteAttributeUint32("SpliceTimeType", uint32(symbol.SpliceTimeType))
	w.WriteAttributeUint32("SpliceTimeDays", symbol.SpliceTimeDays)
	w.WriteAttributeUint32("IEFlags", uint32(symbol.IEFlags))
	w.WriteAttributeString("Category", symbol.Category)
	w.WriteAttributeString("Exchange", symbol.Exchange)
	w.WriteAttributeString("CFI", symbol.CFI)
	w.WriteAttributeUint32("Sector", uint32(symbol.Sector))
	w.WriteAttributeUint32("Industry", uint32(symbol.Industry))
	w.WriteAttributeString("Country", symbol.Country)
	w.WriteAttributeUint32("SubscriptionsDelay", symbol.SubscriptionsDelay)
	w.WriteEndObject()
	return w.String()
}

// sessionsToJSON serializes sessions as nested JSON arrays: [[{Open,OpenHours,Close,CloseHours},...],...]
func sessionsToJSON(sessions [][]MTConSymbolSession) string {
	w := NewJSONWriter()
	w.writeBeginCollection()
	for i, day := range sessions {
		if i > 0 {
			w.isFirstElement = false
		}
		innerW := NewJSONWriter()
		innerW.writeBeginCollection()
		for j, sess := range day {
			if j > 0 {
				innerW.isFirstElement = false
			}
			innerW.WriteBeginObject()
			innerW.WriteAttributeUint32("Open", sess.Open)
			innerW.WriteAttributeUint32("OpenHours", sess.OpenHours)
			innerW.WriteAttributeUint32("Close", sess.Close)
			innerW.WriteAttributeUint32("CloseHours", sess.CloseHours)
			innerW.WriteEndObject()
			innerW.isFirstElement = false
		}
		innerW.writeEndCollection()
		if !w.isFirstElement {
			w.builder.WriteByte(',')
		}
		w.builder.WriteString(innerW.String())
		w.isFirstElement = false
	}
	w.writeEndCollection()
	return w.String()
}
