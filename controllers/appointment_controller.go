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
	"github.com/gin-gonic/gin"

	"github.com/Shaieb524/web-clinic.git/responses"
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

func BookAppointmentSlot(c *gin.Context) {

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "failed to parse request", Data: map[string]interface{}{"error": "could't validate request body!"}})
		return
	}

	if checkRoleResult := helpers.RoleValidator(requestData.Role, "patient"); checkRoleResult != "allowed" {
		c.JSON(http.StatusUnauthorized,
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
				Data: map[string]interface{}{"problem": "Only patients are allowed to book appointment slots!"}},
		)
		return
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
		c.JSON(http.StatusInternalServerError,
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error finding doctor", Data: map[string]interface{}{"error": err}},
		)
		return
	}

	intSlotNo, err := strconv.Atoi(requestSlotdata.SlotNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "error parsing slotNo to integer!", Data: map[string]interface{}{"error": err}},
		)
		return
	}

	var newSlotData SlotUpdateData
	newSlotData.PatientID = requestSlotdata.PatientID
	newSlotData.Duration, err = strconv.Atoi(requestSlotdata.Duration)
	newSlotData.isBooked = true

	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotdata.AppointmentDay, int32(intSlotNo), newSlotData)

	// insert booked appointment to collection
	bookedAppointmentItem := models.BookedAppointment{
		PatientId: requestSlotdata.PatientID,
		DoctorId:  requestSlotdata.DoctorID,
		SlotNo:    intSlotNo,
		Day:       requestSlotdata.AppointmentDay,
		Duration:  newSlotData.Duration,
	}
	insertItemToBookedAppointmentsCollection(bookedAppointmentItem)

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"bookedSlot": updatedSlot}},
	)
}

func CancelAppointmentSlot(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "failed", Data: map[string]interface{}{"data": "could't validate request body!"}})
		return
	}

	checkRoleResult := requestData.Role
	if checkRoleResult != "doctor" && checkRoleResult != "admin" {
		c.JSON(http.StatusUnauthorized,
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed",
				Data: map[string]interface{}{"problem": "Only Doctor & Admins are allowed to cancel an appointment!"}},
		)
		return
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
		c.JSON(http.StatusInternalServerError,
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: map[string]interface{}{"problem": err}},
		)
		return
	}

	intSlotNo, err := strconv.Atoi(requestSlotdata.SlotNo)
	if err != nil {
		fmt.Println("error parsing slotNo to integer!", err)
		c.JSON(http.StatusInternalServerError,
			responses.UserResponse{Status: http.StatusInternalServerError, Message: "failed", Data: map[string]interface{}{"problem": err}},
		)
		return
	}

	var newSlotData SlotUpdateData
	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotdata.AppointmentDay, int32(intSlotNo), newSlotData)

	deleteItemFromBookedAppointmentsCollection(requestSlotdata.DoctorID, requestSlotdata.AppointmentDay, int32(intSlotNo))

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"canceledSlot": updatedSlot}},
	)
}

//TODO should return error handle
func ExtractAppoinmentSlotFromDoctorProfile(doctorProfile primitive.M, slotDay string, slotNo int32) interface{} {
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

func ViewAppointmentDetails(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "failed", Data: map[string]interface{}{"data": "could't validate request body!"}})
		return
	}

	doctorId := requestData["doctorId"]
	doctorProfile := getDoctorProfileByStringId(doctorId.(string))
	slot := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, requestData["appointmentDay"].(string), int32(requestData["slotNo"].(float64)))

	if err := requestData["role"] == "patient" && requestData["patientId"] != slot.(primitive.M)["patientid"]; err {
		c.JSON(http.StatusUnauthorized,
			responses.UserResponse{Status: http.StatusUnauthorized, Message: "failed", Data: map[string]interface{}{"message": "your are not authorized!"}},
		)
		return
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"appointment_details": slot}},
	)
}

func insertItemToBookedAppointmentsCollection(ba models.BookedAppointment) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := bookedAppointmentsCollection.InsertOne(ctx, ba)

	if err != nil {
		fmt.Println("err inserting item from booked appointments : ", err)
	}

	fmt.Println("item inserted to booked appointemtn res : ", result)
}

func deleteItemFromBookedAppointmentsCollection(doctorId string, slotDay string, slotNo int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"$and": []bson.M{
			{"doctorid": doctorId},
			{"day": slotDay},
			{"slotno": slotNo},
		},
	}
	result, err := bookedAppointmentsCollection.DeleteOne(ctx, query)

	if err != nil {
		fmt.Println("err deleting item from booked appointments : ", err)
	}

	fmt.Println("item deleted from booked appointemtn res : ", result)
}

func ViewPatientAppointmentsHistory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	query := bson.M{"patientid": idParam}
	results, err := bookedAppointmentsCollection.Find(ctx, query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)

	var paitentBookedAppointments []models.BookedAppointment
	for results.Next(ctx) {
		var singleAppointment models.BookedAppointment
		if err = results.Decode(&singleAppointment); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		paitentBookedAppointments = append(paitentBookedAppointments, singleAppointment)
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"patient_appointments": paitentBookedAppointments}},
	)
}
