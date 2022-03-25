package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shaieb524/web-clinic.git/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SlotUpdateData struct {
	PatientID string
	Duration  int
	isBooked  bool
}

func BookAppointmentSlot(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	doctorId := data["doctorId"]
	doctorObjId, err := primitive.ObjectIDFromHex(doctorId)
	if err != nil {
		panic(err)
	}

	query := bson.D{{Key: "_id", Value: doctorObjId}}

	var doctorDoc bson.M
	if err := userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println("Error finding doctor : ", err)
	}

	intSlotNo, err := strconv.Atoi(data["slotNo"])
	if err != nil {
		fmt.Println("error parsing slotNo to integer!")
	}

	var newSlotData SlotUpdateData
	newSlotData.PatientID = data["patientId"]
	newSlotData.Duration, err = strconv.Atoi(data["duration"])
	newSlotData.isBooked = true

	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, data["appointmentDay"], intSlotNo, newSlotData)

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"bookedSlot": updatedSlot}},
	)

}

func ExtractAppoinmentSlotFromDoctorProfile(doctorProfile primitive.M, slotDay string, slotNo int) interface{} {

	// break down doctor schedule data structure
	ds := doctorProfile["schedule"]
	ws := ds.(primitive.M)["weeklyschedule"]
	day := ws.(primitive.M)[slotDay]
	appointmentsSlots := day.(primitive.M)["appointments"]
	slot := appointmentsSlots.(primitive.A)[slotNo-1]

	return slot
}

func UpdateAppointmentSlot(doctorObjId primitive.ObjectID, doctorProfile primitive.M,
	slotDay string, slotNo int, newSlotData SlotUpdateData) interface{} {

	slot := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, slotDay, slotNo)
	slot.(primitive.M)["patientid"] = newSlotData.PatientID
	slot.(primitive.M)["duration"] = newSlotData.Duration
	slot.(primitive.M)["isbooked"] = newSlotData.isBooked
	newSchedule := doctorProfile["schedule"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	updatedSchedule, err := userCollection.UpdateOne(
		ctx,
		bson.M{"_id": doctorObjId},
		bson.D{
			{"$set", bson.D{{"schedule", newSchedule}}},
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated %v Documents!\n", updatedSchedule.ModifiedCount)

	return slot
}
