// pkg/scanner/scanner.go
package scanner

import (
	"fmt"
	"net/url"

	"github.com/enenisme/poc_scan/pkg/expression"
	"github.com/enenisme/poc_scan/pkg/model"
	"github.com/enenisme/poc_scan/pkg/utils"
)

type Scanner struct {
	httpClient *utils.HTTPClient
}

func NewScanner() *Scanner {
	return &Scanner{
		httpClient: utils.NewHTTPClient(),
	}
}

func (s *Scanner) ScanWithPOCFile(targetURL string, pocPath string) (*model.ScanResult, error) {
	poc, err := utils.LoadPOCFromFile(pocPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load POC: %v", err)
	}
	return s.ScanWithPOC(targetURL, poc)
}

func (s *Scanner) ScanWithPOC(targetURL string, poc *model.POC) (*model.ScanResult, error) {
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %v", err)
	}

	// Initialize result
	result := &model.ScanResult{
		Vulnerable: false,
		VulnName:   poc.Name,
		Info:       "",
	}

	// Execute rules based on expression
	exprEval := expression.NewEvaluator()

	// Process each rule
	for ruleName, rule := range poc.Rules {
		resp, err := s.httpClient.DoRequest(baseURL, &rule.Request)
		if err != nil {
			return nil, fmt.Errorf("request failed for rule %s: %v", ruleName, err)
		}

		// Evaluate rule expression
		isVuln, err := exprEval.Evaluate(rule.Expression, resp)
		if err != nil {
			return nil, fmt.Errorf("expression evaluation failed for rule %s: %v", ruleName, err)
		}

		// Store rule result for final expression evaluation
		exprEval.SetRuleResult(ruleName, isVuln)
	}

	// Evaluate final expression
	if poc.Expression != "" {
		result.Vulnerable, err = exprEval.EvaluateFinal(poc.Expression)
		if err != nil {
			return nil, fmt.Errorf("final expression evaluation failed: %v", err)
		}

		if result.Vulnerable && poc.Detail.Vulnerability.Proof != nil {
			if info, ok := poc.Detail.Vulnerability.Proof.(map[interface{}]interface{})["info"]; ok {
				result.Info = fmt.Sprintf("%v", info)
			}
		}
	}

	return result, nil
}
