package models

// SignatureStatusSummary represents counts of signatures by status.
type SignatureStatusSummary struct {
	Total    int `json:"total"`
	Signed   int `json:"signed"`
	Pending  int `json:"pending"`
	Declined int `json:"declined"`
}
