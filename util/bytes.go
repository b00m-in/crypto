package util

import (
	"fmt"
        "io/ioutil"
)

var (
        crtext = "-crt.pem"
        keyext = "-key.pem"
)

//
func GetBytes(certDir, keyDir string, domains []string) ([][]byte, [][]byte, error) {
        cbs := make([][]byte, 0)
        kbs := make([][]byte, 0)
        if len(domains) == 0{
                return cbs, kbs, fmt.Errorf("%s \n", "no domains")
        }
        for _, d := range domains {
                cb, err := ioutil.ReadFile(certDir + d + crtext)
                if err != nil {
                        fmt.Errorf("%s %v \n", certDir + d + crtext, err)
                        break //skip the key
                }

                kb, err := ioutil.ReadFile(keyDir + d + keyext)
                if err != nil {
                        fmt.Errorf("%s %v \n", keyDir + d + keyext, err)
                        break //skip adding to slice
                }
                cbs = append(cbs, cb)
                kbs = append(kbs, kb)
        }
        return cbs, kbs, nil
}

