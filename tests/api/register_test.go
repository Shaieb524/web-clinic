package apiTesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert" // add Testify package
)

func TestRegister(t *testing.T) {
	tests := []struct {
		description    string
		method         string
		route          string
		requestBody    map[string]interface{}
		requestHeaders map[string]string
		expectedCode   int
	}{
		{
			description: "POST endpoint to register user",
			method:      http.MethodPost,
			route:       "/register",
			requestBody: map[string]interface{}{
				"name":     "test doctor",
				"email":    "test@gmail.com",
				"password": "123",
				"role":     "doctor",
			},
			requestHeaders: map[string]string{`Content-Type`: `application/json`},
			expectedCode:   200,
		},
		{
			description: "POST endpoint to login user",
			method:      http.MethodPost,
			route:       "/login",
			requestBody: map[string]interface{}{
				"email":    "test@gmail.com",
				"password": "123",
			},
			requestHeaders: map[string]string{`Content-Type`: `application/json`},
			expectedCode:   200,
		},
	}

	app := fiber.New()
	app.Post("/register", func(c *fiber.Ctx) error {
		var requestData map[string]string
		if err := c.BodyParser(&requestData); err != nil {
			return err
		}
		return c.SendString("Hello, World!")
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var requestData map[string]string
		if err := c.BodyParser(&requestData); err != nil {
			return err
		}
		if requestData["email"] != "test@gmail" || requestData["password"] != "123" {
			fmt.Println("ifffff")
			return c.Status(http.StatusUnauthorized).SendString("Login failed!")
		} else {
			fmt.Println("elseee")
			return c.SendString("Login successfuly!")
		}
	})

	for _, test := range tests {
		// Create a new http request with the route from the test case
		rbody, _ := json.Marshal(test.requestBody)
		req := httptest.NewRequest(test.method, test.route, bytes.NewReader(rbody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
