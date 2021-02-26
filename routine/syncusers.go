package routine

import (
	"log"
	"runtime/debug"
	"time"
	"trup/ctx"
	"trup/db"
)

func SyncUsersLoop(ctx *ctx.Context) {
	time.Sleep(5 * time.Minute)

	for {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panicked in SyncUsersLoop with error: %v; Stack: %s\n", err, debug.Stack())
				}
			}()

			members, err := ctx.Members()
			if err != nil {
				log.Printf("Failed to get unique members; Error: %v\n", err)
				return
			}

			err = db.AddUsers(members)
			if err != nil {
				log.Printf("Failed to add users to database; Error: %v\n", err)
			} else {
				log.Println("Successfully added users to database")
			}
		}()

		time.Sleep(24 * time.Hour)
	}
}
