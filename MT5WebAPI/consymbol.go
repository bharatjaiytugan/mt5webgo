package mt5webapi

// MTConSymbolSession represents symbol trade and quotes session config.
type MTConSymbolSession struct {
	Open       uint32 `json:"Open"`
	OpenHours  uint32 `json:"OpenHours"`
	Close      uint32 `json:"Close"`
	CloseHours uint32 `json:"CloseHours"`
}

// EnSectors - economical sectors
type EnSectors uint32

const (
	SectorUndefined            EnSectors = 0
	SectorBasicMaterials       EnSectors = 1
	SectorCommunicationService EnSectors = 2
	SectorConsumerCyclical     EnSectors = 3
	SectorConsumerDefensive    EnSectors = 4
	SectorEnergy               EnSectors = 5
	SectorFinancial            EnSectors = 6
	SectorHealthcare           EnSectors = 7
	SectorIndustrials          EnSectors = 8
	SectorRealEstate           EnSectors = 9
	SectorTechnology           EnSectors = 10
	SectorUtilities            EnSectors = 11
	SectorCurrency             EnSectors = 12
	SectorCurrencyCrypto       EnSectors = 13
	SectorIndexes              EnSectors = 14
	SectorCommodities          EnSectors = 15
	// enumeration borders
	SectorFirst EnSectors = SectorUndefined
	SectorLast  EnSectors = SectorCommodities
)

// EnIndustries - economical industries
type EnIndustries uint32

