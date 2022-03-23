package models

type ScheduleDay struct {
	Appointments []AppointmentSlot `json:"appointments,omitempty" validate:"required"`
}
