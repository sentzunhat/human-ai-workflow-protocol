package normalization

// FixOperation is one detected normalization change.
type FixOperation struct {
	OpID         string       `json:"opId"`
	Type         string       `json:"type"`
	ItemID       string       `json:"itemId"`
	FileToModify string       `json:"fileToModify"`
	LineRange    [2]int       `json:"lineRange"`
	Description  string       `json:"description"`
	Safety       string       `json:"safety"`
	Confidence   float64      `json:"confidence"`
	RuleID       string       `json:"ruleId,omitempty"`
	Blocked      *BlockedInfo `json:"blocked,omitempty"`
}

// BlockedInfo explains why a proposed operation needs human input.
type BlockedInfo struct {
	ID         string   `json:"id"`
	Rule       string   `json:"rule"`
	ItemID     string   `json:"itemId"`
	Confidence float64  `json:"confidence"`
	Candidates []string `json:"candidates"`
	Reason     string   `json:"reason"`
	Recovery   string   `json:"recovery"`
}
