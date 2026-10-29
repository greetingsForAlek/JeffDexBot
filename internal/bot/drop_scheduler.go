package bot

import (
	"fmt"
	"time"
)

func (b *Bot) startDropScheduler() chan struct{} {
	stop := make(chan struct{})
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Starting scheduled character drop...")
				b.sendDrop(b.Session, b.DropChannelID)

			case <- stop:
				return
			}
		}
	}()

	return stop
}