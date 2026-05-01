//+------------------------------------------------------------------+
//|                                             MetaTrader 5 Web API |
//|                             Copyright 2000-2025, MetaQuotes Ltd. |
//|                                               www.metaquotes.net |
//+------------------------------------------------------------------+
package mt5webapi

import (
	"fmt"
)

// WEB_API_DATE is the date of the current Web API version.
const WEB_API_DATE = "24 Apr 2026"

// EnPumpModes represents pump modes.
type EnPumpModes int

const (
	PUMP_MODE_NONE EnPumpModes = 0x00000000
)

// CONNECTION_TIMEOUT is the default connection timeout in milliseconds.
const CONNECTION_TIMEOUT = 5000

// MT5WebAPI is the main class of the MetaTrader 5 Web API.
type MT5WebAPI struct {
	agent        string
	crypt        EnCryptModes
	asyncConnect *mtAsyncConnect
	timeBase     *timeHandler
	commonBase   *commonHandler
	groupBase    *groupHandler
	symbolBase   *symbolHandler
	orderBase    *orderHandler
	positionBase *positionHandler
	dealBase     *dealHandler
	historyBase  *historyHandler
	tickBase     *tickHandler
	userBase     *userHandler
	mailBase     *mailHandler
	newsBase     *newsHandler
	pingBase     *pingHandler
	customBase   *customHandler
	tradeBase    *tradeHandler
	dealerBase       *dealerHandler
	serverConfigBase *serverConfigHandler
	batchBase        *batchHandler
	backupBase       *backupHandler
	tradeCalcBase    *tradeCalcHandler
	clientBase       *clientHandler
	chartBase        *chartHandler
	dailyBase        *dailyHandler
	settingsBase       *settingsHandler
	subscriptionBase   *subscriptionHandler
	configCrudBase     *configCrudHandler
	infraBase          *infraHandler
	miscBase           *miscHandler
}

// NewMT5WebAPI creates a new MT5WebAPI instance with logging callback.
func NewMT5WebAPI(agent string, logWriteFunc CallbackLogWrite) *MT5WebAPI {
	api := &MT5WebAPI{
		agent: agent,
		crypt: CRYPT_MODE_AES,
	}
	logInit(agent, logWriteFunc)
	return api
}

// NewMT5WebAPIDefault creates a new MT5WebAPI instance with default settings.
func NewMT5WebAPIDefault() *MT5WebAPI {
	return &MT5WebAPI{
		agent: "WEBAPI",
		crypt: CRYPT_MODE_AES,
	}
}

// IsConnected returns true if connected to the MT5 server.
func (api *MT5WebAPI) IsConnected() bool {
	return api.asyncConnect != nil && api.asyncConnect.IsConnected()
}

// GetWriteDebugLog returns whether debug logging is enabled.
func (api *MT5WebAPI) GetWriteDebugLog() bool {
	return IsWriteDebugLog()
}

// EnableWriteDebugLog enables or disables debug logging.
func (api *MT5WebAPI) EnableWriteDebugLog(value bool) {
	SetWriteDebugLog(value)
}

// Connect establishes a connection to the MetaTrader 5 server.
func (api *MT5WebAPI) Connect(server string, port int, login uint64, password string, pumpModes EnPumpModes, crypt EnCryptModes, timeout int) MTRetCode {
	// check server
	if server == "" {
		logWrite(LogTypeError, "server name is empty")
		return MT_RET_ERR_PARAMS
	}
	// check port
	if port == 0 {
		logWrite(LogTypeError, "port number is invalid")
		return MT_RET_ERR_PARAMS
	}
	// check password
	if password == "" {
		logWrite(LogTypeError, "password is empty")
		return MT_RET_ERR_PARAMS
	}
	// stop existing connection
	if api.asyncConnect != nil {
		api.asyncConnect.Stop()
	}
	api.crypt = crypt
	// sync connection, only for authorization on MT server
	conn := newMTConnect(server, port, timeout)
	result := conn.Connect()
	if result != MT_RET_OK {
		return result
	}
	// authorization to MetaTrader 5 server
	auth := newMTAuth(conn, api.agent)
	cryptRand, result := auth.Auth(login, password, api.crypt)
	if result != MT_RET_OK {
		conn.Disconnect()
		api.Disconnect()
		return result
	}
	// create async connection
	api.asyncConnect = newMTAsyncConnect()
	if api.asyncConnect.Start(conn.Socket(), crypt, timeout, cryptRand, password) {
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "connect to MT started")
		}
		api.changeConnection()
	} else {
		if IsWriteDebugLog() {
			logWrite(LogTypeDebug, "connect to MT failed")
		}
		return MT_RET_ERR_NETWORK
	}
	return MT_RET_OK
}

// ConnectDefault establishes a connection with default encryption and timeout.
func (api *MT5WebAPI) ConnectDefault(server string, port int, login uint64, password string, pumpModes EnPumpModes) MTRetCode {
	return api.Connect(server, port, login, password, pumpModes, CRYPT_MODE_AES, CONNECTION_TIMEOUT)
}

// Disconnect disconnects from the MetaTrader 5 server and stops all goroutines.
func (api *MT5WebAPI) Disconnect() {
	if api.asyncConnect != nil {
		api.asyncConnect.Stop()
	}
}

// TimeGet gets current time configuration from server.
func (api *MT5WebAPI) TimeGet() (*MTConTime, MTRetCode) {
	if api.timeBase == nil {
		api.timeBase = newTimeHandler(api.asyncConnect)
	}
	return api.timeBase.TimeGet()
}

// TimeServer gets current server time in unix format.
func (api *MT5WebAPI) TimeServer() int64 {
	if api.timeBase == nil {
		api.timeBase = newTimeHandler(api.asyncConnect)
	}
	return api.timeBase.TimeServer()
}

// CommonGet gets common server information.
func (api *MT5WebAPI) CommonGet() (*MTConCommon, MTRetCode) {
	if api.commonBase == nil {
		api.commonBase = newCommonHandler(api.asyncConnect)
	}
	return api.commonBase.CommonGet()
}

// GroupTotal gets the total count of groups.
func (api *MT5WebAPI) GroupTotal() (int, MTRetCode) {
	if api.groupBase == nil {
		api.groupBase = newGroupHandler(api.asyncConnect)
	}
	return api.groupBase.GroupTotal()
}

// GroupNext gets the next group by position.
func (api *MT5WebAPI) GroupNext(pos uint32) (*MTConGroup, MTRetCode) {
	if api.groupBase == nil {
		api.groupBase = newGroupHandler(api.asyncConnect)
	}
	return api.groupBase.GroupNext(pos)
}

// GroupGet gets a group by name.
func (api *MT5WebAPI) GroupGet(name string) (*MTConGroup, MTRetCode) {
	if api.groupBase == nil {
		api.groupBase = newGroupHandler(api.asyncConnect)
	}
	return api.groupBase.GroupGet(name)
}

// GroupAdd adds a new group.
func (api *MT5WebAPI) GroupAdd(group *MTConGroup) (*MTConGroup, MTRetCode) {
	if api.groupBase == nil {
		api.groupBase = newGroupHandler(api.asyncConnect)
	}
	return api.groupBase.GroupAdd(group)
}

// GroupDelete deletes a group by name.
func (api *MT5WebAPI) GroupDelete(name string) MTRetCode {
	if api.groupBase == nil {
		api.groupBase = newGroupHandler(api.asyncConnect)
	}
	return api.groupBase.GroupDelete(name)
}

// SymbolTotal gets the total count of symbols.
func (api *MT5WebAPI) SymbolTotal() (int, MTRetCode) {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolTotal()
}

// SymbolNext gets the next symbol by position.
func (api *MT5WebAPI) SymbolNext(pos int) (*MTConSymbol, MTRetCode) {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolNext(pos)
}

// SymbolAdd adds a new symbol.
func (api *MT5WebAPI) SymbolAdd(symbol *MTConSymbol) (*MTConSymbol, MTRetCode) {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolAdd(symbol)
}

// SymbolDelete deletes a symbol by name.
func (api *MT5WebAPI) SymbolDelete(name string) MTRetCode {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolDelete(name)
}

// SymbolGet gets a symbol by name.
func (api *MT5WebAPI) SymbolGet(name string) (*MTConSymbol, MTRetCode) {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolGet(name)
}

// SymbolGetGroup gets a symbol configuration for a specific group.
func (api *MT5WebAPI) SymbolGetGroup(name, group string) (*MTConSymbol, MTRetCode) {
	if api.symbolBase == nil {
		api.symbolBase = newSymbolHandler(api.asyncConnect)
	}
	return api.symbolBase.SymbolGetGroup(name, group)
}

