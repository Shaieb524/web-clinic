package models

type Slot struct {
	DoctorId     string        `json:"doctorId,omitempty" validate:"required"`
	Day          string        `json:"day,omitempty" validate:"required"`
	Appointments []Appointment `json:"appointments,omitempty" validate:"required"`
}
