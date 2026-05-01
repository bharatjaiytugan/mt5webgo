//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// tickHandler implements tick-related protocol commands.
type tickHandler struct {
	apiBase
}

// TickLast retrieves the last ticks for a symbol.
func (h *tickHandler) TickLast(symbol string) ([]*MTTick, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbol,
	}
	answer := h.send(WEB_CMD_TICK_LAST, params)
	if answer == nil {
		logWrite(LogTypeError, "send ticks last failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TICK_LAST, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse ticks last failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	ticks := parseTicksFromJSON(jsonStr)
	if ticks == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return ticks, MT_RET_OK
}

// TickLastGroup retrieves the last ticks for a symbol filtered by group.
func (h *tickHandler) TickLastGroup(symbol, group string) ([]*MTTick, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbol,
		WEB_PARAM_GROUP:  group,
	}
	answer := h.send(WEB_CMD_TICK_LAST_GROUP, params)
	if answer == nil {
		logWrite(LogTypeError, "send ticks last group failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TICK_LAST_GROUP, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse ticks group last failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	ticks := parseTicksFromJSON(jsonStr)
	if ticks == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return ticks, MT_RET_OK
}

// TickStat retrieves tick statistics for a symbol.
func (h *tickHandler) TickStat(symbol string) ([]*MTTickStat, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbol,
	}
	answer := h.send(WEB_CMD_TICK_STAT, params)
	if answer == nil {
		logWrite(LogTypeError, "send tick stat failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TICK_STAT, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse tick stat failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	stats := parseTickStatsFromJSON(jsonStr)
	if stats == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return stats, MT_RET_OK
}

// parseTicksFromJSON parses a list of MTTick from a JSON array string.
func parseTicksFromJSON(jsonStr string) []*MTTick {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing ticks from json failed: %v", err))
		return nil
	}
	var ticks []*MTTick
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			t := parseTickFromMap(m)
			if t != nil {
				ticks = append(ticks, t)
			}
		}
	}
	return ticks
}

func parseTickFromMap(m map[string]interface{}) *MTTick {
	if m == nil {
		return nil
	}
	t := &MTTick{}
	t.Symbol = jsonGetString(m, "Symbol")
	t.Digits = jsonGetUint32(m, "Digits")
	t.Bid = jsonGetFloat64(m, "Bid")
	t.Ask = jsonGetFloat64(m, "Ask")
	t.Last = jsonGetFloat64(m, "Last")
	t.Volume = jsonGetUint64(m, "Volume")
	if _, ok := m["VolumeReal"]; ok {
		t.VolumeReal = jsonGetFloat64(m, "VolumeReal")
	} else {
		t.VolumeReal = float64(t.Volume)
	}
	t.Datetime = jsonGetUint64(m, "Datetime")
	if _, ok := m["DatetimeMsc"]; ok {
		t.DatetimeMsc = jsonGetUint64(m, "DatetimeMsc")
	} else {
		t.DatetimeMsc = t.Datetime * 1000
	}
	return t
}

// parseTickStatsFromJSON parses a list of MTTickStat from a JSON array string.
func parseTickStatsFromJSON(jsonStr string) []*MTTickStat {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing tick stats from json failed: %v", err))
		return nil
	}
	var stats []*MTTickStat
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			s := parseTickStatFromMap(m)
			if s != nil {
				stats = append(stats, s)
			}
		}
	}
	return stats
}

func parseTickStatFromMap(m map[string]interface{}) *MTTickStat {
	if m == nil {
		return nil
	}
	s := &MTTickStat{}
	s.Symbol = jsonGetString(m, "Symbol")
	s.Digits = jsonGetUint32(m, "Digits")
	s.Bid = jsonGetFloat64(m, "Bid")
	s.BidLow = jsonGetFloat64(m, "BidLow")
	s.BidHigh = jsonGetFloat64(m, "BidHigh")
	s.BidDir = EnTickStatDirection(jsonGetUint32(m, "BidDir"))
	s.Ask = jsonGetFloat64(m, "Ask")
	s.AskLow = jsonGetFloat64(m, "AskLow")
	s.AskHigh = jsonGetFloat64(m, "AskHigh")
	s.AskDir = EnTickStatDirection(jsonGetUint32(m, "AskDir"))
	s.Last = jsonGetFloat64(m, "Last")
	s.LastLow = jsonGetFloat64(m, "LastLow")
	s.LastHigh = jsonGetFloat64(m, "LastHigh")
	s.LastDir = EnTickStatDirection(jsonGetUint32(m, "LastDir"))
	s.Volume = jsonGetUint64(m, "Volume")
	if _, ok := m["VolumeReal"]; ok {
		s.VolumeReal = jsonGetFloat64(m, "VolumeReal")
	} else {
		s.VolumeReal = float64(s.Volume)
	}
	s.VolumeLow = jsonGetUint64(m, "VolumeLow")
	if _, ok := m["VolumeLowReal"]; ok {
		s.VolumeLowReal = jsonGetFloat64(m, "VolumeLowReal")
	} else {
		s.VolumeLowReal = float64(s.VolumeLow)
	}
	s.VolumeHigh = jsonGetUint64(m, "VolumeHigh")
	if _, ok := m["VolumeHighReal"]; ok {
		s.VolumeHighReal = jsonGetFloat64(m, "VolumeHighReal")
	} else {
		s.VolumeHighReal = float64(s.VolumeHigh)
	}
	s.VolumeDir = EnTickStatDirection(jsonGetUint32(m, "VolumeDir"))
	s.TradeDeals = jsonGetUint64(m, "TradeDeals")
	s.TradeVolume = jsonGetUint64(m, "TradeVolume")
	if _, ok := m["TradeVolumeReal"]; ok {
		s.TradeVolumeReal = jsonGetFloat64(m, "TradeVolumeReal")
	} else {
		s.TradeVolumeReal = float64(s.TradeVolume)
	}
	s.TradeTurnover = jsonGetUint64(m, "TradeTurnover")
	s.TradeInterest = jsonGetUint64(m, "TradeInterest")
	s.TradeBuyOrders = jsonGetUint64(m, "TradeBuyOrders")
	s.TradeBuyVolume = jsonGetUint64(m, "TradeBuyVolume")
	if _, ok := m["TradeBuyVolumeReal"]; ok {
		s.TradeBuyVolumeReal = jsonGetFloat64(m, "TradeBuyVolumeReal")
	} else {
		s.TradeBuyVolumeReal = float64(s.TradeBuyVolume)
	}
	s.TradeSellOrders = jsonGetUint64(m, "TradeSellOrders")
	s.TradeSellVolume = jsonGetUint64(m, "TradeSellVolume")
	if _, ok := m["TradeSellVolumeReal"]; ok {
		s.TradeSellVolumeReal = jsonGetFloat64(m, "TradeSellVolumeReal")
	} else {
		s.TradeSellVolumeReal = float64(s.TradeSellVolume)
	}
	s.PriceOpen = jsonGetFloat64(m, "PriceOpen")
	s.PriceClose = jsonGetFloat64(m, "PriceClose")
	s.PriceChange = jsonGetFloat64(m, "PriceChange")
	s.PriceVolatility = jsonGetFloat64(m, "PriceVolatility")
	s.PriceTheoretical = jsonGetFloat64(m, "PriceTheoretical")
	return s
}

