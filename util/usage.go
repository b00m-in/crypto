package util

var usageStr = `
Usage: certman [options]

Common Options:
    -h, --help                        Show this message
    -c,  --config <file>              Configuration file (eg. b00m.config) required

serverx Options:
    -httpPort, --port <port>          Use port to listen for http connection requests (default "1883")
    --host <host>                     Network host to listen on. (default "0.0.0.0")

sds Options:
    -sdsPort <port>                   Use port to listen on for envoy grpc connections
    -nodeId <test-id>                 Use test-id as node (eg. -nodeId "test-id")

clientx Options:
    -w, -watchDir                     Directory to watch for cert expiries
        -domains                      List of domains for which to request new/renewed certs from letsencrypt


Logging Options:
    -d, --debug <bool>                Enable debugging output (default false)
    -D                                Debug and trace

`
