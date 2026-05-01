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

// Client protocol command constants.
const (
	WEB_CMD_CLIENT_ADD         = "CLIENT_ADD"
	WEB_CMD_CLIENT_UPDATE      = "CLIENT_UPDATE"
	WEB_CMD_CLIENT_DELETE      = "CLIENT_DELETE"
	WEB_CMD_CLIENT_GET         = "CLIENT_GET"
	WEB_CMD_CLIENT_IDS         = "CLIENT_IDS"
	WEB_CMD_CLIENT_HISTORY_GET = "CLIENT_HISTORY_GET"
	WEB_CMD_CLIENT_USER_ADD    = "CLIENT_USER_ADD"
	WEB_CMD_CLIENT_USER_DELETE = "CLIENT_USER_DELETE"
	WEB_CMD_CLIENT_USER_LOGINS = "CLIENT_USER_LOGINS"
)

// MTClient represents a client record in the MetaTrader 5 server.
type MTClient struct {
	RecordID         uint64 `json:"RecordID"`
	ClientType       uint32 `json:"ClientType"`
	ClientStatus     uint32 `json:"ClientStatus"`
	ClientExternalID string `json:"ClientExternalID"`
	AssignedManager  uint64 `json:"AssignedManager"`
	DateCreated      int64  `json:"DateCreated"`
	DateModified     int64  `json:"DateModified"`
	Comment          string `json:"Comment"`
	// Personal info
	Name       string `json:"Name"`
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	MiddleName string `json:"MiddleName"`
	Country    string `json:"Country"`
	City       string `json:"City"`
	State      string `json:"State"`
	ZipCode    string `json:"ZipCode"`
	Address    string `json:"Address"`
	Phone      string `json:"Phone"`
	Email      string `json:"Email"`
	ID         string `json:"ID"`
}

// clientHandler implements client-related protocol commands.
type clientHandler struct {
	apiBase
}

// ClientAdd adds a new client to the server.
// clientJSON is the JSON of ONE client. It will be wrapped in array format [{}] as required by the server.
// Returns the created client with RecordID.
func (h *clientHandler) ClientAdd(clientJSON string) (*MTClient, MTRetCode) {
	if clientJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	// server expects array format: [{client}]
	bodyJSON := "[" + clientJSON + "]"
	params := map[string]string{
		WEB_PARAM_BODYTEXT: bodyJSON,
	}
	answer := h.send(WEB_CMD_CLIENT_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send client add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// response: CLIENT_ADD|RETCODE=0 Done|\r\n[{"id":"1278","retcode":"0 Done"}]
	// parse command and retcode first
	_, pos := getCommand(answerStr)
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
		if param == nil {
			break
		}
		if param.Name == WEB_PARAM_RETCODE {
			retCode = param.Value
		}
	}
	errCode := getRetCode(retCode)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("client add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	// parse JSON array response to get id
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" {
		if idx := strings.Index(answerStr, "["); idx >= 0 {
			jsonStr = strings.TrimSpace(answerStr[idx:])
		}
	}
	if jsonStr != "" {
		var results []map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &results); err == nil && len(results) > 0 {
			client := &MTClient{}
			// id field - may be string "1002" or number 1002
			switch v := results[0]["id"].(type) {
			case string:
				if id, err := strconv.ParseUint(v, 10, 64); err == nil {
					client.RecordID = id
				}
			case float64:
				client.RecordID = uint64(v)
			}
			return client, MT_RET_OK
		}
	}
	return nil, MT_RET_OK
}

// ClientUpdate updates an existing client on the server.
// clientJSON is the JSON representation of the client to update (must include RecordID).
func (h *clientHandler) ClientUpdate(clientJSON string) (*MTClient, MTRetCode) {
	if clientJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: clientJSON,
	}
	answer := h.send(WEB_CMD_CLIENT_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send client update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" {
		if idx := strings.Index(answerStr, "{"); idx >= 0 {
			jsonStr = strings.TrimSpace(answerStr[idx:])
		}
	}
	if jsonStr == "" {
		return nil, MT_RET_OK
	}
	client := parseClientFromJSON(jsonStr)
	return client, MT_RET_OK
}

// ClientDelete deletes a client by record ID.
func (h *clientHandler) ClientDelete(recordID uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(recordID),
	}
	answer := h.send(WEB_CMD_CLIENT_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send client delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// ClientGet retrieves a client by record ID.
func (h *clientHandler) ClientGet(recordID uint64) (*MTClient, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_ID: formatUint64(recordID),
	}
	answer := h.send(WEB_CMD_CLIENT_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send client get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" {
		if idx := strings.Index(answerStr, "{"); idx >= 0 {
			jsonStr = strings.TrimSpace(answerStr[idx:])
		}
	}
	if jsonStr == "" {
		return nil, MT_RET_OK
	}
	client := parseClientFromJSON(jsonStr)
	return client, MT_RET_OK
}

