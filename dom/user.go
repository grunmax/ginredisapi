package dom

//type TodoItem struct {
//	Id        string `redis:"id" 		json:"id"` //`json:"-"`
//	Title     string `redis:"title" 	json:"title"`
//	Completed bool   `redis:"completed"	json:"completed"`
//	Order     int    `redis:"order"		json:"order"`
//}

type UserSignupForm struct {
	Id              string
	Email           string `validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password        string `validate:"nonzero"`
	PasswordConfirm string `validate:"nonzero"`
}

type UserLoginForm struct {
	Id       string
	Email    string `validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password string `validate:"nonzero"`
}
