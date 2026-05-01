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

// WebSocket command constants for infrastructure configuration.
const (
	// Email
	WEB_CMD_EMAIL_TOTAL  = "EMAIL_TOTAL"
	WEB_CMD_EMAIL_NEXT   = "EMAIL_NEXT"
	WEB_CMD_EMAIL_GET    = "EMAIL_GET"
	WEB_CMD_EMAIL_ADD    = "EMAIL_ADD"
	WEB_CMD_EMAIL_DELETE = "EMAIL_DELETE"
	WEB_CMD_EMAIL_SHIFT  = "EMAIL_SHIFT"
	WEB_CMD_EMAIL_SEND   = "EMAIL_SEND"
	// Messenger
	WEB_CMD_MESSENGER_TOTAL  = "MESSENGER_TOTAL"
	WEB_CMD_MESSENGER_NEXT   = "MESSENGER_NEXT"
	WEB_CMD_MESSENGER_GET    = "MESSENGER_GET"
	WEB_CMD_MESSENGER_ADD    = "MESSENGER_ADD"
	WEB_CMD_MESSENGER_DELETE = "MESSENGER_DELETE"
	WEB_CMD_MESSENGER_SHIFT  = "MESSENGER_SHIFT"
	WEB_CMD_MESSENGER_SEND   = "MESSENGER_SEND"
	// Gateway
	WEB_CMD_GATEWAY_TOTAL        = "GATEWAY_TOTAL"
	WEB_CMD_GATEWAY_NEXT         = "GATEWAY_NEXT"
	WEB_CMD_GATEWAY_GET          = "GATEWAY_GET"
	WEB_CMD_GATEWAY_ADD          = "GATEWAY_ADD"
	WEB_CMD_GATEWAY_DELETE       = "GATEWAY_DELETE"
	WEB_CMD_GATEWAY_SHIFT        = "GATEWAY_SHIFT"
	WEB_CMD_GATEWAY_RESTART      = "GATEWAY_RESTART"
	WEB_CMD_GATEWAY_GET_POSITION = "GATEWAY_GET_POSITION"
	WEB_CMD_GATEWAY_MODULE_TOTAL = "GATEWAY_MODULE_TOTAL"
	WEB_CMD_GATEWAY_MODULE_NEXT  = "GATEWAY_MODULE_NEXT"
	WEB_CMD_GATEWAY_MODULE_GET   = "GATEWAY_MODULE_GET"
	// Feeder
	WEB_CMD_FEEDER_TOTAL       = "FEEDER_TOTAL"
	WEB_CMD_FEEDER_NEXT        = "FEEDER_NEXT"
	WEB_CMD_FEEDER_GET         = "FEEDER_GET"
	WEB_CMD_FEEDER_ADD         = "FEEDER_ADD"
	WEB_CMD_FEEDER_DELETE      = "FEEDER_DELETE"
	WEB_CMD_FEEDER_SHIFT       = "FEEDER_SHIFT"
	WEB_CMD_FEEDER_RESTART     = "FEEDER_RESTART"
	WEB_CMD_FEEDER_MODULE_TOTAL = "FEEDER_MODULE_TOTAL"
	WEB_CMD_FEEDER_MODULE_NEXT  = "FEEDER_MODULE_NEXT"
	WEB_CMD_FEEDER_MODULE_GET   = "FEEDER_MODULE_GET"
	// Report
	WEB_CMD_REPORT_TOTAL       = "REPORT_TOTAL"
	WEB_CMD_REPORT_NEXT        = "REPORT_NEXT"
	WEB_CMD_REPORT_GET         = "REPORT_GET"
	WEB_CMD_REPORT_ADD         = "REPORT_ADD"
	WEB_CMD_REPORT_DELETE      = "REPORT_DELETE"
	WEB_CMD_REPORT_SHIFT       = "REPORT_SHIFT"
	WEB_CMD_REPORT_MODULE_TOTAL = "REPORT_MODULE_TOTAL"
	WEB_CMD_REPORT_MODULE_NEXT  = "REPORT_MODULE_NEXT"
	WEB_CMD_REPORT_MODULE_GET   = "REPORT_MODULE_GET"
	// Plugin
	WEB_CMD_PLUGIN_TOTAL       = "PLUGIN_TOTAL"
	WEB_CMD_PLUGIN_NEXT        = "PLUGIN_NEXT"
	WEB_CMD_PLUGIN_GET         = "PLUGIN_GET"
	WEB_CMD_PLUGIN_ADD         = "PLUGIN_ADD"
	WEB_CMD_PLUGIN_DELETE      = "PLUGIN_DELETE"
	WEB_CMD_PLUGIN_SHIFT       = "PLUGIN_SHIFT"
	WEB_CMD_PLUGIN_MODULE_TOTAL = "PLUGIN_MODULE_TOTAL"
	WEB_CMD_PLUGIN_MODULE_NEXT  = "PLUGIN_MODULE_NEXT"
	WEB_CMD_PLUGIN_MODULE_GET   = "PLUGIN_MODULE_GET"
)

