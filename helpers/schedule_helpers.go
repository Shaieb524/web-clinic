package helpers

import (
	"github.com/Shaieb524/web-clinic.git/models"
)

// create specified No. of appointments item for a specified docter
func GenerateAppointementSlots(slotsNo int, doctorId string) []models.AppointmentSlot {
	apptpointmentSlots := make([]models.AppointmentSlot, slotsNo)

	// assign incremental no. to slots
	for i := range apptpointmentSlots {
		apptpointmentSlots[i].SlotNo = i + 1
	}

	return apptpointmentSlots
}

// create a weelky schedule for a specified doctor
func GenerateWeekDoctorSchedule(doctorId string) models.DoctorSchedule {
	weekDays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thusday", "Friday", "Saturday"}
	schedule := make(map[string]models.ScheduleDay)
	appointmentSlots := GenerateAppointementSlots(3, doctorId)

	for i := range weekDays {
		sd := models.ScheduleDay{}
		sd.AppointmentSlots = appointmentSlots
		schedule[weekDays[i]] = sd
	}

	var ds models.DoctorSchedule
	ds.WeeklySchedule = schedule

	return ds
}
