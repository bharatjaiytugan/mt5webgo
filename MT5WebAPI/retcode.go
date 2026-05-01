//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

// MTRetCode represents error codes returned by the MetaTrader 5 Web API.
type MTRetCode int

const (
	//--- successfully codes
	MT_RET_OK      MTRetCode = 0 // ok
	MT_RET_OK_NONE MTRetCode = 1 // ok, no data

	//--- common errors
	MT_RET_ERROR           MTRetCode = 2  // Common error
	MT_RET_ERR_PARAMS      MTRetCode = 3  // Invalid parameters
	MT_RET_ERR_DATA        MTRetCode = 4  // Invalid data
	MT_RET_ERR_DISK        MTRetCode = 5  // Disk error
	MT_RET_ERR_MEM         MTRetCode = 6  // Memory error
	MT_RET_ERR_NETWORK     MTRetCode = 7  // Network error
	MT_RET_ERR_PERMISSIONS MTRetCode = 8  // Not enough permissions
	MT_RET_ERR_TIMEOUT     MTRetCode = 9  // Operation timeout
	MT_RET_ERR_CONNECTION  MTRetCode = 10 // No connection
	MT_RET_ERR_NOSERVICE   MTRetCode = 11 // Service is not available
	MT_RET_ERR_FREQUENT    MTRetCode = 12 // Too frequent requests
	MT_RET_ERR_NOTFOUND    MTRetCode = 13 // Not found
	MT_RET_ERR_PARTIAL     MTRetCode = 14 // Partial error
	MT_RET_ERR_SHUTDOWN    MTRetCode = 15 // Server shutdown in progress
	MT_RET_ERR_CANCEL      MTRetCode = 16 // Operation has been canceled
	MT_RET_ERR_DUPLICATE   MTRetCode = 17 // Duplicate data

	//--- authentication retcodes
	MT_RET_AUTH_CLIENT_INVALID  MTRetCode = 1000 // Invalid terminal type
	MT_RET_AUTH_ACCOUNT_INVALID MTRetCode = 1001 // Invalid account
	MT_RET_AUTH_ACCOUNT_DISABLED MTRetCode = 1002 // Account disabled
	MT_RET_AUTH_ADVANCED        MTRetCode = 1003 // Advanced authorization necessary
	MT_RET_AUTH_CERTIFICATE     MTRetCode = 1004 // Certificate required
	MT_RET_AUTH_CERTIFICATE_BAD MTRetCode = 1005 // Invalid certificate
	MT_RET_AUTH_NOTCONFIRMED    MTRetCode = 1006 // Certificate is not confirmed
	MT_RET_AUTH_SERVER_INTERNAL MTRetCode = 1007 // Attempt to connect to non-access server
	MT_RET_AUTH_SERVER_BAD      MTRetCode = 1008 // Server isn't authenticated
	MT_RET_AUTH_UPDATE_ONLY     MTRetCode = 1009 // Only updates available
	MT_RET_AUTH_CLIENT_OLD      MTRetCode = 1010 // Client has old version
	MT_RET_AUTH_MANAGER_NOCONFIG MTRetCode = 1011 // Manager account doesn't have manager config
	MT_RET_AUTH_MANAGER_IPBLOCK MTRetCode = 1012 // IP address unallowed for manager
	MT_RET_AUTH_GROUP_INVALID   MTRetCode = 1013 // Group is not initialized (server restart neccesary)
	MT_RET_AUTH_CA_DISABLED     MTRetCode = 1014 // Certificate generation disabled
	MT_RET_AUTH_INVALID_ID      MTRetCode = 1015 // Invalid or disabled server id [check server's id]
	MT_RET_AUTH_INVALID_IP      MTRetCode = 1016 // Unallowed address [check server's ip address]
	MT_RET_AUTH_INVALID_TYPE    MTRetCode = 1017 // Invalid server type [check server's id and type]
	MT_RET_AUTH_SERVER_BUSY     MTRetCode = 1018 // Server is busy
	MT_RET_AUTH_SERVER_CERT     MTRetCode = 1019 // Invalid server certificate
	MT_RET_AUTH_ACCOUNT_UNKNOWN MTRetCode = 1020 // Unknown account
	MT_RET_AUTH_SERVER_OLD      MTRetCode = 1021 // Old server version
	MT_RET_AUTH_SERVER_LIMIT    MTRetCode = 1022 // Server cannot be connected due to license limitation
	MT_RET_AUTH_MOBILE_DISABLED MTRetCode = 1023 // Mobile connection aren't allowed in server license

	//--- config management retcodes
	MT_RET_CFG_LAST_ADMIN       MTRetCode = 2000 // Last admin config deleting
	MT_RET_CFG_LAST_ADMIN_GROUP MTRetCode = 2001 // Last admin group cannot be deleted
	MT_RET_CFG_NOT_EMPTY        MTRetCode = 2003 // Accounts or trades in group
	MT_RET_CFG_INVALID_RANGE    MTRetCode = 2004 // Invalid accounts or trades ranges
	MT_RET_CFG_NOT_MANAGER_LOGIN MTRetCode = 2005 // Manager account is not from manager group
	MT_RET_CFG_BUILTIN          MTRetCode = 2006 // Built-in protected config
	MT_RET_CFG_DUPLICATE        MTRetCode = 2007 // Configuration duplicate
	MT_RET_CFG_LIMIT_REACHED    MTRetCode = 2008 // Configuration limit reached
	MT_RET_CFG_NO_ACCESS_TO_MAIN MTRetCode = 2009 // Invalid network configuration
	MT_RET_CFG_DEALER_ID_EXIST  MTRetCode = 2010 // Dealer with same ID exists
	MT_RET_CFG_BIND_ADDR_EXIST  MTRetCode = 2011 // Bind address already exists
	MT_RET_CFG_WORKING_TRADE    MTRetCode = 2012 // Attempt to delete working trade server

	//--- client management retcodes
	MT_RET_USR_LAST_ADMIN        MTRetCode = 3001 // Last admin account deleting
	MT_RET_USR_LOGIN_EXHAUSTED   MTRetCode = 3002 // Logins range exhausted
	MT_RET_USR_LOGIN_PROHIBITED  MTRetCode = 3003 // Login reserved at another server
	MT_RET_USR_LOGIN_EXIST       MTRetCode = 3004 // Account already exists
	MT_RET_USR_SUICIDE           MTRetCode = 3005 // Attempt of self-deletion
	MT_RET_USR_INVALID_PASSWORD  MTRetCode = 3006 // Invalid account password
	MT_RET_USR_LIMIT_REACHED     MTRetCode = 3007 // Users limit reached
	MT_RET_USR_HAS_TRADES        MTRetCode = 3008 // Account has open trades
	MT_RET_USR_DIFFERENT_SERVERS MTRetCode = 3009 // Attempt to move account to different server
	MT_RET_USR_DIFFERENT_CURRENCY MTRetCode = 3010 // Attempt to move account to different currency group
	MT_RET_USR_IMPORT_BALANCE    MTRetCode = 3011 // Account balance import error
	MT_RET_USR_IMPORT_GROUP      MTRetCode = 3012 // Account import with invalid group

	//--- trades management retcodes
	MT_RET_TRADE_LIMIT_REACHED  MTRetCode = 4001 // Orders or deals limit reached
	MT_RET_TRADE_ORDER_EXIST    MTRetCode = 4002 // Order already exists
	MT_RET_TRADE_ORDER_EXHAUSTED MTRetCode = 4003 // Orders range exhausted
	MT_RET_TRADE_DEAL_EXHAUSTED MTRetCode = 4004 // Deals range exhausted
	MT_RET_TRADE_MAX_MONEY      MTRetCode = 4005 // Money limit reached

	//--- report generation retcodes
	MT_RET_REPORT_SNAPSHOT      MTRetCode = 5001 // Base snapshot error
	MT_RET_REPORT_NOTSUPPORTED  MTRetCode = 5002 // Method doesn't support for this report
	MT_RET_REPORT_NODATA        MTRetCode = 5003 // No report data
	MT_RET_REPORT_TEMPLATE_BAD  MTRetCode = 5004 // Bad template
	MT_RET_REPORT_TEMPLATE_END  MTRetCode = 5005 // End of template (template success processed)
	MT_RET_REPORT_INVALID_ROW   MTRetCode = 5006 // Invalid row size
	MT_RET_REPORT_LIMIT_REPEAT  MTRetCode = 5007 // Tag repeat limit reached
	MT_RET_REPORT_LIMIT_REPORT  MTRetCode = 5008 // Report size limit reached

	//--- price history reports retcodes
	MT_RET_HST_SYMBOL_NOTFOUND MTRetCode = 6001 // Symbol not found, try to restart history server

	//--- trade request retcodes
	MT_RET_REQUEST_INWAY              MTRetCode = 10001 // Request on the way
	MT_RET_REQUEST_ACCEPTED           MTRetCode = 10002 // Request accepted
	MT_RET_REQUEST_PROCESS            MTRetCode = 10003 // Request processed
	MT_RET_REQUEST_REQUOTE            MTRetCode = 10004 // Request Requoted
	MT_RET_REQUEST_PRICES             MTRetCode = 10005 // Request Prices
	MT_RET_REQUEST_REJECT             MTRetCode = 10006 // Request rejected
	MT_RET_REQUEST_CANCEL             MTRetCode = 10007 // Request canceled
	MT_RET_REQUEST_PLACED             MTRetCode = 10008 // Order from requestplaced
	MT_RET_REQUEST_DONE               MTRetCode = 10009 // Request executed
	MT_RET_REQUEST_DONE_PARTIAL       MTRetCode = 10010 // Request executed partially
	MT_RET_REQUEST_ERROR              MTRetCode = 10011 // Request common error
	MT_RET_REQUEST_TIMEOUT            MTRetCode = 10012 // Request timeout
	MT_RET_REQUEST_INVALID            MTRetCode = 10013 // Invalid request
	MT_RET_REQUEST_INVALID_VOLUME     MTRetCode = 10014 // Invalid volume
	MT_RET_REQUEST_INVALID_PRICE      MTRetCode = 10015 // Invalid price
	MT_RET_REQUEST_INVALID_STOPS      MTRetCode = 10016 // Invalid stops or price
	MT_RET_REQUEST_TRADE_DISABLED     MTRetCode = 10017 // Trade disabled
	MT_RET_REQUEST_MARKET_CLOSED      MTRetCode = 10018 // Market closed
	MT_RET_REQUEST_NO_MONEY           MTRetCode = 10019 // Not enough money
	MT_RET_REQUEST_PRICE_CHANGED      MTRetCode = 10020 // Price changed
	MT_RET_REQUEST_PRICE_OFF          MTRetCode = 10021 // No prices
	MT_RET_REQUEST_INVALID_EXP        MTRetCode = 10022 // Invalid order expiration
	MT_RET_REQUEST_ORDER_CHANGED      MTRetCode = 10023 // Order has been changed already
	MT_RET_REQUEST_TOO_MANY           MTRetCode = 10024 // Too many trade requests
	MT_RET_REQUEST_NO_CHANGES         MTRetCode = 10025 // Request doesn't contain changes
	MT_RET_REQUEST_AT_DISABLED_SERVER MTRetCode = 10026 // AutoTrading disabled by server
	MT_RET_REQUEST_AT_DISABLED_CLIENT MTRetCode = 10027 // AutoTrading disabled by client
	MT_RET_REQUEST_LOCKED             MTRetCode = 10028 // Request locked by dealer
	MT_RET_REQUEST_FROZED             MTRetCode = 10029 // Order or position frozen
	MT_RET_REQUEST_INVALID_FILL       MTRetCode = 10030 // Unsupported filling mode
	MT_RET_REQUEST_CONNECTION         MTRetCode = 10031 // No connection
	MT_RET_REQUEST_ONLY_REAL          MTRetCode = 10032 // Allowed for real accounts only
	MT_RET_REQUEST_LIMIT_ORDERS       MTRetCode = 10033 // Orders limit reached
	MT_RET_REQUEST_LIMIT_VOLUME       MTRetCode = 10034 // Volume limit reached

	//--- dealer retcodes
	MT_RET_REQUEST_RETURN        MTRetCode = 11000 // Request returned in queue
	MT_RET_REQUEST_DONE_CANCEL   MTRetCode = 11001 // Request partially filled, remainder has been canceled
	MT_RET_REQUEST_REQUOTE_RETURN MTRetCode = 11002 // Request requoted and returned in queue with new prices

	//--- API retcodes
	MT_RET_ERR_NOTIMPLEMENT MTRetCode = 12000 // Not implement yet
	MT_RET_ERR_NOTMAIN      MTRetCode = 12001 // Operation must be performed on main server
	MT_RET_ERR_NOTSUPPORTED MTRetCode = 12002 // Command doesn't supported
	MT_RET_ERR_DEADLOCK     MTRetCode = 12003 // Operation canceled due possible deadlock
	MT_RET_ERR_LOCKED       MTRetCode = 12004 // Operation on locked entity

	//--- Messengers retcodes
	MT_RET_MESSENGER_INVALID_PHONE MTRetCode = 14000 // Invalid phone number
	MT_RET_MESSENGER_NOT_MOBILE    MTRetCode = 14001 // Phone number isn't mobile

	//--- Subscriptions retcodes
	MT_RET_SUBS_NOT_FOUND              MTRetCode = 15000 // Subscription is not found
	MT_RET_SUBS_NOT_FOUND_CFG          MTRetCode = 15001 // Subscription config is not found
	MT_RET_SUBS_NOT_FOUND_USER         MTRetCode = 15002 // User for subscription is not found
	MT_RET_SUBS_DISABLED               MTRetCode = 15003 // Subscription is disabled
	MT_RET_SUBS_PERMISSION_USER        MTRetCode = 15004 // Subscription is not allowed for user
	MT_RET_SUBS_PERMISSION_SUBSCRIBE   MTRetCode = 15005 // Subscribe is not allowed
	MT_RET_SUBS_PERMISSION_UNSUBSCRIBE MTRetCode = 15006 // Unsubscribe is not allowed
	MT_RET_SUBS_REAL_ONLY              MTRetCode = 15007 // Subscriptions are available for real users only
)
