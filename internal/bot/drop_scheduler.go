package bot

import (
	"fmt"
	"time"
)

func (b *Bot) startDropScheduler() {
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		for range ticker.C {
			fmt.Println("Starting scheduled character drop...")
		}
	}()
}