package models

type User struct {
	Name      string         `json:"name,omitempty" validate:"required"`
	Email     string         `json:"email,omitempty" validate:"required"`
	Password  []byte         `json:"password,omitempty" validate:"required"`
	Role      string         `json:"role,omitempty" validate:"required"`
	Schedule  DoctorSchedule `json:"slot,omitempty" validate:"required"`
	Available bool           `json:"availability,omitempty" validate:"required"`
}
