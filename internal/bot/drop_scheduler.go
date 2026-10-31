package bot

import (
	"fmt"
	"time"
)

// startDropScheduler function
func (b *Bot) startDropScheduler() chan struct{} {
	stop := make(chan struct{}) // stop channel
	ticker := time.NewTicker(60 * time.Minute) // new ticker to count time

	// new goroutine
	go func() {
		defer ticker.Stop()

		// loop
		for {
			select {
			case <-ticker.C:
				// start scheduled character drop
				fmt.Println("Starting scheduled character drop...")
				b.sendDrop(b.Session, b.DropChannelID)

			case <- stop:
				return // return if stopped
			}
		}
	}()

	return stop
}