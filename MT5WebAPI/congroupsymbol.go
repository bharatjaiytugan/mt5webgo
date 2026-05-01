package mt5webapi

import "math"

// Default value constants for MTConGroupSymbol
const (
	DefaultValueUint   uint32  = 0xffffffff
	DefaultValueInt    int32   = 0x7fffffff
	DefaultValueUint64 uint64  = 0xffffffffffffffff
	DefaultValueInt64  int64   = 0x7fffffffffffffff
	DefaultValueDouble float64 = math.MaxFloat64
)

// EnGroupSymbolREFlags - Requests Execution flags (group symbol level)
type EnGroupSymbolREFlags uint32

const (
	GSREFlagsNone  EnGroupSymbolREFlags = 0
	GSREFlagsOrder EnGroupSymbolREFlags = 1
	GSREFlagsAll   EnGroupSymbolREFlags = GSREFlagsOrder
)

// EnGroupSymbolPermissionsFlags - permissions flags (group symbol level)
type EnGroupSymbolPermissionsFlags uint32

const (
	GSPermissionNone    EnGroupSymbolPermissionsFlags = 0
	GSPermissionBook    EnGroupSymbolPermissionsFlags = 1
	GSPermissionDefault EnGroupSymbolPermissionsFlags = GSPermissionBook
	GSPermissionAll     EnGroupSymbolPermissionsFlags = GSPermissionBook
)

// MTConGroupSymbol represents symbols configuration for clients group.
type MTConGroupSymbol struct {
	Path                string                        `json:"Path"`
	TradeMode           EnTradeMode                   `json:"TradeMode"`
	ExecMode            EnCalcMode                    `json:"ExecMode"`
	FillFlags           EnFillingFlags                `json:"FillFlags"`
	ExpirFlags          EnExpirationFlags             `json:"ExpirFlags"`
	OrderFlags          EnOrderFlags                  `json:"OrderFlags"`
	SpreadDiff          int32                         `json:"SpreadDiff"`
	SpreadDiffBalance   int32                         `json:"SpreadDiffBalance"`
	StopsLevel          int32                         `json:"StopsLevel"`
	FreezeLevel         int32                         `json:"FreezeLevel"`
	VolumeMin           uint64                        `json:"VolumeMin"`
	VolumeMinExt        uint64                        `json:"VolumeMinExt"`
	VolumeMax           uint64                        `json:"VolumeMax"`
	VolumeMaxExt        uint64                        `json:"VolumeMaxExt"`
	VolumeStep          uint64                        `json:"VolumeStep"`
	VolumeStepExt       uint64                        `json:"VolumeStepExt"`
	VolumeLimit         uint64                        `json:"VolumeLimit"`
	VolumeLimitExt      uint64                        `json:"VolumeLimitExt"`
	MarginCheckMode     EnMarginFlags                 `json:"MarginCheckMode"` // deprecated
	MarginFlags         EnMarginFlags                 `json:"MarginFlags"`
	MarginInitial       float64                       `json:"MarginInitial"`
	MarginMaintenance   float64                       `json:"MarginMaintenance"`
	MarginRateInitial   map[EnMarginRateTypes]float64 `json:"-"`
	MarginRateMaintenance map[EnMarginRateTypes]float64 `json:"-"`
	MarginRateLiquidity float64                       `json:"MarginLiquidity"`
	MarginHedged        float64                       `json:"MarginHedged"`
	MarginRateCurrency  float64                       `json:"MarginCurrency"`
	MarginLong          float64                       `json:"MarginLong"`      // deprecated
	MarginShort         float64                       `json:"MarginShort"`     // deprecated
	MarginLimit         float64                       `json:"MarginLimit"`     // deprecated
	MarginStop          float64                       `json:"MarginStop"`      // deprecated
	MarginStopLimit     float64                       `json:"MarginStopLimit"` // deprecated
	SwapMode            EnSwapMode                    `json:"SwapMode"`
	SwapLong            float64                       `json:"SwapLong"`
	SwapShort           float64                       `json:"SwapShort"`
	Swap3Day            int32                         `json:"Swap3Day"`
	REFlags             EnGroupSymbolREFlags          `json:"REFlags"`
	RETimeout           uint32                        `json:"RETimeout"`
	IEFlags             uint32                        `json:"IEFlags"`
	IECheckMode         uint32                        `json:"IECheckMode"`
	IETimeout           uint32                        `json:"IETimeout"`
	IESlipProfit        uint32                        `json:"IESlipProfit"`
	IESlipLosing        uint32                        `json:"IESlipLosing"`
	IEVolumeMax         uint64                        `json:"IEVolumeMax"`
	IEVolumeMaxExt      uint64                        `json:"IEVolumeMaxExt"`
	PermissionsFlags    EnGroupSymbolPermissionsFlags `json:"PermissionsFlags"`
	BookDepthLimit      uint32                        `json:"PermissionsBookdepth"`
	// Margin rate fields serialized as individual JSON keys
	MarginInitialBuy            float64 `json:"MarginInitialBuy"`
	MarginInitialSell           float64 `json:"MarginInitialSell"`
	MarginInitialBuyLimit       float64 `json:"MarginInitialBuyLimit"`
	MarginInitialSellLimit      float64 `json:"MarginInitialSellLimit"`
	MarginInitialBuyStop        float64 `json:"MarginInitialBuyStop"`
	MarginInitialSellStop       float64 `json:"MarginInitialSellStop"`
	MarginInitialBuyStopLimit   float64 `json:"MarginInitialBuyStopLimit"`
	MarginInitialSellStopLimit  float64 `json:"MarginInitialSellStopLimit"`
	MarginMaintenanceBuy        float64 `json:"MarginMaintenanceBuy"`
	MarginMaintenanceSell       float64 `json:"MarginMaintenanceSell"`
	MarginMaintenanceBuyLimit   float64 `json:"MarginMaintenanceBuyLimit"`
	MarginMaintenanceSellLimit  float64 `json:"MarginMaintenanceSellLimit"`
	MarginMaintenanceBuyStop    float64 `json:"MarginMaintenanceBuyStop"`
	MarginMaintenanceSellStop   float64 `json:"MarginMaintenanceSellStop"`
	MarginMaintenanceBuyStopLimit  float64 `json:"MarginMaintenanceBuyStopLimit"`
	MarginMaintenanceSellStopLimit float64 `json:"MarginMaintenanceSellStopLimit"`
}

