package command

import (
	"trup/ctx"
	"trup/db"
)

const (
	syncUsage = "sync <showcase>"
)

func sync(ctx *ctx.MessageContext, args []string) {
	_ = ctx.Session.ChannelTyping(ctx.Message.ChannelID)

	if len(args) < 1 {
		ctx.ReportUserError("Usage: " + syncUsage)
		return
	}

	if args[1] == "showcase" {
		ctx.Reply("Syncing... This may take a while")

		var entries []db.ShowcaseEntry

		var beforeID string
		for {
			msgs, err := ctx.Session.ChannelMessages(ctx.Env.ChannelShowcase, 100, beforeID, "", "")
			if err != nil {
				panic(err)
			}

			if len(msgs) == 0 {
				break
			}

			for _, msg := range msgs {
				createDate, err := msg.Timestamp.Parse()
				if err != nil {
					ctx.ReportError("Aborting, can't get timestamp from message id "+ctx.Message.ID, err)
					return
				}

				var score int
				for _, r := range msg.Reactions {
					if r.Emoji.Name == "❤" {
						score = r.Count
						break
					}
				}

				entries = append(entries, db.ShowcaseEntry{
					MessageID:  msg.ID,
					UserID:     msg.Author.ID,
					Score:      score,
					CreateDate: createDate,
				})
			}

			beforeID = msgs[len(msgs)-1].ID
		}

		err := db.AddShowcaseEntries(entries)
		if err != nil {
			ctx.ReportError("Failed to add showcase entries to the database", err)
			return
		}

		ctx.Reply("Success")
		return
	}

	ctx.ReportUserError("Usage: " + syncUsage)
}
