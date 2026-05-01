//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
)

// WEB_API_VERSION is the current version of the Web API.
const WEB_API_VERSION = 5830

// EnCryptModes defines the encryption mode for the connection.
type EnCryptModes int

const (
	CRYPT_MODE_NONE EnCryptModes = 0 // no encryption
	CRYPT_MODE_AES  EnCryptModes = 1 // AES-256 encryption
)

// authStartAnswer holds the parsed response from AUTH_START.
type authStartAnswer struct {
	RetCode string
	SrvRand string
}

func newAuthStartAnswer() *authStartAnswer {
	return &authStartAnswer{
		RetCode: "-1",
		SrvRand: "none",
	}
}

// authAnswer holds the parsed response from AUTH_ANSWER.
type authAnswer struct {
	RetCode   string
	CliRand   string
	CryptRand string
}

func newAuthAnswer() *authAnswer {
	return &authAnswer{
		RetCode: "-1",
		CliRand: "none",
	}
}

// mtAuth handles authentication with the MetaTrader 5 server.
type mtAuth struct {
	connect *mtConnect
	agent   string
}

// newMTAuth creates a new authentication handler.
func newMTAuth(connect *mtConnect, agent string) *mtAuth {
	return &mtAuth{
		connect: connect,
		agent:   agent,
	}
}

// Auth performs full authentication against the MetaTrader 5 server.
// Returns the crypt random bytes and the result code.
func (a *mtAuth) Auth(login uint64, password string, crypt EnCryptModes) ([]byte, MTRetCode) {
	// send auth start
	startAnswer, result := a.sendAuthStart(login, crypt)
	if result != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("auth start failed: [%d] %s", result, GetError(result)))
		return nil, result
	}
	// decode server random from hex
	randCode := GetFromHex(startAnswer.SrvRand)
	// generate client random
	randomCliCode := GetRandomHex(16)
	// hash password with server random
	hash := GetHashFromPassword(password, randCode)
	// send auth answer
	answerResult, result := a.sendAuthAnswer(hash, randomCliCode)
	if result != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("auth answer failed: [%d] %s", result, GetError(result)))
		return nil, result
	}
	// verify server's password hash
	hashPassword := GetHashFromPassword(password, randomCliCode)
	cliRandBytes := encodeUTF16LE(answerResult.CliRand)
	if CompareBytes(hashPassword, cliRandBytes) {
		logWrite(LogTypeError, fmt.Sprintf("server sent incorrect password hash: is: %s, my: %s", answerResult.CliRand, GetHex(hashPassword)))
		return nil, MT_RET_AUTH_SERVER_BAD
	}
	// get crypt rand from server
	cryptRand := GetFromHex(answerResult.CryptRand)
	return cryptRand, MT_RET_OK
}

// sendAuthStart sends the AUTH_START request.
func (a *mtAuth) sendAuthStart(login uint64, crypt EnCryptModes) (*authStartAnswer, MTRetCode) {
	data := map[string]string{
		WEB_PARAM_VERSION: fmt.Sprintf("%d", WEB_API_VERSION),
		WEB_PARAM_AGENT:   a.agent,
		WEB_PARAM_LOGIN:   fmt.Sprintf("%d", login),
		WEB_PARAM_TYPE:    MANAGER,
	}
	if crypt == CRYPT_MODE_NONE {
		data[WEB_PARAM_CRYPT_METHOD] = WEB_VAL_CRYPT_NONE
	} else {
		data[WEB_PARAM_CRYPT_METHOD] = WEB_VAL_CRYPT_AES256OFB
	}
	if !a.connect.Send(WEB_CMD_AUTH_START, data, true) {
		return nil, MT_RET_ERR_NETWORK
	}
	// read answer
	answer := a.connect.Read(true, false)
	if answer == nil {
		logWrite(LogTypeError, "answer auth start is empty")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result: %s", answerStr))
	}
	// parse answer
	startAnswer, result := a.parseAuthStart(answerStr)
	if result != MT_RET_OK {
		logWrite(LogTypeError, fmt.Sprintf("parse auth start failed: [%d] %s", result, GetError(result)))
		return nil, result
	}
	return startAnswer, MT_RET_OK
}

