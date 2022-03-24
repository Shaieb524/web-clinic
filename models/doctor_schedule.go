package models

type DoctorSchedule struct {
	// DoctorId       string
	WeeklySchedule map[string]ScheduleDay
}
