package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shaieb524/web-clinic.git/customsturctures"
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

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BodyParser(&requestData); err != nil {
		return err
	}

	requestSlotdata := requestData.Slotdata
	doctorId := requestSlotdata.DoctorID
	doctorObjId, err := primitive.ObjectIDFromHex(doctorId)
	if err != nil {
		panic(err)
	}

	query := bson.D{{Key: "_id", Value: doctorObjId}}

	var doctorDoc bson.M
	if err := userCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println("Error finding doctor : ", err)

		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{Status: http.StatusOK, Message: "failed", Data: &fiber.Map{"bookedSlot": err}},
		)
	}

	intSlotNo, err := strconv.Atoi(requestSlotdata.SlotNo)
	if err != nil {
		fmt.Println("error parsing slotNo to integer!", err)
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{Status: http.StatusOK, Message: "failed", Data: &fiber.Map{"bookedSlot": err}},
		)
	}

	var newSlotData SlotUpdateData
	newSlotData.PatientID = requestSlotdata.PatientID
	newSlotData.Duration, err = strconv.Atoi(requestSlotdata.Duration)
	newSlotData.isBooked = true

	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotdata.AppointmentDay, intSlotNo, newSlotData)

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

	fmt.Println("doctorObjId : ", doctorObjId)

	slot := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, slotDay, slotNo)

	// update booked slot with new data
	slot.(primitive.M)["patientid"] = newSlotData.PatientID
	slot.(primitive.M)["duration"] = newSlotData.Duration
	slot.(primitive.M)["isbooked"] = newSlotData.isBooked
	newSchedule := doctorProfile["schedule"]

	// update doctor document in database
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
