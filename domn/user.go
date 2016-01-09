package domn

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
