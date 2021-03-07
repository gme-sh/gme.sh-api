package db

import (
	"log"
	"time"
)

type ExpirationCheck struct {
	Interval            time.Duration
	LastExpirationCheck time.Time
	DB                  PersistentDatabase
	DryRun              bool
}

func NewExpirationCheck(interval time.Duration, dryRun bool, database PersistentDatabase) *ExpirationCheck {
	return &ExpirationCheck{
		Interval: interval,
		DB:       database,
		DryRun:   dryRun,
	}
}

func (e *ExpirationCheck) Check() {
	// find
	expired, err := e.DB.FindExpiredURLs()
	if err != nil {
		log.Println("WARN: Error checking for expiration:", err)
		return
	}
	for _, ex := range expired {
		log.Println("💔 Would delete expired url ::", *ex)
		if !e.DryRun {
			if err := e.DB.DeleteShortenedURL(&ex.ID); err != nil {
				log.Println("⚠️ Error deleting expired url #", ex.ID, ":", err)
				continue
			}
		}
	}
}

func (e *ExpirationCheck) Start(cancel chan bool) {
	t := time.NewTicker(e.Interval)
	for {
		select {
		case <-cancel:
			log.Println("(Cancel) cancelled expiration check")
			return
		case <-t.C:
			///
			// check database for last expiration
			check := e.DB.GetLastExpirationCheck()
			sub := time.Now().Sub(check.LastCheck.Add(-2 * time.Second)) // 2s grace
			if sub <= e.Interval {
				break
			}
			e.DB.UpdateLastExpirationCheck(time.Now().Add(-2 * time.Second)) // now + 2s grace
			///

			e.Check()
			break
		}
	}
}