// infraHandler implements infrastructure configuration protocol commands
// (email, messenger, gateway, feeder, report, plugin).
type infraHandler struct {
	apiBase
}

// ---------------------------------------------------------------------------
// Generic internal methods
// ---------------------------------------------------------------------------

// infraTotal sends a TOTAL command and returns the count.
func (h *infraHandler) infraTotal(cmd string) (int, MTRetCode) {
	answer := h.send(cmd, nil)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return 0, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	totalStr, errCode := parseTotalAnswer(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
		return 0, errCode
	}
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return 0, MT_RET_REPORT_NODATA
	}
	return total, MT_RET_OK
}

// infraNext sends a NEXT command with an index and returns the JSON map.
func (h *infraHandler) infraNext(cmd string, index int) (map[string]interface{}, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", strings.ToLower(cmd), err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// infraGet sends a GET command with arbitrary params and returns the JSON map.
func (h *infraHandler) infraGet(cmd string, params map[string]string) (map[string]interface{}, MTRetCode) {
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", strings.ToLower(cmd), err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// infraGetAll iterates TOTAL + NEXT to retrieve all records for a given entity.
func (h *infraHandler) infraGetAll(totalCmd, nextCmd string) ([]map[string]interface{}, MTRetCode) {
	total, ret := h.infraTotal(totalCmd)
	if ret != MT_RET_OK {
		return nil, ret
	}
	var results []map[string]interface{}
	for i := 0; i < total; i++ {
		item, ret2 := h.infraNext(nextCmd, i)
		if ret2 != MT_RET_OK {
			return results, ret2
		}
		results = append(results, item)
	}
	return results, MT_RET_OK
}

// infraAdd sends an ADD command with a JSON body and returns the resulting record.
func (h *infraHandler) infraAdd(cmd string, bodyJSON string) (map[string]interface{}, MTRetCode) {
	if bodyJSON == "" {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: bodyJSON,
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
		return nil, errCode
	}
	result, err := jsonToMap(jsonStr)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("parse %s json failed: %v", strings.ToLower(cmd), err))
		return nil, MT_RET_REPORT_NODATA
	}
	return result, MT_RET_OK
}

// infraDelete sends a DELETE command with the given params (index or name).
func (h *infraHandler) infraDelete(cmd string, params map[string]string) MTRetCode {
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
	}
	return errCode
}

// infraShift sends a SHIFT command with index and shift offset.
func (h *infraHandler) infraShift(cmd string, index, shift int) MTRetCode {
	params := map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
		WEB_PARAM_SHIFT: formatInt(shift),
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
	}
	return errCode
}

// infraSendBody sends a command with a JSON body (fire-and-confirm pattern, e.g. SEND).
func (h *infraHandler) infraSendBody(cmd string, bodyJSON string) MTRetCode {
	if bodyJSON == "" {
		return MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_BODYTEXT: bodyJSON,
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
	}
	return errCode
}

// infraRestart sends a RESTART command with a name parameter.
func (h *infraHandler) infraRestart(cmd string, name string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_NAME: name,
	}
	answer := h.send(cmd, params)
	if answer == nil {
		logWrite(LogTypeError, fmt.Sprintf("send %s failed", strings.ToLower(cmd)))
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(cmd, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse %s failed: %s", strings.ToLower(cmd), GetErrorStandard(errCode)))
	}
	return errCode
}

// ---------------------------------------------------------------------------
// Email (8 methods)
// ---------------------------------------------------------------------------

// EmailTotal gets the total count of email configurations.
func (h *infraHandler) EmailTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_EMAIL_TOTAL)
}

// EmailNext gets the email configuration at the given index.
func (h *infraHandler) EmailNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_EMAIL_NEXT, index)
}

// EmailGet gets an email configuration by index.
func (h *infraHandler) EmailGet(index int) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_EMAIL_GET, map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	})
}

// EmailGetAll gets all email configurations.
func (h *infraHandler) EmailGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_EMAIL_TOTAL, WEB_CMD_EMAIL_NEXT)
}

// EmailAdd adds a new email configuration from JSON.
func (h *infraHandler) EmailAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_EMAIL_ADD, bodyJSON)
}

// EmailDelete deletes an email configuration by index.
func (h *infraHandler) EmailDelete(index int) MTRetCode {
	return h.infraDelete(WEB_CMD_EMAIL_DELETE, map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	})
}

