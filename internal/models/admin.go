package models

type Admin struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
}

type AdminAuth struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type AdminsUsers struct {
	AdminID string   `json:"userID" bson:"adminID"`
	Users   []string `json:"admin" bson:"users"`
}

type AssignmentStatus struct {
	AdminID string `json:"adminID" bson:"adminID"`
	UserID  string `json:"userID" bson:"userID"`
	Task    string `json:"task" bson:"task"`
	Status  bool   `json:"status" bson:"status"`
}
