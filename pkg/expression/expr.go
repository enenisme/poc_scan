// pkg/expression/expr.go
package expression

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"example.com/project/pkg/utils"
)

type Evaluator struct {
	ruleResults map[string]bool
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		ruleResults: make(map[string]bool),
	}
}

func (e *Evaluator) SetRuleResult(ruleName string, result bool) {
	e.ruleResults[ruleName] = result
}

func (e *Evaluator) Evaluate(expression string, resp *utils.HTTPResponse) (bool, error) {
	// Handle common response checks
	if strings.Contains(expression, "response.status") {
		statusCheck := fmt.Sprintf("response.status==%d", resp.Status)
		if !strings.Contains(expression, statusCheck) {
			return false, nil
		}
	}

	// Handle body contains check
	if strings.Contains(expression, "response.body.bcontains") {
		matches := regexp.MustCompile(`bcontains\(b'([^']+)'\)`).FindAllStringSubmatch(expression, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			searchStr := match[1]
			if !bytes.Contains(resp.Body, []byte(searchStr)) {
				return false, nil
			}
		}
	}

	// Handle header checks
	if strings.Contains(expression, "response.headers") {
		matches := regexp.MustCompile(`headers\["([^"]+)"\]\.contains\("([^"]+)"\)`).FindAllStringSubmatch(expression, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}
			headerName := match[1]
			headerValue := match[2]
			if !strings.Contains(resp.Headers.Get(headerName), headerValue) {
				return false, nil
			}
		}
	}

	// Handle regex matches
	if strings.Contains(expression, "bmatches") {
		matches := regexp.MustCompile(`'([^']+)'\.bmatches\(response\.body\)`).FindAllStringSubmatch(expression, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			pattern := match[1]
			matched, err := regexp.Match(pattern, resp.Body)
			if err != nil {
				return false, fmt.Errorf("invalid regex pattern: %v", err)
			}
			if !matched {
				return false, nil
			}
		}
	}

	return true, nil
}

func (e *Evaluator) EvaluateFinal(expression string) (bool, error) {
	// Handle rule combination logic (AND/OR)
	parts := strings.Split(expression, "&&")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasSuffix(part, "()") {
			ruleName := strings.TrimSuffix(part, "()")
			if result, exists := e.ruleResults[ruleName]; !exists || !result {
				return false, nil
			}
		}
	}
	return true, nil
}
