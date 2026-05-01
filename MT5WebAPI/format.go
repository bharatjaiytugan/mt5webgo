//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import "fmt"

// GetError returns a human-readable error string for the given MTRetCode.
func GetError(code MTRetCode) string {
	switch code {
	case MT_RET_OK:
		return "ok"
	case MT_RET_OK_NONE:
		return "ok; no data"
	case MT_RET_ERROR:
		return "Common error"
	case MT_RET_ERR_PARAMS:
		return "Invalid parameters"
	case MT_RET_ERR_DATA:
		return "Invalid data"
	case MT_RET_ERR_DISK:
		return "Disk error"
	case MT_RET_ERR_MEM:
		return "Memory error"
	case MT_RET_ERR_NETWORK:
		return "Network error"
	case MT_RET_ERR_PERMISSIONS:
		return "Not enough permissions"
	case MT_RET_ERR_TIMEOUT:
		return "Operation timeout"
	case MT_RET_ERR_CONNECTION:
		return "No connection"
	case MT_RET_ERR_NOSERVICE:
		return "Service is not available"
	case MT_RET_ERR_FREQUENT:
		return "Too frequent requests"
	case MT_RET_ERR_NOTFOUND:
		return "Not found"
	case MT_RET_ERR_PARTIAL:
		return "Partial error"
	case MT_RET_ERR_SHUTDOWN:
		return "Server shutdown in progress"
	case MT_RET_ERR_CANCEL:
		return "Operation has been canceled"
	case MT_RET_ERR_DUPLICATE:
		return "Duplicate data"
	//--- authentication retcodes
	case MT_RET_AUTH_CLIENT_INVALID:
		return "Invalid terminal type"
	case MT_RET_AUTH_ACCOUNT_INVALID:
		return "Invalid account"
	case MT_RET_AUTH_ACCOUNT_DISABLED:
		return "Account disabled"
	case MT_RET_AUTH_ADVANCED:
		return "Advanced authorization necessary"
	case MT_RET_AUTH_CERTIFICATE:
		return "Certificate required"
	case MT_RET_AUTH_CERTIFICATE_BAD:
		return "Invalid certificate"
	case MT_RET_AUTH_NOTCONFIRMED:
		return "Certificate is not confirmed"
	case MT_RET_AUTH_SERVER_INTERNAL:
		return "Attempt to connect to non-access server"
	case MT_RET_AUTH_SERVER_BAD:
		return "Server is not authenticated"
	case MT_RET_AUTH_UPDATE_ONLY:
		return "Only updates available"
	case MT_RET_AUTH_CLIENT_OLD:
		return "Client has old version"
	case MT_RET_AUTH_MANAGER_NOCONFIG:
		return "Manager account does not have manager config"
	case MT_RET_AUTH_MANAGER_IPBLOCK:
		return "IP address unallowed for manager"
	case MT_RET_AUTH_GROUP_INVALID:
		return "Group is not initialized (server restart neccesary)"
	case MT_RET_AUTH_CA_DISABLED:
		return "Certificate generation disabled"
	case MT_RET_AUTH_INVALID_ID:
		return "Invalid or disabled server id [check server's id]"
	case MT_RET_AUTH_INVALID_IP:
		return "Unallowed address [check server's ip address]"
	case MT_RET_AUTH_INVALID_TYPE:
		return "Invalid server type [check server's id and type]"
	case MT_RET_AUTH_SERVER_BUSY:
		return "Server is busy"
	case MT_RET_AUTH_SERVER_CERT:
		return "Invalid server certificate"
	case MT_RET_AUTH_ACCOUNT_UNKNOWN:
		return "Unknown account"
	case MT_RET_AUTH_SERVER_OLD:
		return "Old server version"
	case MT_RET_AUTH_SERVER_LIMIT:
		return "Server cannot be connected due to license limitation"
	case MT_RET_AUTH_MOBILE_DISABLED:
		return "Mobile connection aren't allowed in server license"
	//--- config management retcodes
	case MT_RET_CFG_LAST_ADMIN:
		return "Last admin config deleting"
	case MT_RET_CFG_LAST_ADMIN_GROUP:
		return "Last admin group cannot be deleted"
	case MT_RET_CFG_NOT_EMPTY:
		return "Accounts or trades in group"
	case MT_RET_CFG_INVALID_RANGE:
		return "Invalid accounts or trades ranges"
	case MT_RET_CFG_NOT_MANAGER_LOGIN:
		return "Manager account is not from manager group"
	case MT_RET_CFG_BUILTIN:
		return "Built-in protected config"
	case MT_RET_CFG_DUPLICATE:
		return "Configuration duplicate"
	case MT_RET_CFG_LIMIT_REACHED:
		return "Configuration limit reached"
	case MT_RET_CFG_NO_ACCESS_TO_MAIN:
		return "Invalid network configuration"
	case MT_RET_CFG_DEALER_ID_EXIST:
		return "Dealer with same ID exists"
	case MT_RET_CFG_BIND_ADDR_EXIST:
		return "Bind address already exists"
	case MT_RET_CFG_WORKING_TRADE:
		return "Attempt to delete working trade server"
	//--- client management retcodes
	case MT_RET_USR_LAST_ADMIN:
		return "Last admin account deleting"
	case MT_RET_USR_LOGIN_EXHAUSTED:
		return "Logins range exhausted"
	case MT_RET_USR_LOGIN_PROHIBITED:
		return "Login reserved at another server"
	case MT_RET_USR_LOGIN_EXIST:
		return "Account already exists"
	case MT_RET_USR_SUICIDE:
		return "Attempt of self-deletion"
	case MT_RET_USR_INVALID_PASSWORD:
		return "Invalid account password"
	case MT_RET_USR_LIMIT_REACHED:
		return "Users limit reached"
	case MT_RET_USR_HAS_TRADES:
		return "Account has open trades"
	case MT_RET_USR_DIFFERENT_SERVERS:
		return "Attempt to move account to different server"
	case MT_RET_USR_DIFFERENT_CURRENCY:
		return "Attempt to move account to different currency group"
	case MT_RET_USR_IMPORT_BALANCE:
		return "Account balance import error"
	case MT_RET_USR_IMPORT_GROUP:
		return "Account import with invalid group"
	//--- trades management retcodes
	case MT_RET_TRADE_LIMIT_REACHED:
		return "Orders or deals limit reached"
	case MT_RET_TRADE_ORDER_EXIST:
		return "Order already exists"
	case MT_RET_TRADE_ORDER_EXHAUSTED:
		return "Orders range exhausted"
	case MT_RET_TRADE_DEAL_EXHAUSTED:
		return "Deals range exhausted"
	case MT_RET_TRADE_MAX_MONEY:
		return "Money limit reached"
	//--- report generation retcodes
	case MT_RET_REPORT_SNAPSHOT:
		return "Base snapshot error"
	case MT_RET_REPORT_NOTSUPPORTED:
		return "Method doesn't support for this report"
	case MT_RET_REPORT_NODATA:
		return "No report data"
	case MT_RET_REPORT_TEMPLATE_BAD:
		return "Bad template"
	case MT_RET_REPORT_TEMPLATE_END:
		return "End of template (template success processed)"
	case MT_RET_REPORT_INVALID_ROW:
		return "Invalid row size"
	case MT_RET_REPORT_LIMIT_REPEAT:
		return "Tag repeat limit reached"
	case MT_RET_REPORT_LIMIT_REPORT:
		return "Report size limit reached"
	//--- price history reports retcodes
	case MT_RET_HST_SYMBOL_NOTFOUND:
		return "Symbol not found; try to restart history server"
	//--- trade request retcodes
	case MT_RET_REQUEST_INWAY:
		return "Request on the way"
	case MT_RET_REQUEST_ACCEPTED:
		return "Request accepted"
	case MT_RET_REQUEST_PROCESS:
		return "Request processed"
	case MT_RET_REQUEST_REQUOTE:
		return "Request Requoted"
	case MT_RET_REQUEST_PRICES:
		return "Request Prices"
	case MT_RET_REQUEST_REJECT:
		return "Request rejected"
	case MT_RET_REQUEST_CANCEL:
		return "Request canceled"
	case MT_RET_REQUEST_PLACED:
		return "Order from requestplaced"
	case MT_RET_REQUEST_DONE:
		return "Request executed"
	case MT_RET_REQUEST_DONE_PARTIAL:
		return "Request executed partially"
	case MT_RET_REQUEST_ERROR:
		return "Request common error"
	case MT_RET_REQUEST_TIMEOUT:
		return "Request timeout"
	case MT_RET_REQUEST_INVALID:
		return "Invalid request"
	case MT_RET_REQUEST_INVALID_VOLUME:
		return "Invalid volume"
	case MT_RET_REQUEST_INVALID_PRICE:
		return "Invalid price"
	case MT_RET_REQUEST_INVALID_STOPS:
		return "Invalid stops or price"
	case MT_RET_REQUEST_TRADE_DISABLED:
		return "Trade disabled"
	case MT_RET_REQUEST_MARKET_CLOSED:
		return "Market closed"
	case MT_RET_REQUEST_NO_MONEY:
		return "Not enough money"
	case MT_RET_REQUEST_PRICE_CHANGED:
		return "Price changed"
	case MT_RET_REQUEST_PRICE_OFF:
		return "No prices"
	case MT_RET_REQUEST_INVALID_EXP:
		return "Invalid order expiration"
	case MT_RET_REQUEST_ORDER_CHANGED:
		return "Order has been changed already"
	case MT_RET_REQUEST_TOO_MANY:
		return "Too many trade requests"
	case MT_RET_REQUEST_NO_CHANGES:
		return "Request doesn't contain changes"
	case MT_RET_REQUEST_AT_DISABLED_SERVER:
		return "AutoTrading disabled by server"
	case MT_RET_REQUEST_AT_DISABLED_CLIENT:
		return "AutoTrading disabled by client"
	case MT_RET_REQUEST_LOCKED:
		return "Request locked by dealer"
	case MT_RET_REQUEST_FROZED:
		return "Order or position frozen"
	case MT_RET_REQUEST_INVALID_FILL:
		return "Unsupported filling mode"
	case MT_RET_REQUEST_CONNECTION:
		return "No connection"
	case MT_RET_REQUEST_ONLY_REAL:
		return "Allowed for real accounts only"
	case MT_RET_REQUEST_LIMIT_ORDERS:
		return "Orders limit reached"
	case MT_RET_REQUEST_LIMIT_VOLUME:
		return "Volume limit reached"
	//--- dealer retcodes
	case MT_RET_REQUEST_RETURN:
		return "Request returned in queue"
	case MT_RET_REQUEST_DONE_CANCEL:
		return "Request partially filled; remainder has been canceled"
	case MT_RET_REQUEST_REQUOTE_RETURN:
		return "Request requoted and returned in queue with new prices"
	//--- API retcodes
	case MT_RET_ERR_NOTIMPLEMENT:
		return "Not implement yet"
	case MT_RET_ERR_NOTMAIN:
		return "Operation must be performed on main server"
	case MT_RET_ERR_NOTSUPPORTED:
		return "Command doesn't supported"
	case MT_RET_ERR_DEADLOCK:
		return "Operation canceled due possible deadlock"
	case MT_RET_ERR_LOCKED:
		return "Operation on locked entity"
	//--- Messengers retcodes
	case MT_RET_MESSENGER_INVALID_PHONE:
		return "Invalid phone number"
	case MT_RET_MESSENGER_NOT_MOBILE:
		return "Phone number isn't mobile"
	//--- Subscriptions retcodes
	case MT_RET_SUBS_NOT_FOUND:
		return "Subscription is not found"
	case MT_RET_SUBS_NOT_FOUND_CFG:
		return "Subscription config is not found"
	case MT_RET_SUBS_NOT_FOUND_USER:
		return "User for subscription is not found"
	case MT_RET_SUBS_DISABLED:
		return "Subscription is disabled"
	case MT_RET_SUBS_PERMISSION_USER:
		return "Subscription is not allowed for user"
	case MT_RET_SUBS_PERMISSION_SUBSCRIBE:
		return "Subscribe is not allowed"
	case MT_RET_SUBS_PERMISSION_UNSUBSCRIBE:
		return "Unsubscribe is not allowed"
	case MT_RET_SUBS_REAL_ONLY:
		return "Subscriptions are available for real users only"
	}
	return "unknown error"
}

// GetErrorStandard returns a formatted error string like "message (CODE)" for logging.
func GetErrorStandard(code MTRetCode) string {
	return fmt.Sprintf("%s (%d)", GetError(code), code)
}