const (
	IndustryUndefined EnIndustries = 0
	// Basic Materials
	IndustryAgriculturalInputs EnIndustries = 1
	IndustryAluminium          EnIndustries = 2
	IndustryBuildingMaterials  EnIndustries = 3
	IndustryChemicals          EnIndustries = 4
	IndustryCokingCoal         EnIndustries = 5
	IndustryCopper             EnIndustries = 6
	IndustryGold               EnIndustries = 7
	IndustryLumberWood         EnIndustries = 8
	IndustryIndustrialMetals   EnIndustries = 9
	IndustryPreciousMetals     EnIndustries = 10
	IndustryPaper              EnIndustries = 11
	IndustrySilver             EnIndustries = 12
	IndustrySpecialtyChemicals EnIndustries = 13
	IndustrySteel              EnIndustries = 14
	// Basic Materials borders
	IndustryBasicMaterialsFirst EnIndustries = 1
	IndustryBasicMaterialsLast  EnIndustries = 14
	IndustryBasicMaterialsEnd   EnIndustries = 50
	// Communication Services
	IndustryAdvertising     EnIndustries = 51
	IndustryBroadcasting    EnIndustries = 52
	IndustryGamingMultimedia EnIndustries = 53
	IndustryEntertainment   EnIndustries = 54
	IndustryInternetContent EnIndustries = 55
	IndustryPublishing      EnIndustries = 56
	IndustryTelecom         EnIndustries = 57
	// Communication borders
	IndustryCommunicationFirst EnIndustries = 51
	IndustryCommunicationLast  EnIndustries = 57
	IndustryCommunicationEnd   EnIndustries = 100
	// Consumer Cyclical
	IndustryApparelManufacturing EnIndustries = 101
	IndustryApparelRetail        EnIndustries = 102
	IndustryAutoManufacturers    EnIndustries = 103
	IndustryAutoParts            EnIndustries = 104
	IndustryAutoDealership       EnIndustries = 105
	IndustryDepartmentStores     EnIndustries = 106
	IndustryFootwearAccessories  EnIndustries = 107
	IndustryFurnishings          EnIndustries = 108
	IndustryGambling             EnIndustries = 109
	IndustryHomeImprovRetail     EnIndustries = 110
	IndustryInternetRetail       EnIndustries = 111
	IndustryLeisure              EnIndustries = 112
	IndustryLodging              EnIndustries = 113
	IndustryLuxuryGoods          EnIndustries = 114
	IndustryPackagingContainers  EnIndustries = 115
	IndustryPersonalServices     EnIndustries = 116
	IndustryRecreationalVehicles EnIndustries = 117
	IndustryResidentConstruction EnIndustries = 118
	IndustryResortsCasinos       EnIndustries = 119
	IndustryRestaurants          EnIndustries = 120
	IndustrySpecialtyRetail      EnIndustries = 121
	IndustryTextileManufacturing EnIndustries = 122
	IndustryTravelServices       EnIndustries = 123
	// Consumer Cyclical borders
	IndustryConsumerCyclFirst EnIndustries = 101
	IndustryConsumerCyclLast  EnIndustries = 123
	IndustryConsumerCyclEnd   EnIndustries = 150
	// Consumer Defensive
	IndustryBeveragesBrewers  EnIndustries = 151
	IndustryBeveragesNonAlco  EnIndustries = 152
	IndustryBeveragesWineries EnIndustries = 153
	IndustryConfectioners     EnIndustries = 154
	IndustryDiscountStores    EnIndustries = 155
	IndustryEducationTraining EnIndustries = 156
	IndustryFarmProducts      EnIndustries = 157
	IndustryFoodDistribution  EnIndustries = 158
	IndustryGroceryStores     EnIndustries = 159
	IndustryHouseholdProducts EnIndustries = 160
	IndustryPackagedFoods     EnIndustries = 161
	IndustryTobacco           EnIndustries = 162
	// Consumer Defensive borders
	IndustryConsumerDefFirst EnIndustries = 151
	IndustryConsumerDefLast  EnIndustries = 162
	IndustryConsumerDefEnd   EnIndustries = 200
	// Energy
	IndustryOilGasDrilling    EnIndustries = 201
	IndustryOilGasEP          EnIndustries = 202
	IndustryOilGasEquipment   EnIndustries = 203
	IndustryOilGasIntegrated  EnIndustries = 204
	IndustryOilGasMidstream   EnIndustries = 205
	IndustryOilGasRefining    EnIndustries = 206
	IndustryThermalCoal       EnIndustries = 207
	IndustryUranium           EnIndustries = 208
	// Energy borders
	IndustryEnergyFirst EnIndustries = 201
	IndustryEnergyLast  EnIndustries = 208
	IndustryEnergyEnd   EnIndustries = 250
	// Financial
	IndustryExchangeTradedFund    EnIndustries = 251
	IndustryAssetsManagement      EnIndustries = 252
	IndustryBanksDiversified      EnIndustries = 253
	IndustryBanksRegional         EnIndustries = 254
	IndustryCapitalMarkets        EnIndustries = 255
	IndustryCloseEndFundDebt      EnIndustries = 256
	IndustryCloseEndFundEquity    EnIndustries = 257
	IndustryCloseEndFundForeign   EnIndustries = 258
	IndustryCreditServices        EnIndustries = 259
	IndustryFinancialConglomerate EnIndustries = 260
	IndustryFinancialDataExchange EnIndustries = 261
	IndustryInsuranceBrokers      EnIndustries = 262
	IndustryInsuranceDiversified  EnIndustries = 263
	IndustryInsuranceLife         EnIndustries = 264
	IndustryInsuranceProperty     EnIndustries = 265
	IndustryInsuranceReinsurance  EnIndustries = 266
	IndustryInsuranceSpecialty    EnIndustries = 267
	IndustryMortgageFinance       EnIndustries = 268
	IndustryShellCompanies        EnIndustries = 269
	// Financial borders
	IndustryFinancialFirst EnIndustries = 251
	IndustryFinancialLast  EnIndustries = 269
	IndustryFinancialEnd   EnIndustries = 300
	// Healthcare
	IndustryBiotechnology          EnIndustries = 301
	IndustryDiagnosticsResearch    EnIndustries = 302
	IndustryDrugsManufacturers     EnIndustries = 303
	IndustryDrugsManufacturersSpec EnIndustries = 304
	IndustryHealthcarePlans        EnIndustries = 305
	IndustryHealthInformation      EnIndustries = 306
	IndustryMedicalFacilities      EnIndustries = 307
	IndustryMedicalDevices         EnIndustries = 308
	IndustryMedicalDistribution    EnIndustries = 309
	IndustryMedicalInstruments     EnIndustries = 310
	IndustryPharmRetailers         EnIndustries = 311
	// Healthcare borders
	IndustryHealthcareFirst EnIndustries = 301
	IndustryHealthcareLast  EnIndustries = 311
	IndustryHealthcareEnd   EnIndustries = 350
	// Industrials
	IndustryAerospaceDefense         EnIndustries = 351
	IndustryAirlines                 EnIndustries = 352
	IndustryAirportsServices         EnIndustries = 353
	IndustryBuildingProducts         EnIndustries = 354
	IndustryBusinessEquipment        EnIndustries = 355
	IndustryConglomerates            EnIndustries = 356
	IndustryConsultingServices       EnIndustries = 357
	IndustryElectricalEquipment      EnIndustries = 358
	IndustryEngineeringConstruction  EnIndustries = 359
	IndustryFarmHeavyMachinery       EnIndustries = 360
	IndustryIndustrialDistribution   EnIndustries = 361
	IndustryInfrastructureOperations EnIndustries = 362
	IndustryFreightLogistics         EnIndustries = 363
	IndustryMarineShipping           EnIndustries = 364
	IndustryMetalFabrication         EnIndustries = 365
	IndustryPollutionControl         EnIndustries = 366
	IndustryRailroads                EnIndustries = 367
	IndustryRentalLeasing            EnIndustries = 368
	IndustrySecurityProtection       EnIndustries = 369
	IndustrySpecialtyBusinessServices EnIndustries = 370
	IndustrySpecialtyMachinery       EnIndustries = 371
	IndustryStuffingEmployment       EnIndustries = 372
	IndustryToolsAccessories         EnIndustries = 373
	IndustryTrucking                 EnIndustries = 374
	IndustryWasteManagement          EnIndustries = 375
	// Industrials borders
	IndustryIndustrialsFirst EnIndustries = 351
	IndustryIndustrialsLast  EnIndustries = 375
	IndustryIndustrialsEnd   EnIndustries = 400
	// Real Estate
	IndustryRealEstateDevelopment EnIndustries = 401
	IndustryRealEstateDiversified EnIndustries = 402
	IndustryRealEstateServices    EnIndustries = 403
	IndustryREITDiversified       EnIndustries = 404
	IndustryREITHealthcare        EnIndustries = 405
	IndustryREITHotelMotel        EnIndustries = 406
	IndustryREITIndustrial        EnIndustries = 407
	IndustryREITMortgage          EnIndustries = 408
	IndustryREITOffice            EnIndustries = 409
	IndustryREITResidential       EnIndustries = 410
	IndustryREITRetail            EnIndustries = 411
	IndustryREITSpecialty         EnIndustries = 412
	// Real Estate borders
	IndustryRealEstateFirst EnIndustries = 401
	IndustryRealEstateLast  EnIndustries = 412
	IndustryRealEstateEnd   EnIndustries = 450
	// Technology
	IndustryCommunicationEquipment EnIndustries = 451
	IndustryComputerHardware       EnIndustries = 452
	IndustryConsumerElectronics    EnIndustries = 453
	IndustryElectronicComponents   EnIndustries = 454
	IndustryElectronicDistribution EnIndustries = 455
	IndustryITServices             EnIndustries = 456
	IndustryScientificInstruments  EnIndustries = 457
	IndustrySemiconductorEquipment EnIndustries = 458
	IndustrySemiconductors         EnIndustries = 459
	IndustrySoftwareApplication    EnIndustries = 460
	IndustrySoftwareInfrastructure EnIndustries = 461
	IndustrySolar                  EnIndustries = 462
	// Technology borders
	IndustryTechnologyFirst EnIndustries = 451
	IndustryTechnologyLast  EnIndustries = 462
	IndustryTechnologyEnd   EnIndustries = 500
	// Utilities
	IndustryUtilitiesDiversified      EnIndustries = 501
	IndustryUtilitiesPowerproducers   EnIndustries = 502
	IndustryUtilitiesRenewable        EnIndustries = 503
	IndustryUtilitiesRegulatedElectric EnIndustries = 504
	IndustryUtilitiesRegulatedGas     EnIndustries = 505
	IndustryUtilitiesRegulatedWater   EnIndustries = 506
	// Utilities borders
	IndustryUtilitiesFirst EnIndustries = 501
	IndustryUtilitiesLast  EnIndustries = 506
	IndustryUtilitiesEnd   EnIndustries = 550
	// Commodities
	IndustryCommoditiesAgricultural EnIndustries = 551
	IndustryCommoditiesEnergy      EnIndustries = 552
	IndustryCommoditiesMetals      EnIndustries = 553
	IndustryCommoditiesPrecious    EnIndustries = 554
	// Commodities borders
	IndustryCommoditiesFirst EnIndustries = 551
	IndustryCommoditiesLast  EnIndustries = 554
	IndustryCommoditiesEnd   EnIndustries = 600
	// enumeration borders
	IndustryFirst EnIndustries = 0
	IndustryLast  EnIndustries = IndustryCommoditiesLast
)

