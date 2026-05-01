//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// newsHandler implements news-related protocol commands.
type newsHandler struct {
	apiBase
}

// NewsSend sends a news item to users.
func (h *newsHandler) NewsSend(subject, category string, language, priority uint32, text string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_SUBJECT:  subject,
		WEB_PARAM_CATEGORY: category,
		WEB_PARAM_LANGUAGE: formatUint32(language),
		WEB_PARAM_PRIORITY: formatUint32(priority),
		WEB_PARAM_BODYTEXT: text,
	}
	answer := h.send(WEB_CMD_NEWS_SEND, params)
	if answer == nil {
		logWrite(LogTypeError, "send news failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_NEWS_SEND, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse news failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}