// --- Stub types for MTTick and MTTickStat if not yet defined elsewhere ---
// These are used by the tick handler. If already defined in the project,
// remove these and use the existing definitions.

// MTTick represents a price tick.
type MTTick struct {
	Symbol      string  `json:"Symbol"`
	Digits      uint32  `json:"Digits"`
	Bid         float64 `json:"Bid"`
	Ask         float64 `json:"Ask"`
	Last        float64 `json:"Last"`
	Volume      uint64  `json:"Volume"`
	VolumeReal  float64 `json:"VolumeReal"`
	Datetime    uint64  `json:"Datetime"`
	DatetimeMsc uint64  `json:"DatetimeMsc"`
}

// EnTickStatDirection represents tick price direction.
type EnTickStatDirection uint32

const (
	TickStatDirectionNone  EnTickStatDirection = 0
	TickStatDirectionUp    EnTickStatDirection = 1
	TickStatDirectionDown  EnTickStatDirection = 2
	TickStatDirectionFirst                     = TickStatDirectionNone
	TickStatDirectionLast                      = TickStatDirectionDown
)

// MTTickStat represents tick statistics.
type MTTickStat struct {
	Symbol           string              `json:"Symbol"`
	Digits           uint32              `json:"Digits"`
	Bid              float64             `json:"Bid"`
	BidLow           float64             `json:"BidLow"`
	BidHigh          float64             `json:"BidHigh"`
	BidDir           EnTickStatDirection `json:"BidDir"`
	Ask              float64             `json:"Ask"`
	AskLow           float64             `json:"AskLow"`
	AskHigh          float64             `json:"AskHigh"`
	AskDir           EnTickStatDirection `json:"AskDir"`
	Last             float64             `json:"Last"`
	LastLow          float64             `json:"LastLow"`
	LastHigh         float64             `json:"LastHigh"`
	LastDir          EnTickStatDirection `json:"LastDir"`
	Volume           uint64              `json:"Volume"`
	VolumeReal       float64             `json:"VolumeReal"`
	VolumeLow        uint64              `json:"VolumeLow"`
	VolumeLowReal    float64             `json:"VolumeLowReal"`
	VolumeHigh       uint64              `json:"VolumeHigh"`
	VolumeHighReal   float64             `json:"VolumeHighReal"`
	VolumeDir        EnTickStatDirection `json:"VolumeDir"`
	TradeDeals       uint64              `json:"TradeDeals"`
	TradeVolume      uint64              `json:"TradeVolume"`
	TradeVolumeReal  float64             `json:"TradeVolumeReal"`
	TradeTurnover    uint64              `json:"TradeTurnover"`
	TradeInterest    uint64              `json:"TradeInterest"`
	TradeBuyOrders   uint64              `json:"TradeBuyOrders"`
	TradeBuyVolume   uint64              `json:"TradeBuyVolume"`
	TradeBuyVolumeReal  float64          `json:"TradeBuyVolumeReal"`
	TradeSellOrders  uint64              `json:"TradeSellOrders"`
	TradeSellVolume  uint64              `json:"TradeSellVolume"`
	TradeSellVolumeReal float64          `json:"TradeSellVolumeReal"`
	PriceOpen        float64             `json:"PriceOpen"`
	PriceClose       float64             `json:"PriceClose"`
	PriceChange      float64             `json:"PriceChange"`
	PriceVolatility  float64             `json:"PriceVolatility"`
	PriceTheoretical float64             `json:"PriceTheoretical"`
}

