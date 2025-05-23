package delivery

import "github.com/asticode/go-astiav"

const (
	DefaultVideoClockRate uint32 = 90000
	DefaultVideoWidth     uint16 = 1280
	DefaultVideoHeight    uint16 = 720
	DefaultPixelFormat           = astiav.PixelFormatYuv420P
	DefaultVideoFPS       uint8  = 25
)

const (
	DefaultSPSBase64 = "AAAAAWdCwB+2gFAFuhAAAAMAEAAAAwMo8YMq"
	DefaultPPSBase64 = "AAAAAWjOPIA="
)