// EmailShift shifts an email configuration position.
func (h *infraHandler) EmailShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_EMAIL_SHIFT, index, shift)
}

// EmailSend sends an email using the given JSON body.
func (h *infraHandler) EmailSend(bodyJSON string) MTRetCode {
	return h.infraSendBody(WEB_CMD_EMAIL_SEND, bodyJSON)
}

// ---------------------------------------------------------------------------
// Messenger (8 methods)
// ---------------------------------------------------------------------------

// MessengerTotal gets the total count of messenger configurations.
func (h *infraHandler) MessengerTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_MESSENGER_TOTAL)
}

// MessengerNext gets the messenger configuration at the given index.
func (h *infraHandler) MessengerNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_MESSENGER_NEXT, index)
}

// MessengerGet gets a messenger configuration by index.
func (h *infraHandler) MessengerGet(index int) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_MESSENGER_GET, map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	})
}

// MessengerGetAll gets all messenger configurations.
func (h *infraHandler) MessengerGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_MESSENGER_TOTAL, WEB_CMD_MESSENGER_NEXT)
}

// MessengerAdd adds a new messenger configuration from JSON.
func (h *infraHandler) MessengerAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_MESSENGER_ADD, bodyJSON)
}

// MessengerDelete deletes a messenger configuration by index.
func (h *infraHandler) MessengerDelete(index int) MTRetCode {
	return h.infraDelete(WEB_CMD_MESSENGER_DELETE, map[string]string{
		WEB_PARAM_INDEX: formatInt(index),
	})
}

// MessengerShift shifts a messenger configuration position.
func (h *infraHandler) MessengerShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_MESSENGER_SHIFT, index, shift)
}

// MessengerSend sends a messenger message using the given JSON body.
func (h *infraHandler) MessengerSend(bodyJSON string) MTRetCode {
	return h.infraSendBody(WEB_CMD_MESSENGER_SEND, bodyJSON)
}

// ---------------------------------------------------------------------------
// Gateway (12 methods)
// ---------------------------------------------------------------------------

// GatewayTotal gets the total count of gateway configurations.
func (h *infraHandler) GatewayTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_GATEWAY_TOTAL)
}

// GatewayNext gets the gateway configuration at the given index.
func (h *infraHandler) GatewayNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_GATEWAY_NEXT, index)
}

// GatewayGet gets a gateway configuration by name.
func (h *infraHandler) GatewayGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_GATEWAY_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// GatewayGetAll gets all gateway configurations.
func (h *infraHandler) GatewayGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_GATEWAY_TOTAL, WEB_CMD_GATEWAY_NEXT)
}

// GatewayAdd adds a new gateway configuration from JSON.
func (h *infraHandler) GatewayAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_GATEWAY_ADD, bodyJSON)
}

// GatewayDelete deletes a gateway configuration by name.
func (h *infraHandler) GatewayDelete(name string) MTRetCode {
	return h.infraDelete(WEB_CMD_GATEWAY_DELETE, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// GatewayShift shifts a gateway configuration position.
func (h *infraHandler) GatewayShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_GATEWAY_SHIFT, index, shift)
}

// GatewayRestart restarts a gateway by name.
func (h *infraHandler) GatewayRestart(name string) MTRetCode {
	return h.infraRestart(WEB_CMD_GATEWAY_RESTART, name)
}

// GatewayGetPosition gets position data for a gateway by name.
func (h *infraHandler) GatewayGetPosition(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_GATEWAY_GET_POSITION, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// GatewayModuleTotal gets the total count of gateway modules.
func (h *infraHandler) GatewayModuleTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_GATEWAY_MODULE_TOTAL)
}

// GatewayModuleNext gets the gateway module at the given index.
func (h *infraHandler) GatewayModuleNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_GATEWAY_MODULE_NEXT, index)
}

// GatewayModuleGet gets a gateway module by name.
func (h *infraHandler) GatewayModuleGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_GATEWAY_MODULE_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// ---------------------------------------------------------------------------
// Feeder (11 methods)
// ---------------------------------------------------------------------------

// FeederTotal gets the total count of feeder configurations.
func (h *infraHandler) FeederTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_FEEDER_TOTAL)
}

// FeederNext gets the feeder configuration at the given index.
func (h *infraHandler) FeederNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_FEEDER_NEXT, index)
}

// FeederGet gets a feeder configuration by name.
func (h *infraHandler) FeederGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_FEEDER_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// FeederGetAll gets all feeder configurations.
func (h *infraHandler) FeederGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_FEEDER_TOTAL, WEB_CMD_FEEDER_NEXT)
}

// FeederAdd adds a new feeder configuration from JSON.
func (h *infraHandler) FeederAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_FEEDER_ADD, bodyJSON)
}

