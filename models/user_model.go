package models

type User struct {
	Name          string  `json:"name,omitempty" validate:"required"`
	Email         string  `json:"email,omitempty" validate:"required"`
	Password      []byte  `json:"password,omitempty" validate:"required"`
	Role          string  `json:"role,omitempty" validate:"required"`
	Token         *string `json:"token"`
	Refresh_token *string `json:"refresh_token"`
}