// OrderGet gets an order by ticket.
func (api *MT5WebAPI) OrderGet(ticket uint64) (*MTOrder, MTRetCode) {
	if api.orderBase == nil {
		api.orderBase = newOrderHandler(api.asyncConnect)
	}
	return api.orderBase.OrderGet(ticket)
}

// OrderGetTotal gets the total count of orders for a user.
func (api *MT5WebAPI) OrderGetTotal(login uint64) (uint32, MTRetCode) {
	if api.orderBase == nil {
		api.orderBase = newOrderHandler(api.asyncConnect)
	}
	return api.orderBase.OrderGetTotal(login)
}

// OrderGetPage gets orders by page.
func (api *MT5WebAPI) OrderGetPage(login uint64, offset, total uint32) ([]*MTOrder, MTRetCode) {
	if api.orderBase == nil {
		api.orderBase = newOrderHandler(api.asyncConnect)
	}
	return api.orderBase.OrderGetPage(login, offset, total)
}

// PositionGet gets a position by login and symbol.
func (api *MT5WebAPI) PositionGet(login uint64, symbol string) (*MTPosition, MTRetCode) {
	if api.positionBase == nil {
		api.positionBase = newPositionHandler(api.asyncConnect)
	}
	return api.positionBase.PositionGet(login, symbol)
}

// PositionGetTotal gets the total count of positions for a user.
func (api *MT5WebAPI) PositionGetTotal(login uint64) (uint32, MTRetCode) {
	if api.positionBase == nil {
		api.positionBase = newPositionHandler(api.asyncConnect)
	}
	return api.positionBase.PositionGetTotal(login)
}

// PositionGetPage gets positions by page.
func (api *MT5WebAPI) PositionGetPage(login uint64, offset, total uint32) ([]*MTPosition, MTRetCode) {
	if api.positionBase == nil {
		api.positionBase = newPositionHandler(api.asyncConnect)
	}
	return api.positionBase.PositionGetPage(login, offset, total)
}

// DealGet gets a deal by ticket.
func (api *MT5WebAPI) DealGet(ticket uint64) (*MTDeal, MTRetCode) {
	if api.dealBase == nil {
		api.dealBase = newDealHandler(api.asyncConnect)
	}
	return api.dealBase.DealGet(ticket)
}

// DealGetTotal gets the total count of deals for a user in a date range.
func (api *MT5WebAPI) DealGetTotal(login uint64, from, to int64) (uint32, MTRetCode) {
	if api.dealBase == nil {
		api.dealBase = newDealHandler(api.asyncConnect)
	}
	return api.dealBase.DealGetTotal(login, from, to)
}

// DealGetPage gets deals by page.
func (api *MT5WebAPI) DealGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTDeal, MTRetCode) {
	if api.dealBase == nil {
		api.dealBase = newDealHandler(api.asyncConnect)
	}
	return api.dealBase.DealGetPage(login, from, to, offset, total)
}

// HistoryGet gets a historical order by ticket.
func (api *MT5WebAPI) HistoryGet(ticket uint64) (*MTOrder, MTRetCode) {
	if api.historyBase == nil {
		api.historyBase = newHistoryHandler(api.asyncConnect)
	}
	return api.historyBase.HistoryGet(ticket)
}

// HistoryGetTotal gets the total count of historical orders.
func (api *MT5WebAPI) HistoryGetTotal(login uint64, from, to int64) (uint32, MTRetCode) {
	if api.historyBase == nil {
		api.historyBase = newHistoryHandler(api.asyncConnect)
	}
	return api.historyBase.HistoryGetTotal(login, from, to)
}

// HistoryGetPage gets historical orders by page.
func (api *MT5WebAPI) HistoryGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTOrder, MTRetCode) {
	if api.historyBase == nil {
		api.historyBase = newHistoryHandler(api.asyncConnect)
	}
	return api.historyBase.HistoryGetPage(login, from, to, offset, total)
}

// TickLast gets last ticks for a symbol.
func (api *MT5WebAPI) TickLast(symbol string) ([]*MTTick, MTRetCode) {
	if api.tickBase == nil {
		api.tickBase = newTickHandler(api.asyncConnect)
	}
	return api.tickBase.TickLast(symbol)
}

// TickLastGroup gets last ticks for a symbol and group.
func (api *MT5WebAPI) TickLastGroup(symbol, group string) ([]*MTTick, MTRetCode) {
	if api.tickBase == nil {
		api.tickBase = newTickHandler(api.asyncConnect)
	}
	return api.tickBase.TickLastGroup(symbol, group)
}

// TickStat gets tick statistics for a symbol.
func (api *MT5WebAPI) TickStat(symbol string) ([]*MTTickStat, MTRetCode) {
	if api.tickBase == nil {
		api.tickBase = newTickHandler(api.asyncConnect)
	}
	return api.tickBase.TickStat(symbol)
}

// UserAdd adds a new user to the server.
func (api *MT5WebAPI) UserAdd(user *MTUser) (*MTUser, MTRetCode) {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.Add(user)
}

// UserUpdate updates a user on the server.
func (api *MT5WebAPI) UserUpdate(user *MTUser) (*MTUser, MTRetCode) {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.Update(user)
}

// UserDelete deletes a user from the server.
func (api *MT5WebAPI) UserDelete(login uint64) MTRetCode {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.Delete(login)
}

// UserGet gets a user by login.
func (api *MT5WebAPI) UserGet(login uint64) (*MTUser, MTRetCode) {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.Get(login)
}

// UserPasswordCheck checks a user's password.
func (api *MT5WebAPI) UserPasswordCheck(login uint64, password string, passType EnUsersPasswords) MTRetCode {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.PasswordCheck(login, password, passType)
}

// UserPasswordCheckMain checks a user's main password.
func (api *MT5WebAPI) UserPasswordCheckMain(login uint64, password string) MTRetCode {
	return api.UserPasswordCheck(login, password, UserPassMain)
}

// UserPasswordChange changes a user's password.
func (api *MT5WebAPI) UserPasswordChange(login uint64, newPassword string, passType EnUsersPasswords) MTRetCode {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.PasswordChange(login, newPassword, passType)
}

// UserPasswordChangeMain changes a user's main password.
func (api *MT5WebAPI) UserPasswordChangeMain(login uint64, newPassword string) MTRetCode {
	return api.UserPasswordChange(login, newPassword, UserPassMain)
}

// UserDepositChange changes a user's deposit.
func (api *MT5WebAPI) UserDepositChange(login uint64, newDeposit float64, comment string, dealType EnDealAction) MTRetCode {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.DepositChange(login, newDeposit, comment, dealType)
}

// UserAccountGet gets account information for a user.
func (api *MT5WebAPI) UserAccountGet(login uint64) (*MTAccount, MTRetCode) {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.AccountGet(login)
}

// UserLogins gets a list of user logins for a group.
func (api *MT5WebAPI) UserLogins(group string) ([]uint64, MTRetCode) {
	if api.userBase == nil {
		api.userBase = newUserHandler(api.asyncConnect)
	}
	return api.userBase.UserLogins(group)
}

// MailSend sends mail to a user.
func (api *MT5WebAPI) MailSend(to, subject, text string) MTRetCode {
	if api.mailBase == nil {
		api.mailBase = newMailHandler(api.asyncConnect)
	}
	return api.mailBase.MailSend(to, subject, text)
}

// NewsSend sends news to users.
func (api *MT5WebAPI) NewsSend(subject, category string, language, priority uint32, text string) MTRetCode {
	if api.newsBase == nil {
		api.newsBase = newNewsHandler(api.asyncConnect)
	}
	return api.newsBase.NewsSend(subject, category, language, priority, text)
}

// TradeBalance performs a trade balance operation.
func (api *MT5WebAPI) TradeBalance(login uint64, dealType EnDealAction, balance float64, comment string) MTRetCode {
	if api.tradeBase == nil {
		api.tradeBase = newTradeHandler(api.asyncConnect)
	}
	return api.tradeBase.TradeBalance(login, dealType, balance, comment)
}

// TradeBalanceWithTicket performs a trade balance operation and returns the ticket.
func (api *MT5WebAPI) TradeBalanceWithTicket(login uint64, dealType EnDealAction, balance float64, comment string) (uint64, MTRetCode) {
	if api.tradeBase == nil {
		api.tradeBase = newTradeHandler(api.asyncConnect)
	}
	return api.tradeBase.TradeBalanceWithTicket(login, dealType, balance, comment)
}

