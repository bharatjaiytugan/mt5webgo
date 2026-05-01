//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
)

// Settings command constants.
const (
	WEB_CMD_SETTING_GET    = "SETTING_GET"
	WEB_CMD_SETTING_SET    = "SETTING_SET"
	WEB_CMD_SETTING_DELETE = "SETTING_DELETE"
)

// Settings parameter constants.
const (
	WEB_PARAM_SECTION = "SECTION"
	WEB_PARAM_KEY     = "KEY"
)

// settingsHandler implements settings protocol commands.
type settingsHandler struct {
	apiBase
}

// SettingGet gets a settings file from the server.
// section is the directory, key is the filename.
// Returns the parsed JSON content as a map and a return code.
func (h *settingsHandler) SettingGet(section, key string) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_SECTION: section,
		WEB_PARAM_KEY:     key,
	}
	answer := h.send(WEB_CMD_SETTING_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send setting get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer setting get: %s", answerStr))
	}
	// parse with tolerant command name
	jsonStr, errCode := parseTolerantJSON(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse setting get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	m, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse setting get json failed: %v", err))
		return nil, MT_RET_ERR_DATA
	}
	return m, MT_RET_OK
}

// SettingSet writes a settings file to the server.
// section is the directory, key is the filename, contentJSON is the JSON content to write.
// Returns a return code.
func (h *settingsHandler) SettingSet(section, key, contentJSON string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_SECTION:  section,
		WEB_PARAM_KEY:      key,
		WEB_PARAM_BODYTEXT: contentJSON,
	}
	answer := h.send(WEB_CMD_SETTING_SET, params)
	if answer == nil {
		logWrite(LogTypeError, "send setting set failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer setting set: %s", answerStr))
	}
	// try tolerant parsing first, fall back to parseEmptyResult
	errCode := parseEmptyResultTolerant(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse setting set failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// SettingDelete deletes a settings file from the server.
// section is the directory, key is the filename.
// Returns a return code.
func (h *settingsHandler) SettingDelete(section, key string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_SECTION: section,
		WEB_PARAM_KEY:     key,
	}
	answer := h.send(WEB_CMD_SETTING_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send setting delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer setting delete: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_SETTING_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse setting delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// parseEmptyResultTolerant checks a simple answer that only contains RETCODE (no JSON body),
// accepting any command name returned by the server.
func parseEmptyResultTolerant(answer string) MTRetCode {
	_, pos := getCommand(answer)
	if pos < 0 {
		return MT_RET_ERR_DATA
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
	return getRetCode(retCode)
}