// EnFillingFlags - allowed filling modes flags
type EnFillingFlags uint32

const (
	FillFlagsNone  EnFillingFlags = 0
	FillFlagsFOK   EnFillingFlags = 1
	FillFlagsIOC   EnFillingFlags = 2
	FillFlagsFirst EnFillingFlags = FillFlagsFOK
	FillFlagsAll   EnFillingFlags = FillFlagsFOK | FillFlagsIOC
)

// EnExpirationFlags - allowed order expiration modes flags
type EnExpirationFlags uint32

const (
	TimeFlagsNone         EnExpirationFlags = 0
	TimeFlagsGTC          EnExpirationFlags = 1
	TimeFlagsDay          EnExpirationFlags = 2
	TimeFlagsSpecified    EnExpirationFlags = 4
	TimeFlagsSpecifiedDay EnExpirationFlags = 8
	TimeFlagsFirst        EnExpirationFlags = TimeFlagsGTC
	TimeFlagsAll          EnExpirationFlags = TimeFlagsGTC | TimeFlagsDay | TimeFlagsSpecified | TimeFlagsSpecifiedDay
)

// EnTradeMode - allowed trade modes
type EnTradeMode uint32

const (
	TradeDisabled  EnTradeMode = 0
	TradeLongOnly  EnTradeMode = 1
	TradeShortOnly EnTradeMode = 2
	TradeCloseOnly EnTradeMode = 3
	TradeFull      EnTradeMode = 4
	TradeFirst     EnTradeMode = TradeDisabled
	TradeLast      EnTradeMode = TradeFull
)

