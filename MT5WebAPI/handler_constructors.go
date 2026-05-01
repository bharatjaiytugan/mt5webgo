package mt5webapi

func newTimeHandler(connect *mtAsyncConnect) *timeHandler {
	return &timeHandler{apiBase: newAPIBase(connect)}
}

func newCommonHandler(connect *mtAsyncConnect) *commonHandler {
	return &commonHandler{apiBase: newAPIBase(connect)}
}

func newGroupHandler(connect *mtAsyncConnect) *groupHandler {
	return &groupHandler{apiBase: newAPIBase(connect)}
}

func newSymbolHandler(connect *mtAsyncConnect) *symbolHandler {
	return &symbolHandler{apiBase: newAPIBase(connect)}
}

func newOrderHandler(connect *mtAsyncConnect) *orderHandler {
	return &orderHandler{apiBase: newAPIBase(connect)}
}

func newPositionHandler(connect *mtAsyncConnect) *positionHandler {
	return &positionHandler{apiBase: newAPIBase(connect)}
}

func newDealHandler(connect *mtAsyncConnect) *dealHandler {
	return &dealHandler{apiBase: newAPIBase(connect)}
}

func newHistoryHandler(connect *mtAsyncConnect) *historyHandler {
	return &historyHandler{apiBase: newAPIBase(connect)}
}

func newTickHandler(connect *mtAsyncConnect) *tickHandler {
	return &tickHandler{apiBase: newAPIBase(connect)}
}

func newUserHandler(connect *mtAsyncConnect) *userHandler {
	return &userHandler{apiBase: newAPIBase(connect)}
}

func newMailHandler(connect *mtAsyncConnect) *mailHandler {
	return &mailHandler{apiBase: newAPIBase(connect)}
}

func newNewsHandler(connect *mtAsyncConnect) *newsHandler {
	return &newsHandler{apiBase: newAPIBase(connect)}
}

func newPingHandler(connect *mtAsyncConnect) *pingHandler {
	return &pingHandler{apiBase: newAPIBase(connect)}
}

func newCustomHandler(connect *mtAsyncConnect) *customHandler {
	return &customHandler{apiBase: newAPIBase(connect)}
}

func newTradeHandler(connect *mtAsyncConnect) *tradeHandler {
	return &tradeHandler{apiBase: newAPIBase(connect)}
}

func newServerHandler(connect *mtAsyncConnect) *serverHandler {
	return &serverHandler{apiBase: newAPIBase(connect)}
}

func newDealerHandler(connect *mtAsyncConnect) *dealerHandler {
	return &dealerHandler{apiBase: newAPIBase(connect)}
}

func newServerConfigHandler(connect *mtAsyncConnect) *serverConfigHandler {
	return &serverConfigHandler{apiBase: newAPIBase(connect)}
}

func newBatchHandler(connect *mtAsyncConnect) *batchHandler {
	return &batchHandler{apiBase: newAPIBase(connect)}
}

func newBackupHandler(connect *mtAsyncConnect) *backupHandler {
	return &backupHandler{apiBase: newAPIBase(connect)}
}

func newTradeCalcHandler(connect *mtAsyncConnect) *tradeCalcHandler {
	return &tradeCalcHandler{apiBase: newAPIBase(connect)}
}

func newClientHandler(connect *mtAsyncConnect) *clientHandler {
	return &clientHandler{apiBase: newAPIBase(connect)}
}

func newChartHandler(connect *mtAsyncConnect) *chartHandler {
	return &chartHandler{apiBase: newAPIBase(connect)}
}

func newDailyHandler(connect *mtAsyncConnect) *dailyHandler {
	return &dailyHandler{apiBase: newAPIBase(connect)}
}

func newSettingsHandler(connect *mtAsyncConnect) *settingsHandler {
	return &settingsHandler{apiBase: newAPIBase(connect)}
}

func newSubscriptionHandler(connect *mtAsyncConnect) *subscriptionHandler {
	return &subscriptionHandler{apiBase: newAPIBase(connect)}
}

func newConfigCrudHandler(connect *mtAsyncConnect) *configCrudHandler {
	return &configCrudHandler{apiBase: newAPIBase(connect)}
}

func newInfraHandler(connect *mtAsyncConnect) *infraHandler {
	return &infraHandler{apiBase: newAPIBase(connect)}
}

func newMiscHandler(connect *mtAsyncConnect) *miscHandler {
	return &miscHandler{apiBase: newAPIBase(connect)}
}
