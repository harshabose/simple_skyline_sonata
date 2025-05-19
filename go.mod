module github.com/harshabose/simple_webrtc_comm

go 1.23.3

require (
	github.com/asticode/go-astiav v0.33.1
	github.com/harshabose/simple_webrtc_comm/client v0.0.0
	github.com/harshabose/simple_webrtc_comm/datachannel v0.0.0
	github.com/harshabose/simple_webrtc_comm/mediasource v0.0.0
	github.com/harshabose/simple_webrtc_comm/transcode v0.0.0
	github.com/pion/interceptor v0.1.37
	github.com/pion/webrtc/v4 v4.0.10
)

require (
	cloud.google.com/go v0.117.0 // indirect
	cloud.google.com/go/auth v0.14.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.7 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	cloud.google.com/go/firestore v1.18.0 // indirect
	cloud.google.com/go/iam v1.2.2 // indirect
	cloud.google.com/go/longrunning v0.6.2 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	firebase.google.com/go v3.13.0+incompatible // indirect
	github.com/aler9/gomavlib v1.3.0 // indirect
	github.com/asticode/go-astikit v0.52.0 // indirect
	github.com/bluenviron/gortsplib/v4 v4.12.3 // indirect
	github.com/bluenviron/mediacommon v1.14.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.1 // indirect
	github.com/harshabose/simple_webrtc_comm/mediasink v0.0.0 // indirect
	github.com/harshabose/tools/buffer v0.0.0 // indirect
	github.com/pion/datachannel v1.5.10 // indirect
	github.com/pion/dtls/v3 v3.0.4 // indirect
	github.com/pion/ice/v4 v4.0.6 // indirect
	github.com/pion/logging v0.2.3 // indirect
	github.com/pion/mdns/v2 v2.0.7 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.15 // indirect
	github.com/pion/rtp v1.8.11 // indirect
	github.com/pion/sctp v1.8.36 // indirect
	github.com/pion/sdp/v3 v3.0.10 // indirect
	github.com/pion/srtp/v3 v3.0.4 // indirect
	github.com/pion/stun/v3 v3.0.0 // indirect
	github.com/pion/transport/v3 v3.0.7 // indirect
	github.com/pion/turn/v4 v4.0.0 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	github.com/wlynxg/anet v0.0.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.58.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/oauth2 v0.26.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.10.0 // indirect
	google.golang.org/api v0.222.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20241118233622-e639e219e697 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250212204824-5a70512c5d8b // indirect
	google.golang.org/grpc v1.70.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)

replace (
	github.com/harshabose/simple_webrtc_comm/client => ./dependencies/client
	github.com/harshabose/simple_webrtc_comm/datachannel => ./dependencies/datachannel
	github.com/harshabose/simple_webrtc_comm/mediasink => ./dependencies/mediasink
	github.com/harshabose/simple_webrtc_comm/mediasource => ./dependencies/mediasource
	github.com/harshabose/simple_webrtc_comm/transcode => ./dependencies/transcode
	github.com/harshabose/tools/buffer => ./dependencies/tools/buffer
)
