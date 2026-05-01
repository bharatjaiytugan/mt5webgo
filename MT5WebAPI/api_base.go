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

// apiBase is the base type embedded by all protocol handlers.
// It holds a reference to the async connection for communication with the MT server.
type apiBase struct {
	asyncConnect *mtAsyncConnect
	sendFunc     func(command string, params map[string]string) ([]byte, error)
}

// newAPIBase creates a new apiBase with the given async connection.
func newAPIBase(connect *mtAsyncConnect) apiBase {
	return apiBase{asyncConnect: connect}
}

// sendViaAsync sends a command with parameters via the async connection and returns the response.
func (a *apiBase) sendViaAsync(command string, params map[string]string, isNeedReceived bool) []byte {
	asyncSend := newMTAsyncSend()
	if !asyncSend.BeginSend(a.asyncConnect, command, params, isNeedReceived) {
		logWrite(LogTypeError, fmt.Sprintf("command '%s' did not add to queue", command))
		return nil
	}
	return asyncSend.EndSend()
}

// send sends a command with parameters to the MT5 server and returns the raw response bytes.
// If an asyncConnect is set, it uses the async path; otherwise falls back to sendFunc.
func (a *apiBase) send(command string, params map[string]string) []byte {
	if a.asyncConnect != nil {
		return a.sendViaAsync(command, params, true)
	}
	if a.sendFunc == nil {
		logWrite(LogTypeError, "send function not set")
		return nil
	}
	result, err := a.sendFunc(command, params)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed: %v", command, err))
		return nil
	}
	return result
}

// sendNoRecv sends a command with parameters without expecting a response (fire and forget).
func (a *apiBase) sendNoRecv(command string, params map[string]string) {
	if a.asyncConnect != nil {
		a.sendViaAsync(command, params, false)
		return
	}
	if a.sendFunc == nil {
		logWrite(LogTypeError, "send function not set")
		return
	}
	_, _ = a.sendFunc(command, params)
}

// setConnection sets a new async connection for the API base.
func (a *apiBase) setConnection(connect *mtAsyncConnect) {
	a.asyncConnect = connect
}

// parseEmptyResult checks a simple answer that only contains RETCODE (no JSON body).
func (a *apiBase) parseEmptyResult(command, answer string) MTRetCode {
	cmdReal, pos := getCommand(answer)
	if cmdReal != command {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, command))
		return MT_RET_ERR_DATA
	}
	ba := newBaseAnswer()
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
	return getRetCode(ba.RetCode)
}

// parseSingleJSON parses a standard answer: check command, extract retcode, extract JSON body.
func parseSingleJSON(command, answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != command {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, command))
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
	jsonStr := getJSON(answer, posEnd)
	if jsonStr == "" {
		return "", MT_RET_REPORT_NODATA
	}
	return jsonStr, MT_RET_OK
}

// parseTotalAnswer parses an answer that returns RETCODE and TOTAL.
func parseTotalAnswer(command, answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != command {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, command))
		return "", MT_RET_ERR_DATA
	}
	retCode := ""
	total := ""
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			retCode = param.Value
		case WEB_PARAM_TOTAL:
			total = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return "", errCode
	}
	return total, MT_RET_OK
}

// parseTicketAnswer parses an answer that returns RETCODE and TICKET.
func parseTicketAnswer(command, answer string) (uint64, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != command {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, command))
		return 0, MT_RET_ERR_DATA
	}
	retCode := ""
	var ticket uint64
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			retCode = param.Value
		case WEB_PARAM_TICKET:
			t, err := strconv.ParseUint(param.Value, 10, 64)
			if err == nil {
				ticket = t
			}
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return 0, errCode
	}
	return ticket, MT_RET_OK
}

