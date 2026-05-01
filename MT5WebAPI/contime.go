package mt5webapi

// EnTimeTableMode - day working mode
type EnTimeTableMode int32

const (
	TimeModeDisabled EnTimeTableMode = 0 // work enabled
	TimeModeEnabled  EnTimeTableMode = 1 // work disabled
	TimeModeFirst    EnTimeTableMode = TimeModeDisabled
	TimeModeLast     EnTimeTableMode = TimeModeEnabled
)

// MTConTime represents time configuration on server.
type MTConTime struct {
	DaylightState int32             `json:"DaylightState"`
	Daylight      int32             `json:"Daylight"`
	TimeZone      int32             `json:"TimeZone"`
	TimeServer    string            `json:"TimeServer"`
	Days          [][]EnTimeTableMode `json:"Days"`
}
