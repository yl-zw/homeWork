package domain

type ReqSingUpUser struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}
type ReqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Profile struct {
	Email           string `json:"email"`
	UserName        string `json:"userName"`
	Birthday        string `json:"birthday"`
	PersonalProfile string `json:"personalProfile"`
}
