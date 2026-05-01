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

// Protocol command constants.
const (
	// authorization
	WEB_CMD_AUTH_START  = "AUTH_START"  // begin of authorization on server
	WEB_CMD_AUTH_ANSWER = "AUTH_ANSWER" // end of authorization on server
	// config
	WEB_CMD_COMMON_GET = "COMMON_GET" // get config
	// time
	WEB_CMD_TIME_SERVER = "TIME_SERVER" // get time of server
	WEB_CMD_TIME_GET    = "TIME_GET"    // get config of times
	// api
	WEB_PREFIX_WEBAPI  = "MT5WEBAPI%04x%04x%s" // format of package for api
	WEB_PACKET_FORMAT  = "%04x%04x%s"           // format of packages second or more request
	WEB_API_WORD       = "WebAPI"
	// params
	WEB_PARAM_VERSION        = "VERSION"        // version of authorization
	WEB_PARAM_RETCODE        = "RETCODE"        // code answer
	WEB_PARAM_LOGIN          = "LOGIN"          // login
	WEB_PARAM_TYPE           = "TYPE"           // type of connection, type of data, type of operation
	WEB_PARAM_AGENT          = "AGENT"          // agent name
	WEB_PARAM_SRV_RAND       = "SRV_RAND"       // server random string
	WEB_PARAM_SRV_RAND_ANSWER = "SRV_RAND_ANSWER" // answer on server random string
	WEB_PARAM_CLI_RAND       = "CLI_RAND"       // client's random string
	WEB_PARAM_CLI_RAND_ANSWER = "CLI_RAND_ANSWER" // answer to clients random string
	WEB_PARAM_TIME           = "TIME"           // time param
	WEB_PARAM_TOTAL          = "TOTAL"          // total
	WEB_PARAM_INDEX          = "INDEX"          // index
	WEB_PARAM_GROUP          = "GROUP"          // group
	WEB_PARAM_SYMBOL         = "SYMBOL"         // symbol
	WEB_PARAM_NAME           = "NAME"           // name
	WEB_PARAM_COMPANY        = "COMPANY"        // company
	WEB_PARAM_LANGUAGE       = "LANGUAGE"        // language (LANGID)
	WEB_PARAM_COUNTRY        = "COUNTRY"        // country
	WEB_PARAM_CITY           = "CITY"           // city
	WEB_PARAM_STATE          = "STATE"          // state
	WEB_PARAM_ZIPCODE        = "ZIPCODE"        // zipcode
	WEB_PARAM_ADDRESS        = "ADDRESS"        // address
	WEB_PARAM_PHONE          = "PHONE"          // phone
	WEB_PARAM_EMAIL          = "EMAIL"          // email
	WEB_PARAM_ID             = "ID"             // id
	WEB_PARAM_STATUS         = "STATUS"         // status
	WEB_PARAM_COMMENT        = "COMMENT"        // comment
	WEB_PARAM_COLOR          = "COLOR"          // color
	WEB_PARAM_PASS_MAIN      = "PASS_MAIN"      // main password
	WEB_PARAM_PASS_INVESTOR  = "PASS_INVESTOR"  // invest password
	WEB_PARAM_PASS_API       = "PASS_API"       // API password
	WEB_PARAM_PASS_PHONE     = "PASS_PHONE"     // phone password
	WEB_PARAM_LEVERAGE       = "LEVERAGE"       // leverage
	WEB_PARAM_RIGHTS         = "RIGHTS"         // rights
	WEB_PARAM_BALANCE        = "BALANCE"        // balance
	WEB_PARAM_PASSWORD       = "PASSWORD"       // password
	WEB_PARAM_TICKET         = "TICKET"         // ticket
	WEB_PARAM_OFFSET         = "OFFSET"         // offset for page requests
	WEB_PARAM_FROM           = "FROM"           // from time
	WEB_PARAM_TO             = "TO"             // to time
	WEB_PARAM_TRANS_ID       = "TRANS_ID"       // trans id
	WEB_PARAM_SUBJECT        = "SUBJECT"        // subject
	WEB_PARAM_CATEGORY       = "CATEGORY"       // category
	WEB_PARAM_PRIORITY       = "PRIORITY"       // priority
	WEB_PARAM_BODYTEXT       = "BODY_TEXT"       // big text
	WEB_PARAM_CHECK_MARGIN   = "CHECK_MARGIN"   // check margin
	// crypt
	WEB_PARAM_CRYPT_METHOD = "CRYPT_METHOD" // method of crypt
	WEB_PARAM_CRYPT_RAND   = "CRYPT_RAND"   // random string for crypt
	// group
	WEB_CMD_GROUP_TOTAL  = "GROUP_TOTAL"  // get count groups
	WEB_CMD_GROUP_NEXT   = "GROUP_NEXT"   // get next group
	WEB_CMD_GROUP_GET    = "GROUP_GET"    // get info about group
	WEB_CMD_GROUP_DELETE = "GROUP_DELETE" // group delete
	WEB_CMD_GROUP_ADD    = "GROUP_ADD"    // group add
	// symbols
	WEB_CMD_SYMBOL_TOTAL     = "SYMBOL_TOTAL"     // get count symbols
	WEB_CMD_SYMBOL_NEXT      = "SYMBOL_NEXT"      // get next symbol
	WEB_CMD_SYMBOL_GET       = "SYMBOL_GET"       // get info about symbol
	WEB_CMD_SYMBOL_GET_GROUP = "SYMBOL_GET_GROUP" // get info about symbol group
	WEB_CMD_SYMBOL_DELETE    = "SYMBOL_DELETE"    // symbol delete
	WEB_CMD_SYMBOL_ADD       = "SYMBOL_ADD"       // symbol add
	// user
	WEB_CMD_USER_ADD         = "USER_ADD"         // add new user
	WEB_CMD_USER_UPDATE      = "USER_UPDATE"      // update user
	WEB_CMD_USER_DELETE      = "USER_DELETE"      // delete user
	WEB_CMD_USER_GET         = "USER_GET"         // get user information
	WEB_CMD_USER_PASS_CHECK  = "USER_PASS_CHECK"  // user check
	WEB_CMD_USER_PASS_CHANGE = "USER_PASS_CHANGE" // password change
	WEB_CMD_USER_ACCOUNT_GET = "USER_ACCOUNT_GET" // account info get
	WEB_CMD_USER_USER_LOGINS = "USER_LOGINS"      // users logins get
	// password type
	WEB_VAL_USER_PASS_MAIN     = "MAIN"
	WEB_VAL_USER_PASS_INVESTOR = "INVESTOR"
	WEB_VAL_USER_PASS_API      = "API"
	// crypt modes
	WEB_VAL_CRYPT_NONE      = "NONE"
	WEB_VAL_CRYPT_AES256OFB = "AES256OFB"
	// trade command
	WEB_CMD_USER_DEPOSIT_CHANGE = "USER_DEPOSIT_CHANGE" // deposit change
	// work with order
	WEB_CMD_ORDER_GET       = "ORDER_GET"       // get order
	WEB_CMD_ORDER_GET_TOTAL = "ORDER_GET_TOTAL" // get count orders
	WEB_CMD_ORDER_GET_PAGE  = "ORDER_GET_PAGE"  // get order from history
	// work with position
	WEB_CMD_POSITION_GET       = "POSITION_GET"       // get position
	WEB_CMD_POSITION_GET_TOTAL = "POSITION_GET_TOTAL" // get count positions
	WEB_CMD_POSITION_GET_PAGE  = "POSITION_GET_PAGE"  // get positions
	// work with deal
	WEB_CMD_DEAL_GET       = "DEAL_GET"       // get deal
	WEB_CMD_DEAL_GET_TOTAL = "DEAL_GET_TOTAL" // get count deals
	WEB_CMD_DEAL_GET_PAGE  = "DEAL_GET_PAGE"  // get list of deals
	// work with history
	WEB_CMD_HISTORY_GET       = "HISTORY_GET"       // get history
	WEB_CMD_HISTORY_GET_TOTAL = "HISTORY_GET_TOTAL" // get count of history order
	WEB_CMD_HISTORY_GET_PAGE  = "HISTORY_GET_PAGE"  // get list of history
	// work with ticks
	WEB_CMD_TICK_LAST       = "TICK_LAST"       // get tick
	WEB_CMD_TICK_LAST_GROUP = "TICK_LAST_GROUP" // get tick by group name
	WEB_CMD_TICK_STAT       = "TICK_STAT"       // tick stat
	// mail
	WEB_CMD_MAIL_SEND = "MAIL_SEND"
	// news
	WEB_CMD_NEWS_SEND = "NEWS_SEND"
	// chart / book
	WEB_CMD_CHART_GET = "CHART_GET" // get 1-minute bar history
	WEB_CMD_BOOK_GET  = "BOOK_GET"  // get market depth (DOM)
	// data param for chart
	WEB_PARAM_DATA = "DATA" // data flags: d=datetime, o=open, h=high, l=low, c=close, t=tick volume, v=real volume, s=spread
	// ping
	WEB_CMD_PING = "PING"
	// trade
	WEB_CMD_TRADE_BALANCE = "TRADE_BALANCE"
	// server restart
	WEB_CMD_SERVER_RESTART = "SERVER_RESTART"
	// Manager
	MANAGER = "Manager"
)

