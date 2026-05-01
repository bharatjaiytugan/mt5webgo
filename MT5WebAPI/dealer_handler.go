//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// Dealer protocol command constants.
const (
	WEB_CMD_DEALER_SEND       = "DEALER_SEND"
	WEB_CMD_DEALER_UPDATES = "DEALER_UPDATES"
)

// EnTradeAction represents dealer trade action types.
type EnTradeAction uint32

const (
	TA_DEALER_POS_EXECUTE  EnTradeAction = 200 // execute position (market buy/sell)
	TA_DEALER_ORD_PENDING  EnTradeAction = 201 // place pending order
	TA_DEALER_POS_MODIFY   EnTradeAction = 202 // modify position
	TA_DEALER_ORD_MODIFY   EnTradeAction = 203 // modify pending order
	TA_DEALER_ORD_REMOVE   EnTradeAction = 204 // remove pending order
	TA_DEALER_ORD_ACTIVATE EnTradeAction = 205 // activate pending order
	TA_DEALER_BALANCE      EnTradeAction = 206 // balance operation
	TA_DEALER_ORD_SLIMIT   EnTradeAction = 207 // stop-limit order
	TA_DEALER_CLOSE_BY     EnTradeAction = 208 // close by
)

// TradeRequest represents a dealer trade request.
type TradeRequest struct {
	Action         EnTradeAction
	Login          uint64
	Symbol         string
	Volume         uint64  // in 1/10000 lots (100 = 0.01 lot)
	Type           uint32  // EnOrderType value
	TypeFill       uint32  // EnOrderFilling value
	TypeTime       uint32  // EnOrderTime value
	PriceOrder     float64
	PriceTrigger   float64 // for stop-limit orders
	PriceSL        float64
	PriceTP        float64
	PriceDeviation uint32
	TimeExpiration int64
	Order          uint64 // ticket for modify/cancel
	Position       uint64 // position ticket
	Comment        string
	Digits         uint32
}

// TradeResult represents the result of a dealer trade operation.
type TradeResult struct {
	ID            uint64
	RetCode       string
	ResultRetcode string
	ResultDeal    uint64
	ResultOrder   uint64
	ResultVolume  uint64
	ResultPrice   float64
	ResultComment string
}

// dealerHandler implements dealer trade protocol commands.
type dealerHandler struct {
	apiBase
}

// DealerSend sends a dealer trade request to the MT5 server.
func (h *dealerHandler) DealerSend(req *TradeRequest) (*TradeResult, MTRetCode) {
	if req == nil {
		logWrite(LogTypeError, "dealer send: request is nil")
		return nil, MT_RET_ERR_PARAMS
	}
	// serialize request to JSON
	body := tradeRequestToJSON(req)
	params := map[string]string{
		WEB_PARAM_BODYTEXT: body,
	}
	answer := h.send(WEB_CMD_DEALER_SEND, params)
	if answer == nil {
		logWrite(LogTypeError, "send dealer request failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer dealer send: %s", answerStr))
	}
	// parse response
	result, errCode := parseDealerSendAnswer(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse dealer send failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	return result, MT_RET_OK
}

// DealerGetResult retrieves the result of a previously submitted dealer request.
func (h *dealerHandler) DealerGetResult(id uint64) (*TradeResult, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(id),
	}
	answer := h.send(WEB_CMD_DEALER_UPDATES, params)
	if answer == nil {
		logWrite(LogTypeError, "send dealer get result failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer dealer get result: %s", answerStr))
	}
	// parse response
	result, errCode := parseDealerResultAnswer(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse dealer get result failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	return result, MT_RET_OK
}

