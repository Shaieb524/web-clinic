package models

type ScheduleDay struct {
	Day          string            `json:"day,omitempty" validate:"required"`
	Appointments []AppointmentSlot `json:"appointments,omitempty" validate:"required"`
}
