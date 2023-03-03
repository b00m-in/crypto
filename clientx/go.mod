module b00m.in/crypto/clientx

require (
	github.com/golang/glog v1.0.0
	golang.org/x/crypto v0.0.0-20200406173513-056763e48d71
)

require (
	b00m.in/crypto/util v0.0.0
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 // indirect
	golang.org/x/text v0.3.0 // indirect
)

replace b00m.in/crypto/util => ../util

go 1.19
