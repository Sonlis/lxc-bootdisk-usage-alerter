package lxc

import (
    "os/exec"
    "fmt"
    //"os"
    "strconv"
)

type Lxc struct {
    Name string
    id string
    status string
}

func List() ([]Lxc, error) {
    cmd := exec.Command("pct", "list")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("Failed to run pct list: %v", err)
    }
    return getLxcsFromOutput(output), nil
}

func (l *Lxc) GetStorageUsage() (float64, error) {
    cmd := exec.Command("pct", "df", l.id)
    output, err := cmd.Output()
    if err != nil {
        return 0.0, fmt.Errorf("Failed to run pct df: %v", err)
    }
    words := splitWords(string(output))
    storageUsage, err := calculateStorageUsage(words[10], words[11])
    return storageUsage, err
}

// The output of pct list is as follows:
// VMID       Status     Lock         Name
// 102        running                 oui
// 103        running                 non
// Remove the first row, and return a list of lxc, one per row.
func getLxcsFromOutput(output []byte) []Lxc {
    words := splitWords(string(output))
    runningLxcs := words[4:]
    return splitLxcStatuses(runningLxcs)
}

func splitLxcStatuses(runningLxcs []string) []Lxc {
    if len(runningLxcs) == 0 {
        return nil
    }
    lxcs := []Lxc{}
    for i := 0; i < len(runningLxcs); i += 3 {
        lxcStatus := Lxc{
            id: runningLxcs[i],
            status: runningLxcs[i+1],
            Name: runningLxcs[i+2],
        }
        lxcs = append(lxcs, lxcStatus)
    }
    return lxcs
}

func splitWords(out string) []string {
    words := []string{}
    letters := []rune{}
    for _, value := range out {
        // UTF-8 values for space and line feed respectively
        if value != 32 && value != 10 {
           letters = append(letters, value) 
        } else {
            if len(letters) > 0 {
                words = append(words, string(letters))
                letters = nil
            }
        }
    }
    return words
}

func calculateStorageUsage(used, available string) (float64, error) {
    if used == "" || available == "" {
        return 0.0, fmt.Errorf("Storage used or available is empty")
    }

    if available == "0" {
        return 100.0, nil
    }

    // Remove the last element of the string as it is the unit
    usedValue, err := strconv.ParseFloat(used[:len(used) - 1], 64)
    if err != nil {
        return 0.0, fmt.Errorf("Error parsing used storage value: %v", err)
    }

    availableValue, err := strconv.ParseFloat(available[:len(available) - 1], 64)
    if err != nil {
        return 0.0, fmt.Errorf("Error parsing available storage value: %v", err)
    }

    if !isSameUnit(used, available) {
        usedMultiplier, availableMultiplier := getUnitMultiplier(used, available)
        usedValue = usedValue * usedMultiplier
        availableValue = availableValue * availableMultiplier
    }
    percentageUsed := usedValue / (usedValue + availableValue) * 100
    return percentageUsed, nil
}

func isSameUnit(used, available string) bool {
    if used[len(used) - 1] == available[len(available) - 1] {
        return true
    }
    return false
}

func getUnitMultiplier(used, available string) (float64, float64) {
    multipliers := map[string]float64 {
        "K": 1000,
        "M": 1000000,
        "G": 1000000000,
    }
    return multipliers[returnUnit(used)], multipliers[returnUnit(available)]
}

func returnUnit(measure string) string {
    return measure[len(measure) - 1 :]
}