// SetVolumeMin sets VolumeMin with default-aware conversion.
func (gs *MTConGroupSymbol) SetVolumeMin(v uint64) {
	if v == DefaultValueUint64 {
		gs.VolumeMinExt = v
	} else {
		gs.VolumeMinExt = v * 10000
	}
	gs.VolumeMin = v
}

// SetVolumeMax sets VolumeMax with default-aware conversion.
func (gs *MTConGroupSymbol) SetVolumeMax(v uint64) {
	if v == DefaultValueUint64 {
		gs.VolumeMaxExt = v
	} else {
		gs.VolumeMaxExt = v * 10000
	}
	gs.VolumeMax = v
}

// SetVolumeStep sets VolumeStep with default-aware conversion.
func (gs *MTConGroupSymbol) SetVolumeStep(v uint64) {
	if v == DefaultValueUint64 {
		gs.VolumeStepExt = v
	} else {
		gs.VolumeStepExt = v * 10000
	}
	gs.VolumeStep = v
}

// SetVolumeLimit sets VolumeLimit with default-aware conversion.
func (gs *MTConGroupSymbol) SetVolumeLimit(v uint64) {
	if v == DefaultValueUint64 {
		gs.VolumeLimitExt = v
	} else {
		gs.VolumeLimitExt = v * 10000
	}
	gs.VolumeLimit = v
}

// SetIEVolumeMax sets IEVolumeMax with default-aware conversion.
func (gs *MTConGroupSymbol) SetIEVolumeMax(v uint64) {
	if v == DefaultValueUint64 {
		gs.IEVolumeMaxExt = v
	} else {
		gs.IEVolumeMaxExt = v * 10000
	}
	gs.IEVolumeMax = v
}

// GetVolumeMin returns VolumeMinExt converted to old volume format, preserving default.
func (gs *MTConGroupSymbol) GetVolumeMin() uint64 {
	if gs.VolumeMinExt == DefaultValueUint64 {
		return DefaultValueUint64
	}
	return gs.VolumeMinExt / 10000
}

// GetVolumeMax returns VolumeMaxExt converted to old volume format, preserving default.
func (gs *MTConGroupSymbol) GetVolumeMax() uint64 {
	if gs.VolumeMaxExt == DefaultValueUint64 {
		return DefaultValueUint64
	}
	return gs.VolumeMaxExt / 10000
}

// GetVolumeStep returns VolumeStepExt converted to old volume format, preserving default.
func (gs *MTConGroupSymbol) GetVolumeStep() uint64 {
	if gs.VolumeStepExt == DefaultValueUint64 {
		return DefaultValueUint64
	}
	return gs.VolumeStepExt / 10000
}