// TradeBalanceWithMarginCheck performs a trade balance operation with margin check.
func (api *MT5WebAPI) TradeBalanceWithMarginCheck(login uint64, dealType EnDealAction, balance float64, comment string, checkMargin bool) (uint64, MTRetCode) {
	if api.tradeBase == nil {
		api.tradeBase = newTradeHandler(api.asyncConnect)
	}
	return api.tradeBase.TradeBalanceWithMarginCheck(login, dealType, balance, comment, checkMargin)
}

// Ping sends a ping to the server.
func (api *MT5WebAPI) Ping() MTRetCode {
	if api.pingBase == nil {
		api.pingBase = newPingHandler(api.asyncConnect)
	}
	return api.pingBase.PingSend()
}

// CustomSend sends a custom command to the MT server.
func (api *MT5WebAPI) CustomSend(command string, params map[string]string, body string) ([]byte, MTRetCode) {
	if api.customBase == nil {
		api.customBase = newCustomHandler(api.asyncConnect)
	}
	return api.customBase.CustomSend(command, params, body)
}

// DealerSend sends a dealer trade request to the server.
func (api *MT5WebAPI) DealerSend(req *TradeRequest) (*TradeResult, MTRetCode) {
	if api.dealerBase == nil {
		api.dealerBase = newDealerHandler(api.asyncConnect)
	}
	return api.dealerBase.DealerSend(req)
}

// DealerGetResult retrieves the result of a previously submitted dealer request.
func (api *MT5WebAPI) DealerGetResult(id uint64) (*TradeResult, MTRetCode) {
	if api.dealerBase == nil {
		api.dealerBase = newDealerHandler(api.asyncConnect)
	}
	return api.dealerBase.DealerGetResult(id)
}

// DealerBuyMarket sends a market buy order via dealer.
func (api *MT5WebAPI) DealerBuyMarket(login uint64, symbol string, volume uint64, comment string) (*TradeResult, MTRetCode) {
	// get current price
	ticks, ret := api.TickLast(symbol)
	if ret != MT_RET_OK || len(ticks) == 0 {
		logWrite(LogTypeError, fmt.Sprintf("dealer buy market: failed to get tick for %s", symbol))
		return nil, ret
	}
	// get symbol digits
	sym, ret := api.SymbolGet(symbol)
	if ret != MT_RET_OK || sym == nil {
		logWrite(LogTypeError, fmt.Sprintf("dealer buy market: failed to get symbol %s", symbol))
		return nil, ret
	}
	req := &TradeRequest{
		Action:     TA_DEALER_POS_EXECUTE,
		Login:      login,
		Symbol:     symbol,
		Volume:     volume,
		Type:       uint32(OpBuy),
		TypeFill:   uint32(OrderFillFOK),
		PriceOrder: ticks[0].Ask,
		Digits:     sym.Digits,
		Comment:    comment,
	}
	return api.DealerSend(req)
}

// DealerSellMarket sends a market sell order via dealer.
func (api *MT5WebAPI) DealerSellMarket(login uint64, symbol string, volume uint64, comment string) (*TradeResult, MTRetCode) {
	// get current price
	ticks, ret := api.TickLast(symbol)
	if ret != MT_RET_OK || len(ticks) == 0 {
		logWrite(LogTypeError, fmt.Sprintf("dealer sell market: failed to get tick for %s", symbol))
		return nil, ret
	}
	// get symbol digits
	sym, ret := api.SymbolGet(symbol)
	if ret != MT_RET_OK || sym == nil {
		logWrite(LogTypeError, fmt.Sprintf("dealer sell market: failed to get symbol %s", symbol))
		return nil, ret
	}
	req := &TradeRequest{
		Action:     TA_DEALER_POS_EXECUTE,
		Login:      login,
		Symbol:     symbol,
		Volume:     volume,
		Type:       uint32(OpSell),
		TypeFill:   uint32(OrderFillFOK),
		PriceOrder: ticks[0].Bid,
		Digits:     sym.Digits,
		Comment:    comment,
	}
	return api.DealerSend(req)
}

// DealerPendingOrder sends a pending order (limit/stop) via dealer.
func (api *MT5WebAPI) DealerPendingOrder(login uint64, symbol string, orderType EnOrderType, volume uint64, price float64, sl, tp float64, comment string) (*TradeResult, MTRetCode) {
	// get symbol digits
	sym, ret := api.SymbolGet(symbol)
	if ret != MT_RET_OK || sym == nil {
		logWrite(LogTypeError, fmt.Sprintf("dealer pending order: failed to get symbol %s", symbol))
		return nil, ret
	}
	req := &TradeRequest{
		Action:     TA_DEALER_ORD_PENDING,
		Login:      login,
		Symbol:     symbol,
		Volume:     volume,
		Type:       uint32(orderType),
		TypeFill:   uint32(OrderFillReturn),
		PriceOrder: price,
		PriceSL:    sl,
		PriceTP:    tp,
		Digits:     sym.Digits,
		Comment:    comment,
	}
	return api.DealerSend(req)
}

// DealerModifyOrder modifies an existing pending order via dealer.
func (api *MT5WebAPI) DealerModifyOrder(login uint64, order uint64, price float64, sl, tp float64) (*TradeResult, MTRetCode) {
	req := &TradeRequest{
		Action:     TA_DEALER_ORD_MODIFY,
		Login:      login,
		Order:      order,
		PriceOrder: price,
		PriceSL:    sl,
		PriceTP:    tp,
	}
	return api.DealerSend(req)
}

// DealerCancelOrder cancels a pending order via dealer.
func (api *MT5WebAPI) DealerCancelOrder(login uint64, order uint64) (*TradeResult, MTRetCode) {
	req := &TradeRequest{
		Action: TA_DEALER_ORD_REMOVE,
		Login:  login,
		Order:  order,
	}
	return api.DealerSend(req)
}

// CommonSet updates the common server configuration.
func (api *MT5WebAPI) CommonSet(config *MTConCommon) (*MTConCommon, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.CommonSet(config)
}

// AccessServerTotal gets the total count of access servers.
func (api *MT5WebAPI) AccessServerTotal() (int, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerTotal()
}

// AccessServerNext gets the access server config at the given index.
func (api *MT5WebAPI) AccessServerNext(index int) (map[string]interface{}, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerNext(index)
}

// AccessServerGet gets an access server config by ID.
func (api *MT5WebAPI) AccessServerGet(id int) (map[string]interface{}, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerGet(id)
}

// AccessServerAdd adds a new access server configuration.
func (api *MT5WebAPI) AccessServerAdd(serverJSON string) (map[string]interface{}, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerAdd(serverJSON)
}

// AccessServerDelete deletes an access server by ID.
func (api *MT5WebAPI) AccessServerDelete(id int) MTRetCode {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerDelete(id)
}

// AccessServerShift shifts an access server position.
func (api *MT5WebAPI) AccessServerShift(index, shift int) MTRetCode {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.AccessServerShift(index, shift)
}

// TLSCertTotal gets the total count of TLS certificates.
func (api *MT5WebAPI) TLSCertTotal() (int, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TLSCertTotal()
}

// TLSCertNext gets the TLS certificate config at the given index.
func (api *MT5WebAPI) TLSCertNext(index int) (map[string]interface{}, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TLSCertNext(index)
}

// TLSCertAdd adds a new TLS certificate.
func (api *MT5WebAPI) TLSCertAdd(certJSON string) MTRetCode {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TLSCertAdd(certJSON)
}

// TLSCertDelete deletes a TLS certificate by index.
func (api *MT5WebAPI) TLSCertDelete(index int) MTRetCode {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TLSCertDelete(index)
}

// TLSCertShift shifts a TLS certificate position.
func (api *MT5WebAPI) TLSCertShift(index, shift int) MTRetCode {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TLSCertShift(index, shift)
}

// TimeSet updates the server time configuration.
func (api *MT5WebAPI) TimeSet(config *MTConTime) (*MTConTime, MTRetCode) {
	if api.serverConfigBase == nil {
		api.serverConfigBase = newServerConfigHandler(api.asyncConnect)
	}
	return api.serverConfigBase.TimeSet(config)
}

// GroupList retrieves a list of group names matching the given mask.
func (api *MT5WebAPI) GroupList(mask string) ([]string, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.GroupList(mask)
}

// GroupAddBatch adds multiple groups at once from a JSON array of group configs.
func (api *MT5WebAPI) GroupAddBatch(groupsJSON string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.GroupAddBatch(groupsJSON)
}