// parseTimeServerAnswer parses the TIME_SERVER response, extracting RETCODE and TIME.
func parseTimeServerAnswer(answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != WEB_CMD_TIME_SERVER {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", cmdReal, WEB_CMD_TIME_SERVER))
		return "", MT_RET_ERR_DATA
	}
	retCode := ""
	timeVal := "none"
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			retCode = param.Value
		case WEB_PARAM_TIME:
			timeVal = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		return "", errCode
	}
	if timeVal == "" || timeVal == "none" {
		return "", MT_RET_ERR_PARAMS
	}
	return timeVal, MT_RET_OK
}

// jsonToMap parses a JSON string into a map[string]interface{}.
func jsonToMap(data string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// jsonToSlice parses a JSON array string into a []interface{}.
func jsonToSlice(data string) ([]interface{}, error) {
	var s []interface{}
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, err
	}
	return s, nil
}

// --- JSON field extraction helpers ---

func jsonGetString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

func jsonGetUint64(m map[string]interface{}, key string) uint64 {
	v, ok := m[key]
	if !ok || v == nil {
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

func jsonGetInt64(m map[string]interface{}, key string) int64 {
	v, ok := m[key]
	if !ok || v == nil {
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

func jsonGetUint32(m map[string]interface{}, key string) uint32 {
	return uint32(jsonGetUint64(m, key))
}

func jsonGetInt32(m map[string]interface{}, key string) int32 {
	return int32(jsonGetInt64(m, key))
}

func jsonGetFloat64(m map[string]interface{}, key string) float64 {
	v, ok := m[key]
	if !ok || v == nil {
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

func jsonGetInt(m map[string]interface{}, key string) int {
	return int(jsonGetInt64(m, key))
}

// userToJSON serializes an MTUser to a JSON string for BODY_TEXT params.
func userToJSON(user *MTUser) string {
	if user == nil {
		return "{}"
	}
	b, err := json.Marshal(user)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// parseUserFromJSON parses an MTUser from a JSON string.
func parseUserFromJSON(jsonStr string) *MTUser {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing user from json failed: %v", err))
		return nil
	}
	return parseUserFromMap(m)
}

// parseUserFromMap parses an MTUser from a map[string]interface{}.
func parseUserFromMap(m map[string]interface{}) *MTUser {
	if m == nil {
		return nil
	}
	user := &MTUser{}
	user.Login = jsonGetUint64(m, "Login")
	user.Group = jsonGetString(m, "Group")
	user.CertSerialNumber = jsonGetUint64(m, "CertSerialNumber")
	user.Rights = EnUsersRights(jsonGetUint64(m, "Rights"))
	user.MQID = jsonGetString(m, "MQID")
	user.Registration = jsonGetInt64(m, "Registration")
	user.LastAccess = jsonGetInt64(m, "LastAccess")
	user.LastPassChange = jsonGetInt64(m, "LastPassChange")
	user.LastIP = jsonGetString(m, "LastIP")
	user.Name = jsonGetString(m, "Name")
	user.Company = jsonGetString(m, "Company")
	user.Account = jsonGetString(m, "Account")
	user.Country = jsonGetString(m, "Country")
	user.Language = jsonGetUint32(m, "Language")
	user.ClientID = jsonGetUint64(m, "ClientID")
	user.City = jsonGetString(m, "City")
	user.State = jsonGetString(m, "State")
	user.ZIPCode = jsonGetString(m, "ZipCode")
	user.Address = jsonGetString(m, "Address")
	user.Phone = jsonGetString(m, "Phone")
	user.Email = jsonGetString(m, "Email")
	user.ID = jsonGetString(m, "ID")
	user.Status = jsonGetString(m, "Status")
	user.Comment = jsonGetString(m, "Comment")
	user.Color = jsonGetUint32(m, "Color")
	user.PhonePassword = jsonGetString(m, "PhonePassword")
	user.Leverage = jsonGetUint32(m, "Leverage")
	user.Agent = jsonGetUint64(m, "Agent")
	user.Balance = jsonGetFloat64(m, "Balance")
	user.Credit = jsonGetFloat64(m, "Credit")
	user.InterestRate = jsonGetFloat64(m, "InterestRate")
	user.CommissionDaily = jsonGetFloat64(m, "CommissionDaily")
	user.CommissionMonthly = jsonGetFloat64(m, "CommissionMonthly")
	user.CommissionAgentDaily = jsonGetFloat64(m, "CommissionAgentDaily")
	user.CommissionAgentMonthly = jsonGetFloat64(m, "CommissionAgentMonthly")
	user.BalancePrevDay = jsonGetFloat64(m, "BalancePrevDay")
	user.BalancePrevMonth = jsonGetFloat64(m, "BalancePrevMonth")
	user.EquityPrevDay = jsonGetFloat64(m, "EquityPrevDay")
	user.EquityPrevMonth = jsonGetFloat64(m, "EquityPrevMonth")
	user.TradeAccounts = jsonGetString(m, "TradeAccounts")
	user.LeadCampaign = jsonGetString(m, "LeadCampaign")
	user.LeadSource = jsonGetString(m, "LeadSource")
	return user
}

// parseUsersFromJSON parses a list of MTUser from a JSON array string.
func parseUsersFromJSON(jsonStr string) []*MTUser {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing users from json failed: %v", err))
		return nil
	}
	var users []*MTUser
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			if u := parseUserFromMap(m); u != nil {
				users = append(users, u)
			}
		}
	}
	return users
}

// parseAccountFromJSON parses an MTAccount from a JSON string.
func parseAccountFromJSON(jsonStr string) *MTAccount {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing account from json failed: %v", err))
		return nil
	}
	acc := &MTAccount{}
	acc.Login = jsonGetUint64(m, "Login")
	acc.CurrencyDigits = jsonGetUint32(m, "CurrencyDigits")
	acc.Balance = jsonGetFloat64(m, "Balance")
	acc.Credit = jsonGetFloat64(m, "Credit")
	acc.Margin = jsonGetFloat64(m, "Margin")
	acc.MarginFree = jsonGetFloat64(m, "MarginFree")
	acc.MarginLevel = jsonGetFloat64(m, "MarginLevel")
	acc.MarginLeverage = jsonGetUint32(m, "MarginLeverage")
	acc.Profit = jsonGetFloat64(m, "Profit")
	acc.Storage = jsonGetFloat64(m, "Storage")
	acc.Commission = jsonGetFloat64(m, "Commission")
	acc.Floating = jsonGetFloat64(m, "Floating")
	acc.Equity = jsonGetFloat64(m, "Equity")
	acc.SOActivation = EnSoActivation(jsonGetUint32(m, "SOActivation"))
	acc.SOTime = jsonGetInt64(m, "SOTime")
	acc.SOLevel = jsonGetFloat64(m, "SOLevel")
	acc.SOEquity = jsonGetFloat64(m, "SOEquity")
	acc.SOMargin = jsonGetFloat64(m, "SOMargin")
	acc.MarginInitial = jsonGetFloat64(m, "MarginInitial")
	acc.MarginMaintenance = jsonGetFloat64(m, "MarginMaintenance")
	acc.Assets = jsonGetFloat64(m, "Assets")
	acc.Liabilities = jsonGetFloat64(m, "Liabilities")
	acc.BlockedCommission = jsonGetFloat64(m, "BlockedCommission")
	acc.BlockedProfit = jsonGetFloat64(m, "BlockedProfit")
	return acc
}

// parseLoginsFromJSON parses a list of login IDs from a JSON array string.
func parseLoginsFromJSON(jsonStr string) []uint64 {
	var raw []interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing logins from json failed: %v", err))
		return nil
	}
	logins := make([]uint64, 0, len(raw))
	for _, v := range raw {
		switch val := v.(type) {
		case float64:
			logins = append(logins, uint64(val))
		case string:
			if n, err := strconv.ParseUint(val, 10, 64); err == nil {
				logins = append(logins, n)
			}
		}
	}
	return logins
}

// parseOrderFromJSON parses an MTOrder from a JSON string.
func parseOrderFromJSON(jsonStr string) *MTOrder {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing order from json failed: %v", err))
		return nil
	}
	return parseOrderFromMap(m)
}