// EnExecutionMode - order execution modes
type EnExecutionMode uint32

const (
	ExecutionRequest  EnExecutionMode = 0
	ExecutionInstant  EnExecutionMode = 1
	ExecutionMarket   EnExecutionMode = 2
	ExecutionExchange EnExecutionMode = 3
	ExecutionFirst    EnExecutionMode = ExecutionRequest
	ExecutionLast     EnExecutionMode = ExecutionExchange
)

// EnCalcMode - profit and margin calculation modes
type EnCalcMode uint32

const (
	// market maker modes
	TradeModeForex          EnCalcMode = 0
	TradeModeFutures        EnCalcMode = 1
	TradeModeCFD            EnCalcMode = 2
	TradeModeCFDIndex       EnCalcMode = 3
	TradeModeCFDLeverage    EnCalcMode = 4
	TradeModeForexNoLeverage EnCalcMode = 5
	// market maker borders
	TradeModeMMFirst EnCalcMode = TradeModeForex
	TradeModeMMLast  EnCalcMode = TradeModeForexNoLeverage
	// exchange modes
	TradeModeExchStocks      EnCalcMode = 32
	TradeModeExchFutures     EnCalcMode = 33
	TradeModeExchFuturesForts EnCalcMode = 34
	TradeModeExchOptions     EnCalcMode = 35
	TradeModeExchOptionsMargin EnCalcMode = 36
	TradeModeExchBonds       EnCalcMode = 37
	TradeModeExchStocksMOEX  EnCalcMode = 38
	TradeModeExchBondsMOEX   EnCalcMode = 39
	// exchange borders
	TradeModeExchFirst EnCalcMode = TradeModeExchStocks
	TradeModeExchLast  EnCalcMode = TradeModeExchBondsMOEX
	// service modes
	TradeModeServCollateral EnCalcMode = 64
	// service borders
	TradeModeServFirst EnCalcMode = TradeModeServCollateral
	TradeModeServLast  EnCalcMode = TradeModeServCollateral
	// enumeration borders
	TradeModeCalcFirst EnCalcMode = TradeModeForex
	TradeModeCalcLast  EnCalcMode = TradeModeServCollateral
)

// EnGTCMode - orders expiration modes
type EnGTCMode uint32

const (
	OrdersGTC          EnGTCMode = 0
	OrdersDaily        EnGTCMode = 1
	OrdersDailyNoStops EnGTCMode = 2
	OrdersGTCFirst     EnGTCMode = OrdersGTC
	OrdersGTCLast      EnGTCMode = OrdersDailyNoStops
)

// EnTickFlags - tick collection flags
type EnTickFlags uint64

const (
	TickRealtime       EnTickFlags = 1
	TickCollectRaw     EnTickFlags = 2
	TickFeedStats      EnTickFlags = 4
	TickNegativePrices EnTickFlags = 8
	TickNone           EnTickFlags = 0
	TickAll            EnTickFlags = TickRealtime | TickCollectRaw | TickFeedStats | TickNegativePrices
)

// EnChartMode - chart mode
type EnChartMode uint32

const (
	ChartModeBidPrice  EnChartMode = 0
	ChartModeLastPrice EnChartMode = 1
	ChartModeOld       EnChartMode = 255
	ChartModeFirst     EnChartMode = ChartModeBidPrice
	ChartModeLast      EnChartMode = ChartModeOld
)

// EnMarginFlags - margin check modes (symbol-level)
type EnMarginFlags uint32

const (
	SymbolMarginFlagsNone          EnMarginFlags = 0
	SymbolMarginFlagsCheckProcess  EnMarginFlags = 1
	SymbolMarginFlagsCheckSLTP     EnMarginFlags = 2
	SymbolMarginFlagsHedgeLargeLeg EnMarginFlags = 4
	SymbolMarginFlagsFirst         EnMarginFlags = SymbolMarginFlagsNone
	SymbolMarginFlagsLast          EnMarginFlags = SymbolMarginFlagsHedgeLargeLeg
)

