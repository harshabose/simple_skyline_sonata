package delivery

const (
	DefaultVideoClockRate uint32 = 90000
	DefaultVideoWidth     uint16 = 854
	DefaultVideoHeight    uint16 = 480
	DefaultVideoFPS       uint8  = 25
)

const (
	DefaultSPSBase64 = "AAAAAWdCwB/aA2D3m4QAAAMABAAAAwDLkYAMNQAYa4pIAeMGVA=="
	DefaultPPSBase64 = "AAAAAWjOPIA="
)

const (
	InitialBitrate              int64 = 500_000
	MinimumBitrate              int64 = 100_000
	MaximumBitrate              int64 = 5_000_000
	CutVideoBelowMinimumBitrate       = false
)