func parseOrderFromMap(m map[string]interface{}) *MTOrder {
	if m == nil {
		return nil
	}
	o := &MTOrder{}
	o.Order = jsonGetUint64(m, "Order")
	o.ExternalID = jsonGetString(m, "ExternalID")
	o.Login = jsonGetUint64(m, "Login")
	o.Dealer = jsonGetUint64(m, "Dealer")
	o.Symbol = jsonGetString(m, "Symbol")
	o.Digits = jsonGetUint32(m, "Digits")
	o.DigitsCurrency = jsonGetUint32(m, "DigitsCurrency")
	o.ContractSize = jsonGetFloat64(m, "ContractSize")
	o.State = EnOrderState(jsonGetUint32(m, "State"))
	o.Reason = EnOrderReason(jsonGetUint32(m, "Reason"))
	o.TimeSetup = jsonGetInt64(m, "TimeSetup")
	o.TimeExpiration = jsonGetInt64(m, "TimeExpiration")
	o.TimeDone = jsonGetInt64(m, "TimeDone")
	o.TimeSetupMsc = jsonGetInt64(m, "TimeSetupMsc")
	o.TimeDoneMsc = jsonGetInt64(m, "TimeDoneMsc")
	o.ModifyFlags = EnOrderTradeModifyFlags(jsonGetUint32(m, "ModifyFlags"))
	o.Type = EnOrderType(jsonGetUint32(m, "Type"))
	o.TypeFill = EnOrderFilling(jsonGetUint32(m, "TypeFill"))
	o.TypeTime = EnOrderTime(jsonGetUint32(m, "TypeTime"))
	o.PriceOrder = jsonGetFloat64(m, "PriceOrder")
	o.PriceTrigger = jsonGetFloat64(m, "PriceTrigger")
	o.PriceCurrent = jsonGetFloat64(m, "PriceCurrent")
	o.PriceSL = jsonGetFloat64(m, "PriceSL")
	o.PriceTP = jsonGetFloat64(m, "PriceTP")
	o.SetVolumeInitialExt(jsonGetUint64(m, "VolumeInitialExt"))
	if o.VolumeInitialExt() == 0 {
		o.SetVolumeInitial(jsonGetUint64(m, "VolumeInitial"))
	}
	o.SetVolumeCurrentExt(jsonGetUint64(m, "VolumeCurrentExt"))
	if o.VolumeCurrentExt() == 0 {
		o.SetVolumeCurrent(jsonGetUint64(m, "VolumeCurrent"))
	}
	o.ExpertID = jsonGetUint64(m, "ExpertID")
	o.ExpertPositionID = jsonGetUint64(m, "PositionID")
	o.PositionByID = jsonGetUint64(m, "PositionByID")
	o.Comment = jsonGetString(m, "Comment")
	o.ActivationMode = EnOrderActivation(jsonGetUint32(m, "ActivationMode"))
	o.ActivationTime = jsonGetInt64(m, "ActivationTime")
	o.ActivationPrice = jsonGetFloat64(m, "ActivationPrice")
	o.ActivationFlags = EnOrderTradeActivationFlags(jsonGetUint32(m, "ActivationFlags"))
	return o
}

