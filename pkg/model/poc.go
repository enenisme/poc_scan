// pkg/model/poc.go
package model

type POC struct {
	ID         string            `yaml:"id"`
	Name       string            `yaml:"name"`
	Tags       []string          `yaml:"tags"`
	Transport  string            `yaml:"transport"`
	Rules      map[string]*Rule  `yaml:"rules"`
	Set        map[string]string `yaml:"set"`
	Expression string            `yaml:"expression"`
	Detail     Detail            `yaml:"detail"`
}

type Rule struct {
	Request    Request     `yaml:"request"`
	Expression string      `yaml:"expression"`
	Output     interface{} `yaml:"output"`
}

type Request struct {
	Cache           bool              `yaml:"cache"`
	Method          string            `yaml:"method"`
	Path            string            `yaml:"path"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Body            string            `yaml:"body,omitempty"`
	FollowRedirects bool              `yaml:"follow_redirects"`
}

type Detail struct {
	Fingerprint   Fingerprint   `yaml:"fingerprint"`
	Vulnerability Vulnerability `yaml:"vulnerability"`
}

type Fingerprint struct {
	Softhard       string `yaml:"softhard"`
	Product        string `yaml:"product"`
	Company        string `yaml:"company"`
	Category       string `yaml:"category"`
	ParentCategory string `yaml:"parent_category"`
}

type Vulnerability struct {
	Proof interface{} `yaml:"proof"`
}

type ScanResult struct {
	Vulnerable bool
	VulnName   string
	Info       string
	Error      error
}
