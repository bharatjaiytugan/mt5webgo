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

// WebSocket commands for daily reports.
const (
	WEB_CMD_DAILY_GET       = "DAILY_GET"
	WEB_CMD_DAILY_GET_BATCH = "DAILY_GET_BATCH"
)

// MTDailyReport represents a daily report record from the server.
type MTDailyReport struct {
	Login             uint64
	Name              string
	Group             string
	Currency          string
	CurrencyDigits    uint32
	Company           string
	Email             string
	Balance           float64
	Credit            float64
	InterestRate      float64
	CommissionDaily   float64
	CommissionMonthly float64
	AgentDaily        float64
	AgentMonthly      float64
	BalancePrevDay    float64
	BalancePrevMonth  float64
	EquityPrevDay     float64
	EquityPrevMonth   float64
	Margin            float64
	MarginFree        float64
	Equity            float64
	Profit            float64
	ProfitEquity      float64
	DailyProfit       float64
	DailyBalance      float64
	DailyCredit       float64
	DailyCharge       float64
	DailyCorrection   float64
	DailyBonus        float64
	DailyStorage      float64
	DailyCommission   float64
	DailyAgent        float64
	Timestamp         string
	DatetimePrev      int64
}

// dailyHandler implements daily report protocol commands.
type dailyHandler struct {
	apiBase
}

// DailyGet retrieves daily reports for a login in a date range.
func (h *dailyHandler) DailyGet(login uint64, from, to int64) ([]*MTDailyReport, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
		WEB_PARAM_FROM:  formatInt64(from),
		WEB_PARAM_TO:    formatInt64(to),
	}
	answer := h.send(WEB_CMD_DAILY_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send daily get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseDailyResponse(WEB_CMD_DAILY_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse daily get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	dailies := parseDailiesFromJSON(jsonStr)
	if dailies == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return dailies, MT_RET_OK
}

// DailyGetBatch retrieves daily reports for multiple logins in a date range.
// Implemented by calling DailyGet for each login since DAILY_GET_BATCH is not supported on WebSocket.
func (h *dailyHandler) DailyGetBatch(logins string, from, to int64) ([]*MTDailyReport, MTRetCode) {
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
		reports, ret := h.DailyGet(loginNum, from, to)
		if ret == MT_RET_OK && reports != nil {
			allReports = append(allReports, reports...)
		}
	}
	if allReports == nil {
		return []*MTDailyReport{}, MT_RET_OK
	}
	return allReports, MT_RET_OK
}

// parseDailyResponse parses daily report responses, accepting any server command format.
func parseDailyResponse(expectedCmd, answer string) (string, MTRetCode) {
	cmdReal, pos := getCommand(answer)
	if cmdReal != expectedCmd {
		// accept any command name the server returns (may use slash format)
		logWrite(LogTypeDebug, fmt.Sprintf("daily: server returned command '%s', expected '%s' - accepting", cmdReal, expectedCmd))
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

// parseDailyFromMap parses an MTDailyReport from a map[string]interface{}.
func parseDailyFromMap(m map[string]interface{}) *MTDailyReport {
	if m == nil {
		return nil
	}
	d := &MTDailyReport{}
	d.Login = jsonGetUint64(m, "Login")
	d.Name = jsonGetString(m, "Name")
	d.Group = jsonGetString(m, "Group")
	d.Currency = jsonGetString(m, "Currency")
	d.CurrencyDigits = jsonGetUint32(m, "CurrencyDigits")
	d.Company = jsonGetString(m, "Company")
	d.Email = jsonGetString(m, "Email")
	d.Balance = jsonGetFloat64(m, "Balance")
	d.Credit = jsonGetFloat64(m, "Credit")
	d.InterestRate = jsonGetFloat64(m, "InterestRate")
	d.CommissionDaily = jsonGetFloat64(m, "CommissionDaily")
	d.CommissionMonthly = jsonGetFloat64(m, "CommissionMonthly")
	d.AgentDaily = jsonGetFloat64(m, "AgentDaily")
	d.AgentMonthly = jsonGetFloat64(m, "AgentMonthly")
	d.BalancePrevDay = jsonGetFloat64(m, "BalancePrevDay")
	d.BalancePrevMonth = jsonGetFloat64(m, "BalancePrevMonth")
	d.EquityPrevDay = jsonGetFloat64(m, "EquityPrevDay")
	d.EquityPrevMonth = jsonGetFloat64(m, "EquityPrevMonth")
	d.Margin = jsonGetFloat64(m, "Margin")
	d.MarginFree = jsonGetFloat64(m, "MarginFree")
	d.Equity = jsonGetFloat64(m, "Equity")
	d.Profit = jsonGetFloat64(m, "Profit")
	d.ProfitEquity = jsonGetFloat64(m, "ProfitEquity")
	d.DailyProfit = jsonGetFloat64(m, "DailyProfit")
	d.DailyBalance = jsonGetFloat64(m, "DailyBalance")
	d.DailyCredit = jsonGetFloat64(m, "DailyCredit")
	d.DailyCharge = jsonGetFloat64(m, "DailyCharge")
	d.DailyCorrection = jsonGetFloat64(m, "DailyCorrection")
	d.DailyBonus = jsonGetFloat64(m, "DailyBonus")
	d.DailyStorage = jsonGetFloat64(m, "DailyStorage")
	d.DailyCommission = jsonGetFloat64(m, "DailyCommission")
	d.DailyAgent = jsonGetFloat64(m, "DailyAgent")
	d.Timestamp = jsonGetString(m, "Timestamp")
	d.DatetimePrev = jsonGetInt64(m, "DatetimePrev")
	return d
}

// parseDailiesFromJSON parses daily reports from a JSON string.
// The response may be a JSON array of objects or a single object.
func parseDailiesFromJSON(jsonStr string) []*MTDailyReport {
	jsonStr = trimSpaces(jsonStr)
	if len(jsonStr) == 0 {
		return nil
	}
	// try array first
	if jsonStr[0] == '[' {
		var arr []interface{}
		if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
			logWrite(LogTypeError, fmt.Sprintf("parsing daily reports from json array failed: %v", err))
			return nil
		}
		var dailies []*MTDailyReport
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				if d := parseDailyFromMap(m); d != nil {
					dailies = append(dailies, d)
				}
			}
		}
		return dailies
	}
	// single object
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing daily report from json failed: %v", err))
		return nil
	}
	d := parseDailyFromMap(m)
	if d == nil {
		return nil
	}
	return []*MTDailyReport{d}
}

// trimSpaces trims leading and trailing whitespace from a string.
func trimSpaces(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	j := len(s)
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