// parseOrdersFromJSON parses a list of MTOrder from a JSON array string.
func parseOrdersFromJSON(jsonStr string) []*MTOrder {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing orders page from json failed: %v", err))
		return nil
	}
	var orders []*MTOrder
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			if o := parseOrderFromMap(m); o != nil {
				orders = append(orders, o)
			}
		}
	}
	return orders
}

// parseDealFromJSON parses an MTDeal from a JSON string.
func parseDealFromJSON(jsonStr string) *MTDeal {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing deal from json failed: %v", err))
		return nil
	}
	return parseDealFromMap(m)
}

func parseDealFromMap(m map[string]interface{}) *MTDeal {
	if m == nil {
		return nil
	}
	d := &MTDeal{}
	d.Deal = jsonGetUint64(m, "Deal")
	d.ExternalID = jsonGetString(m, "ExternalID")
	d.Login = jsonGetUint64(m, "Login")
	d.Dealer = jsonGetUint64(m, "Dealer")
	d.Order = jsonGetUint64(m, "Order")
	d.Action = EnDealAction(jsonGetUint32(m, "Action"))
	d.Entry = EnEntryFlags(jsonGetUint32(m, "Entry"))
	d.Reason = EnDealReason(jsonGetUint32(m, "Reason"))
	d.Digits = jsonGetUint32(m, "Digits")
	d.DigitsCurrency = jsonGetUint32(m, "DigitsCurrency")
	d.ContractSize = jsonGetFloat64(m, "ContractSize")
	d.Time = jsonGetInt64(m, "Time")
	d.TimeMsc = jsonGetInt64(m, "TimeMsc")
	d.Symbol = jsonGetString(m, "Symbol")
	d.Price = jsonGetFloat64(m, "Price")
	d.SetVolumeExt(jsonGetUint64(m, "VolumeExt"))
	if d.VolumeExt() == 0 {
		d.SetVolume(jsonGetUint64(m, "Volume"))
	}
	d.Profit = jsonGetFloat64(m, "Profit")
	d.Storage = jsonGetFloat64(m, "Storage")
	d.Commission = jsonGetFloat64(m, "Commission")
	d.CommissionAgent = jsonGetFloat64(m, "CommissionAgent")
	d.RateProfit = jsonGetFloat64(m, "RateProfit")
	d.RateMargin = jsonGetFloat64(m, "RateMargin")
	d.ExpertID = jsonGetUint64(m, "ExpertID")
	d.PositionID = jsonGetUint64(m, "PositionID")
	d.Comment = jsonGetString(m, "Comment")
	d.ProfitRaw = jsonGetFloat64(m, "ProfitRaw")
	d.PricePosition = jsonGetFloat64(m, "PricePosition")
	d.SetVolumeClosedExt(jsonGetUint64(m, "VolumeClosedExt"))
	if d.VolumeClosedExt() == 0 {
		d.SetVolumeClosed(jsonGetUint64(m, "VolumeClosed"))
	}
	d.TickValue = jsonGetFloat64(m, "TickValue")
	d.TickSize = jsonGetFloat64(m, "TickSize")
	d.Flags = jsonGetUint64(m, "Flags")
	d.Gateway = jsonGetString(m, "Gateway")
	d.PriceGateway = jsonGetFloat64(m, "PriceGateway")
	d.ModifyFlags = EnDealTradeModifyFlags(jsonGetUint32(m, "ModifyFlags"))
	d.PriceSL = jsonGetFloat64(m, "PriceSL")
	d.PriceTP = jsonGetFloat64(m, "PriceTP")
	return d
}

