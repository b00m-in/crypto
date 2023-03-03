package main

import (
        //"context"
	"fmt"
        "os"
        "time"
        "b00m.in/crypto/sds"
        "b00m.in/crypto/serverx"
        "b00m.in/crypto/clientx"
        "b00m.in/crypto/util"
)

var (
        //configfile string
)

func main() {
        cfg, err := util.ConfigureConfig(os.Args[1:])
        if err != nil {
                fmt.Printf("%s %v \n", "exiting", err)
        }

        //run the sds server and the routine that refreshes the snapshot
	//ctx := context.Background()
        errch := make(chan error, 5)
        ch := make(chan int)
	ss := sds.NewServer(cfg.Sds)
        if ss != nil {
                srv, err := ss.NewSDSWithCache()
                if err != nil {
                        fmt.Printf("%v \n", err)
                }
                //sds.RunServer(ctx, srv, ss.Port)
                go sds.RunServerGo(errch, srv, ss.Port)
                go ss.RefreshSc(ch)
        }

        //run serverx
        s := serverx.NewServerx(cfg)
        if s == nil {
                fmt.Printf("%s \n", "exiting")
        }
        /*if err := s.Run(); err != nil {
                fmt.Printf("%s %v \n", "exiting", err)
        }*/
        go s.RunGo(errch)
        //run clientx - it depends on serverx already running
        c := clientx.NewClient(cfg.Cx)
        expiries := make(chan error, 5)
        go c.CheckExpiries(expiries)
        go c.MakeHttpsCalls(expiries)

        go serverx.StatCertsCfg(cfg.Wd, expiries)

        for {
                select {
                case <-time.After(1 * time.Hour):
                        ch<-1 // refresh the snapshot every 1 hour
                        //fmt.Printf("%v \n", err)
                case err := <-errch:
                        fmt.Printf("%v \n", err)

                }
        }

}

