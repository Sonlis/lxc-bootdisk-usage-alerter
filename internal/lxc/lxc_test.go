package lxc

import (
    "testing"
)

// For the 2 following tests, their respective calling function checks
// beforehand if the string is empty, so we should not hit out of bonds
// errors when indexing based on the string length.
func TestReturnUnit(t *testing.T) {
    str := "wesh"
    want := "h"
    if got := returnUnit(str); got != want {
        t.Errorf("Want %s, got %s", want, got)
    }
}

func TestIsSameUnit(t *testing.T) {
    if !isSameUnit("haha", "olala") {
        t.Errorf("haha and olala should return the same last character")
    }

    if isSameUnit("this is not the same", "trust mE") {
        t.Errorf("same and mE should bot return the same last character")
    }
}

func TestCalculateStorageUsage(t *testing.T) {
    used := "1G"
    available := "1000M"
    want := 50.0
    percentageUsed, err := calculateStorageUsage(used, available)
    if err != nil {
        t.Errorf("Error calculating storage usage: %v", err)
    }
    if percentageUsed != want {
        t.Errorf("We should have %f%% usage, got %f%%", want, percentageUsed)
    }

    used = "1G"
    available = "1000000K"
    want = 50.0
    percentageUsed, err = calculateStorageUsage(used, available)
    if err != nil {
        t.Errorf("Error calculating storage usage: %v", err)
    }
    if percentageUsed != want {
        t.Errorf("We should have %f%% usage, got %f%%", want, percentageUsed)
    }

    used = "1G"
    available = "1G"
    want = 50.0
    percentageUsed, err = calculateStorageUsage(used, available)
    if err != nil {
        t.Errorf("Error calculating storage usage: %v", err)
    }
    if percentageUsed != want {
        t.Errorf("We should have %f%% usage, got %f%%", want, percentageUsed)
    }

    used = "1G"
    available = "0"
    want = 100.0
    percentageUsed, err = calculateStorageUsage(used, available)
    if err != nil {
        t.Errorf("Error calculating storage usage: %v", err)
    }
    if percentageUsed != want {
        t.Errorf("We should have %f%% usage, got %f%%", want, percentageUsed)
    }

    used = "1G"
    available = ""
    percentageUsed, err = calculateStorageUsage(used, available)
    if err == nil {
        t.Errorf("Should throw error when available is empty, but didn't")
    }

    used = "0"
    available = "over 9 thousand"
    percentageUsed, err = calculateStorageUsage(used, available)
    if err == nil {
        t.Errorf("Should throw error when used is empty, but didn't")
    }
}

func testGetLxcsFromOutput(t *testing.T) {
    output := "VMID Status Lock Name\n100 running 0 test\n101 running 0 test2\n"
    want := []Lxc{
        {
            id: "100",
            status: "running",
            Name: "test",
        },
        {
            id: "101",
            status: "running",
            Name: "test2",
        },
    }
    got := getLxcsFromOutput([]byte(output))
    for i := 0; i < len(want); i++ {
        if got[i] != want[i] {
            t.Errorf("Want %v, got %v", want[i], got[i])
        }
    }
    output = "VMID Status Lock Name\n"
    got = getLxcsFromOutput([]byte(output))
    if got != nil {
        t.Errorf("Should return nil when no lxc is running")
    }
}

func testSplitWords(t *testing.T) {
    out := "VMID Status Lock Name\n100 running 0 test\n101 running 0 test2\n"
    want := []string{"VMID", "Status", "Lock", "Name", "100", "running", "0", "test", "101", "running", "0", "test2"}
    got := splitWords(out)
    for i := 0; i < len(want); i++ {
        if got[i] != want[i] {
            t.Errorf("Want %s, got %s", want[i], got[i])
        }
    }
}