// EnSwapMode - swaps calculation modes
type EnSwapMode uint32

const (
	SwapDisabled            EnSwapMode = 0
	SwapByPoints            EnSwapMode = 1
	SwapBySymbolCurrency    EnSwapMode = 2
	SwapByMarginCurrency    EnSwapMode = 3
	SwapByGroupCurrency     EnSwapMode = 4
	SwapByInterestCurrent   EnSwapMode = 5
	SwapByInterestOpen      EnSwapMode = 6
	SwapReopenByClosePrice  EnSwapMode = 7
	SwapReopenByBid         EnSwapMode = 8
	SwapByProfitCurrency    EnSwapMode = 9
	SwapFirst               EnSwapMode = SwapDisabled
	SwapLast                EnSwapMode = SwapByProfitCurrency
)

// EnInstantMode - Instant Execution Modes
type EnInstantMode uint32

const (
	InstantCheckNormal EnInstantMode = 0
	InstantCheckFirst  EnInstantMode = InstantCheckNormal
	InstantCheckLast   EnInstantMode = InstantCheckNormal
)

// EnRequestFlags - Request Execution Flags
type EnRequestFlags uint32

const (
	RequestFlagsNone  EnRequestFlags = 0
	RequestFlagsOrder EnRequestFlags = 1
	RequestFlagsAll   EnRequestFlags = RequestFlagsOrder
)

// EnTradeInstantFlags - Instant Execution Flags
type EnTradeInstantFlags uint32

const (
	InstantFlagsNone             EnTradeInstantFlags = 0
	InstantFlagsFastConfirmation EnTradeInstantFlags = 1
	InstantFlagsAll              EnTradeInstantFlags = InstantFlagsFastConfirmation
)

// EnMarginRateTypes - Margin Rate Types
type EnMarginRateTypes uint32

const (
	MarginRateBuy           EnMarginRateTypes = 0
	MarginRateSell          EnMarginRateTypes = 1
	MarginRateBuyLimit      EnMarginRateTypes = 2
	MarginRateSellLimit     EnMarginRateTypes = 3
	MarginRateBuyStop       EnMarginRateTypes = 4
	MarginRateSellStop      EnMarginRateTypes = 5
	MarginRateBuyStopLimit  EnMarginRateTypes = 6
	MarginRateSellStopLimit EnMarginRateTypes = 7
	MarginRateFirst         EnMarginRateTypes = MarginRateBuy
	MarginRateLast          EnMarginRateTypes = MarginRateSellStopLimit
)

// EnOrderFlags - allowed order flags
type EnOrderFlags uint32

const (
	OrderFlagsNone      EnOrderFlags = 0
	OrderFlagsMarket    EnOrderFlags = 1
	OrderFlagsLimit     EnOrderFlags = 2
	OrderFlagsStop      EnOrderFlags = 4
	OrderFlagsStopLimit EnOrderFlags = 8
	OrderFlagsSL        EnOrderFlags = 16
	OrderFlagsTP        EnOrderFlags = 32
	OrderFlagsCloseBy   EnOrderFlags = 64
	OrderFlagsFirst     EnOrderFlags = OrderFlagsMarket
	OrderFlagsAll       EnOrderFlags = OrderFlagsMarket | OrderFlagsLimit | OrderFlagsStop | OrderFlagsStopLimit | OrderFlagsSL | OrderFlagsTP | OrderFlagsCloseBy
)

// EnSymbolTradeFlags - common trade flags (symbol-level)
type EnSymbolTradeFlags uint64

const (
	SymbolTradeFlagsNone           EnSymbolTradeFlags = 0
	SymbolTradeFlagsProfitByMarket EnSymbolTradeFlags = 1
	SymbolTradeFlagsAllowSignals   EnSymbolTradeFlags = 2
	SymbolTradeFlagsAll            EnSymbolTradeFlags = SymbolTradeFlagsProfitByMarket | SymbolTradeFlagsAllowSignals
	SymbolTradeFlagsDefault        EnSymbolTradeFlags = SymbolTradeFlagsAllowSignals
)

// EnOptionMode - Options Mode
type EnOptionMode uint32

const (
	OptionModeEuropeanCall EnOptionMode = 0
	OptionModeEuropeanPut  EnOptionMode = 1
	OptionModeAmericanCall EnOptionMode = 2
	OptionModeAmericanPut  EnOptionMode = 3
	OptionModeFirst        EnOptionMode = OptionModeEuropeanCall
	OptionModeLast         EnOptionMode = OptionModeAmericanPut
)

