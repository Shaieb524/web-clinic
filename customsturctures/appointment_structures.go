package customsturctures

type SlotUpdateData_2 struct {
	PatientID string
	Duration  int
	isBooked  bool
}

type NewSlotDataRequest struct {
	DoctorID       string
	PatientID      string
	AppointmentDay string
	SlotNo         string
	Duration       string
}

type BookSlotRequest struct {
	Role     string
	Slotdata NewSlotDataRequest
}
