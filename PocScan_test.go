package poc_scan

import (
	"fmt"
	"testing"
)

func TestPocScan(t *testing.T) {
	pocScan := NewPocScan("https://example.com")
	pocScan.Scan()
	for _, result := range pocScan.Results {
		fmt.Printf("VulnName: %s\nInfo: %s\n", result.VulnName, result.Info)
	}
}