// parseAuthStart parses the AUTH_START response.
func (a *mtAuth) parseAuthStart(answer string) (*authStartAnswer, MTRetCode) {
	command, pos := getCommand(answer)
	if command != WEB_CMD_AUTH_START {
		logWrite(LogTypeError, fmt.Sprintf("answer command '%s' is incorrect, wait %s", command, WEB_CMD_AUTH_START))
		return nil, MT_RET_ERR_DATA
	}
	startAnswer := newAuthStartAnswer()
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			startAnswer.RetCode = param.Value
		case WEB_PARAM_SRV_RAND:
			startAnswer.SrvRand = param.Value
		}
	}
	// check ret code
	result := getRetCode(startAnswer.RetCode)
	if result != MT_RET_OK {
		return nil, result
	}
	// check SrvRand
	if startAnswer.SrvRand == "" || startAnswer.SrvRand == "none" {
		logWrite(LogTypeError, "srv rand answer incorrect")
		return nil, MT_RET_ERR_PARAMS
	}
	return startAnswer, MT_RET_OK
}

// sendAuthAnswer sends the AUTH_ANSWER request with password hash and client random.
func (a *mtAuth) sendAuthAnswer(hash, randomCliCode []byte) (*authAnswer, MTRetCode) {
	data := map[string]string{
		WEB_PARAM_SRV_RAND_ANSWER: GetHex(hash),
		WEB_PARAM_CLI_RAND:        GetHex(randomCliCode),
	}
	if !a.connect.Send(WEB_CMD_AUTH_ANSWER, data, false) {
		logWrite(LogTypeError, "send auth answer failed")
		return nil, MT_RET_ERR_NETWORK
	}
	// read answer
	answer := a.connect.Read(true, false)
	if answer == nil {
		logWrite(LogTypeError, "answer auth answer is empty")
		return nil, MT_RET_ERR_NETWORK
	}
	answerStr := GetString(answer)
	if IsWriteDebugLog() {
		logWrite(LogTypeDebug, fmt.Sprintf("result: %s", answerStr))
	}
	// parse answer
	authAns, result := parseAuthAnswer(answerStr)
	if result != MT_RET_OK {
		logWrite(LogTypeError, "parse auth answer failed")
		return nil, result
	}
	return authAns, MT_RET_OK
}

// parseAuthAnswer parses the AUTH_ANSWER response.
func parseAuthAnswer(answer string) (*authAnswer, MTRetCode) {
	command, pos := getCommand(answer)
	if command != WEB_CMD_AUTH_ANSWER {
		logWrite(LogTypeError, fmt.Sprintf("type answer '%s' is incorrect, is not %s", command, WEB_CMD_AUTH_ANSWER))
		return nil, MT_RET_ERR_DATA
	}
	ans := newAuthAnswer()
	posEnd := -1
	for {
		param := getNextParam(answer, &pos, &posEnd)
		if param == nil {
			break
		}
		switch param.Name {
		case WEB_PARAM_RETCODE:
			ans.RetCode = param.Value
		case WEB_PARAM_CLI_RAND_ANSWER:
			ans.CliRand = param.Value
		case WEB_PARAM_CRYPT_RAND:
			ans.CryptRand = param.Value
		}
	}
	// check ret code
	result := getRetCode(ans.RetCode)
	if result != MT_RET_OK {
		return nil, result
	}
	// check CliRand
	if ans.CliRand == "" || ans.CliRand == "none" {
		logWrite(LogTypeError, "cli rand answer incorrect")
		return nil, MT_RET_ERR_PARAMS
	}
	return ans, MT_RET_OK
}
