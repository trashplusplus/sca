package target

type Target struct {
	Id       int    `json:"id"`
	Name     string `json:"name" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Notes    string `json:"notes"`
	Complete bool   `json:"complete"`
}
