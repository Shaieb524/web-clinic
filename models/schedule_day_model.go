package models

type ScheduleDay struct {
	AppointmentSlots []AppointmentSlot `json:"appointments_slots,omitempty" validate:"required"`
}
