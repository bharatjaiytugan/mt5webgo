package mt5webapi

// EnGroupTradeFlags - trade rights flags (group-level)
type EnGroupTradeFlags uint64

const (
	TradeFlagsNone             EnGroupTradeFlags = 0
	TradeFlagsSwaps            EnGroupTradeFlags = 1
	TradeFlagsTrailing         EnGroupTradeFlags = 2
	TradeFlagsExperts          EnGroupTradeFlags = 4
	TradeFlagsExpiration       EnGroupTradeFlags = 8
	TradeFlagsSignalsAll       EnGroupTradeFlags = 16
	TradeFlagsSignalsOwn       EnGroupTradeFlags = 32
	TradeFlagsSOCompensation   EnGroupTradeFlags = 64
	TradeFlagsDefault          EnGroupTradeFlags = TradeFlagsSwaps | TradeFlagsTrailing | TradeFlagsExperts | TradeFlagsExpiration | TradeFlagsSignalsAll
	TradeFlagsAll              EnGroupTradeFlags = TradeFlagsSwaps | TradeFlagsTrailing | TradeFlagsExperts | TradeFlagsExpiration | TradeFlagsSignalsAll | TradeFlagsSignalsOwn | TradeFlagsSOCompensation
)

// EnGroupPermissionsFlags - group permissions flags
type EnGroupPermissionsFlags uint64

const (
	GroupPermissionNone              EnGroupPermissionsFlags = 0
	GroupPermissionCertConfirm       EnGroupPermissionsFlags = 1
	GroupPermissionEnableConnection  EnGroupPermissionsFlags = 2
	GroupPermissionResetPassword     EnGroupPermissionsFlags = 4
	GroupPermissionForcedOTPUsage    EnGroupPermissionsFlags = 8
	GroupPermissionRiskWarning       EnGroupPermissionsFlags = 16
	GroupPermissionRegulationProtect EnGroupPermissionsFlags = 32
	GroupPermissionAll               EnGroupPermissionsFlags = GroupPermissionCertConfirm | GroupPermissionEnableConnection | GroupPermissionResetPassword | GroupPermissionForcedOTPUsage | GroupPermissionRiskWarning | GroupPermissionRegulationProtect
)

// EnAuthMode - authorization mode
type EnAuthMode uint32

const (
	AuthStandard  EnAuthMode = 0
	AuthRSA1024   EnAuthMode = 1
	AuthRSA2048   EnAuthMode = 2
	AuthRSACustom EnAuthMode = 3
	AuthFirst     EnAuthMode = AuthStandard
	AuthLast      EnAuthMode = AuthRSACustom
)

// EnAuthOTPMode - One-Time-Password mode
type EnAuthOTPMode uint32

const (
	AuthOTPDisabled      EnAuthOTPMode = 0
	AuthOTPTotpSHA256    EnAuthOTPMode = 1
	AuthOTPTotpSHA256Web EnAuthOTPMode = 2
	AuthOTPFirst         EnAuthOTPMode = AuthOTPDisabled
	AuthOTPLast          EnAuthOTPMode = AuthOTPTotpSHA256Web
)

// EnReportsMode - reports generation mode
type EnReportsMode uint32

const (
	ReportsDisabled EnReportsMode = 0
	ReportsStandard EnReportsMode = 1
	ReportsFirst    EnReportsMode = ReportsDisabled
	ReportsLast     EnReportsMode = ReportsStandard
)

// EnReportsFlags - reports generation flags
type EnReportsFlags uint32

const (
	ReportsFlagsNone       EnReportsFlags = 0
	ReportsFlagsEmail      EnReportsFlags = 1
	ReportsFlagsSupport    EnReportsFlags = 2
	ReportsFlagsStatements EnReportsFlags = 4
	ReportsFlagsAll        EnReportsFlags = ReportsFlagsEmail | ReportsFlagsSupport | ReportsFlagsStatements
)

// EnNewsMode - news modes
type EnNewsMode uint32