// HEADER_LENGTH is the length of the protocol header in bytes.
const HEADER_LENGTH = 9

// protocolHeader represents the header of a protocol packet.
type protocolHeader struct {
	SizeBody     int
	NumberPacket int
	Flag         int
}

// parseHeader parses a protocol header from raw byte data.
func parseHeader(data []byte) *protocolHeader {
	if data == nil {
		return nil
	}
	if len(data) < HEADER_LENGTH {
		return nil
	}
	result := &protocolHeader{}
	// parse size body (4 hex chars)
	sizeBody, err := strconv.ParseInt(string(data[0:4]), 16, 32)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("get size body failed: %v", err))
		return nil
	}
	result.SizeBody = int(sizeBody)
	// parse number of packet (4 hex chars)
	numberPacket, err := strconv.ParseInt(string(data[4:8]), 16, 32)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("get number packet failed: %v", err))
		return nil
	}
	result.NumberPacket = int(numberPacket)
	// parse flag (1 hex char)
	flag, err := strconv.ParseInt(string(data[8:9]), 16, 8)
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("get flag failed: %v", err))
		return nil
	}
	result.Flag = int(flag)
	return result
}

// getCommand extracts the command name from an answer string.
// Returns the command and the position of the '|' separator.
func getCommand(answer string) (string, int) {
	if answer == "" {
		return "", -1
	}
	pos := strings.Index(answer, "|")
	if pos > 0 {
		return answer[:pos], pos
	}
	return "", -1
}

