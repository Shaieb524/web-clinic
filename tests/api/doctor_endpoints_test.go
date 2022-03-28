package apiTesting

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert" // add Testify package
)

func TestDoctorEndpoints(t *testing.T) {
	tests := []struct {
		description    string
		method         string
		route          string
		requestBody    map[string]interface{}
		requestHeaders map[string]string
		expectedCode   int
	}{
		{
			description:    "GET to get a doctor by id",
			method:         http.MethodGet,
			route:          "/doctors/:id",
			requestHeaders: map[string]string{`Content-Type`: `application/json`},
			expectedCode:   200,
		},
	}

	app := fiber.New()
	app.Get("/doctors/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		fmt.Println("idParam : ", idParam)
		return c.SendString("Hello, World!")
	})

	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