// FeederDelete deletes a feeder configuration by name.
func (h *infraHandler) FeederDelete(name string) MTRetCode {
	return h.infraDelete(WEB_CMD_FEEDER_DELETE, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// FeederShift shifts a feeder configuration position.
func (h *infraHandler) FeederShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_FEEDER_SHIFT, index, shift)
}

// FeederRestart restarts a feeder by name.
func (h *infraHandler) FeederRestart(name string) MTRetCode {
	return h.infraRestart(WEB_CMD_FEEDER_RESTART, name)
}

// FeederModuleTotal gets the total count of feeder modules.
func (h *infraHandler) FeederModuleTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_FEEDER_MODULE_TOTAL)
}

// FeederModuleNext gets the feeder module at the given index.
func (h *infraHandler) FeederModuleNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_FEEDER_MODULE_NEXT, index)
}

// FeederModuleGet gets a feeder module by name.
func (h *infraHandler) FeederModuleGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_FEEDER_MODULE_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// ---------------------------------------------------------------------------
// Report (10 methods)
// ---------------------------------------------------------------------------

// ReportTotal gets the total count of report configurations.
func (h *infraHandler) ReportTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_REPORT_TOTAL)
}

// ReportNext gets the report configuration at the given index.
func (h *infraHandler) ReportNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_REPORT_NEXT, index)
}

// ReportGet gets a report configuration by name.
func (h *infraHandler) ReportGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_REPORT_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// ReportGetAll gets all report configurations.
func (h *infraHandler) ReportGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_REPORT_TOTAL, WEB_CMD_REPORT_NEXT)
}

// ReportAdd adds a new report configuration from JSON.
func (h *infraHandler) ReportAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_REPORT_ADD, bodyJSON)
}

// ReportDelete deletes a report configuration by name.
func (h *infraHandler) ReportDelete(name string) MTRetCode {
	return h.infraDelete(WEB_CMD_REPORT_DELETE, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// ReportShift shifts a report configuration position.
func (h *infraHandler) ReportShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_REPORT_SHIFT, index, shift)
}

// ReportModuleTotal gets the total count of report modules.
func (h *infraHandler) ReportModuleTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_REPORT_MODULE_TOTAL)
}

// ReportModuleNext gets the report module at the given index.
func (h *infraHandler) ReportModuleNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_REPORT_MODULE_NEXT, index)
}

// ReportModuleGet gets a report module by name.
func (h *infraHandler) ReportModuleGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_REPORT_MODULE_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// ---------------------------------------------------------------------------
// Plugin (10 methods)
// ---------------------------------------------------------------------------

// PluginTotal gets the total count of plugin configurations.
func (h *infraHandler) PluginTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_PLUGIN_TOTAL)
}

// PluginNext gets the plugin configuration at the given index.
func (h *infraHandler) PluginNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_PLUGIN_NEXT, index)
}

// PluginGet gets a plugin configuration by name.
func (h *infraHandler) PluginGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_PLUGIN_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// PluginGetAll gets all plugin configurations.
func (h *infraHandler) PluginGetAll() ([]map[string]interface{}, MTRetCode) {
	return h.infraGetAll(WEB_CMD_PLUGIN_TOTAL, WEB_CMD_PLUGIN_NEXT)
}

// PluginAdd adds a new plugin configuration from JSON.
func (h *infraHandler) PluginAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	return h.infraAdd(WEB_CMD_PLUGIN_ADD, bodyJSON)
}

// PluginDelete deletes a plugin configuration by name.
func (h *infraHandler) PluginDelete(name string) MTRetCode {
	return h.infraDelete(WEB_CMD_PLUGIN_DELETE, map[string]string{
		WEB_PARAM_NAME: name,
	})
}

// PluginShift shifts a plugin configuration position.
func (h *infraHandler) PluginShift(index, shift int) MTRetCode {
	return h.infraShift(WEB_CMD_PLUGIN_SHIFT, index, shift)
}

// PluginModuleTotal gets the total count of plugin modules.
func (h *infraHandler) PluginModuleTotal() (int, MTRetCode) {
	return h.infraTotal(WEB_CMD_PLUGIN_MODULE_TOTAL)
}

// PluginModuleNext gets the plugin module at the given index.
func (h *infraHandler) PluginModuleNext(index int) (map[string]interface{}, MTRetCode) {
	return h.infraNext(WEB_CMD_PLUGIN_MODULE_NEXT, index)
}

// PluginModuleGet gets a plugin module by name.
func (h *infraHandler) PluginModuleGet(name string) (map[string]interface{}, MTRetCode) {
	return h.infraGet(WEB_CMD_PLUGIN_MODULE_GET, map[string]string{
		WEB_PARAM_NAME: name,
	})
}
