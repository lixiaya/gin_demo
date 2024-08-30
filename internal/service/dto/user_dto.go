package dto

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegister struct {
	Uid      string `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Nickname string `json:"nickname"`
	Captcha  string `json:"captcha"`
}

func (u UserRegister) TableName() string {
	return "users"
}