// answerParam represents a single parameter from the protocol answer.
type answerParam struct {
	Name  string
	Value string
}

// getNextParam extracts the next parameter from the answer string starting at pos.
// posEnd tracks the end of parameters (\r\n). Both pos and posEnd are updated.
func getNextParam(answer string, pos *int, posEnd *int) *answerParam {
	// find end of params on first call
	if *posEnd < 0 {
		idx := strings.Index(answer, "\r\n")
		if idx < 0 {
			*posEnd = len(answer)
		} else {
			*posEnd = idx
		}
	}
	if (*pos + 1) >= len(answer) || (*pos+1) >= *posEnd {
		return nil
	}
	paramName := ""
	paramVal := ""
	currParam := true
	*pos++
	for *pos < len(answer) && *pos < *posEnd {
		symbol := answer[*pos]
		if symbol == '=' {
			currParam = false
			*pos++
			continue
		}
		if symbol == '|' {
			break
		}
		if currParam {
			paramName += string(symbol)
		} else {
			paramVal += string(symbol)
		}
		*pos++
	}
	return &answerParam{Name: paramName, Value: paramVal}
}

// getRetCode parses a return code string into an MTRetCode value.
func getRetCode(retCode string) MTRetCode {
	if retCode == "" {
		return MT_RET_ERROR
	}
	parts := strings.SplitN(retCode, " ", 2)
	result, err := strconv.Atoi(parts[0])
	if err != nil {
		logWrite(LogTypeError, fmt.Sprintf("error parsing answer code from string: %s", retCode))
		return MT_RET_ERROR
	}
	return MTRetCode(result)
}

// getJSON extracts the JSON body from an answer string starting at the given position.
func getJSON(answer string, pos int) string {
	idx := strings.Index(answer[pos:], "\n")
	if idx >= 0 {
		return strings.TrimSpace(answer[pos+idx:])
	}
	return ""
}

// baseAnswer represents a base answer from the MT server.
type baseAnswer struct {
	RetCode string
}

// newBaseAnswer creates a new baseAnswer with default values.
func newBaseAnswer() *baseAnswer {
	return &baseAnswer{RetCode: "-1"}
}

// baseAnswerJSON represents a base answer with JSON data from the MT server.
type baseAnswerJSON struct {
	RetCode    string
	ConfigJSON string
}

// MAX_LENGTH_JSON is the maximum length of JSON data accepted.
const MAX_LENGTH_JSON = 52428800
