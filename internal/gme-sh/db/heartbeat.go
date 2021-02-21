package db

import (
	"log"
	"sync"
	"time"
)

// LastHeartbeatError contains the last thrown error of the backend (TemporaryDatabase).
// If everything is running, this variable should be nil
var LastHeartbeatError error
var mu = &sync.Mutex{}

func createTickerAndCheck(tdb TemporaryDatabase, c chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := tdb.Heartbeat()

			// only send error if new
			if err != LastHeartbeatError {
				if err == nil || LastHeartbeatError == nil || err.Error() != LastHeartbeatError.Error() {
					if err != nil {
						log.Println("💔 Heartbeat failed:", err)
					} else {
						log.Println("💚 Heartbeat successful.")
					}
				}
			}

			mu.Lock()
			LastHeartbeatError = err
			mu.Unlock()

			break
		case <-c:
			log.Println("\U0001FAA6 RIP: Heartbeat Service stopped.")
			ticker.Stop()
			return
		}
	}
}

// CreateHeartbeatService x starts the Heartbeat service and creates a channel to end the service by sending true
func CreateHeartbeatService(tdb TemporaryDatabase) chan bool {
	c := make(chan bool, 1)
	go createTickerAndCheck(tdb, c)
	return c
}
