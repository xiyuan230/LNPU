package model

const (
	GUEST = iota
	ADMIN
)

type User struct {
	OpenID       string `json:"openid" gorm:"primaryKey"`
	Name         string `json:"name"`
	StudentID    string `json:"student_id"`
	SSOPassword  string `json:"sso_password"`
	JwxtPassword string `json:"jwxt_password"`
	ExpPassword  string `json:"exp_password"`
	Role         int    `json:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == ADMIN
}
