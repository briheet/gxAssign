package models

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Class     string `json:"class"`
	Section   string `json:"section"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
