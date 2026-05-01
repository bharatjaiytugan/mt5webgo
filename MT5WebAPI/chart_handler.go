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

// MTBar represents a single 1-minute OHLC bar.
type MTBar struct {
	Datetime int64
	Open     float64
	High     float64
	Low      float64
	Close    float64
	TickVol  uint64
	RealVol  uint64
	Spread   uint32
}

// MTBookItem represents a single entry in the market depth (DOM).
type MTBookItem struct {
	Type   uint32  // 1=sell, 2=buy, etc.
	Price  float64
	Volume uint64
}

// chartHandler implements chart/price-related protocol commands.
type chartHandler struct {
	apiBase
}

// ChartGet retrieves 1-minute bar history for a symbol.
// data specifies which fields to return: d=datetime, o=open, h=high, l=low, c=close,
// t=tick volume, v=real volume, s=spread. Default is "dohlc".
// Response is a JSON array of arrays: [[timestamp,open,high,low,close], ...]
func (h *chartHandler) ChartGet(symbol string, from, to int64, data string) ([]*MTBar, MTRetCode) {
	if data == "" {
		data = "dohlc"
	}
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbol,
		WEB_PARAM_FROM:   formatInt64(from),
		WEB_PARAM_TO:     formatInt64(to),
		WEB_PARAM_DATA:   data,
	}
	answer := h.send(WEB_CMD_CHART_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send chart get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	jsonStr, errCode := parseTolerantJSON(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse chart get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	bars := parseBarsFromJSON(jsonStr, data)
	if bars == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return bars, MT_RET_OK
}

// BookGet retrieves the market depth (DOM) for one or more symbols (comma-separated).
// Response is a JSON map: {"EURUSD": [{"type":"1","price":"1.12345","volume":"100"}, ...]}
func (h *chartHandler) BookGet(symbol string) (map[string][]MTBookItem, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SYMBOL: symbol,
	}
	answer := h.send(WEB_CMD_BOOK_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send book get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	jsonStr, errCode := parseTolerantJSON(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse book get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	book := parseBookFromJSON(jsonStr)
	if book == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return book, MT_RET_OK
}

// parseBarsFromJSON parses bar data from a JSON array of arrays.
// Each inner array contains values corresponding to the data parameter letters.
func parseBarsFromJSON(jsonStr string, data string) []*MTBar {
	var raw [][]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing bars from json failed: %v", err))
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	var bars []*MTBar
	for _, row := range raw {
		bar := parseBarFromArray(row, data)
		if bar != nil {
			bars = append(bars, bar)
		}
	}
	return bars
}

// parseBarFromArray parses a single bar from a JSON array based on the data flags.
func parseBarFromArray(arr []interface{}, data string) *MTBar {
	if len(arr) == 0 {
		return nil
	}
	bar := &MTBar{}
	idx := 0
	for _, ch := range data {
		if idx >= len(arr) {
			break
		}
		val := arr[idx]
		switch ch {
		case 'd':
			bar.Datetime = toInt64(val)
		case 'o':
			bar.Open = toFloat64(val)
		case 'h':
			bar.High = toFloat64(val)
		case 'l':
			bar.Low = toFloat64(val)
		case 'c':
			bar.Close = toFloat64(val)
		case 't':
			bar.TickVol = toUint64(val)
		case 'v':
			bar.RealVol = toUint64(val)
		case 's':
			bar.Spread = uint32(toUint64(val))
		}
		idx++
	}
	return bar
}

// parseBookFromJSON parses market depth data from a JSON map of symbol -> array of book items.
func parseBookFromJSON(jsonStr string) map[string][]MTBookItem {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing book from json failed: %v", err))
		return nil
	}
	if len(raw) == 0 {
		return nil
	}
	result := make(map[string][]MTBookItem)
	for sym, val := range raw {
		arr, ok := val.([]interface{})
		if !ok {
			continue
		}
		var items []MTBookItem
		for _, entry := range arr {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			item := MTBookItem{
				Type:   uint32(toUint64(m["type"])),
				Price:  toFloat64(m["price"]),
				Volume: toUint64(m["volume"]),
			}
			items = append(items, item)
		}
		result[sym] = items
	}
	return result
}

// toFloat64 converts an interface{} value (float64, string, json.Number) to float64.
func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case string:
		r, _ := strconv.ParseFloat(val, 64)
		return r
	case json.Number:
		r, _ := val.Float64()
		return r
	}
	return 0
}

// toInt64 converts an interface{} value to int64.
func toInt64(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int64(val)
	case string:
		r, _ := strconv.ParseInt(val, 10, 64)
		return r
	case json.Number:
		r, _ := val.Int64()
		return r
	}
	return 0
}

// toUint64 converts an interface{} value to uint64.
func toUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return uint64(val)
	case string:
		r, _ := strconv.ParseUint(val, 10, 64)
		return r
	case json.Number:
		r, _ := val.Int64()
		return uint64(r)
	}
	return 0
}
