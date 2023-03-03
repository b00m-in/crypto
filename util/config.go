package util

import (
        "flag"
	"fmt"
        "io/ioutil"
        "os"
        "encoding/json"
)

type Config struct {
        Le LetsEncrypt `json:"letencrypt"`
        Domains []string `json:"domains"`
        Sx Serverx `json:"serverx"`
        Wd Watch `json:"watch"`
        Sds SDS  `json:"sds"`
        Cx Clientx `json:"clientx"`
}

type LetsEncrypt struct {
        Prod2 string `json:"prod2"`
        Stag2 string `json:"stag2"`
        Reg2 string `json:"reg2"`
}

type Serverx struct {
        KeyPath string `json:"keyPath"`
        Email string `json:"email"`
        Cache string `json:"cache"`
        RenewBefore int64 `json:"renewBefore"`
        HttpsPort int `json:"httpsPort"`
        WhiteList []string `json:"whiteList"`
}

type Watch struct {
        BundleDir string `json:"bundleDir"`
        CertDir string `json:"certDir"`
        KeyDir string `json:"keyDir"`
        Stat int `json:"stat"`
        Domains []string `json:"domains"`
}

type SDS struct {
        NodeId string `json:"nodeId"`
        Version string `json:"version"`
        SdsPort int `json:"sdsPort"`
        CertDir string `json:"certDir"`
        KeyDir string `json:"keyDir"`
        Domains []string `json:"domains"`
        Debug bool `json:"debug"`
}

type Clientx struct {
        WatchDir string `json:"watchDir"`
        Domains []string `json:"domains"`
        Prod bool `json:"prod"`
}

func showHelp() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func ConfigureConfig(args []string) (*Config, error) {
        config := &Config{}
        var (
                help bool
                configFile string
        )
        fs := flag.NewFlagSet("crypto", flag.ExitOnError)
        fs.Usage = showHelp
	fs.StringVar(&configFile, "c", "./b00m.config", "Config file required.")
	fs.BoolVar(&help, "h", false, "Show this message.")
	fs.BoolVar(&help, "help", false, "Show this message.")

        fs.StringVar(&config.Sx.Email, "email", "", "letsencrypt eaccount")
        fs.StringVar(&config.Sx.Cache, "cache", "", "letsencrypt cache")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if help {
		showHelp()
		return nil, nil
	}
	if configFile != "" {
		tmpConfig, e := LoadConfig(configFile)
		if e != nil {
			return nil, e
		} else {
			config = tmpConfig
		}
	}
        if config.Sds.CertDir != "" {

        }
        return config, nil
}

func LoadConfig(filename string) (*Config, error) {
        cfgs, err := ioutil.ReadFile(filename)
        if err != nil {
                return nil, err
        }
        var config Config
        err = json.Unmarshal(cfgs, &config)
        if err != nil {
                fmt.Printf("json unmarshl error: %v\n", err)
                return nil, err
        }
        return &config, nil
}
