package handlers

import (
	"github.com/sujayramesh/SMS/model"
	"github.com/sujayramesh/SMS/utils"

	"net/http"
	"sync"
	"testing"
)

func TestSessionCreation(t *testing.T) {
	SessionList.Init()
	defer SessionList.ResetSessionMap()

	status := http.StatusOK

	uuid := utils.GenerateUUID()
	if err := SessionList.CreateSessionEntry(uuid, 100); err != nil {
		t.Fatal("Create session error: ", err)
	}

	var ActiveSessions []model.SessionDetails
	SessionList.GetSessionEntries(&ActiveSessions)
	if len(ActiveSessions) == 0 {
		t.Fatal("Session map is empty.")
	}
	ActiveSessions = nil

	req := model.SessionDetails{UUID: uuid, Duration: 150}
	if status = SessionList.ModifySessionEntry(&req); status != http.StatusOK {
		t.Fatal("modify session status: ", status)
	}

	if status = SessionList.DeleteSessionEntry(uuid); status != http.StatusOK {
		t.Fatal("Delete session status: ", status)
	}

}

func TestSessionCreationParallel(t *testing.T) {

	SessionList.Init()
	defer SessionList.ResetSessionMap()

	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		uuid := utils.GenerateUUID()
		go func(i int, uuid string, wg *sync.WaitGroup) {
			defer wg.Done()
			_ = SessionList.CreateSessionEntry(uuid, 100+i)
		}(i, uuid, &wg)
	}
	wg.Wait()

	var ActiveSessions []model.SessionDetails
	SessionList.GetSessionEntries(&ActiveSessions)

	if len(ActiveSessions) != 100 {
		t.Fatal("Session map does not have 10 entries .")
	}
	ActiveSessions = nil
}
