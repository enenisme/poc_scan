package poc_scan

import (
	"fmt"
	"testing"
)

func TestPocScan(t *testing.T) {
	pocScan := NewPocScan("https://example.com", "C:\\Users\\张裕波\\Desktop\\project\\web\\backend\\data\\poc")
	pocScan.Scan()
	for _, result := range pocScan.Results {
		fmt.Printf("VulnName: %s\nInfo: %s\n", result.VulnName, result.Info)
	}
}
