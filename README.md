# lxc-bootdisk-usage-alerter

Small script to throw alerts to gotify in case an lxc bootdisk is above 89%.

It invokes the `pct` tool, a proxmox wrapper around `lxc-container`, to get all running lxcs, for each of them gets the memory left
on their bootdisk also using pct, and the alert to gotify through HTTP.

## Requirements

[go](https://go.dev/doc/install)

## Installation

Build the binary for the desired architecture and operating system:
```
GOOS=linux GOARCH=amd64 go build # Example for an x86_64 linux machine
```

## Testing
```
go test ./...
```
