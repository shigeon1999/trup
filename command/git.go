package command

import (
	"trup/ctx"
	"trup/db"
	"trup/misc"

	"github.com/jackc/pgx"
)

const (
	gitUsage = "git [url]"
	gitHelp  = "Adds a git link to your fetch."
)

func git(ctx *ctx.MessageContext, args []string) {
	user := ctx.Message.Author.ID

	if len(args) == 1 {
		setItFirstMsg := "You need to set your !git url first"
		profile, err := db.GetProfile(user)
		if err != nil {
			if err.Error() == pgx.ErrNoRows.Error() {
				ctx.ReportUserError(setItFirstMsg)
				return
			} else {
				ctx.ReportError("Failed to fetch your profile", err)
				return
			}
		}

		if profile.Git == "" {
			ctx.ReportUserError(setItFirstMsg)
			return
		}

		_, _ = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, profile.Git)
		return
	}

	url := args[1]

	if !misc.IsValidURL(url) {
		ctx.ReportUserError("You need to provide a valid url")
		return
	}

	profile, err := db.GetProfile(user)
	if err != nil {
		if err.Error() != pgx.ErrNoRows.Error() {
			ctx.ReportError("Failed to fetch your profile", err)
			return
		}
		profile = db.NewProfile(user, url, "", "")
	} else {
		profile.Git = url
	}

	err = profile.Save()
	if err != nil {
		ctx.ReportError("failed to save git url", err)
		return
	}

	ctx.Reply("Success. You can run !git or !fetch to retrieve the url")
}
