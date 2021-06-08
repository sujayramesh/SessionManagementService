package handlers

import (
	"github.com/sujayramesh/SMS/model"

	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type SessionParameters struct {
	StartTime     time.Time
	EndTime       int64
	TimerDuration int
}

// Main map
type SessionMapType map[string]SessionParameters

type SessionData struct {
	ActiveSessionMap SessionMapType
}

//Method to init Session holder data structure.
func (sessionData *SessionData) Init() {
	sessionData.ActiveSessionMap = make(SessionMapType)
}

//Method to init Session holder data structure.
func (sessionData *SessionData) ResetSessionMap() {
	sessionData.ActiveSessionMap = nil
}

// Method to create session entry
func (sessionData *SessionData) CreateSessionEntry(uuid string, duration int) error {

	mutex.Lock()
	var err error
	if _, ok := sessionData.ActiveSessionMap[uuid]; ok {
		fmt.Println("UUID: ", uuid, " already present in active session data")

		// Set appropriate error
		err = errors.New("session could not be stored")

	} else {
		fmt.Println("UUID not present. Creating entry.")
		sp := SessionParameters{StartTime: time.Now(), EndTime: time.Now().Unix() + int64(duration), TimerDuration: duration}
		sessionData.ActiveSessionMap[uuid] = sp
		err = nil
	}
	mutex.Unlock()
	return err
}

// Method to delete aforementioned session
func (sessionData *SessionData) DeleteSessionEntry(uuid string) int {

	ret := http.StatusNotFound
	currUnixTime := time.Now().Unix()

	mutex.Lock()
	if val, ok := sessionData.ActiveSessionMap[uuid]; ok {
		// Delete only if the session is active.
		if val.EndTime > currUnixTime {
			fmt.Println("UUID: ", uuid, " present. Deleting the session.")
			delete(sessionData.ActiveSessionMap, uuid)
			ret = http.StatusOK
		}
	}
	mutex.Unlock()
	return ret
}

// Method to list or retrieve active session entries.
func (sessionData *SessionData) GetSessionEntries(activeSessions *[]model.SessionDetails) {

	// Current time in number of seconds
	currUnixTime := time.Now().Unix()

	//Lock the map for concurrent operations.
	mutex.Lock()
	for uuid, sp := range sessionData.ActiveSessionMap {
		var activeSample model.SessionDetails
		activeSample.UUID = uuid
		// End time - current time is the remaining duration of the session.
		activeSample.Duration = int(sp.EndTime - currUnixTime)

		// Ensuring expired sessions are not included in the list
		if activeSample.Duration >= 0 {
			*activeSessions = append(*activeSessions, activeSample)
		}
	}
	mutex.Unlock()
}

// Method to Modify session entry
func (sessionData *SessionData) ModifySessionEntry(req *model.SessionDetails) int {

	ret := http.StatusNotFound
	currUnixTime := time.Now().Unix()

	mutex.Lock()
	if val, ok := sessionData.ActiveSessionMap[req.UUID]; ok {

		if val.EndTime > currUnixTime {
			// Session active. Allow modify.
			fmt.Println("UUID: ", req.UUID, " present. Extending the session.")
			sp := SessionParameters{StartTime: time.Now(), EndTime: time.Now().Unix() + int64(req.Duration), TimerDuration: req.Duration}
			sessionData.ActiveSessionMap[req.UUID] = sp
			ret = http.StatusOK
		}
	}
	mutex.Unlock()
	return ret
}

// Method to delete expired session entries
func (sessionData *SessionData) DeleteExpiredSessions() {

	// Current time in number of seconds
	currUnixTime := time.Now().Unix()

	//Lock the map for concurrent operations.
	mutex.Lock()
	for uuid, sp := range sessionData.ActiveSessionMap {

		// Expired sessions time will be higher than the original session EndTime.
		if currUnixTime > sp.EndTime {
			fmt.Println("DeleteExpiredSessions: Session with UUID: ", uuid, " has expired.")
			delete(sessionData.ActiveSessionMap, uuid)
		}
	}
	mutex.Unlock()
}

var SessionList SessionData
