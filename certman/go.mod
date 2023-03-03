module b00m.in/crypto/certman

replace b00m.in/crypto/serverx => ../serverx

replace b00m.in/crypto/util => ../util

replace b00m.in/crypto/sds => ../sds

replace b00m.in/crypto/clientx => ../clientx

go 1.19

require (
	b00m.in/crypto/clientx v0.0.0
	b00m.in/crypto/sds v0.0.0
	b00m.in/crypto/serverx v0.0.0-00010101000000-000000000000
	b00m.in/crypto/util v0.0.0
)

require (
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/envoyproxy/go-control-plane v0.10.3-0.20221219165740-8b998257ff09 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.9.1 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20220822174746-9e6da59bd2fc // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
