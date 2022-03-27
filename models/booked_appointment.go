package models

import "time"

type BookedAppointment struct {
	SlotNo    int       `json:"slotNo,omitempty" validate:"required"`
	PatientId string    `json:"patientId,omitempty" validate:"required"`
	DoctorId  string    `json:"doctorId,omitempty" validate:"required"`
	Duration  int       `json:"duration,omitempty" validate:"required"`
	Date      time.Time `json:"date,omitempty" validate:"required"`
	Day       string    `json:"day,omitempty" validate:"required"`
}
