package models

type SignatureQueryOptions struct {
	Search    string
	Status    *int
	Method    string
	SortBy    string
	Order     string
	CursorID  int
	Limit     int
	Direction string
}