// GroupDeleteBatch deletes multiple groups by comma-separated names.
func (api *MT5WebAPI) GroupDeleteBatch(groups string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.GroupDeleteBatch(groups)
}

// SymbolAddBatch adds multiple symbols at once from a JSON array of symbol configs.
func (api *MT5WebAPI) SymbolAddBatch(symbolsJSON string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.SymbolAddBatch(symbolsJSON)
}

// SymbolDeleteBatch deletes multiple symbols by comma-separated names.
func (api *MT5WebAPI) SymbolDeleteBatch(symbols string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.SymbolDeleteBatch(symbols)
}

// OrderGetBatch retrieves orders matching the given logins, groups, and tickets (comma-separated).
func (api *MT5WebAPI) OrderGetBatch(logins, groups, tickets string) ([]*MTOrder, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.OrderGetBatch(logins, groups, tickets)
}

// OrderUpdateDirect updates an order directly (identified by ticket in JSON body).
func (api *MT5WebAPI) OrderUpdateDirect(orderJSON string) (*MTOrder, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.OrderUpdateDirect(orderJSON)
}

// OrderDeleteDirect deletes orders directly by comma-separated tickets.
func (api *MT5WebAPI) OrderDeleteDirect(tickets string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.OrderDeleteDirect(tickets)
}

// OrderCancel cancels orders by comma-separated tickets, moving them to history.
func (api *MT5WebAPI) OrderCancel(tickets string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.OrderCancel(tickets)
}

// HistoryGetBatch retrieves historical orders matching the given filters.
func (api *MT5WebAPI) HistoryGetBatch(logins, groups, tickets string, from, to int64) ([]*MTOrder, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.HistoryGetBatch(logins, groups, tickets, from, to)
}

// HistoryUpdateDirect updates a historical order directly (identified by ticket in JSON body).
func (api *MT5WebAPI) HistoryUpdateDirect(orderJSON string) (*MTOrder, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.HistoryUpdateDirect(orderJSON)
}

// HistoryDeleteDirect deletes historical orders by comma-separated tickets.
func (api *MT5WebAPI) HistoryDeleteDirect(tickets string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.HistoryDeleteDirect(tickets)
}

// DealGetBatch retrieves deals matching the given logins, groups, and tickets (comma-separated) within a time range.
func (api *MT5WebAPI) DealGetBatch(logins, groups, tickets string, from, to int64) ([]*MTDeal, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.DealGetBatch(logins, groups, tickets, from, to)
}

// DealUpdateDirect updates a deal directly (identified by Deal ticket in JSON body).
func (api *MT5WebAPI) DealUpdateDirect(dealJSON string) (*MTDeal, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.DealUpdateDirect(dealJSON)
}

// DealDeleteDirect deletes deals by comma-separated tickets.
func (api *MT5WebAPI) DealDeleteDirect(tickets string) MTRetCode {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.DealDeleteDirect(tickets)
}

// PositionGetBatch gets positions by logins, groups, or tickets.
func (api *MT5WebAPI) PositionGetBatch(logins, groups, tickets string) ([]*MTPosition, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.PositionGetBatch(logins, groups, tickets)
}

// UserGetBatch gets users by logins or groups (comma-separated).
func (api *MT5WebAPI) UserGetBatch(logins, groups string) ([]*MTUser, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.UserGetBatch(logins, groups)
}

// UserBackupList gets user backup dates in a time range.
func (api *MT5WebAPI) UserBackupList(from, to int64) ([]int64, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.UserBackupList(from, to)
}

// UserBackupGet gets users from a specific backup for a given login.
func (api *MT5WebAPI) UserBackupGet(backupDate int64, login uint64) ([]*MTUser, MTRetCode) {
	if api.batchBase == nil {
		api.batchBase = newBatchHandler(api.asyncConnect)
	}
	return api.batchBase.UserBackupGet(backupDate, login)
}

// OrderBackupList gets available backup dates in a time range.
func (api *MT5WebAPI) OrderBackupList(from, to int64) ([]int64, MTRetCode) {
	if api.backupBase == nil {
		api.backupBase = newBackupHandler(api.asyncConnect)
	}
	return api.backupBase.OrderBackupList(from, to)
}

// OrderBackupGet gets orders from a specific backup for a given login.
func (api *MT5WebAPI) OrderBackupGet(backupDate int64, login uint64) ([]*MTOrder, MTRetCode) {
	if api.backupBase == nil {
		api.backupBase = newBackupHandler(api.asyncConnect)
	}
	return api.backupBase.OrderBackupGet(backupDate, login)
}

// OrderBackupRestore restores an order from backup using its JSON representation.
func (api *MT5WebAPI) OrderBackupRestore(orderJSON string) (*MTOrder, MTRetCode) {
	if api.backupBase == nil {
		api.backupBase = newBackupHandler(api.asyncConnect)
	}
	return api.backupBase.OrderBackupRestore(orderJSON)
}

// OrderReopen reopens a pending order from history by ticket.
func (api *MT5WebAPI) OrderReopen(ticket uint64) (*MTOrder, MTRetCode) {
	if api.backupBase == nil {
		api.backupBase = newBackupHandler(api.asyncConnect)
	}
	return api.backupBase.OrderReopen(ticket)
}

// CalcRateBuy calculates the buy conversion rate between two currencies.
func (api *MT5WebAPI) CalcRateBuy(base, currency string) (float64, MTRetCode) {
	if api.tradeCalcBase == nil {
		api.tradeCalcBase = newTradeCalcHandler(api.asyncConnect)
	}
	return api.tradeCalcBase.CalcRateBuy(base, currency)
}

// CalcRateSell calculates the sell conversion rate between two currencies.
func (api *MT5WebAPI) CalcRateSell(base, currency string) (float64, MTRetCode) {
	if api.tradeCalcBase == nil {
		api.tradeCalcBase = newTradeCalcHandler(api.asyncConnect)
	}
	return api.tradeCalcBase.CalcRateSell(base, currency)
}

// CheckMargin checks the margin required for an order.
// orderType: 0=buy, 1=sell. volume is in 1/10000 lots (e.g. 10000 = 1.0 lot).
func (api *MT5WebAPI) CheckMargin(login uint64, symbol string, orderType uint32, volume uint64, price float64) (map[string]interface{}, MTRetCode) {
	if api.tradeCalcBase == nil {
		api.tradeCalcBase = newTradeCalcHandler(api.asyncConnect)
	}
	return api.tradeCalcBase.CheckMargin(login, symbol, orderType, volume, price)
}

// CalcProfit calculates the profit for a trade.
// orderType: 0=buy, 1=sell. volume is in 1/10000 lots (e.g. 10000 = 1.0 lot).
func (api *MT5WebAPI) CalcProfit(group, symbol string, orderType uint32, volume uint64, priceOpen, priceClose float64) (float64, MTRetCode) {
	if api.tradeCalcBase == nil {
		api.tradeCalcBase = newTradeCalcHandler(api.asyncConnect)
	}
	return api.tradeCalcBase.CalcProfit(group, symbol, orderType, volume, priceOpen, priceClose)
}

// ClientAdd adds a new client record to the server.
func (api *MT5WebAPI) ClientAdd(clientJSON string) (*MTClient, MTRetCode) {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientAdd(clientJSON)
}

// ClientUpdate updates an existing client record on the server.
func (api *MT5WebAPI) ClientUpdate(clientJSON string) (*MTClient, MTRetCode) {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientUpdate(clientJSON)
}

// ClientDelete deletes a client record by record ID.
func (api *MT5WebAPI) ClientDelete(recordID uint64) MTRetCode {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientDelete(recordID)
}

// ClientGet retrieves a client record by record ID.
func (api *MT5WebAPI) ClientGet(recordID uint64) (*MTClient, MTRetCode) {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientGet(recordID)
}

// ClientGetIDs retrieves client record IDs available to the manager.
// group is optional filter (e.g. "demo\\*"). Pass "" for all.
func (api *MT5WebAPI) ClientGetIDs(group string) ([]uint64, MTRetCode) {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientGetIDs(group)
}

// ClientUserAdd links a user login to a client record.
func (api *MT5WebAPI) ClientUserAdd(clientID, login uint64) MTRetCode {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientUserAdd(clientID, login)
}

// ClientUserDelete unlinks a user login from a client record.
func (api *MT5WebAPI) ClientUserDelete(clientID, login uint64) MTRetCode {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientUserDelete(clientID, login)
}

