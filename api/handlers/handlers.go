package handlers

import (
	"github.com/sujayramesh/SMS/model"
	"github.com/sujayramesh/SMS/utils"

	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler for CreateSession
func CreateSession(c echo.Context) (err error) {

	req := new(model.SessionRequest)
	var resp model.SessionResponse
	status := http.StatusOK

	if err = c.Bind(req); err != nil {
		fmt.Println("Json unmarshal failed. Error:", err)
		status = http.StatusBadRequest
		return c.JSON(status, resp)
	}

	if req.Expiry == nil {
		fmt.Println("Expiry is 0. Defaulting to 30s")
		i := 30
		req.Expiry = &i
	} else {
		fmt.Println("Expiry time present in request. Value: ", *req.Expiry)
	}

	// Generate UUID
	uid := utils.GenerateUUID()

	if err := SessionList.CreateSessionEntry(uid, *req.Expiry); err != nil {
		fmt.Println("Error encountered: ", err)
		status = http.StatusInternalServerError
	}

	return c.JSON(status, resp)
}

// Handler for DeleteSession
func DeleteSession(c echo.Context) (err error) {

	req := new(model.SessionUUID)
	var resp model.SessionResponse
	status := http.StatusOK

	if err = c.Bind(req); err != nil {
		fmt.Println("Json unmarshal failed. Error:", err)
		status = http.StatusBadRequest
		return c.JSON(status, resp)
	}

	status = SessionList.DeleteSessionEntry(req.UUID)

	return c.JSON(status, resp)
}

// Handler for ModifySession
func ModifySession(c echo.Context) (err error) {

	req := new(model.SessionDetails)
	var resp model.SessionResponse
	status := http.StatusOK

	if err = c.Bind(req); err != nil {
		fmt.Println("Json unmarshal failed. Error:", err)
		status = http.StatusBadRequest
		return c.JSON(status, resp)
	}

	if req.Duration > 300 {
		fmt.Println("Requested duration is greater than 300. Limiting to 300s")
		req.Duration = 300
	}

	status = SessionList.ModifySessionEntry(req)

	return c.JSON(status, resp)
}

// Handler for Listing all sessions
func GetSession(c echo.Context) (err error) {

	status := http.StatusOK

	var ActiveSessions []model.SessionDetails
	SessionList.GetSessionEntries(&ActiveSessions)

	return c.JSON(status, ActiveSessions)
}

// Handler for deleting expired sessions
func DeleteExpiredSessions() {
	SessionList.DeleteExpiredSessions()
}
