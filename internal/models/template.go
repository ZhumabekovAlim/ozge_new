package models

type Template struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"company_id"`
	Name      string `json:"name"`
	FilePath  string `json:"file_path"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}
