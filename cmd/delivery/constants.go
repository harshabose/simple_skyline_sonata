package delivery

import "github.com/asticode/go-astiav"

const (
	DefaultVideoClockRate uint32 = 90000
	DefaultVideoWidth     uint16 = 854
	DefaultVideoHeight    uint16 = 480
	DefaultPixelFormat           = astiav.PixelFormatYuv420P
	DefaultVideoFPS       uint8  = 25
)

const (
	DefaultSPSBase64 = "AAAAAWdCwB+mgNg95uEAAAMAAQAAAwAy5GADDUAGGuKSAHjBlQ=="
	DefaultPPSBase64 = "AAAAAWjOPoA="
)

const (
	InitialBitrate              int64 = 800_000
	MinimumBitrate              int64 = 100_000
	MaximumBitrate              int64 = 2_000_000
	CutVideoBelowMinimumBitrate       = false
)
