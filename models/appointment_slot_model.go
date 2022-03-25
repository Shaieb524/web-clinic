package models

import "time"

type AppointmentSlot struct {
	SlotNo    int    `json:"slotNo,omitempty" validate:"required"`
	PatientId string `json:"patientId,omitempty" validate:"required"`
	// DoctorId  string    `json:"doctorId,omitempty" validate:"required"`
	Duration int       `json:"duration,omitempty" validate:"required"`
	IsBooked bool      `json:"isBooked,omitempty" validate:"required"`
	Date     time.Time `json:"date,omitempty" validate:"required"`
}
