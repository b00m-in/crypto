package serverx

import (
        "context"
	"crypto"
	"crypto/x509"
	"crypto/ecdsa"
	"crypto/rsa"
        "encoding/pem"
        "errors"
	"fmt"
        "golang.org/x/crypto/acme"
        "golang.org/x/crypto/acme/autocert"
        "io/ioutil"
        "net/http"
        "time"
)

var (
        httpPort int
        httpsPort int
        servicename string
        origCrt string
        dstCrt string
        dstKey string
        stat int
        domains = []string{"example.com", "example.in", "ex.ample.in", "example.com+rsa"}
        keyName = "./path_to/acme_key.pem"
        keyGen = false
        privKey *rsa.PrivateKey
        man *autocert.Manager
        c *acme.Client
        prod1 = "https://acme-v01.api.letsencrypt.org/directory"
        prod2 = "https://acme-v02.api.letsencrypt.org/directory"
        reg1 = "https://acme-v01.api.letsencrypt.org/acme/reg"
        stag2 = "https://acme-staging-v02.api.letsencrypt.org/directory"
        stag1 = "https://acme-staging.api.letsencrypt.org/directory"
)

func startHttp() {
        mux := http.NewServeMux()
        mux.Handle("/", http.HandlerFunc(handleRoot))
        //mux.Handle("/", http.HandlerFunc(RedirectHttp))
        hs := http.Server{
                ReadTimeout: time.Duration(5) * time.Second,
                WriteTimeout: time.Duration(5) * time.Second,
                Addr: fmt.Sprintf(":%d", httpPort),
                Handler: mux,
        }
        err := hs.ListenAndServe()
        if err != nil {
                fmt.Printf("Oops: %v \n", err)
        }
}

func startHttps() {
        mux := http.NewServeMux()
        /*mux.Handle("/debug/vars", expvar.Handler())
        mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
        mux.Handle("/admin/packets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                toks := strings.Split(r.URL.Path, "/")
                if len(toks) <= 3 {
                        glog.Infof("Nothing to see at %s \n", r.URL.Path)
                        epbs := make([]*data.Pub, 3)
                        render := Render {Message: "Nothing to see here", Pubs: epbs, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                id, err := strconv.ParseInt(toks[3], 10, 64)
                if err != nil {
                        glog.Infof("strconv: %v \n", err)
                        render := Render1 {"Nothing to see here", &data.Packet{}, dflt_ctgrs}
                        _ = tmpl_adm_pck_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                pk, err := data.GetLastPacket(id)
                if err != nil {
                        glog.Infof("Https %v \n", err)
                        render := Render1 {"Packets", &data.Packet{}, dflt_ctgrs}
                        _ = tmpl_adm_pck_one.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render1 {"Packets", pk, dflt_ctgrs}
                err = tmpl_adm_pck_one.ExecuteTemplate(w, "base", render)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        return
                }
                return
        }))
        mux.Handle("/api/", http.HandlerFunc(handleAPI))
        mux.Handle("/subs/", http.HandlerFunc(handleSubs))
        mux.Handle("/pubs/", http.HandlerFunc(handlePubs))
        mux.Handle("/admin/subs/", http.HandlerFunc(handleAdmin))
        mux.Handle("/admin/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                //pbs := make([]*data.Pub, 0)
                //render := Render {"Pubs", pbs, dflt_ctgrs}
                //err = tmpl_adm_gds_lst.ExecuteTemplate(w, "admin", s0)
                pbs, err := data.GetPubs(10)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        epbs := make([]*data.Pub, 3)
                        render := Render {Message: "Pubs", Pubs: epbs, Categories: dflt_ctgrs}
                        _ = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                        return
                }
                render := Render {Message: "Pubs", Pubs: pbs, Categories: dflt_ctgrs}
                err = tmpl_adm_pbs_lst.ExecuteTemplate(w, "base", render)
                if err != nil {
                        fmt.Printf("Https %v \n", err)
                        return
                }
                return
        }))*/
        mux.Handle("/", http.HandlerFunc(handleRoot))
        hs := http.Server{
                ReadTimeout: time.Duration(2) * time.Second,
                WriteTimeout: time.Duration(2) * time.Second,
                Addr: fmt.Sprintf(":%d", httpsPort), //":https", //
                TLSConfig: man.TLSConfig(),
                Handler: mux,
        }
        err := hs.ListenAndServeTLS("", "")
        if err != nil {
                fmt.Printf("Https %v \n", err)
                return
        }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
        // SERVICE_NAME used with Docker
        //fmt.Printf("%s %v \n", "Hello from ", os.Getenv("SERVICE_NAME")) 
        //fmt.Fprintf(w, "Hello from %s @ %d \n", os.Getenv("SERVICE_NAME"), httpPort)
        fmt.Printf("%s %v \n", "Hello from ", servicename) 
        fmt.Fprintf(w, "Hello from %s @ %d \n", servicename, httpPort)
}

func RedirectHttp(w http.ResponseWriter, r *http.Request) {
        if r.TLS != nil || r.Host == "" {
                http.Error(w, "Not Found", 404)
        }
        u := r.URL
        u.Host = r.Host
        u.Scheme = "https"
        switch r.Method {
        case "GET":
                http.Redirect(w, r, u.String(), 302)
        case "POST":
                http.Redirect(w, r, u.String(), 307)
        }
}

func GetKey(path string) (*ecdsa.PrivateKey, error) {
        keybs, err := ioutil.ReadFile(path)
        if err != nil {
                return nil, err
        }
        d, _ := pem.Decode(keybs)
        if d == nil {
                return nil, fmt.Errorf("pem no block found")
        }
        k, err := x509.ParseECPrivateKey(d.Bytes)
        if err != nil {
                return nil, err
        }
        return k, nil
}

// Attempt to parse the given private key DER block. OpenSSL 0.9.8 generates
// PKCS#1 private keys by default, while OpenSSL 1.0.0 generates PKCS#8 keys.
// OpenSSL ecparam generates SEC1 EC private keys for ECDSA. We try all three.
//
// Inspired by parsePrivateKey in crypto/tls/tls.go.
func parsePrivateKey(der []byte) (crypto.Signer, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			return key, nil
		case *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("acme/autocert: unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}

	return nil, errors.New("acme/autocert: failed to parse private key")
}

func discover() {
        ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
        defer cancel()

        dir, err := man.Client.Discover(ctx)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println(dir)
}

func register() {
        ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
        defer cancel()

        var contact []string
        contact = []string{"mailto:" + man.Email}
        a := &acme.Account{Contact: contact}
        x, err := man.Client.Register(ctx, a, man.Prompt)
        if err == nil || isAccountAlreadyExist(err) {
                fmt.Printf("%v %v \n", err, x)
        } else {
                fmt.Printf("Error registering: %v %v \n", err, x)
        }
}

func isAccountAlreadyExist(err error) bool {
        if err == acme.ErrAccountAlreadyExists {
                fmt.Println("Account already exists")
                return true
        }
        ae, ok := err.(*acme.Error)
        return ok && ae.StatusCode == http.StatusConflict
}

