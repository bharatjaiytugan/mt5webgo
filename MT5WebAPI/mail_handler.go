//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// mailHandler implements mail-related protocol commands.
type mailHandler struct {
	apiBase
}

// MailSend sends a mail to a user or group of users.
func (h *mailHandler) MailSend(to, subject, text string) MTRetCode {
	params := map[string]string{
		WEB_PARAM_TO:       to,
		WEB_PARAM_SUBJECT:  subject,
		WEB_PARAM_BODYTEXT: text,
	}
	answer := h.send(WEB_CMD_MAIL_SEND, params)
	if answer == nil {
		logWrite(LogTypeError, "send email failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_MAIL_SEND, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse email failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}
