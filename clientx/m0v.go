package clientx

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
        "net/http"
        "path/filepath"
        "strings"
        "strconv"
        "time"
        "b00m.in/crypto/util"
        //"golang.org/x/net/http2"
)

var (
        //url = "https://e.m0v.in:8000/"
        url = "https://127.0.0.1:8030"
        //sni = "ev.b00m.in"
        //url = "https://www.b00m.in" //"https://b00m.in/service/1"
        //sni = "www.b00m.in"
        s1 = "https://e.m0v.in:8000/service/2"
        boom = "https://b00m.in"
)

type Client struct {
        WatchDir string
        Domains []string
        Prod bool
}

func NewClient(cfg util.Clientx) *Client {
        if cfg.WatchDir != "" {
                return &Client{WatchDir: cfg.WatchDir, Prod: cfg.Prod}
        }
        return &Client{Prod: false}
}

func (c* Client) MakeHttpsCalls(expiries <-chan error) {
        for {
                select {
                case <-time.After(time.Hour * 1):
                case err := <-expiries:
                        fmt.Printf("read channel %v \n", err)
                        // get sni
                        sni := err.Error()
                        fmt.Printf("making https call to %s\n", sni)
                        //check if rsa cert needed
                        rsa := false
                        if strings.Contains(sni, "=") {
                                ss := strings.Split(sni, "=")
                                if len(ss) > 0 {
                                        sni = ss[0]
                                        one, err := strconv.Atoi(ss[1])
                                        if err != nil {
                                                fmt.Printf("strconv rsa: %v \n", err)
                                        } else {
                                                if x509.RSA == x509.PublicKeyAlgorithm(one) {
                                                        rsa = true
                                                }
                                        }
                                }
                        }
                        if c.Prod {
                                if err := c.makeHttpsCall(sni, rsa); err != nil {
                                        fmt.Errorf("%v\n", err)
                                }
                        }
                }
        }
}

func (c *Client) makeHttpsCall(sni string, rsa bool) error {
        var tr *http.Transport
        if rsa {
                tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false, ServerName: sni, NextProtos: []string{"http/1.1", "acme-tls/1"/* "h2"  */}, CipherSuites: []uint16{tls.TLS_RSA_WITH_AES_256_CBC_SHA}}}
        } else {
	        tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false, ServerName: sni, NextProtos: []string{"http/1.1", "acme-tls/1"/* "h2"  */}}}
        }
        /*err := http2.ConfigureTransport(tr)
        if err != nil {
                fmt.Printf("req %v \n", err)
        }*/
	client := &http.Client{Transport: tr}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                fmt.Printf("req %v \n", err)
        }
        resp, err := client.Do(req)
        //resp, err := http.Get(url)
        if err != nil {
                fmt.Printf("error %v \n", err)
                return err
        }
        if resp != nil {
                defer resp.Body.Close()
        } else {
                return nil
        }
        fmt.Printf("%s \n",resp.Status)

        var bs []byte
        if resp.ContentLength > 0 {
                bs = make([]byte, resp.ContentLength)
        } else {
                bs = make([]byte, 1000)
        }

        n, err := resp.Body.Read(bs)
        if err != nil {
                fmt.Printf("error %v \n", err)
                return err
        }
        fmt.Printf("%v %v \n", n, string(bs))
        return nil
}

//CheckExpiries walks the filepath provided in cfg.WatchDir and calls util.CheckExpiry on all the certificates it encounters in the filepath. If it finds an error (an expired certificate) it adds an entry to the expired channel. It can't read from this channel. 
func (c *Client) CheckExpiries(expired chan<- error) {
        for {
                t := time.NewTicker(time.Minute * 2)
                select {
                case <-t.C:
                        fmt.Printf("walk: %s \n", c.WatchDir)
                        err := filepath.Walk(c.WatchDir, util.CheckExpiry)
                        if err != nil {
                                //add to renew queue
                                fmt.Printf("%v\n", err)
                                expired<-err
                        }
                        t.Stop()
                case <-time.After(time.Hour * 1): // check expiries every 1 hour
                        fmt.Printf("walk: %s \n", c.WatchDir)
                        err := filepath.Walk(c.WatchDir, util.CheckExpiry)
                        if err != nil {
                                //add to renew queue
                                fmt.Printf("%v\n", err)
                                expired<-err
                        }
                        //return nil
                }
        }
}