// ClientUserLogins retrieves user logins linked to a client record.
func (api *MT5WebAPI) ClientUserLogins(clientID uint64) ([]uint64, MTRetCode) {
	if api.clientBase == nil {
		api.clientBase = newClientHandler(api.asyncConnect)
	}
	return api.clientBase.ClientUserLogins(clientID)
}

// ChartGet retrieves 1-minute bar history for a symbol.
// data specifies fields: d=datetime, o=open, h=high, l=low, c=close, t=tick volume, v=real volume, s=spread.
func (api *MT5WebAPI) ChartGet(symbol string, from, to int64, data string) ([]*MTBar, MTRetCode) {
	if api.chartBase == nil {
		api.chartBase = newChartHandler(api.asyncConnect)
	}
	return api.chartBase.ChartGet(symbol, from, to, data)
}

// BookGet retrieves the market depth (DOM) for one or more symbols (comma-separated).
func (api *MT5WebAPI) BookGet(symbol string) (map[string][]MTBookItem, MTRetCode) {
	if api.chartBase == nil {
		api.chartBase = newChartHandler(api.asyncConnect)
	}
	return api.chartBase.BookGet(symbol)
}

// DailyGet retrieves daily reports for a login in a date range.
func (api *MT5WebAPI) DailyGet(login uint64, from, to int64) ([]*MTDailyReport, MTRetCode) {
	if api.dailyBase == nil {
		api.dailyBase = newDailyHandler(api.asyncConnect)
	}
	return api.dailyBase.DailyGet(login, from, to)
}

// DailyGetBatch retrieves daily reports for multiple logins in a date range.
func (api *MT5WebAPI) DailyGetBatch(logins string, from, to int64) ([]*MTDailyReport, MTRetCode) {
	if api.dailyBase == nil {
		api.dailyBase = newDailyHandler(api.asyncConnect)
	}
	return api.dailyBase.DailyGetBatch(logins, from, to)
}

// SettingGet gets a settings file from the server.
// section is the directory, key is the filename.
func (api *MT5WebAPI) SettingGet(section, key string) (map[string]interface{}, MTRetCode) {
	if api.settingsBase == nil {
		api.settingsBase = newSettingsHandler(api.asyncConnect)
	}
	return api.settingsBase.SettingGet(section, key)
}

// SettingSet writes a settings file to the server.
// section is the directory, key is the filename, contentJSON is the JSON content.
func (api *MT5WebAPI) SettingSet(section, key, contentJSON string) MTRetCode {
	if api.settingsBase == nil {
		api.settingsBase = newSettingsHandler(api.asyncConnect)
	}
	return api.settingsBase.SettingSet(section, key, contentJSON)
}

// SettingDelete deletes a settings file from the server.
// section is the directory, key is the filename.
func (api *MT5WebAPI) SettingDelete(section, key string) MTRetCode {
	if api.settingsBase == nil {
		api.settingsBase = newSettingsHandler(api.asyncConnect)
	}
	return api.settingsBase.SettingDelete(section, key)
}

// SubscriptionCfgTotal gets the total count of subscription configs.
func (api *MT5WebAPI) SubscriptionCfgTotal() (int, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgTotal()
}

// SubscriptionCfgNext gets the subscription config at the given index.
func (api *MT5WebAPI) SubscriptionCfgNext(index int) (map[string]interface{}, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgNext(index)
}

// SubscriptionCfgGet gets a subscription config by ID.
func (api *MT5WebAPI) SubscriptionCfgGet(id int) (map[string]interface{}, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgGet(id)
}

// SubscriptionCfgAdd adds a new subscription config.
func (api *MT5WebAPI) SubscriptionCfgAdd(cfgJSON string) (map[string]interface{}, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgAdd(cfgJSON)
}

// SubscriptionCfgDelete deletes a subscription config by ID.
func (api *MT5WebAPI) SubscriptionCfgDelete(id int) MTRetCode {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgDelete(id)
}

// SubscriptionCfgShift shifts a subscription config position.
func (api *MT5WebAPI) SubscriptionCfgShift(index, shift int) MTRetCode {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCfgShift(index, shift)
}

// SubscriptionGet gets subscriptions for a user login.
func (api *MT5WebAPI) SubscriptionGet(login uint64) (map[string]interface{}, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionGet(login)
}

// SubscriptionExist checks whether a user is subscribed to a subscription config.
func (api *MT5WebAPI) SubscriptionExist(login uint64, id int) (bool, MTRetCode) {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionExist(login, id)
}

// SubscriptionJoin joins a user to a subscription config.
func (api *MT5WebAPI) SubscriptionJoin(login uint64, id int) MTRetCode {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionJoin(login, id)
}

// SubscriptionCancel cancels a user's subscription.
func (api *MT5WebAPI) SubscriptionCancel(login uint64, id int) MTRetCode {
	if api.subscriptionBase == nil {
		api.subscriptionBase = newSubscriptionHandler(api.asyncConnect)
	}
	return api.subscriptionBase.SubscriptionCancel(login, id)
}

// --- Leverage Config CRUD ---

// LeverageTotal gets the total count of leverage configs.
func (api *MT5WebAPI) LeverageTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageTotal()
}

// LeverageNext gets the leverage config at the given index.
func (api *MT5WebAPI) LeverageNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageNext(index)
}

// LeverageGet gets a leverage config by name.
func (api *MT5WebAPI) LeverageGet(name string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageGet(name)
}

// LeverageAdd adds a new leverage config.
func (api *MT5WebAPI) LeverageAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageAdd(bodyJSON)
}

// LeverageDelete deletes a leverage config by name.
func (api *MT5WebAPI) LeverageDelete(name string) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageDelete(name)
}

// LeverageShift shifts a leverage config position.
func (api *MT5WebAPI) LeverageShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.LeverageShift(index, shift)
}

// --- Firewall Config CRUD ---

// FirewallTotal gets the total count of firewall configs.
func (api *MT5WebAPI) FirewallTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallTotal()
}

// FirewallNext gets the firewall config at the given index.
func (api *MT5WebAPI) FirewallNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallNext(index)
}

// FirewallGetAll gets all firewall configs.
func (api *MT5WebAPI) FirewallGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallGetAll()
}

// FirewallAdd adds a new firewall config.
func (api *MT5WebAPI) FirewallAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallAdd(bodyJSON)
}

// FirewallDelete deletes a firewall config by index.
func (api *MT5WebAPI) FirewallDelete(index int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallDelete(index)
}

// FirewallShift shifts a firewall config position.
func (api *MT5WebAPI) FirewallShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.FirewallShift(index, shift)
}

// --- Holiday Config CRUD ---

// HolidayTotal gets the total count of holiday configs.
func (api *MT5WebAPI) HolidayTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayTotal()
}

// HolidayNext gets the holiday config at the given index.
func (api *MT5WebAPI) HolidayNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayNext(index)
}

// HolidayGetAll gets all holiday configs.
func (api *MT5WebAPI) HolidayGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayGetAll()
}

// HolidayAdd adds a new holiday config.
func (api *MT5WebAPI) HolidayAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayAdd(bodyJSON)
}

// HolidayDelete deletes a holiday config by index.
func (api *MT5WebAPI) HolidayDelete(index int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayDelete(index)
}

// HolidayShift shifts a holiday config position.
func (api *MT5WebAPI) HolidayShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HolidayShift(index, shift)
}

// --- Manager Config CRUD ---

// ManagerTotal gets the total count of manager configs.
func (api *MT5WebAPI) ManagerTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerTotal()
}

// ManagerNext gets the manager config at the given index.
func (api *MT5WebAPI) ManagerNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerNext(index)
}

// ManagerGet gets a manager config by login.
func (api *MT5WebAPI) ManagerGet(login uint64) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerGet(login)
}

// ManagerGetAll gets all manager configs.
func (api *MT5WebAPI) ManagerGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerGetAll()
}

// ManagerAdd adds a new manager config.
func (api *MT5WebAPI) ManagerAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerAdd(bodyJSON)
}

// ManagerDelete deletes a manager config by login.
func (api *MT5WebAPI) ManagerDelete(login uint64) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerDelete(login)
}

// ManagerShift shifts a manager config position.
func (api *MT5WebAPI) ManagerShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.ManagerShift(index, shift)
}

// --- Route Config CRUD ---

// RouteTotal gets the total count of route configs.
func (api *MT5WebAPI) RouteTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteTotal()
}

// RouteNext gets the route config at the given index.
func (api *MT5WebAPI) RouteNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteNext(index)
}

