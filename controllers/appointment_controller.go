package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Shaieb524/web-clinic.git/configs"
	"github.com/Shaieb524/web-clinic.git/customsturctures"
	"github.com/Shaieb524/web-clinic.git/helpers"
	"github.com/Shaieb524/web-clinic.git/models"

	"github.com/Shaieb524/web-clinic.git/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookedAppointmentsCollection *mongo.Collection = configs.GetCollection(configs.DB, "bookedAppointments")

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

	if checkRoleResult := helpers.RoleValidator(requestData.Role, "patient"); checkRoleResult != "allowed" {
		return c.Status(http.StatusUnauthorized).JSON(
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
				Data: &fiber.Map{"problem": "Only patients are allowed to book appointment slots!"}},
		)
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
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: &fiber.Map{"problem": err}},
		)
	}

	intSlotNo, err := strconv.Atoi(requestSlotdata.SlotNo)
	if err != nil {
		fmt.Println("error parsing slotNo to integer!", err)
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: &fiber.Map{"problem": err}},
		)
	}

	var newSlotData SlotUpdateData
	newSlotData.PatientID = requestSlotdata.PatientID
	newSlotData.Duration, err = strconv.Atoi(requestSlotdata.Duration)
	newSlotData.isBooked = true

	// fmt.Println("doc id : ", doctorObjId)
	// fmt.Printf("doc id %T\n: ", doctorObjId)

	fmt.Println("newSlotData : ", newSlotData)

	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotdata.AppointmentDay, int32(intSlotNo), newSlotData)

	// insert booked appointment to collection
	bookedAppointmentItem := models.BookedAppointment{
		PatientId: requestSlotdata.PatientID,
		DoctorId:  requestSlotdata.DoctorID,
		SlotNo:    intSlotNo,
		Day:       requestSlotdata.AppointmentDay,
		Duration:  newSlotData.Duration,
	}
	insertToBookedAppointmentsCollection(bookedAppointmentItem)

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"bookedSlot": updatedSlot}},
	)
}

func CancelAppointmentSlot(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BodyParser(&requestData); err != nil {
		return err
	}

	checkRoleResult := requestData.Role
	if checkRoleResult != "doctor" && checkRoleResult != "admin" {
		return c.Status(http.StatusUnauthorized).JSON(
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
				Data: &fiber.Map{"problem": "Only Doctor & Admins are allowed to cancel an appointment!"}},
		)
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
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: &fiber.Map{"problem": err}},
		)
	}

	intSlotNo, err := strconv.Atoi(requestSlotdata.SlotNo)
	if err != nil {
		fmt.Println("error parsing slotNo to integer!", err)
		return c.Status(http.StatusInternalServerError).JSON(
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: &fiber.Map{"problem": err}},
		)
	}

	var newSlotData SlotUpdateData
	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotdata.AppointmentDay, int32(intSlotNo), newSlotData)
	fmt.Println("newSlot : ", newSlotData)

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"canceledSlot": updatedSlot}},
	)
}

func ExtractAppoinmentSlotFromDoctorProfile(doctorProfile primitive.M, slotDay string, slotNo int32) interface{} {
	fmt.Println("---------------------Extract-------------------------: ")

	// break down doctor schedule data structure
	ds := doctorProfile["schedule"]
	ws := ds.(primitive.M)["weeklyschedule"]
	day := ws.(primitive.M)[slotDay]
	appointmentsSlots := day.(primitive.M)["appointmentslots"]
	slot := appointmentsSlots.(primitive.A)[slotNo-1]

	return slot
}

func UpdateAppointmentSlot(doctorObjId primitive.ObjectID, doctorProfile primitive.M,
	slotDay string, slotNo int32, newSlotData SlotUpdateData) interface{} {

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

func ViewAppointmentDetails(c *fiber.Ctx) error {
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return err
	}

	doctorId := requestData["doctorId"]
	doctorProfile := getDoctorProfileByStringId(doctorId.(string))
	slot := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, requestData["slotDay"].(string), int32(requestData["slotNo"].(float64)))

	if err := requestData["role"] == "patient" && requestData["patientId"] != slot.(primitive.M)["patientid"]; err {
		return c.Status(http.StatusUnauthorized).JSON(
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed", Data: &fiber.Map{"message": "your are not authorized!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"appointment_details": slot}},
	)
}

func insertToBookedAppointmentsCollection(ba models.BookedAppointment) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := bookedAppointmentsCollection.InsertOne(ctx, ba)

	if err != nil {
		fmt.Println("erro : ", err)
	}

	fmt.Println("insert res /L ", result)
}

// func ViewPatientAppointmentsHistory(c *fiber.Ctx) error {

// }