// parseDealsFromJSON parses a list of MTDeal from a JSON array string.
func parseDealsFromJSON(jsonStr string) []*MTDeal {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing deals page from json failed: %v", err))
		return nil
	}
	var deals []*MTDeal
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			if d := parseDealFromMap(m); d != nil {
				deals = append(deals, d)
			}
		}
	}
	return deals
}

// parsePositionFromJSON parses an MTPosition from a JSON string.
func parsePositionFromJSON(jsonStr string) *MTPosition {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing position from json failed: %v", err))
		return nil
	}
	return parsePositionFromMap(m)
}

func parsePositionFromMap(m map[string]interface{}) *MTPosition {
	if m == nil {
		return nil
	}
	p := &MTPosition{}
	p.Position = jsonGetUint64(m, "Position")
	p.ExternalID = jsonGetString(m, "ExternalID")
	p.Login = jsonGetUint64(m, "Login")
	p.Dealer = jsonGetUint64(m, "Dealer")
	p.Symbol = jsonGetString(m, "Symbol")
	p.Action = EnPositionAction(jsonGetUint32(m, "Action"))
	p.Digits = jsonGetUint32(m, "Digits")
	p.DigitsCurrency = jsonGetUint32(m, "DigitsCurrency")
	p.Reason = EnPositionReason(jsonGetUint32(m, "Reason"))
	p.ContractSize = jsonGetFloat64(m, "ContractSize")
	p.TimeCreate = jsonGetInt64(m, "TimeCreate")
	p.TimeUpdate = jsonGetInt64(m, "TimeUpdate")
	p.ModifyFlags = EnPositionTradeModifyFlags(jsonGetUint32(m, "ModifyFlags"))
	p.PriceOpen = jsonGetFloat64(m, "PriceOpen")
	p.PriceCurrent = jsonGetFloat64(m, "PriceCurrent")
	p.PriceSL = jsonGetFloat64(m, "PriceSL")
	p.PriceTP = jsonGetFloat64(m, "PriceTP")
	p.SetVolumeExt(jsonGetUint64(m, "VolumeExt"))
	if p.VolumeExt() == 0 {
		p.SetVolume(jsonGetUint64(m, "Volume"))
	}
	p.Profit = jsonGetFloat64(m, "Profit")
	p.Storage = jsonGetFloat64(m, "Storage")
	p.Commission = jsonGetFloat64(m, "Commission")
	p.RateProfit = jsonGetFloat64(m, "RateProfit")
	p.RateMargin = jsonGetFloat64(m, "RateMargin")
	p.ExpertID = jsonGetUint64(m, "ExpertID")
	p.ExpertPositionID = jsonGetUint64(m, "ExpertPositionID")
	p.Comment = jsonGetString(m, "Comment")
	p.ActivationMode = EnPositionActivation(jsonGetUint32(m, "ActivationMode"))
	p.ActivationTime = jsonGetInt64(m, "ActivationTime")
	p.ActivationPrice = jsonGetFloat64(m, "ActivationPrice")
	p.ActivationFlags = EnPositionTradeActivationFlags(jsonGetUint32(m, "ActivationFlags"))
	return p
}

