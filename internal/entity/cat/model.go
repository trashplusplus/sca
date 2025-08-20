package cat

type Cat struct {
	Id         int    `json:"id"omitempty`
	Name       string `json:"name" binding:"required"`
	Experience int    `json:"experience" binding:"required"`
	Breed      string `json:"breed" binding:"required"`
	Salary     int    `json:"salary" binding:"required"`
}
