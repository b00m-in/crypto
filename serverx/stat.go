package serverx

import (
        "bytes"
        "fmt"
        "io/ioutil"
	"os"
        "strings"
        "time"
        "b00m.in/crypto/util"
)

var (
        crtext = "-crt.pem"
        keyext = "-key.pem"
)

// formatForProxy splits a single file containing both EC private key and certificates into 2 files containing the private key and certificates respectively
func formatForProxy(inPath, outKeyPath, outCrtPath string) bool {
        /*crt, err := os.OpenFile(inPath, os.O_RDONLY, 0600)
        if err != nil {
                fmt.Printf("Error opening crt file: %v \n", err)
                return false
        }*/

        crtOnly, err := os.OpenFile(outCrtPath, os.O_WRONLY|os.O_CREATE/*|os.O_TRUNC*/, 0600)
        if err != nil {
                fmt.Printf("Error opening crtOnly file: %v \n", err)
                return false
        }

        keyOnly, err := os.OpenFile(outKeyPath, os.O_WRONLY|os.O_CREATE/*|os.O_TRUNC*/, 0600)
        if err != nil {
                fmt.Printf("Error opening keyOnly file: %v \n", err)
                return false
        }
        defer crtOnly.Close()
        defer keyOnly.Close()

        /*r, w, err := os.Pipe()
        if err != nil {
                fmt.Printf("Error piping: %v \n", err)
                return false
        }
        r = crt
        w = crtOnly*/

        crtb, err := ioutil.ReadFile(inPath)
        if err != nil {
                fmt.Printf("Error reading: %v \n", err)
                return false
        }
        bs := bytes.Split(crtb, []byte("END EC PRIVATE KEY-----\n"))
        if _, err = keyOnly.Write(bs[0]); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        if _, err = keyOnly.WriteAt([]byte("END EC PRIVATE KEY-----\n"), int64(len(bs[0]))); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        if _, err = crtOnly.Write(bs[1]); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        return true
}

// formatForProxyRSA splits a single file containing both RSA private key and certificates into 2 files containing the private key and certificates respectively
func formatForProxyRSA(inPath, outKeyPath, outCrtPath string) bool {
        /*crt, err := os.OpenFile(inPath, os.O_RDONLY, 0600)
        if err != nil {
                fmt.Printf("Error opening crt file: %v \n", err)
                return false
        }*/

        crtOnly, err := os.OpenFile(outCrtPath, os.O_WRONLY|os.O_CREATE/*|os.O_TRUNC*/, 0600)
        if err != nil {
                fmt.Printf("Error opening crtOnly file: %v \n", err)
                return false
        }

        keyOnly, err := os.OpenFile(outKeyPath, os.O_WRONLY|os.O_CREATE/*|os.O_TRUNC*/, 0600)
        if err != nil {
                fmt.Printf("Error opening keyOnly file: %v \n", err)
                return false
        }
        defer crtOnly.Close()
        defer keyOnly.Close()

        /*r, w, err := os.Pipe()
        if err != nil {
                fmt.Printf("Error piping: %v \n", err)
                return false
        }
        r = crt
        w = crtOnly*/

        crtb, err := ioutil.ReadFile(inPath)
        if err != nil {
                fmt.Printf("Error reading: %v \n", err)
                return false
        }
        bs := bytes.Split(crtb, []byte("END RSA PRIVATE KEY-----\n"))
        if _, err = keyOnly.Write(bs[0]); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        if _, err = keyOnly.WriteAt([]byte("END RSA PRIVATE KEY-----\n"), int64(len(bs[0]))); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        if _, err = crtOnly.Write(bs[1]); err != nil {
                fmt.Printf("Error writing: %v \n", err)
                return false
        }
        return true
}

//StatCertsCfg validates cfg before calling StatCerts
func StatCertsCfg(cfg util.Watch, errch chan<- error) {
        if cfg.Stat == 0 {
                cfg.Stat = 60
        }
        if cfg.BundleDir != "" && cfg.CertDir != "" && cfg.KeyDir != "" && len(cfg.Domains) > 0 {
                StatCerts(cfg, errch)
        }
}

//StatCerts periodically checks stat on cert bundles to see if they've changed and calls formatForProxy/formatForProxyRSA to overwrite the certs chains / keys for the proxy in required format
func StatCerts(cfg util.Watch, errch chan<- error) {
        crtOrig := cfg.BundleDir
        crtDest := cfg.CertDir
        keyDest := cfg.KeyDir
        domains := cfg.Domains
        stat := cfg.Stat
        lastMods := make([]time.Time, len(domains))
        for i:=0; i<len(domains); i++ {
                crtInfo, err := os.Lstat(crtOrig + domains[i])
                if err != nil || crtInfo == nil {
                        fmt.Printf("Couldn't stat file %v %v \n", crtOrig + domains[i], err)
                        lastMods[i] = time.Now()
                        errch<-fmt.Errorf("%s", domains[i])
                        //break
                } else {
                        lastMods[i] = crtInfo.ModTime()
                }
        }
        for {
                select {
                case <-time.After(time.Duration(stat) * time.Second):
                        for i:=0; i<len(domains); i++ {
                                crtInfo, err := os.Lstat(crtOrig + domains[i])
                                if err != nil {
                                        fmt.Printf("Couldn't stat file %v %v \n", crtOrig + domains[i], err)
                                }

                                if crtInfo == nil {
                                        break
                                }

                                if crtInfo.ModTime() != lastMods[i] {
                                        lastMods[i] = crtInfo.ModTime()
                                        fmt.Printf("file changed %v %v \n", crtOrig + domains[i], lastMods[i])
                                        if strings.Contains(domains[i], "rsa") { // the private key is RSA
                                                if !formatForProxyRSA(crtOrig + domains[i], keyDest + domains[i] + keyext, crtDest + domains[i] + crtext) {
                                                        fmt.Printf("couldn't format %v to %v %v\n", crtOrig + domains[i], crtDest+ domains[i] + crtext, keyDest + domains[i] + keyext)
                                                } else {
                                                        fmt.Printf("formatted %v to %v %v\n", crtOrig + domains[i], crtDest + domains[i] + crtext, keyDest + domains[i] + keyext)
                                                }
                                        } else { // the private key is EC
                                                if !formatForProxy(crtOrig + domains[i], keyDest + domains[i] + keyext, crtDest + domains[i] + crtext) {
                                                        fmt.Printf("couldn't format %v to %v %v\n", crtOrig + domains[i], crtDest+ domains[i] + crtext, keyDest + domains[i] + keyext)
                                                } else {
                                                        fmt.Printf("formatted %v to %v %v\n", crtOrig + domains[i], crtDest + domains[i] + crtext, keyDest + domains[i] + keyext)
                                                }
                                        }
                                } else {
                                        //fmt.Printf("file unchanged %v %v \n", crtOrig, lastMod)
                                }
                        }
                }
        }
}