// RouteGet gets a route config by name.
func (api *MT5WebAPI) RouteGet(name string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteGet(name)
}

// RouteGetAll gets all route configs.
func (api *MT5WebAPI) RouteGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteGetAll()
}

// RouteAdd adds a new route config.
func (api *MT5WebAPI) RouteAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteAdd(bodyJSON)
}

// RouteDelete deletes a route config by name.
func (api *MT5WebAPI) RouteDelete(name string) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteDelete(name)
}

// RouteShift shifts a route config position.
func (api *MT5WebAPI) RouteShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.RouteShift(index, shift)
}

// --- Spread Config CRUD ---

// SpreadTotal gets the total count of spread configs.
func (api *MT5WebAPI) SpreadTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadTotal()
}

// SpreadNext gets the spread config at the given index.
func (api *MT5WebAPI) SpreadNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadNext(index)
}

// SpreadGetAll gets all spread configs.
func (api *MT5WebAPI) SpreadGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadGetAll()
}

// SpreadAdd adds a new spread config.
func (api *MT5WebAPI) SpreadAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadAdd(bodyJSON)
}

// SpreadDelete deletes a spread config by index.
func (api *MT5WebAPI) SpreadDelete(index int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadDelete(index)
}

// SpreadShift shifts a spread config position.
func (api *MT5WebAPI) SpreadShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SpreadShift(index, shift)
}

// --- HistorySync Config CRUD ---

// HistorySyncTotal gets the total count of history sync configs.
func (api *MT5WebAPI) HistorySyncTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncTotal()
}

// HistorySyncNext gets the history sync config at the given index.
func (api *MT5WebAPI) HistorySyncNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncNext(index)
}

// HistorySyncGetAll gets all history sync configs.
func (api *MT5WebAPI) HistorySyncGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncGetAll()
}

// HistorySyncAdd adds a new history sync config.
func (api *MT5WebAPI) HistorySyncAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncAdd(bodyJSON)
}

// HistorySyncDelete deletes a history sync config by index.
func (api *MT5WebAPI) HistorySyncDelete(index int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncDelete(index)
}

// HistorySyncShift shifts a history sync config position.
func (api *MT5WebAPI) HistorySyncShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncShift(index, shift)
}

// HistorySyncStart triggers a history synchronization start.
func (api *MT5WebAPI) HistorySyncStart() MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.HistorySyncStart()
}

// --- SymbolGroup Config CRUD ---

// SymbolGroupTotal gets the total count of symbol group configs.
func (api *MT5WebAPI) SymbolGroupTotal() (int, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupTotal()
}

// SymbolGroupNext gets the symbol group config at the given index.
func (api *MT5WebAPI) SymbolGroupNext(index int) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupNext(index)
}

// SymbolGroupList retrieves a list of symbol group names.
func (api *MT5WebAPI) SymbolGroupList() ([]string, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupList()
}

// SymbolGroupAdd adds a new symbol group config.
func (api *MT5WebAPI) SymbolGroupAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupAdd(bodyJSON)
}

// SymbolGroupDelete deletes a symbol group config by name.
func (api *MT5WebAPI) SymbolGroupDelete(name string) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupDelete(name)
}

// SymbolGroupShift shifts a symbol group config position.
func (api *MT5WebAPI) SymbolGroupShift(index, shift int) MTRetCode {
	if api.configCrudBase == nil {
		api.configCrudBase = newConfigCrudHandler(api.asyncConnect)
	}
	return api.configCrudBase.SymbolGroupShift(index, shift)
}

// ServerRestart restarts the MT5 server.
func (api *MT5WebAPI) ServerRestart() MTRetCode {
	handler := newServerHandler(api.asyncConnect)
	return handler.ServerRestart()
}

// ---------------------------------------------------------------------------
// Email infrastructure (8 methods)
// ---------------------------------------------------------------------------

// EmailTotal gets the total count of email configurations.
func (api *MT5WebAPI) EmailTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailTotal()
}

// EmailNext gets the email configuration at the given index.
func (api *MT5WebAPI) EmailNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailNext(index)
}

// EmailGet gets an email configuration by index.
func (api *MT5WebAPI) EmailGet(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailGet(index)
}

// EmailGetAll gets all email configurations.
func (api *MT5WebAPI) EmailGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailGetAll()
}

// EmailAdd adds a new email configuration from JSON.
func (api *MT5WebAPI) EmailAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailAdd(bodyJSON)
}

// EmailDelete deletes an email configuration by index.
func (api *MT5WebAPI) EmailDelete(index int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailDelete(index)
}

// EmailShift shifts an email configuration position.
func (api *MT5WebAPI) EmailShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailShift(index, shift)
}

// EmailSend sends an email using the given JSON body.
func (api *MT5WebAPI) EmailSend(bodyJSON string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.EmailSend(bodyJSON)
}

// ---------------------------------------------------------------------------
// Messenger infrastructure (8 methods)
// ---------------------------------------------------------------------------

// MessengerTotal gets the total count of messenger configurations.
func (api *MT5WebAPI) MessengerTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerTotal()
}

// MessengerNext gets the messenger configuration at the given index.
func (api *MT5WebAPI) MessengerNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerNext(index)
}

// MessengerGet gets a messenger configuration by index.
func (api *MT5WebAPI) MessengerGet(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerGet(index)
}

// MessengerGetAll gets all messenger configurations.
func (api *MT5WebAPI) MessengerGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerGetAll()
}

// MessengerAdd adds a new messenger configuration from JSON.
func (api *MT5WebAPI) MessengerAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerAdd(bodyJSON)
}

// MessengerDelete deletes a messenger configuration by index.
func (api *MT5WebAPI) MessengerDelete(index int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerDelete(index)
}

// MessengerShift shifts a messenger configuration position.
func (api *MT5WebAPI) MessengerShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerShift(index, shift)
}

// MessengerSend sends a messenger message using the given JSON body.
func (api *MT5WebAPI) MessengerSend(bodyJSON string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.MessengerSend(bodyJSON)
}

// ---------------------------------------------------------------------------
// Gateway infrastructure (12 methods)
// ---------------------------------------------------------------------------

// GatewayTotal gets the total count of gateway configurations.
func (api *MT5WebAPI) GatewayTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayTotal()
}

// GatewayNext gets the gateway configuration at the given index.
func (api *MT5WebAPI) GatewayNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayNext(index)
}

// GatewayGet gets a gateway configuration by name.
func (api *MT5WebAPI) GatewayGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayGet(name)
}

// GatewayGetAll gets all gateway configurations.
func (api *MT5WebAPI) GatewayGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayGetAll()
}

// GatewayAdd adds a new gateway configuration from JSON.
func (api *MT5WebAPI) GatewayAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayAdd(bodyJSON)
}

// GatewayDelete deletes a gateway configuration by name.
func (api *MT5WebAPI) GatewayDelete(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayDelete(name)
}

// GatewayShift shifts a gateway configuration position.
func (api *MT5WebAPI) GatewayShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayShift(index, shift)
}

// GatewayRestart restarts a gateway by name.
func (api *MT5WebAPI) GatewayRestart(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayRestart(name)
}

// GatewayGetPosition gets position data for a gateway by name.
func (api *MT5WebAPI) GatewayGetPosition(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayGetPosition(name)
}

// GatewayModuleTotal gets the total count of gateway modules.
func (api *MT5WebAPI) GatewayModuleTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayModuleTotal()
}

// GatewayModuleNext gets the gateway module at the given index.
func (api *MT5WebAPI) GatewayModuleNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayModuleNext(index)
}

// GatewayModuleGet gets a gateway module by name.
func (api *MT5WebAPI) GatewayModuleGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.GatewayModuleGet(name)
}

// ---------------------------------------------------------------------------
// Feeder infrastructure (11 methods)
// ---------------------------------------------------------------------------

// FeederTotal gets the total count of feeder configurations.
func (api *MT5WebAPI) FeederTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederTotal()
}

// FeederNext gets the feeder configuration at the given index.
func (api *MT5WebAPI) FeederNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederNext(index)
}

// FeederGet gets a feeder configuration by name.
func (api *MT5WebAPI) FeederGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederGet(name)
}

// FeederGetAll gets all feeder configurations.
func (api *MT5WebAPI) FeederGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederGetAll()
}

// FeederAdd adds a new feeder configuration from JSON.
func (api *MT5WebAPI) FeederAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederAdd(bodyJSON)
}

