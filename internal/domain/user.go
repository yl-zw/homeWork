package domain

import "encoding/json"

type ReqSingUpUser struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Phone           string `json:"phone"`
}
type ReqLoginUser struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	LoginMethod int    `json:"loginMethod"`
}
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

type Profile struct {
	Email           string `json:"email"`
	UserName        string `json:"userName"`
	Birthday        string `json:"birthday"`
	PersonalProfile string `json:"personalProfile"`
}

func (u *Profile) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *Profile) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
