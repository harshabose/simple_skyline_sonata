module github.com/harshabose/simple_webrtc_comm

go 1.23

replace (
    github.com/harshabose/simple_webrtc_comm/client => ./dependencies/client
    github.com/harshabose/simple_webrtc_comm/datachannel => ./dependencies/datachannel
    github.com/harshabose/simple_webrtc_comm/mediasink => ./dependencies/mediasink
    github.com/harshabose/simple_webrtc_comm/mediasource => ./dependencies/mediasource
    github.com/harshabose/simple_webrtc_comm/transcode => ./dependencies/transcode
    github.com/harshabose/tools/buffer => ./dependencies/tools/buffer
)