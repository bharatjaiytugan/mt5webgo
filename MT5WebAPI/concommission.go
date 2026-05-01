package mt5webapi

// EnCommMode - commission mode
type EnCommMode uint32

const (
	CommStandard EnCommMode = 0
	CommAgent    EnCommMode = 1
	CommFirst    EnCommMode = CommStandard
	CommLast     EnCommMode = CommAgent
)

// EnCommRangeMode - commission range mode
type EnCommRangeMode uint32

const (
	CommRangeVolume         EnCommRangeMode = 0
	CommRangeOverturnMoney  EnCommRangeMode = 1
	CommRangeOverturnVolume EnCommRangeMode = 2
	CommRangeFirst          EnCommRangeMode = CommRangeVolume
	CommRangeLast           EnCommRangeMode = CommRangeOverturnVolume
)

// EnCommChargeMode - commission charge modes
type EnCommChargeMode uint32

const (
	CommChargeDaily   EnCommChargeMode = 0
	CommChargeMonthly EnCommChargeMode = 1
	CommChargeInstant EnCommChargeMode = 2
	CommChargeFirst   EnCommChargeMode = CommChargeDaily
	CommChargeLast    EnCommChargeMode = CommChargeInstant
)

// EnCommEntryMode - deal entry mode
type EnCommEntryMode uint32

const (
	CommEntryAll   EnCommEntryMode = 0
	CommEntryIn    EnCommEntryMode = 1
	CommEntryOut   EnCommEntryMode = 2
	CommEntryFirst EnCommEntryMode = CommEntryAll
	CommEntryLast  EnCommEntryMode = CommEntryOut
)

// EnCommissionMode - commission charge mode (tier level)
type EnCommissionMode uint32

const (
	CommMoneyDeposit      EnCommissionMode = 0
	CommMoneySymbolBase   EnCommissionMode = 1
	CommMoneySymbolProfit EnCommissionMode = 2
	CommMoneySymbolMargin EnCommissionMode = 3
	CommPips              EnCommissionMode = 4
	CommPercent           EnCommissionMode = 5
	CommMoneySpecified    EnCommissionMode = 6
	CommTierModeFirst     EnCommissionMode = CommMoneyDeposit
	CommTierModeLast      EnCommissionMode = CommMoneySpecified
)

// EnCommissionVolumeType - commission type by volume
type EnCommissionVolumeType uint32

const (
	CommTypeDeal     EnCommissionVolumeType = 0
	CommTypeVolume   EnCommissionVolumeType = 1
	CommTypeFirst    EnCommissionVolumeType = CommTypeDeal
	CommTypeLast     EnCommissionVolumeType = CommTypeVolume
)

// MTConCommission represents commission config.
type MTConCommission struct {
	Name             string           `json:"Name"`
	Description      string           `json:"Description"`
	Path             string           `json:"Path"`
	Mode             EnCommMode       `json:"Mode"`
	RangeMode        EnCommRangeMode  `json:"RangeMode"`
	ChargeMode       EnCommChargeMode `json:"ChargeMode"`
	TurnoverCurrency string           `json:"TurnoverCurrency"`
	EntryMode        EnCommEntryMode  `json:"EntryMode"`
	Tiers            []MTConCommTier  `json:"Tiers"`
}

// MTConCommTier represents commission tier config.
type MTConCommTier struct {
	Mode      EnCommissionMode       `json:"Mode"`
	Type      EnCommissionVolumeType `json:"Type"`
	Value     float64                `json:"Value"`
	Minimal   float64                `json:"Minimal"`
	Maximal   float64                `json:"Maximal"`
	RangeFrom float64                `json:"RangeFrom"`
	RangeTo   float64                `json:"RangeTo"`
	Currency  string                 `json:"Currency"`
}
