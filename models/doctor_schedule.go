package models

type DoctorSchedule struct {
	DoctorId string
	Schedule map[string]ScheduleDay
}
