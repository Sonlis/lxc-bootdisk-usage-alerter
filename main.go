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
    var gotifyToken, gotifyHost string
    
    if gotifyToken = os.Getenv("GOTIFY_TOKEN"); gotifyToken == "" {
        log.Fatalf("GOTIFY_TOKEN env variable not set")
    }

    if gotifyHost = os.Getenv("GOTIFY_HOST"); gotifyHost == "" {
        log.Fatalf("GOTIFY_HOST env variable not set")
    }

    lxcs, err := lxc.List()
    if err != nil {
        log.Fatalf("Error listing lxcs: %v", err)
    }

    if lxcs == nil {
        log.Printf("No lxc running")
        os.Exit(0)
    }

    var wg sync.WaitGroup
    for _, lxContainer := range lxcs {
        wg.Add(1)
        go func(lxContainer lxc.Lxc) {
            defer wg.Done()
            usage, err := lxContainer.GetStorageUsage()
            if err != nil {
                log.Printf("Error getting storage usage for lxc %s: %v", lxContainer.Name, err)
            }
            if usage > 89 {
                message := fmt.Sprintf("Used storage is at %f%%", usage)
                err = alerting.AlertServiceUnealthy(lxContainer.Name, message, gotifyToken, gotifyHost)
                if err != nil {
                    log.Fatalf("Error sending alert for lxc %s: %v", lxContainer.Name, err)
                }
            }
        }(lxContainer)
    }
    wg.Wait()
}
