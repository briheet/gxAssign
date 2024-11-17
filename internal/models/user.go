package models

type User struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Age       int    `json:"age" bson:"age"`
	Class     string `json:"class" bson:"class"`
	Section   string `json:"section" bson:"section"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type UserAuth struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type UserDocument struct {
	UserID string `json:"userID" bson:"userID"`
	Task   string `json:"task" bson:"task"`
	Admin  string `json:"admin" bson:"admin"`
	Status string `json:"status" bson:"status"`
}

type UserAdmins struct {
	UserID string   `json:"userID" bson:"userID"`
	Admins []string `json:"admin" bson:"admin"`
}
