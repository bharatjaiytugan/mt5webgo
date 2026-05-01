//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"encoding/json"
	"fmt"
)

// timeHandler implements time-related protocol commands.
type timeHandler struct {
	apiBase
}

// TimeServer retrieves the current server time as a unix timestamp.
func (h *timeHandler) TimeServer() int64 {
	answer := h.send(WEB_CMD_TIME_SERVER, nil)
	if answer == nil {
		logWrite(LogTypeError, "send time server failed")
		return 0
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer time: %s", answerStr))
	}
	timeVal, errCode := parseTimeServerAnswer(answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse time server failed: [%d] %s", errCode, GetError(errCode)))
		return 0
	}
	return parseTimeFromString(timeVal)
}

// TimeGet retrieves the server time configuration.
func (h *timeHandler) TimeGet() (*MTConTime, MTRetCode) {
	answer := h.send(WEB_CMD_TIME_GET, nil)
	if answer == nil {
		logWrite(LogTypeError, "send time get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_TIME_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse time get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	timeCon := parseConTimeFromJSON(jsonStr)
	if timeCon == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return timeCon, MT_RET_OK
}

// parseConTimeFromJSON parses an MTConTime from a JSON string.
// MT5 server sends all numeric values as strings, so we parse manually.
func parseConTimeFromJSON(jsonStr string) *MTConTime {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parsing time config from json failed: %v", err))
		return nil
	}
	t := &MTConTime{
		DaylightState: jsonGetInt32(raw, "DaylightState"),
		Daylight:      jsonGetInt32(raw, "Daylight"),
		TimeZone:      jsonGetInt32(raw, "TimeZone"),
		TimeServer:    jsonGetString(raw, "TimeServer"),
	}
	// parse Days array
	if daysRaw, ok := raw["Days"]; ok {
		if daysArr, ok := daysRaw.([]interface{}); ok {
			t.Days = make([][]EnTimeTableMode, len(daysArr))
			for i, dayRaw := range daysArr {
				if dayArr, ok := dayRaw.([]interface{}); ok {
					t.Days[i] = make([]EnTimeTableMode, len(dayArr))
					for j, v := range dayArr {
						switch val := v.(type) {
						case float64:
							t.Days[i][j] = EnTimeTableMode(int(val))
						case string:
							t.Days[i][j] = EnTimeTableMode(jsonGetInt32(map[string]interface{}{"v": val}, "v"))
						}
					}
				}
			}
		}
	}
	return t
}
