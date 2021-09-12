package models

type Users struct {
	Id        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email" gorm:"unique"`
	Password  []byte `json:"-"`
	UserName  string `json:"username" gorm:"unique"`
}