// tradeRequestToJSON serializes a TradeRequest to JSON using JSONWriter.
// Only non-zero/non-empty fields are included. Numeric values are written as strings (MT5 convention).
func tradeRequestToJSON(req *TradeRequest) string {
	w := NewJSONWriter()
	w.WriteBeginObject()
	// Action is always required
	w.WriteAttributeString("Action", formatUint32(uint32(req.Action)))
	// Login is always required
	if req.Login != 0 {
		w.WriteAttributeString("Login", formatUint64(req.Login))
	}
	// Symbol
	if req.Symbol != "" {
		w.WriteAttributeString("Symbol", req.Symbol)
	}
	// Volume
	if req.Volume != 0 {
		w.WriteAttributeString("Volume", formatUint64(req.Volume))
	}
	// Type
	w.WriteAttributeString("Type", formatUint32(req.Type))
	// TypeFill
	w.WriteAttributeString("TypeFill", formatUint32(req.TypeFill))
	// TypeTime
	if req.TypeTime != 0 {
		w.WriteAttributeString("TypeTime", formatUint32(req.TypeTime))
	}
	// PriceOrder
	if req.PriceOrder != 0 {
		w.WriteAttributeString("PriceOrder", formatFloat64(req.PriceOrder))
	}
	// PriceTrigger
	if req.PriceTrigger != 0 {
		w.WriteAttributeString("PriceTrigger", formatFloat64(req.PriceTrigger))
	}
	// PriceSL
	if req.PriceSL != 0 {
		w.WriteAttributeString("PriceSL", formatFloat64(req.PriceSL))
	}
	// PriceTP
	if req.PriceTP != 0 {
		w.WriteAttributeString("PriceTP", formatFloat64(req.PriceTP))
	}
	// PriceDeviation
	if req.PriceDeviation != 0 {
		w.WriteAttributeString("PriceDeviation", formatUint32(req.PriceDeviation))
	}
	// TimeExpiration
	if req.TimeExpiration != 0 {
		w.WriteAttributeString("TimeExpiration", formatInt64(req.TimeExpiration))
	}
	// Order (ticket for modify/cancel)
	if req.Order != 0 {
		w.WriteAttributeString("Order", formatUint64(req.Order))
	}
	// Position
	if req.Position != 0 {
		w.WriteAttributeString("Position", formatUint64(req.Position))
	}
	// Comment
	if req.Comment != "" {
		w.WriteAttributeString("Comment", req.Comment)
	}
	// Digits
	w.WriteAttributeString("Digits", formatUint32(req.Digits))
	w.WriteEndObject()
	return w.String()
}

// parseDealerSendAnswer parses the DEALER_SEND response.
func parseDealerSendAnswer(answer string) (*TradeResult, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != WEB_CMD_DEALER_SEND {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, WEB_CMD_DEALER_SEND))
		return nil, MT_RET_ERR_DATA
	}
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return nil, errCode
	}
	// parse JSON body
	jsonStr := getJSON(answer, posEnd)
	if jsonStr == "" {
		return nil, MT_RET_REPORT_NODATA
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing dealer send result json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	result := &TradeResult{}
	result.ID = jsonGetUint64(m, "id")
	if result.ID == 0 {
		result.ID = jsonGetUint64(m, "ID")
	}
	result.RetCode = jsonGetString(m, "RetCode")
	result.ResultRetcode = jsonGetString(m, "ResultRetcode")
	result.ResultDeal = jsonGetUint64(m, "ResultDeal")
	result.ResultOrder = jsonGetUint64(m, "ResultOrder")
	result.ResultVolume = jsonGetUint64(m, "ResultVolume")
	result.ResultPrice = jsonGetFloat64(m, "ResultPrice")
	result.ResultComment = jsonGetString(m, "ResultComment")
	return result, MT_RET_OK
}

// parseDealerResultAnswer parses the DEALER_GET_RESULT response.
func parseDealerResultAnswer(answer string) (*TradeResult, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != WEB_CMD_DEALER_UPDATES {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, WEB_CMD_DEALER_UPDATES))
		return nil, MT_RET_ERR_DATA
	}
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return nil, errCode
	}
	// parse JSON body
	jsonStr := getJSON(answer, posEnd)
	if jsonStr == "" {
		return nil, MT_RET_REPORT_NODATA
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing dealer get result json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	result := &TradeResult{}
	result.ID = jsonGetUint64(m, "id")
	if result.ID == 0 {
		result.ID = jsonGetUint64(m, "ID")
	}
	result.RetCode = jsonGetString(m, "RetCode")
	result.ResultRetcode = jsonGetString(m, "ResultRetcode")
	result.ResultDeal = jsonGetUint64(m, "ResultDeal")
	result.ResultOrder = jsonGetUint64(m, "ResultOrder")
	result.ResultVolume = jsonGetUint64(m, "ResultVolume")
	result.ResultPrice = jsonGetFloat64(m, "ResultPrice")
	result.ResultComment = jsonGetString(m, "ResultComment")
	return result, MT_RET_OK
}
