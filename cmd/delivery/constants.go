package delivery

import "github.com/asticode/go-astiav"

const (
	DefaultVideoClockRate uint32 = 90000
	DefaultVideoWidth     uint16 = 1920
	DefaultVideoHeight    uint16 = 1080
	DefaultPixelFormat           = astiav.PixelFormatYuv420P
	DefaultVideoFPS       uint8  = 25
)

const (
	DefaultSPSBase64 = "AAAAAWdCwCmmgFAFuhAAAAMAEAAAAwPI8YMq"
	DefaultPPSBase64 = "AAAAAWjOPIA="
)
