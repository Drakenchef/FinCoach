package models

import (
	"FinCoach/internal/app/role"
	"time"
)

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	Role        role.Role     `json:"role"`
	Username    string        `json:"userName"`
	UserId      int           `json:"userid"`
}

type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserName string `json:"user_name"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}
