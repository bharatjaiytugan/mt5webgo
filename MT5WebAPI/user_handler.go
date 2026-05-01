//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// userHandler implements user-related protocol commands.
type userHandler struct {
	apiBase
}

// Add adds a new user to the server.
func (h *userHandler) Add(user *MTUser) (*MTUser, MTRetCode) {
	if user == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_LOGIN:         formatUint64(user.Login),
		WEB_PARAM_PASS_MAIN:     user.MainPassword,
		WEB_PARAM_PASS_INVESTOR: user.InvestPassword,
		WEB_PARAM_RIGHTS:        formatUint64(uint64(user.Rights)),
		WEB_PARAM_GROUP:         user.Group,
		WEB_PARAM_NAME:          user.Name,
		WEB_PARAM_COMPANY:       user.Company,
		WEB_PARAM_LANGUAGE:      formatUint32(user.Language),
		WEB_PARAM_COUNTRY:       user.Country,
		WEB_PARAM_CITY:          user.City,
		WEB_PARAM_STATE:         user.State,
		WEB_PARAM_ZIPCODE:       user.ZIPCode,
		WEB_PARAM_ADDRESS:       user.Address,
		WEB_PARAM_PHONE:         user.Phone,
		WEB_PARAM_EMAIL:         user.Email,
		WEB_PARAM_ID:            user.ID,
		WEB_PARAM_STATUS:        user.Status,
		WEB_PARAM_COMMENT:       user.Comment,
		WEB_PARAM_PASS_PHONE:    user.PhonePassword,
		WEB_PARAM_COLOR:         formatUint32(user.Color),
		WEB_PARAM_LEVERAGE:      formatUint32(user.Leverage),
		WEB_PARAM_AGENT:         formatUint64(user.Agent),
		WEB_PARAM_BODYTEXT:      userToJSON(user),
	}
	answer := h.send(WEB_CMD_USER_ADD, params)
	if answer == nil {
		logWrite(LogTypeError, "send user add failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_USER_ADD, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user add failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	newUser := parseUserFromJSON(jsonStr)
	if newUser == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return newUser, MT_RET_OK
}

// Update updates an existing user on the server.
func (h *userHandler) Update(user *MTUser) (*MTUser, MTRetCode) {
	if user == nil {
		return nil, MT_RET_ERR_PARAMS
	}
	params := map[string]string{
		WEB_PARAM_LOGIN:         formatUint64(user.Login),
		WEB_PARAM_PASS_MAIN:     user.MainPassword,
		WEB_PARAM_PASS_INVESTOR: user.InvestPassword,
		WEB_PARAM_RIGHTS:        formatUint64(uint64(user.Rights)),
		WEB_PARAM_GROUP:         user.Group,
		WEB_PARAM_NAME:          user.Name,
		WEB_PARAM_COMPANY:       user.Company,
		WEB_PARAM_LANGUAGE:      formatUint32(user.Language),
		WEB_PARAM_COUNTRY:       user.Country,
		WEB_PARAM_CITY:          user.City,
		WEB_PARAM_STATE:         user.State,
		WEB_PARAM_ZIPCODE:       user.ZIPCode,
		WEB_PARAM_ADDRESS:       user.Address,
		WEB_PARAM_PHONE:         user.Phone,
		WEB_PARAM_EMAIL:         user.Email,
		WEB_PARAM_ID:            user.ID,
		WEB_PARAM_STATUS:        user.Status,
		WEB_PARAM_COMMENT:       user.Comment,
		WEB_PARAM_PASS_PHONE:    user.PhonePassword,
		WEB_PARAM_COLOR:         formatUint32(user.Color),
		WEB_PARAM_LEVERAGE:      formatUint32(user.Leverage),
		WEB_PARAM_AGENT:         formatUint64(user.Agent),
		WEB_PARAM_BODYTEXT:      userToJSON(user),
	}
	answer := h.send(WEB_CMD_USER_UPDATE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user update failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_USER_UPDATE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user update failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	newUser := parseUserFromJSON(jsonStr)
	if newUser == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return newUser, MT_RET_OK
}

// Delete deletes a user from the server.
func (h *userHandler) Delete(login uint64) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_DELETE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user delete failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_DELETE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user delete failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// Get retrieves user information by login.
func (h *userHandler) Get(login uint64) (*MTUser, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send user get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_USER_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	user := parseUserFromJSON(jsonStr)
	if user == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return user, MT_RET_OK
}

// PasswordCheck checks a user's password.
func (h *userHandler) PasswordCheck(login uint64, password string, passType EnUsersPasswords) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN:    formatUint64(login),
		WEB_PARAM_TYPE:     passwordTypeString(passType),
		WEB_PARAM_PASSWORD: password,
	}
	answer := h.send(WEB_CMD_USER_PASS_CHECK, params)
	if answer == nil {
		logWrite(LogTypeError, "send user password check failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_PASS_CHECK, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user password check failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// PasswordChange changes a user's password.
func (h *userHandler) PasswordChange(login uint64, newPassword string, passType EnUsersPasswords) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN:    formatUint64(login),
		WEB_PARAM_TYPE:     passwordTypeString(passType),
		WEB_PARAM_PASSWORD: newPassword,
	}
	answer := h.send(WEB_CMD_USER_PASS_CHANGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user password change failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_PASS_CHANGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user password change failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// DepositChange changes a user's deposit/balance.
func (h *userHandler) DepositChange(login uint64, deposit float64, comment string, dealType EnDealAction) MTRetCode {
	params := map[string]string{
		WEB_PARAM_LOGIN:   formatUint64(login),
		WEB_PARAM_TYPE:    formatUint32(uint32(dealType)),
		WEB_PARAM_BALANCE: formatFloat64(deposit),
		WEB_PARAM_COMMENT: comment,
	}
	answer := h.send(WEB_CMD_USER_DEPOSIT_CHANGE, params)
	if answer == nil {
		logWrite(LogTypeError, "send user deposit change failed")
		return MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	errCode := h.parseEmptyResult(WEB_CMD_USER_DEPOSIT_CHANGE, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse user deposit change failed: %s", GetErrorStandard(errCode)))
	}
	return errCode
}

// AccountGet retrieves account information for a user.
func (h *userHandler) AccountGet(login uint64) (*MTAccount, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_LOGIN: formatUint64(login),
	}
	answer := h.send(WEB_CMD_USER_ACCOUNT_GET, params)
	if answer == nil {
		logWrite(LogTypeError, "send account get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_USER_ACCOUNT_GET, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse account get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	acc := parseAccountFromJSON(jsonStr)
	if acc == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return acc, MT_RET_OK
}

// UserLogins retrieves user logins for a given group.
func (h *userHandler) UserLogins(group string) ([]uint64, MTRetCode) {
	params := map[string]string{
		WEB_PARAM_GROUP: group,
	}
	answer := h.send(WEB_CMD_USER_USER_LOGINS, params)
	if answer == nil {
		logWrite(LogTypeError, "send logins get failed")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result answer: %s", answerStr))
	}
	jsonStr, errCode := parseSingleJSON(WEB_CMD_USER_USER_LOGINS, answerStr)
	if errCode != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse logins get failed: %s", GetErrorStandard(errCode)))
		return nil, errCode
	}
	logins := parseLoginsFromJSON(jsonStr)
	if logins == nil {
		return nil, MT_RET_REPORT_NODATA
	}
	return logins, MT_RET_OK
}