// GetVolumeLimit returns VolumeLimitExt converted to old volume format, preserving default.
func (gs *MTConGroupSymbol) GetVolumeLimit() uint64 {
	if gs.VolumeLimitExt == DefaultValueUint64 {
		return DefaultValueUint64
	}
	return gs.VolumeLimitExt / 10000
}

// GetIEVolumeMax returns IEVolumeMaxExt converted to old volume format, preserving default.
func (gs *MTConGroupSymbol) GetIEVolumeMax() uint64 {
	if gs.IEVolumeMaxExt == DefaultValueUint64 {
		return DefaultValueUint64
	}
	return gs.IEVolumeMaxExt / 10000
}

// CreateDefaultConGroupSymbol creates an MTConGroupSymbol with default values.
func CreateDefaultConGroupSymbol() *MTConGroupSymbol {
	gs := &MTConGroupSymbol{}
	gs.TradeMode = EnTradeMode(DefaultValueUint)
	gs.ExecMode = EnCalcMode(DefaultValueUint)
	gs.FillFlags = EnFillingFlags(DefaultValueUint)
	gs.ExpirFlags = EnExpirationFlags(DefaultValueUint)
	gs.OrderFlags = OrderFlagsNone
	gs.SpreadDiff = DefaultValueInt
	gs.SpreadDiffBalance = DefaultValueInt
	gs.StopsLevel = DefaultValueInt
	gs.FreezeLevel = DefaultValueInt
	gs.VolumeMin = DefaultValueUint64
	gs.VolumeMinExt = DefaultValueUint64
	gs.VolumeMax = DefaultValueUint64
	gs.VolumeMaxExt = DefaultValueUint64
	gs.VolumeStep = DefaultValueUint64
	gs.VolumeStepExt = DefaultValueUint64
	gs.VolumeLimit = DefaultValueUint64
	gs.VolumeLimitExt = DefaultValueUint64
	gs.MarginFlags = EnMarginFlags(DefaultValueUint)
	gs.MarginCheckMode = EnMarginFlags(DefaultValueUint)
	gs.MarginInitial = DefaultValueDouble
	gs.MarginMaintenance = DefaultValueDouble
	gs.MarginRateInitial = make(map[EnMarginRateTypes]float64)
	gs.MarginRateMaintenance = make(map[EnMarginRateTypes]float64)
	gs.MarginRateLiquidity = DefaultValueDouble
	gs.MarginHedged = DefaultValueDouble
	gs.MarginRateCurrency = DefaultValueDouble
	// deprecated
	gs.MarginLong = DefaultValueDouble
	gs.MarginShort = DefaultValueDouble
	gs.MarginLimit = DefaultValueDouble
	gs.MarginStop = DefaultValueDouble
	gs.MarginStopLimit = DefaultValueDouble
	gs.SwapMode = EnSwapMode(DefaultValueUint)
	gs.SwapLong = DefaultValueDouble
	gs.SwapShort = DefaultValueDouble
	gs.Swap3Day = DefaultValueInt
	gs.REFlags = EnGroupSymbolREFlags(DefaultValueUint)
	gs.RETimeout = DefaultValueUint
	gs.IEFlags = DefaultValueUint
	gs.IECheckMode = DefaultValueUint
	gs.IETimeout = DefaultValueUint
	gs.IESlipProfit = DefaultValueUint
	gs.IESlipLosing = DefaultValueUint
	gs.IEVolumeMax = DefaultValueUint64
	gs.IEVolumeMaxExt = DefaultValueUint64
	gs.PermissionsFlags = GSPermissionDefault
	gs.BookDepthLimit = 0
	// Margin rate initial/maintenance fields default
	gs.MarginInitialBuy = DefaultValueDouble
	gs.MarginInitialSell = DefaultValueDouble
	gs.MarginInitialBuyLimit = DefaultValueDouble
	gs.MarginInitialSellLimit = DefaultValueDouble
	gs.MarginInitialBuyStop = DefaultValueDouble
	gs.MarginInitialSellStop = DefaultValueDouble
	gs.MarginInitialBuyStopLimit = DefaultValueDouble
	gs.MarginInitialSellStopLimit = DefaultValueDouble
	gs.MarginMaintenanceBuy = DefaultValueDouble
	gs.MarginMaintenanceSell = DefaultValueDouble
	gs.MarginMaintenanceBuyLimit = DefaultValueDouble
	gs.MarginMaintenanceSellLimit = DefaultValueDouble
	gs.MarginMaintenanceBuyStop = DefaultValueDouble
	gs.MarginMaintenanceSellStop = DefaultValueDouble
	gs.MarginMaintenanceBuyStopLimit = DefaultValueDouble
	gs.MarginMaintenanceSellStopLimit = DefaultValueDouble
	return gs
}
