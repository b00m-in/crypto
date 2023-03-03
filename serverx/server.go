package serverx

import (
        "crypto/rand"
        "crypto/rsa"
        "crypto/x509"
        "encoding/pem"
        "fmt"
	"github.com/golang/glog"
        "golang.org/x/crypto/acme"
        "golang.org/x/crypto/acme/autocert"
        "io/ioutil"
        "net/http"
        "os"
        "strings"
        "time"
        "b00m.in/crypto/util"
)

type Serverx struct {
        Key *rsa.PrivateKey
        c *acme.Client
        man *autocert.Manager
        cfg *util.Config
}

func NewServerx(ucfg *util.Config) *Serverx {
        data, err := ioutil.ReadFile(ucfg.Sx.KeyPath)
        if err != nil {
                glog.Errorf("%s %v \n", "Generate rsa key", err)
                keyGen = true
        }
        var privKey *rsa.PrivateKey
        if keyGen {
                privKey, err = rsa.GenerateKey(rand.Reader, 2048)
                if err != nil {
                        glog.Errorf("%s \n", "Generating rsa key")
                        return nil //os.Exit(1) // no other option but to exit
                }
                // write key to file
                // open or create file
                keyFile, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
                if err != nil {
                        glog.Errorf("Error opening key file: %v \n", err)
                }
                // marshal private key to bytes
                privBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
                if err != nil {
                        glog.Errorf("Error marshaling key to bytes: %v \n", err)
                }
                // encode bytes to file
                if err := pem.Encode(keyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
                        glog.Errorf("Error encoding pem to file: %v \n", err)
                }
        } else { // using key from file
                // private key
                priv, _ := pem.Decode(data) // ignore public key
                if priv == nil || !strings.Contains(priv.Type, "PRIVATE") {
                        glog.Errorf("%s \n", "Nil rsa key")
                        return nil //os.Exit(1) // no other option but to exit
                        /*if key == nil {
                                key, err = rsa.GenerateKey(rand.Reader, 2048)
                                if err != nil {
                                        fmt.Printf("%s \n", "Generating rsa key")
                                        os.Exit(1) // no other option but to exit
                                }
                        }*/
                }
                signer, err := parsePrivateKey(priv.Bytes)
                if err != nil {
                        glog.Errorf("%s \n", "Parsing rsa key")
                        return nil //os.Exit(1) // no other option but to exit
                }
                privKey = signer.(*rsa.PrivateKey)
        }
        var leprod2 string
        if ucfg.Le.Prod2 != "" {
                leprod2 = ucfg.Le.Prod2
        } else {
                leprod2 = "https://acme-v02.api.letsencrypt.org/directory"
        }
        sc := &acme.Client{DirectoryURL: leprod2, Key: privKey}
        sman := &autocert.Manager{
                Client: sc,
                Email: ucfg.Sx.Email,
                Prompt: autocert.AcceptTOS,
                Cache: autocert.DirCache(ucfg.Sx.Cache),
                RenewBefore: time.Duration(ucfg.Sx.RenewBefore)*time.Hour,
                HostPolicy: autocert.HostWhitelist(ucfg.Sx.WhiteList...),
        }
        return &Serverx{c: sc, man: sman, cfg: ucfg}
}

func (s *Serverx) Run() error {
        mux := http.NewServeMux()
        mux.Handle("/", http.HandlerFunc(handleRoot))
        hs := http.Server{
                ReadTimeout: time.Duration(2) * time.Second,
                WriteTimeout: time.Duration(2) * time.Second,
                Addr: fmt.Sprintf(":%d", s.cfg.Sx.HttpsPort), //httpsPort), //":https", //
                TLSConfig: s.man.TLSConfig(),
                Handler: mux,
        }
        err := hs.ListenAndServeTLS("", "")
        if err != nil {
                glog.Errorf("Https %v \n", err)
                return err
        }
        return nil
}

func (s *Serverx) RunGo(errch chan<- error) {
        mux := http.NewServeMux()
        mux.Handle("/", http.HandlerFunc(handleRoot))
        hs := http.Server{
                ReadTimeout: time.Duration(2) * time.Second,
                WriteTimeout: time.Duration(2) * time.Second,
                Addr: fmt.Sprintf(":%d", s.cfg.Sx.HttpsPort), //httpsPort), //":https", //
                TLSConfig: s.man.TLSConfig(),
                Handler: mux,
        }
        err := hs.ListenAndServeTLS("", "")
        if err != nil {
                glog.Errorf("Https %v \n", err)
                errch<-err //return err
        }
        errch<-nil //return nil
}