// parsePositionsFromJSON parses a list of MTPosition from a JSON array string.
func parsePositionsFromJSON(jsonStr string) []*MTPosition {
	arr, err := jsonToSlice(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing positions page from json failed: %v", err))
		return nil
	}
	var positions []*MTPosition
	for _, item := range arr {
		if m, ok := item.(map[string]interface{}); ok {
			if p := parsePositionFromMap(m); p != nil {
				positions = append(positions, p)
			}
		}
	}
	return positions
}

// passwordTypeString converts a password type enum to its protocol string value.
func passwordTypeString(passType EnUsersPasswords) string {
	switch passType {
	case UserPassMain:
		return WEB_VAL_USER_PASS_MAIN
	case UserPassInvestor:
		return WEB_VAL_USER_PASS_INVESTOR
	case UserPassAPI:
		return WEB_VAL_USER_PASS_API
	default:
		return WEB_VAL_USER_PASS_MAIN
	}
}

// formatUint64 converts a uint64 to its string representation.
func formatUint64(v uint64) string {
	return strconv.FormatUint(v, 10)
}

// formatUint32 converts a uint32 to its string representation.
func formatUint32(v uint32) string {
	return strconv.FormatUint(uint64(v), 10)
}

// formatInt64 converts an int64 to its string representation.
func formatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

// formatInt converts an int to its string representation.
func formatInt(v int) string {
	return strconv.Itoa(v)
}

// formatFloat64 converts a float64 to a string with enough precision (no trailing zeros issue).
func formatFloat64(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// parseTimeFromString extracts unix time from a time server string like "1234567890 Mon...".
func parseTimeFromString(timeStr string) int64 {
	if timeStr == "" || timeStr == "none" {
		return 0
	}
	parts := strings.SplitN(timeStr, " ", 2)
	t, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0
	}
	return t
}
