package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/sonlis/lxc-bootdisdk-usage-alerter/internal/alerting"
	"github.com/sonlis/lxc-bootdisdk-usage-alerter/internal/lxc"
)

func main() {
    gotifyToken := os.Getenv("GOTIFY_TOKEN")
    gotifyHost := os.Getenv("GOTIFY_HOST")
    lxcs, err := lxc.List()
    if err != nil {
        log.Fatalf("%v", err)
    }
    var wg sync.WaitGroup
    for _, lxContainer := range lxcs {
        wg.Add(1)
        go func(lxContainer lxc.Lxc) {
            defer wg.Done()
            usage, err := lxContainer.GetStorageUsage()
            if err != nil {
                log.Fatalf("%v", err)
            }
            if usage > 89 {
                message := fmt.Sprintf("Used storage is at %f%%", usage)
                err = alerting.AlertServiceUnealthy(lxContainer.Name, message, gotifyToken, gotifyHost)
                if err != nil {
                    log.Fatalf("%v", err)
                }
            }
        }(lxContainer)
    }
    wg.Wait()
}
