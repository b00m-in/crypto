# 

certman consists of the following 3 components:

+ serverx has an acme client and responds to acme/tls-1 (tls-alpn-01) challenges. It also has a goroutine that stats the cert bundles every 600s (configurable) and splits newly generated/received bundles into certs and keys. 

+ clientx watches the certs directory and checks for expired or close-to-expiry (within the next 10 days non-configurable set in util/x509.go) certificates and makes requests to serverx with the sni of these certs. It does this every hour (currently non-configurable). If the cert's `NotAfter` time before 10 days into the future `if cert.NotAfter.Before(time.Now().AddDate(0,0,10))` the cert is expired or within 10 days to expiry.

This forces serverx's acme client to make a request to letsencrypt to renew the cert because serverx is run with `RenewBefore` set to 480 hours (20 days). serverx then handles the alpn (acme/tls-s) challenge from letsencrypt and downloads the cert bundles. A serverx goroutine stats the cert directory every 600 seconds (10 minutues configurable) and watches for changes in any of the cert bundles and splits the bundle into cert chains and keys and saves them to new separate directories/files.

+ sds_server reads these certs and keys from file, creates a snapshot cache and makes them available to envoy via SDS. It also refreshes the snapshot cache every 1 hour (non-configurable) with a new version string so that envoy updates it's config every 10 minutes. This is probably overkill for certs/keys which change only once every 90 days but whatever.

envoy is configured to use certs/keys from sds_server's SDS and also forward acme/tls-1 requests to serverx. 

`b00m.in/gin` uses sds_server to also create a snapshot cache of routes with a new version string every time it is rebuilt (potentially having altered routes). Envoy then uses the new version of cluster/endpoint/routes on the fly.  

Bug:
http: TLS handshake error from 127.0.0.1:60146: acme/autocert: no token cert for