const (
	NewsModeDisabled EnNewsMode = 0
	NewsModeHeaders  EnNewsMode = 1
	NewsModeFull     EnNewsMode = 2
	NewsModeFirst    EnNewsMode = NewsModeDisabled
	NewsModeLast     EnNewsMode = NewsModeFull
)

// EnMailMode - internal email modes
type EnMailMode uint32

const (
	MailModeDisabled EnMailMode = 0
	MailModeFull     EnMailMode = 1
	MailModeFirst    EnMailMode = MailModeDisabled
	MailModeLast     EnMailMode = MailModeFull
)

// EnHistoryLimit - client history limits
type EnHistoryLimit uint32

const (
	TradeHistoryAll      EnHistoryLimit = 0
	TradeHistoryMonths1  EnHistoryLimit = 1
	TradeHistoryMonths3  EnHistoryLimit = 2
	TradeHistoryMonths6  EnHistoryLimit = 3
	TradeHistoryYear1    EnHistoryLimit = 4
	TradeHistoryYear2    EnHistoryLimit = 5
	TradeHistoryYear3    EnHistoryLimit = 6
	TradeHistoryFirst    EnHistoryLimit = TradeHistoryAll
	TradeHistoryLast     EnHistoryLimit = TradeHistoryYear3
)

// EnFreeMarginMode - free margin calculation modes
type EnFreeMarginMode uint32

const (
	FreeMarginNotUsePL EnFreeMarginMode = 0
	FreeMarginUsePL    EnFreeMarginMode = 1
	FreeMarginProfit   EnFreeMarginMode = 2
	FreeMarginLoss     EnFreeMarginMode = 3
	FreeMarginFirst    EnFreeMarginMode = FreeMarginNotUsePL
	FreeMarginLast     EnFreeMarginMode = FreeMarginLoss
)

// EnTransferMode - deposit transfer mode
type EnTransferMode uint32

const (
	TransferModeDisabled  EnTransferMode = 0
	TransferModeName      EnTransferMode = 1
	TransferModeGroup     EnTransferMode = 2
	TransferModeNameGroup EnTransferMode = 3
	TransferModeFirst     EnTransferMode = TransferModeDisabled
	TransferModeLast      EnTransferMode = TransferModeNameGroup
)

// EnStopOutMode - stop-out mode
type EnStopOutMode uint32

const (
	StopOutPercent EnStopOutMode = 0
	StopOutMoney   EnStopOutMode = 1
	StopOutFirst   EnStopOutMode = StopOutPercent
	StopOutLast    EnStopOutMode = StopOutMoney
)

// EnMarginFreeProfitMode - mode of calculation of the free margin of the fixed income
type EnMarginFreeProfitMode uint32

const (
	FreeMarginProfitPL   EnMarginFreeProfitMode = 0
	FreeMarginProfitLoss EnMarginFreeProfitMode = 1
	FreeMarginProfitFirst EnMarginFreeProfitMode = FreeMarginProfitPL
	FreeMarginProfitLast  EnMarginFreeProfitMode = FreeMarginProfitLoss
)

// EnGroupMarginMode - group risk management mode
type EnGroupMarginMode uint32

const (
	MarginModeRetail          EnGroupMarginMode = 0
	MarginModeExchangeDiscount EnGroupMarginMode = 1
	MarginModeRetailHedged    EnGroupMarginMode = 2
	MarginModeFirst           EnGroupMarginMode = MarginModeRetail
	MarginModeLast            EnGroupMarginMode = MarginModeRetailHedged
)

// EnGroupMarginFlags - margin calculation flags (group-level)
type EnGroupMarginFlags uint32

const (
	GroupMarginFlagsNone     EnGroupMarginFlags = 0
	GroupMarginFlagsClearAcc EnGroupMarginFlags = 1
	GroupMarginFlagsAll      EnGroupMarginFlags = GroupMarginFlagsClearAcc
)

