//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
	"strconv"
	"strings"
)

// Trade calculation command constants.
const (
	WEB_CMD_TRADE_RATE_BUY     = "TRADE_RATE_BUY"
	WEB_CMD_TRADE_RATE_SELL    = "TRADE_RATE_SELL"
	WEB_CMD_TRADE_MARGIN_CHECK = "TRADE_MARGIN_CHECK"
	WEB_CMD_TRADE_PROFIT       = "TRADE_PROFIT"
)

// Trade calculation parameter constants.
const (
	WEB_PARAM_BASE        = "BASE"
	WEB_PARAM_CURRENCY    = "CURRENCY"
	WEB_PARAM_PRICE       = "PRICE"
	WEB_PARAM_VOLUME      = "VOLUME"
	WEB_PARAM_PRICE_OPEN  = "PRICE_OPEN"
	WEB_PARAM_PRICE_CLOSE = "PRICE_CLOSE"
)

// tradeCalcHandler implements trade calculation protocol commands.
type tradeCalcHandler struct {
	apiBase
}

// CalcRateBuy calculates the buy conversion rate between two currencies.
// Returns the rate and a return code.
func (h *tradeCalcHandler) CalcRateBuy(base, currency string) (float64, MTRetCode) {
	return h.calcRate(WEB_CMD_TRADE_RATE_BUY, base, currency)
}

// CalcRateSell calculates the sell conversion rate between two currencies.
// Returns the rate and a return code.
func (h *tradeCalcHandler) CalcRateSell(base, currency string) (float64, MTRetCode) {
	return h.calcRate(WEB_CMD_TRADE_RATE_SELL, base, currency)
}

// calcRate is the common implementation for CalcRateBuy and CalcRateSell.
func (h *tradeCalcHandler) calcRate(command, base, currency string) (float64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_BASE:     base,
		WEB_PARAM_CURRENCY: currency,
	}
	answer := h.send(command, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", command))
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer %s: %s", command, answerStr))
	}
	// parse answer with tolerant command name checking
	rate, errCode := parseRateAnswer(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", command, GetErrorStandard(errCode)))
		return 0, errCode
	}
	return rate, MT_RET_OK
}

// CheckMargin checks the margin required for an order.
// orderType: 0=buy, 1=sell. volume is in 1/10000 lots (e.g. 10000 = 1.0 lot).
// Returns a map with "new" and "current" account state objects.
func (h *tradeCalcHandler) CheckMargin(login uint64, symbol string, orderType uint32, volume uint64, price float64) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN:  formatUint64(login),
		WEB_PARAM_SYMBOL: symbol,
		WEB_PARAM_TYPE:   formatUint32(orderType),
		WEB_PARAM_VOLUME: formatUint64(volume),
		WEB_PARAM_PRICE:  fmt.Sprintf("%.6f", price),
	}
	answer := h.send(WEB_CMD_TRADE_MARGIN_CHECK, params)
	if answer == nil {
		logWrite(LogTypeError, "send trade margin check failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer trade margin check: %s", answerStr))
	}
	// parse with tolerant command name
	jsonStr, errCode := parseTolerantJSON(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse trade margin check failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse margin check json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return m, MT_RET_OK
}

// CalcProfit calculates the profit for a trade.
// orderType: 0=buy, 1=sell. volume is in 1/10000 lots (e.g. 10000 = 1.0 lot).
func (h *tradeCalcHandler) CalcProfit(group, symbol string, orderType uint32, volume uint64, priceOpen, priceClose float64) (float64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_GROUP:       group,
		WEB_PARAM_SYMBOL:      symbol,
		WEB_PARAM_TYPE:        formatUint32(orderType),
		WEB_PARAM_VOLUME:      formatUint64(volume),
		WEB_PARAM_PRICE_OPEN:  fmt.Sprintf("%.6f", priceOpen),
		WEB_PARAM_PRICE_CLOSE: fmt.Sprintf("%.6f", priceClose),
	}
	answer := h.send(WEB_CMD_TRADE_PROFIT, params)
	if answer == nil {
		logWrite(LogTypeError, "send trade profit failed")
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer trade profit: %s", answerStr))
	}
	// parse with tolerant command name
	profit, errCode := parseProfitAnswer(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse trade profit failed: %s", GetErrorStandard(errCode)))
		return 0, errCode
	}
	return profit, MT_RET_OK
}

// parseTolerantJSON extracts RETCODE and JSON body from an answer string,
// accepting any command name returned by the server.
func parseTolerantJSON(answer string) (string, MTRetCode) {
	_, pos := getCommand(answer)
	if pos < 0 {
		return "", MT_RET_ERR_DATA
	}
	ba := &baseAnswerJSON{}
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
	errCode := getRetCode(ba.RetCode)
	if errCode != MT_RET_OK {
		return "", errCode
	}
	// try getJSON first, fall back to finding { directly
	jsonStr := getJSON(answer, posEnd)
	if jsonStr == "" {
		// try to find JSON object directly
		idx := strings.Index(answer, "{")
		if idx >= 0 {
			jsonStr = strings.TrimSpace(answer[idx:])
		}
	}
	if jsonStr == "" {
		return "", MT_RET_REPORT_NODATA
	}
	return jsonStr, MT_RET_OK
}

// parseRateAnswer parses a rate response (Rate field from JSON body).
func parseRateAnswer(answer string) (float64, MTRetCode) {
	jsonStr, errCode := parseTolerantJSON(answer)
	if errCode != MT_RET_OK {
		return 0, errCode
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse rate json failed: %v", err))
		return 0, MT_RET_ERR_DATA
	}
	rateStr := jsonGetString(m, "Rate")
	if rateStr == "" {
		rateStr = jsonGetString(m, "rate")
	}
	if rateStr == "" {
		return 0, MT_RET_REPORT_NODATA
	}
	rate, err := strconv.ParseFloat(rateStr, 64)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse rate value failed: %v", err))
		return 0, MT_RET_ERR_DATA
	}
	return rate, MT_RET_OK
}

// parseProfitAnswer parses a profit response (Profit field from JSON body).
func parseProfitAnswer(answer string) (float64, MTRetCode) {
	jsonStr, errCode := parseTolerantJSON(answer)
	if errCode != MT_RET_OK {
		return 0, errCode
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse profit json failed: %v", err))
		return 0, MT_RET_ERR_DATA
	}
	profitStr := jsonGetString(m, "Profit")
	if profitStr == "" {
		profitStr = jsonGetString(m, "profit")
	}
	if profitStr == "" {
		return 0, MT_RET_REPORT_NODATA
	}
	profit, err := strconv.ParseFloat(profitStr, 64)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse profit value failed: %v", err))
		return 0, MT_RET_ERR_DATA
	}
	return profit, MT_RET_OK
}
