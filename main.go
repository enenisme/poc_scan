// main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"example.com/project/pkg/scanner"
)

func main() {
	var (
		pocPath string
		url     string
		dir     string
	)

	flag.StringVar(&pocPath, "poc", "", "Path to single POC file")
	flag.StringVar(&url, "url", "", "Target URL to scan")
	flag.StringVar(&dir, "dir", "", "Directory containing POC files")
	flag.Parse()

	if url == "" {
		log.Fatal("URL is required")
	}

	if pocPath == "" && dir == "" {
		log.Fatal("Either POC file or directory is required")
	}

	scanner := scanner.NewScanner()

	if pocPath != "" {
		// Single POC scan
		result, err := scanner.ScanWithPOCFile(url, pocPath)
		if err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		if result.Vulnerable {
			fmt.Printf("[VULNERABLE] Found vulnerability: %s\nInfo: %s\n", result.VulnName, result.Info)
		}
	}

	if dir != "" {
		// Batch scan with directory
		files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
		}

		for _, file := range files {
			result, err := scanner.ScanWithPOCFile(url, file)
			if err != nil {
				log.Printf("Scan failed for %s: %v", file, err)
				continue
			}
			if result.Vulnerable {
				fmt.Printf("[VULNERABLE] Found vulnerability in %s: %s\nInfo: %s\n", file, result.VulnName, result.Info)
			}
		}
	}
}
