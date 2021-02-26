package command

import (
	"fmt"
	"log"
	"strings"
	"trup/ctx"
	"trup/db"
	"trup/misc"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

const warnUsage = "warn <@user> <reason>"

func warn(ctx *ctx.MessageContext, args []string) {
	if len(args) < 3 {
		ctx.ReportUserError("Usage: " + warnUsage)
		return
	}

	user := misc.ParseUser(args[1])
	if user == "" {
		ctx.ReportUserError("The first argument must be a user mention")
		return
	}

	reason := strings.Join(args[2:], " ")
	warningMessageLink := fmt.Sprintf("[(warning)](%s)", makeMessageLink(ctx.Message.GuildID, ctx.Message))

	err := ctx.RequestUserByName(true, args[1], func(m *discordgo.Member) error {
		user := m.User.ID

		w := db.NewWarn(ctx.Message.Author.ID, user, reason)
		err := w.Save()
		if err != nil {
			ctx.ReportError("Failed to save your warning", err)
			return nil
		}

		var nth string
		warnCount, err := db.CountWarns(user)
		if err != nil {
			log.Printf("Failed to count warns for user %s; Error: %s\n", user, err)
		}
		if warnCount > 0 {
			nth = " for the " + humanize.Ordinal(warnCount) + " time"
		}

		taker := ctx.Message.Author
		err = db.NewNote(taker.ID, user, fmt.Sprintf("User was warned for: %s %s", reason, warningMessageLink), db.ManualNote).Save()
		if err != nil {
			log.Printf("Failed to save warning note. Error: %s\n", err)
		}

		ctx.Reply(fmt.Sprintf("<@%s> Has been warned%s with reason: %s. <a:police:749871644071165974>", user, nth, reason))
		r := ""
		if reason != "" {
			r = " with reason: " + reason
		}

		_, err = ctx.Session.ChannelMessageSend(
			ctx.Env.ChannelModlog,
			fmt.Sprintf("<@%s> was warned by moderator %s%s. They've been warned%s", user, taker.Username, r, nth),
		)
		if err != nil {
			log.Printf("Error sending warning notice into modlog: %s\n", err)
		}
		return nil
	})
	if err != nil {
		ctx.ReportError("Failed to find the user", err)
		return
	}
}

func makeMessageLink(guildID string, m *discordgo.Message) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, m.ChannelID, m.ID)
}