// FeederDelete deletes a feeder configuration by name.
func (api *MT5WebAPI) FeederDelete(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederDelete(name)
}

// FeederShift shifts a feeder configuration position.
func (api *MT5WebAPI) FeederShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederShift(index, shift)
}

// FeederRestart restarts a feeder by name.
func (api *MT5WebAPI) FeederRestart(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederRestart(name)
}

// FeederModuleTotal gets the total count of feeder modules.
func (api *MT5WebAPI) FeederModuleTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederModuleTotal()
}

// FeederModuleNext gets the feeder module at the given index.
func (api *MT5WebAPI) FeederModuleNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederModuleNext(index)
}

// FeederModuleGet gets a feeder module by name.
func (api *MT5WebAPI) FeederModuleGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.FeederModuleGet(name)
}

// ---------------------------------------------------------------------------
// Report infrastructure (10 methods)
// ---------------------------------------------------------------------------

// ReportTotal gets the total count of report configurations.
func (api *MT5WebAPI) ReportTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportTotal()
}

// ReportNext gets the report configuration at the given index.
func (api *MT5WebAPI) ReportNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportNext(index)
}

// ReportGet gets a report configuration by name.
func (api *MT5WebAPI) ReportGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportGet(name)
}

// ReportGetAll gets all report configurations.
func (api *MT5WebAPI) ReportGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportGetAll()
}

// ReportAdd adds a new report configuration from JSON.
func (api *MT5WebAPI) ReportAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportAdd(bodyJSON)
}

// ReportDelete deletes a report configuration by name.
func (api *MT5WebAPI) ReportDelete(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportDelete(name)
}

// ReportShift shifts a report configuration position.
func (api *MT5WebAPI) ReportShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportShift(index, shift)
}

// ReportModuleTotal gets the total count of report modules.
func (api *MT5WebAPI) ReportModuleTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportModuleTotal()
}

// ReportModuleNext gets the report module at the given index.
func (api *MT5WebAPI) ReportModuleNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportModuleNext(index)
}

// ReportModuleGet gets a report module by name.
func (api *MT5WebAPI) ReportModuleGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.ReportModuleGet(name)
}

// ---------------------------------------------------------------------------
// Plugin infrastructure (10 methods)
// ---------------------------------------------------------------------------

// PluginTotal gets the total count of plugin configurations.
func (api *MT5WebAPI) PluginTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginTotal()
}

// PluginNext gets the plugin configuration at the given index.
func (api *MT5WebAPI) PluginNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginNext(index)
}

// PluginGet gets a plugin configuration by name.
func (api *MT5WebAPI) PluginGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginGet(name)
}

// PluginGetAll gets all plugin configurations.
func (api *MT5WebAPI) PluginGetAll() ([]map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginGetAll()
}

// PluginAdd adds a new plugin configuration from JSON.
func (api *MT5WebAPI) PluginAdd(bodyJSON string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginAdd(bodyJSON)
}

// PluginDelete deletes a plugin configuration by name.
func (api *MT5WebAPI) PluginDelete(name string) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginDelete(name)
}

// PluginShift shifts a plugin configuration position.
func (api *MT5WebAPI) PluginShift(index, shift int) MTRetCode {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginShift(index, shift)
}

// PluginModuleTotal gets the total count of plugin modules.
func (api *MT5WebAPI) PluginModuleTotal() (int, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginModuleTotal()
}

// PluginModuleNext gets the plugin module at the given index.
func (api *MT5WebAPI) PluginModuleNext(index int) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginModuleNext(index)
}

// PluginModuleGet gets a plugin module by name.
func (api *MT5WebAPI) PluginModuleGet(name string) (map[string]interface{}, MTRetCode) {
	if api.infraBase == nil {
		api.infraBase = newInfraHandler(api.asyncConnect)
	}
	return api.infraBase.PluginModuleGet(name)
}

// ---------------------------------------------------------------------------
// Miscellaneous operations (misc_handler)
// ---------------------------------------------------------------------------

// --- Logger ---

// LoggerServerRequest retrieves server log entries.
func (api *MT5WebAPI) LoggerServerRequest(mode, logType string, from, to int64, filter string) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.LoggerServerRequest(mode, logType, from, to, filter)
}

// --- Notification ---

// NotificationSend sends a push notification to a user's mobile terminal.
func (api *MT5WebAPI) NotificationSend(metaquotesIDs, text string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.NotificationSend(metaquotesIDs, text)
}

// --- Position extras ---

// PositionUpdate updates a position on the server.
func (api *MT5WebAPI) PositionUpdate(positionJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionUpdate(positionJSON)
}

// PositionDelete deletes positions by comma-separated tickets.
func (api *MT5WebAPI) PositionDelete(tickets string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionDelete(tickets)
}

// PositionCheck checks a position (server-side validation).
func (api *MT5WebAPI) PositionCheck(positionJSON string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionCheck(positionJSON)
}

// PositionFix fixes a position (server-side fixation).
func (api *MT5WebAPI) PositionFix(positionJSON string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionFix(positionJSON)
}

// PositionBackupList gets position backup dates in a time range.
func (api *MT5WebAPI) PositionBackupList(from, to int64) ([]int64, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionBackupList(from, to)
}

// PositionBackupGet gets positions from a specific backup for a given login.
func (api *MT5WebAPI) PositionBackupGet(backupDate int64, login uint64) ([]*MTPosition, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionBackupGet(backupDate, login)
}

// PositionBackupRestore restores a position from backup.
func (api *MT5WebAPI) PositionBackupRestore(positionJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.PositionBackupRestore(positionJSON)
}

// --- Deal backup ---

// DealBackupList gets deal backup dates in a time range.
func (api *MT5WebAPI) DealBackupList(from, to int64) ([]int64, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DealBackupList(from, to)
}

// DealBackupGet gets deals from a specific backup for a given login.
func (api *MT5WebAPI) DealBackupGet(backupDate int64, login uint64) ([]*MTDeal, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DealBackupGet(backupDate, login)
}

// DealBackupRestore restores a deal from backup.
func (api *MT5WebAPI) DealBackupRestore(dealJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DealBackupRestore(dealJSON)
}

// --- Daily extras ---

// DailyGetPage retrieves paged daily reports for a login in a date range.
func (api *MT5WebAPI) DailyGetPage(login uint64, from, to int64, offset, total uint32) ([]*MTDailyReport, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DailyGetPage(login, from, to, offset, total)
}

// DailyGetLight retrieves a lightweight daily report for a login in a date range.
func (api *MT5WebAPI) DailyGetLight(login uint64, from, to int64) ([]*MTDailyReport, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DailyGetLight(login, from, to)
}

// DailyGetLightBatch retrieves lightweight daily reports for multiple logins (comma-separated).
func (api *MT5WebAPI) DailyGetLightBatch(logins string, from, to int64) ([]*MTDailyReport, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DailyGetLightBatch(logins, from, to)
}

// DailyGetLightPage retrieves lightweight paged daily reports.
func (api *MT5WebAPI) DailyGetLightPage(login uint64, from, to int64, offset, total uint32) ([]*MTDailyReport, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DailyGetLightPage(login, from, to, offset, total)
}

// --- User extras ---

// UserTotal gets the total number of users on the server.
func (api *MT5WebAPI) UserTotal() (int, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserTotal()
}

// UserGroup gets the group name of a user by login.
func (api *MT5WebAPI) UserGroup(login uint64) (string, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserGroup(login)
}

// UserGetExternal gets a user by external ID.
func (api *MT5WebAPI) UserGetExternal(externalID string) (*MTUser, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserGetExternal(externalID)
}

// UserSyncExternal synchronizes a user by external ID.
func (api *MT5WebAPI) UserSyncExternal(userJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserSyncExternal(userJSON)
}

// UserCheckBalance checks the balance for a user.
func (api *MT5WebAPI) UserCheckBalance(login uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserCheckBalance(login)
}

// UserRestore restores a user from backup.
func (api *MT5WebAPI) UserRestore(userJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserRestore(userJSON)
}

// UserArchiveAdd archives a user.
func (api *MT5WebAPI) UserArchiveAdd(login uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserArchiveAdd(login)
}

// UserArchiveGet gets an archived user by login.
func (api *MT5WebAPI) UserArchiveGet(login uint64) (*MTUser, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserArchiveGet(login)
}

// UserArchiveGetBatch gets archived users by comma-separated logins.
func (api *MT5WebAPI) UserArchiveGetBatch(logins string) ([]*MTUser, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserArchiveGetBatch(logins)
}

