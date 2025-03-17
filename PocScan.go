package poc_scan

import (
	"log"
	"path/filepath"

	"github.com/enenisme/poc_scan/pkg/scanner"
)

type PocScan struct {
	pocPath string
	url     string
	dir     string

	Results []PocScanResult
}

type PocScanResult struct {
	VulnName string
	Info     string
}

func NewPocScan(url string, dir string) *PocScan {
	return &PocScan{
		url: url,
		dir: dir,
	}
}

func (p *PocScan) Scan() {
	scanner := scanner.NewScanner()
	// Batch scan with directory
	files, err := filepath.Glob(filepath.Join(p.dir, "*.yaml"))
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		result, err := scanner.ScanWithPOCFile(p.url, file)
		if err != nil {
			log.Printf("Scan failed for %s: %v", file, err)
			continue
		}
		if result.Vulnerable {
			p.Results = append(p.Results, PocScanResult{
				VulnName: result.VulnName,
				Info:     result.Info,
			})
		}
	}

}
