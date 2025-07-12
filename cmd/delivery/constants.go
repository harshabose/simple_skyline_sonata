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
	MinimumBitrate = 200_000
	InitialBitrate = 500_000
	MaximumBitrate = 800_000 // 800kbps
)