// UserAccountGetBatch gets account information for multiple logins (comma-separated).
func (api *MT5WebAPI) UserAccountGetBatch(logins string) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserAccountGetBatch(logins)
}

// UserCertificateGet gets a user's certificate.
func (api *MT5WebAPI) UserCertificateGet(login uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserCertificateGet(login)
}

// UserCertificateUpdate updates a user's certificate.
func (api *MT5WebAPI) UserCertificateUpdate(login uint64, certJSON string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserCertificateUpdate(login, certJSON)
}

// UserCertificateDelete deletes a user's certificate.
func (api *MT5WebAPI) UserCertificateDelete(login uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserCertificateDelete(login)
}

// UserCertificateConfirm confirms a user's certificate.
func (api *MT5WebAPI) UserCertificateConfirm(login uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserCertificateConfirm(login)
}

// UserOTPSecretGet gets a user's OTP secret.
func (api *MT5WebAPI) UserOTPSecretGet(login uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserOTPSecretGet(login)
}

// UserOTPSecretSet sets a user's OTP secret.
func (api *MT5WebAPI) UserOTPSecretSet(login uint64, secret string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.UserOTPSecretSet(login, secret)
}

// --- Subscription extras ---

// MiscSubscriptionAdd adds a subscription for a user.
func (api *MT5WebAPI) MiscSubscriptionAdd(subscriptionJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionAdd(subscriptionJSON)
}

// MiscSubscriptionUpdate updates a subscription.
func (api *MT5WebAPI) MiscSubscriptionUpdate(subscriptionJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionUpdate(subscriptionJSON)
}

// MiscSubscriptionDelete deletes a subscription.
func (api *MT5WebAPI) MiscSubscriptionDelete(login uint64, id int) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionDelete(login, id)
}

// SubscriptionHistoryAdd adds a subscription history record.
func (api *MT5WebAPI) SubscriptionHistoryAdd(historyJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionHistoryAdd(historyJSON)
}

// SubscriptionHistoryUpdate updates a subscription history record.
func (api *MT5WebAPI) SubscriptionHistoryUpdate(historyJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionHistoryUpdate(historyJSON)
}

// SubscriptionHistoryDelete deletes a subscription history record.
func (api *MT5WebAPI) SubscriptionHistoryDelete(login uint64, id int) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionHistoryDelete(login, id)
}

// SubscriptionHistoryGet gets subscription history for a user.
func (api *MT5WebAPI) SubscriptionHistoryGet(login uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.SubscriptionHistoryGet(login)
}

// --- Mail extras ---

// MailGet gets mail messages for a login.
func (api *MT5WebAPI) MailGet(login uint64, offset, total uint32) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.MailGet(login, offset, total)
}

// MailGetBody gets the body of a mail message by ticket.
func (api *MT5WebAPI) MailGetBody(ticket uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.MailGetBody(ticket)
}

// --- News extras ---

// NewsGet gets news messages. id is optional (pass "" for all), total is max count.
func (api *MT5WebAPI) NewsGet(id string, total uint32) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.NewsGet(id, total)
}

// NewsGetBody gets the body of a news message by ticket.
func (api *MT5WebAPI) NewsGetBody(ticket uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.NewsGetBody(ticket)
}

// --- Attachment ---

// AttachmentAdd adds an attachment.
func (api *MT5WebAPI) AttachmentAdd(attachmentJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.AttachmentAdd(attachmentJSON)
}

// AttachmentAttach attaches an attachment to an entity.
func (api *MT5WebAPI) AttachmentAttach(attachJSON string) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.AttachmentAttach(attachJSON)
}

// AttachmentGet gets an attachment by ID.
func (api *MT5WebAPI) AttachmentGet(id uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.AttachmentGet(id)
}

// --- Comment ---

// CommentAdd adds a comment.
func (api *MT5WebAPI) CommentAdd(commentJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.CommentAdd(commentJSON)
}

// CommentUpdate updates a comment.
func (api *MT5WebAPI) CommentUpdate(commentJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.CommentUpdate(commentJSON)
}

// CommentDelete deletes a comment by ID.
func (api *MT5WebAPI) CommentDelete(id uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.CommentDelete(id)
}

// CommentGet gets comments by entity ID.
func (api *MT5WebAPI) CommentGet(id uint64) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.CommentGet(id)
}

// --- Document ---

// DocumentAdd adds a document.
func (api *MT5WebAPI) DocumentAdd(documentJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DocumentAdd(documentJSON)
}

// DocumentUpdate updates a document.
func (api *MT5WebAPI) DocumentUpdate(documentJSON string) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DocumentUpdate(documentJSON)
}

// DocumentDelete deletes a document by ID.
func (api *MT5WebAPI) DocumentDelete(id uint64) MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DocumentDelete(id)
}

// DocumentGet gets a document by ID.
func (api *MT5WebAPI) DocumentGet(id uint64) (map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DocumentGet(id)
}

// DocumentHistoryGet gets document history by ID.
func (api *MT5WebAPI) DocumentHistoryGet(id uint64) ([]map[string]interface{}, MTRetCode) {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.DocumentHistoryGet(id)
}

// --- Quit ---

// Quit sends a QUIT command to disconnect the session.
func (api *MT5WebAPI) Quit() MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.Quit()
}

// --- Test ---

// TestAccess tests server access/connectivity.
func (api *MT5WebAPI) TestAccess() MTRetCode {
	if api.miscBase == nil {
		api.miscBase = newMiscHandler(api.asyncConnect)
	}
	return api.miscBase.TestAccess()
}

// changeConnection updates the connection for all handler objects.
func (api *MT5WebAPI) changeConnection() {
	if api.timeBase != nil {
		api.timeBase.setConnection(api.asyncConnect)
	}
	if api.commonBase != nil {
		api.commonBase.setConnection(api.asyncConnect)
	}
	if api.groupBase != nil {
		api.groupBase.setConnection(api.asyncConnect)
	}
	if api.symbolBase != nil {
		api.symbolBase.setConnection(api.asyncConnect)
	}
	if api.orderBase != nil {
		api.orderBase.setConnection(api.asyncConnect)
	}
	if api.positionBase != nil {
		api.positionBase.setConnection(api.asyncConnect)
	}
	if api.dealBase != nil {
		api.dealBase.setConnection(api.asyncConnect)
	}
	if api.historyBase != nil {
		api.historyBase.setConnection(api.asyncConnect)
	}
	if api.userBase != nil {
		api.userBase.setConnection(api.asyncConnect)
	}
	if api.mailBase != nil {
		api.mailBase.setConnection(api.asyncConnect)
	}
	if api.newsBase != nil {
		api.newsBase.setConnection(api.asyncConnect)
	}
	if api.pingBase != nil {
		api.pingBase.setConnection(api.asyncConnect)
	}
	if api.customBase != nil {
		api.customBase.setConnection(api.asyncConnect)
	}
	if api.tradeBase != nil {
		api.tradeBase.setConnection(api.asyncConnect)
	}
	if api.dealerBase != nil {
		api.dealerBase.setConnection(api.asyncConnect)
	}
	if api.serverConfigBase != nil {
		api.serverConfigBase.setConnection(api.asyncConnect)
	}
	if api.batchBase != nil {
		api.batchBase.setConnection(api.asyncConnect)
	}
	if api.backupBase != nil {
		api.backupBase.setConnection(api.asyncConnect)
	}
	if api.tradeCalcBase != nil {
		api.tradeCalcBase.setConnection(api.asyncConnect)
	}
	if api.clientBase != nil {
		api.clientBase.setConnection(api.asyncConnect)
	}
	if api.chartBase != nil {
		api.chartBase.setConnection(api.asyncConnect)
	}
	if api.dailyBase != nil {
		api.dailyBase.setConnection(api.asyncConnect)
	}
	if api.settingsBase != nil {
		api.settingsBase.setConnection(api.asyncConnect)
	}
	if api.subscriptionBase != nil {
		api.subscriptionBase.setConnection(api.asyncConnect)
	}
	if api.configCrudBase != nil {
		api.configCrudBase.setConnection(api.asyncConnect)
	}
	if api.infraBase != nil {
		api.infraBase.setConnection(api.asyncConnect)
	}
	if api.miscBase != nil {
		api.miscBase.setConnection(api.asyncConnect)
	}
}

// String returns the API version information.
func (api *MT5WebAPI) String() string {
	return fmt.Sprintf("MT5 Web API v%d (%s)", WEB_API_VERSION, WEB_API_DATE)
}