// ClientGetIDs retrieves client record IDs available to the manager.
// group is optional - filter by group mask (e.g. "demo\\*"). Pass "" for all.
func (h *clientHandler) ClientGetIDs(group string) ([]uint64, MTRetCode) {
	params := map[string]string{}
	if group != "" {
		params[WEB_PARAM_GROUP] = group
	}
	answer := h.send(WEB_CMD_CLIENT_IDS, params)
	if answer == nil {
		logWrite(LogTypeError, "send client ids failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client ids failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" {
		return nil, MT_RET_REPORT_NODATA
	}
	ids := parseLoginsFromJSON(jsonStr)
	if ids == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return ids, MT_RET_OK
}

// ClientUserAdd links a user login to a client record.
func (h *clientHandler) ClientUserAdd(clientID, login uint64) MTRetCode {
	params := map[string]string{
		"CLIENT": formatUint64(clientID),
		"USER":   formatUint64(login),
	}
	answer := h.send(WEB_CMD_CLIENT_USER_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send client user add failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client user add failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// ClientUserDelete unlinks a user login from a client record.
func (h *clientHandler) ClientUserDelete(clientID, login uint64) MTRetCode {
	params := map[string]string{
		"CLIENT": formatUint64(clientID),
		"USER":   formatUint64(login),
	}
	answer := h.send(WEB_CMD_CLIENT_USER_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send client user delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client user delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// ClientUserLogins retrieves user logins linked to a client record.
func (h *clientHandler) ClientUserLogins(clientID uint64) ([]uint64, MTRetCode) {
	params := map[string]string{
		"CLIENT": formatUint64(clientID),
	}
	answer := h.send(WEB_CMD_CLIENT_USER_LOGINS, params)
	if answer == nil {
		logWrite(LogTypeError, "send client user logins failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	// tolerant parse - accept any command name format
	cmdReal, pos := getCommand(answerStr)
	_ = cmdReal
	retCode := ""
	posEnd := -1
	for {
		param := getNextParam(answerStr, &pos, &posEnd)
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
		logWrite(LogTypeError, fmt.Sprintf("client user logins failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	jsonStr := getJSON(answerStr, posEnd)
	if jsonStr == "" {
		if idx := strings.Index(answerStr, "{"); idx >= 0 {
			jsonStr = strings.TrimSpace(answerStr[idx:])
		}
	}
	if jsonStr == "" {
		return []uint64{}, MT_RET_OK
	}
	// response format: {"clientID": [login1, login2], "clientID2": [login3]}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse client user logins json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	var logins []uint64
	for _, v := range raw {
		if arr, ok := v.([]interface{}); ok {
			for _, item := range arr {
				switch val := item.(type) {
				case float64:
					logins = append(logins, uint64(val))
				case string:
					if n, err := strconv.ParseUint(val, 10, 64); err == nil {
						logins = append(logins, n)
					}
				}
			}
		}
	}
	return logins, MT_RET_OK
}

// parseClientFromJSON parses an MTClient from a JSON string.
func parseClientFromJSON(jsonStr string) *MTClient {
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing client from json failed: %v", err))
		return nil
	}
	return parseClientFromMap(m)
}

// parseClientFromMap parses an MTClient from a map[string]interface{}.
func parseClientFromMap(m map[string]interface{}) *MTClient {
	if m == nil {
		return nil
	}
	c := &MTClient{}
	c.RecordID = jsonGetUint64(m, "RecordID")
	c.ClientType = jsonGetUint32(m, "ClientType")
	c.ClientStatus = jsonGetUint32(m, "ClientStatus")
	c.ClientExternalID = jsonGetString(m, "ClientExternalID")
	c.AssignedManager = jsonGetUint64(m, "AssignedManager")
	c.DateCreated = jsonGetInt64(m, "DateCreated")
	c.DateModified = jsonGetInt64(m, "DateModified")
	c.Comment = jsonGetString(m, "Comment")
	c.Name = jsonGetString(m, "Name")
	c.FirstName = jsonGetString(m, "FirstName")
	c.LastName = jsonGetString(m, "LastName")
	c.MiddleName = jsonGetString(m, "MiddleName")
	c.Country = jsonGetString(m, "Country")
	c.City = jsonGetString(m, "City")
	c.State = jsonGetString(m, "State")
	c.ZipCode = jsonGetString(m, "ZipCode")
	c.Address = jsonGetString(m, "Address")
	c.Phone = jsonGetString(m, "Phone")
	c.Email = jsonGetString(m, "Email")
	c.ID = jsonGetString(m, "ID")
	return c
}