// MTConGroup represents data config of a group.
type MTConGroup struct {
	Group               string                 `json:"Group"`
	Server              uint64                 `json:"Server"`
	PermissionsFlags    EnGroupPermissionsFlags `json:"PermissionsFlags"`
	AuthMode            EnAuthMode             `json:"AuthMode"`
	AuthPasswordMin     uint32                 `json:"AuthPasswordMin"`
	AuthOTPMode         EnAuthOTPMode          `json:"AuthOTPMode"`
	Company             string                 `json:"Company"`
	CompanyPage         string                 `json:"CompanyPage"`
	CompanyEmail        string                 `json:"CompanyEmail"`
	CompanySupportPage  string                 `json:"CompanySupportPage"`
	CompanySupportEmail string                 `json:"CompanySupportEmail"`
	CompanyCatalog      string                 `json:"CompanyCatalog"`
	Currency            string                 `json:"Currency"`
	CurrencyDigits      uint32                 `json:"CurrencyDigits"`
	ReportsMode         EnReportsMode          `json:"ReportsMode"`
	ReportsFlags        EnReportsFlags         `json:"ReportsFlags"`
	ReportsSMTP         string                 `json:"ReportsSMTP"`
	ReportsSMTPLogin    string                 `json:"ReportsSMTPLogin"`
	ReportsSMTPPass     string                 `json:"ReportsSMTPPass"`
	NewsMode            EnNewsMode             `json:"NewsMode"`
	NewsCategory        string                 `json:"NewsCategory"`
	NewsLangs           []uint32               `json:"NewsLangs"`
	MailMode            EnMailMode             `json:"MailMode"`
	TradeFlags          EnGroupTradeFlags      `json:"TradeFlags"`
	TradeTransferMode   EnTransferMode         `json:"TradeTransferMode"`
	TradeInterestrate   float64                `json:"TradeInterestrate"`
	TradeVirtualCredit  float64                `json:"TradeVirtualCredit"`
	MarginMode          EnGroupMarginMode      `json:"MarginMode"`
	MarginFreeMode      EnFreeMarginMode       `json:"MarginFreeMode"`
	MarginSOMode        EnStopOutMode          `json:"MarginSOMode"`
	MarginCall          float64                `json:"MarginCall"`
	MarginStopOut       float64                `json:"MarginStopOut"`
	MarginFreeProfitMode EnMarginFreeProfitMode `json:"MarginFreeProfitMode"`
	DemoLeverage        uint32                 `json:"DemoLeverage"`
	DemoDeposit         float64                `json:"DemoDeposit"`
	LimitHistory        EnHistoryLimit         `json:"LimitHistory"`
	LimitOrders         uint32                 `json:"LimitOrders"`
	LimitSymbols        uint32                 `json:"LimitSymbols"`
	LimitPositions      uint32                 `json:"LimitPositions"`
	Commissions         []MTConCommission      `json:"Commissions"`
	Symbols             []MTConGroupSymbol     `json:"Symbols"`
}

// CreateDefaultConGroup creates an MTConGroup with default values.
func CreateDefaultConGroup() *MTConGroup {
	g := &MTConGroup{}
	g.PermissionsFlags = GroupPermissionEnableConnection
	g.AuthMode = AuthStandard
	g.AuthPasswordMin = 7
	g.ReportsMode = ReportsDisabled
	g.ReportsFlags = ReportsFlagsNone
	g.Currency = "USD"
	g.CurrencyDigits = 2
	g.NewsMode = NewsModeFull
	g.MailMode = MailModeFull
	g.MarginFreeMode = FreeMarginUsePL
	g.MarginCall = 50
	g.MarginStopOut = 30
	g.MarginSOMode = StopOutPercent
	g.TradeVirtualCredit = 0
	g.MarginFreeProfitMode = FreeMarginProfitPL
	g.DemoLeverage = 0
	g.DemoDeposit = 0
	g.LimitSymbols = 0
	g.LimitOrders = 0
	g.LimitHistory = TradeHistoryAll
	g.TradeInterestrate = 0
	g.TradeFlags = TradeFlagsAll
	return g
}