// EnSpliceType - Splice Type
type EnSpliceType uint32

const (
	SpliceNone       EnSpliceType = 0
	SpliceUnadjusted EnSpliceType = 1
	SpliceAdjusted   EnSpliceType = 2
	SpliceFirst      EnSpliceType = SpliceNone
	SpliceLast       EnSpliceType = SpliceAdjusted
)

// EnSpliceTimeType - Splice Time Type
type EnSpliceTimeType uint32

const (
	SpliceTimeExpiration EnSpliceTimeType = 0
	SpliceTimeFirst      EnSpliceTimeType = SpliceTimeExpiration
	SpliceTimeLast       EnSpliceTimeType = SpliceTimeExpiration
)

// MTConSymbol represents symbols configuration.
type MTConSymbol struct {
	Symbol              string             `json:"Symbol"`
	Path                string             `json:"Path"`
	ISIN                string             `json:"ISIN"`
	Description         string             `json:"Description"`
	International       string             `json:"International"`
	Basis               string             `json:"Basis"`
	Source              string             `json:"Source"`
	Page                string             `json:"Page"`
	CurrencyBase        string             `json:"CurrencyBase"`
	CurrencyBaseDigits  uint32             `json:"CurrencyBaseDigits"`
	CurrencyProfit      string             `json:"CurrencyProfit"`
	CurrencyProfitDigits uint32            `json:"CurrencyProfitDigits"`
	CurrencyMargin      string             `json:"CurrencyMargin"`
	CurrencyMarginDigits uint32            `json:"CurrencyMarginDigits"`
	Color               uint32             `json:"Color"`
	ColorBackground     uint32             `json:"ColorBackground"`
	Digits              uint32             `json:"Digits"`
	Point               float64            `json:"Point"`
	Multiply            float64            `json:"Multiply"`
	TickFlags           EnTickFlags        `json:"TickFlags"`
	TickBookDepth       uint32             `json:"TickBookDepth"`
	ChartMode           EnChartMode        `json:"ChartMode"`
	FilterSoft          uint32             `json:"FilterSoft"`
	FilterSoftTicks     uint32             `json:"FilterSoftTicks"`
	FilterHard          uint32             `json:"FilterHard"`
	FilterHardTicks     uint32             `json:"FilterHardTicks"`
	FilterDiscard       uint32             `json:"FilterDiscard"`
	FilterSpreadMax     uint32             `json:"FilterSpreadMax"`
	FilterSpreadMin     uint32             `json:"FilterSpreadMin"`
	FilterGap           uint32             `json:"FilterGap"`
	FilterGapTicks      uint32             `json:"FilterGapTicks"`
	TradeMode           EnTradeMode        `json:"TradeMode"`
	TradeFlags          EnSymbolTradeFlags `json:"TradeFlags"`
	CalcMode            EnCalcMode         `json:"CalcMode"`
	ExecMode            EnExecutionMode    `json:"ExecMode"`
	GTCMode             EnGTCMode          `json:"GTCMode"`
	FillFlags           EnFillingFlags     `json:"FillFlags"`
	OrderFlags          EnOrderFlags       `json:"OrderFlags"`
	ExpirFlags          EnExpirationFlags  `json:"ExpirFlags"`
	Spread              int32              `json:"Spread"`
	SpreadBalance       int32              `json:"SpreadBalance"`
	SpreadDiff          int32              `json:"SpreadDiff"`
	SpreadDiffBalance   int32              `json:"SpreadDiffBalance"`
	TickValue           float64            `json:"TickValue"`
	TickSize            float64            `json:"TickSize"`
	ContractSize        float64            `json:"ContractSize"`
	StopsLevel          int32              `json:"StopsLevel"`
	FreezeLevel         int32              `json:"FreezeLevel"`
	QuotesTimeout       uint32             `json:"QuotesTimeout"`
	VolumeMin           uint64             `json:"VolumeMin"`
	VolumeMinExt        uint64             `json:"VolumeMinExt"`
	VolumeMax           uint64             `json:"VolumeMax"`
	VolumeMaxExt        uint64             `json:"VolumeMaxExt"`
	VolumeStep          uint64             `json:"VolumeStep"`
	VolumeStepExt       uint64             `json:"VolumeStepExt"`
	VolumeLimit         uint64             `json:"VolumeLimit"`
	VolumeLimitExt      uint64             `json:"VolumeLimitExt"`
	MarginCheckMode     EnMarginFlags      `json:"MarginCheckMode"` // deprecated, use MarginFlags
	MarginFlags         EnMarginFlags      `json:"MarginFlags"`
	MarginInitial       float64            `json:"MarginInitial"`
	MarginMaintenance   float64            `json:"MarginMaintenance"`
	MarginRateInitial   map[EnMarginRateTypes]float64 `json:"-"`
	MarginRateMaintenance map[EnMarginRateTypes]float64 `json:"-"`
	MarginRateLiquidity float64            `json:"MarginLiquidity"`
	MarginHedged        float64            `json:"MarginHedged"`
	MarginRateCurrency  float64            `json:"MarginCurrency"`
	MarginLong          float64            `json:"MarginLong"`          // deprecated
	MarginShort         float64            `json:"MarginShort"`         // deprecated
	MarginLimit         float64            `json:"MarginLimit"`         // deprecated
	MarginStop          float64            `json:"MarginStop"`          // deprecated
	MarginStopLimit     float64            `json:"MarginStopLimit"`     // deprecated
	SwapMode            EnSwapMode         `json:"SwapMode"`
	SwapLong            float64            `json:"SwapLong"`
	SwapShort           float64            `json:"SwapShort"`
	Swap3Day            int32              `json:"Swap3Day"`
	TimeStart           int64              `json:"TimeStart"`
	TimeExpiration      int64              `json:"TimeExpiration"`
	SessionsQuotes      [][]MTConSymbolSession `json:"SessionsQuotes"`
	SessionsTrades      [][]MTConSymbolSession `json:"SessionsTrades"`
	REFlags             EnRequestFlags     `json:"REFlags"`
	RETimeout           uint32             `json:"RETimeout"`
	IECheckMode         EnInstantMode      `json:"IECheckMode"`
	IETimeout           uint32             `json:"IETimeout"`
	IESlipProfit        uint32             `json:"IESlipProfit"`
	IESlipLosing        uint32             `json:"IESlipLosing"`
	IEVolumeMax         uint64             `json:"IEVolumeMax"`
	IEVolumeMaxExt      uint64             `json:"IEVolumeMaxExt"`
	PriceSettle         float64            `json:"PriceSettle"`
	PriceLimitMax       float64            `json:"PriceLimitMax"`
	PriceLimitMin       float64            `json:"PriceLimitMin"`
	PriceStrike         float64            `json:"PriceStrike"`
	OptionsMode         EnOptionMode       `json:"OptionsMode"`
	FaceValue           float64            `json:"FaceValue"`
	AccruedInterest     float64            `json:"AccruedInterest"`
	SpliceType          EnSpliceType       `json:"SpliceType"`
	SpliceTimeType      EnSpliceTimeType   `json:"SpliceTimeType"`
	SpliceTimeDays      uint32             `json:"SpliceTimeDays"`
	IEFlags             EnTradeInstantFlags `json:"IEFlags"`
	Category            string             `json:"Category"`
	Exchange            string             `json:"Exchange"`
	CFI                 string             `json:"CFI"`
	Sector              EnSectors          `json:"Sector"`
	Industry            EnIndustries       `json:"Industry"`
	Country             string             `json:"Country"`
	SubscriptionsDelay  uint32             `json:"SubscriptionsDelay"`
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

// SetVolumeMin sets VolumeMin converting from old volume format and stores as VolumeMinExt.
func (s *MTConSymbol) SetVolumeMin(v uint64) {
	s.VolumeMinExt = v * 10000
	s.VolumeMin = v
}

// SetVolumeMax sets VolumeMax converting from old volume format and stores as VolumeMaxExt.
func (s *MTConSymbol) SetVolumeMax(v uint64) {
	s.VolumeMaxExt = v * 10000
	s.VolumeMax = v
}

// SetVolumeStep sets VolumeStep converting from old volume format and stores as VolumeStepExt.
func (s *MTConSymbol) SetVolumeStep(v uint64) {
	s.VolumeStepExt = v * 10000
	s.VolumeStep = v
}

// SetVolumeLimit sets VolumeLimit converting from old volume format and stores as VolumeLimitExt.
func (s *MTConSymbol) SetVolumeLimit(v uint64) {
	s.VolumeLimitExt = v * 10000
	s.VolumeLimit = v
}

// SetIEVolumeMax sets IEVolumeMax converting from old volume format and stores as IEVolumeMaxExt.
func (s *MTConSymbol) SetIEVolumeMax(v uint64) {
	s.IEVolumeMaxExt = v * 10000
	s.IEVolumeMax = v
}

// GetVolumeMin returns VolumeMinExt converted to old volume format.
func (s *MTConSymbol) GetVolumeMin() uint64 { return s.VolumeMinExt / 10000 }

// GetVolumeMax returns VolumeMaxExt converted to old volume format.
func (s *MTConSymbol) GetVolumeMax() uint64 { return s.VolumeMaxExt / 10000 }

// GetVolumeStep returns VolumeStepExt converted to old volume format.
func (s *MTConSymbol) GetVolumeStep() uint64 { return s.VolumeStepExt / 10000 }

// GetVolumeLimit returns VolumeLimitExt converted to old volume format.
func (s *MTConSymbol) GetVolumeLimit() uint64 { return s.VolumeLimitExt / 10000 }

// GetIEVolumeMax returns IEVolumeMaxExt converted to old volume format.
func (s *MTConSymbol) GetIEVolumeMax() uint64 { return s.IEVolumeMaxExt / 10000 }

// CreateDefaultConSymbol creates an MTConSymbol with default values.
func CreateDefaultConSymbol() *MTConSymbol {
	s := &MTConSymbol{}
	s.CurrencyBase = "USD"
	s.CurrencyProfit = "USD"
	s.CurrencyMargin = "USD"
	s.Digits = 4
	s.TickFlags = TickRealtime
	s.TickBookDepth = 0
	s.FilterSoft = 100
	s.FilterSoftTicks = 10
	s.FilterDiscard = 500
	s.FilterHard = 500
	s.FilterHardTicks = 10
	s.FilterSpreadMax = 0
	s.FilterSpreadMin = 0
	s.TradeMode = TradeFull
	s.TradeFlags = SymbolTradeFlagsNone
	s.CalcMode = TradeModeForex
	s.ExecMode = ExecutionInstant
	s.GTCMode = OrdersGTC
	s.FillFlags = FillFlagsFOK
	s.ExpirFlags = TimeFlagsAll
	s.OrderFlags = OrderFlagsNone
	s.Spread = 0
	s.SpreadBalance = 0
	s.SpreadDiff = 0
	s.SpreadDiffBalance = 0
	s.TickValue = 0
	s.TickSize = 0
	s.ContractSize = 100000
	s.StopsLevel = 5
	s.FreezeLevel = 0
	s.QuotesTimeout = 0
	s.VolumeMin = 0
	s.VolumeMax = 100000
	s.VolumeStep = 10000
	s.VolumeLimit = 0
	s.MarginFlags = SymbolMarginFlagsNone
	s.MarginCheckMode = SymbolMarginFlagsNone
	s.MarginInitial = 0
	s.MarginMaintenance = 0
	s.MarginRateInitial = map[EnMarginRateTypes]float64{
		MarginRateBuy:           0.0,
		MarginRateSell:          0.0,
		MarginRateBuyLimit:      0.0,
		MarginRateSellLimit:     0.0,
		MarginRateBuyStop:       0.0,
		MarginRateSellStop:      0.0,
		MarginRateBuyStopLimit:  0.0,
		MarginRateSellStopLimit: 0.0,
	}
	s.MarginRateMaintenance = map[EnMarginRateTypes]float64{
		MarginRateBuy:           0.0,
		MarginRateSell:          0.0,
		MarginRateBuyLimit:      0.0,
		MarginRateSellLimit:     0.0,
		MarginRateBuyStop:       0.0,
		MarginRateSellStop:      0.0,
		MarginRateBuyStopLimit:  0.0,
		MarginRateSellStopLimit: 0.0,
	}
	s.MarginRateLiquidity = 0
	s.MarginHedged = 0
	s.MarginRateCurrency = 0
	s.MarginLong = 1
	s.MarginShort = 1
	s.MarginLimit = 0
	s.MarginStop = 0
	s.MarginStopLimit = 0
	s.SwapMode = SwapDisabled
	s.SwapLong = 0
	s.SwapShort = 0
	s.Swap3Day = 3
	s.TimeStart = 0
	s.TimeExpiration = 0
	s.REFlags = RequestFlagsNone
	s.RETimeout = 7
	s.IECheckMode = InstantCheckNormal
	s.IETimeout = 7
	s.IESlipProfit = 2
	s.IESlipLosing = 2
	s.IEVolumeMax = 0
	s.PriceSettle = 0
	s.PriceLimitMax = 0
	s.PriceLimitMin = 0
	s.PriceStrike = 0
	s.OptionsMode = OptionModeEuropeanCall
	s.FaceValue = 0
	s.AccruedInterest = 0
	s.SpliceType = SpliceNone
	s.SpliceTimeType = SpliceTimeExpiration
	s.SpliceTimeDays = 0
	s.IEFlags = InstantFlagsNone
	s.Category = ""
	s.Exchange = ""
	s.CFI = ""
	s.Sector = SectorUndefined
	s.Industry = IndustryUndefined
	s.Country = ""
	s.SubscriptionsDelay = 15
	return s
}
