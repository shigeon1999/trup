package ctx

import (
	"log"
	"time"
	"trup/misc"

	"github.com/bwmarrin/discordgo"
)

type MessageContext struct {
	Context
	Message *discordgo.Message
}

func (ctx *MessageContext) Reply(msg string) {
	_, err := ctx.ReplyEmbedSimple("", msg)
	if err != nil {
		log.Printf("Failed to reply to message %s; Error: %s\n", ctx.Message.ID, err)
	}
}

func (ctx *MessageContext) ReportError(msg string, err error) {
	log.Printf(
		"Error Message ID: %s; ChannelID: %s; GuildID: %s; Author ID: %s; msg: %s; error: %s\n",
		ctx.Message.ID,
		ctx.Message.ChannelID,
		ctx.Message.GuildID,
		ctx.Message.Author.ID,
		msg,
		err,
	)
	ctx.Reply(msg)
}

func (ctx *MessageContext) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	if embed.Color == 0 {
		embed.Color = ctx.UserColor()
	}
	if embed.Timestamp == "" {
		embed.Timestamp = time.Now().Format(misc.DiscordDateFormat)
	}
	if embed.Footer == nil {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: "Response to " + ctx.Message.Author.String() + "'s message",
		}
	}

	return ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
}

func (ctx *MessageContext) ReplyEmbedSimple(title string, description string) (*discordgo.Message, error) {
	return ctx.ReplyEmbed(&discordgo.MessageEmbed{
		Title:       title,
		Description: description,
	})
}

func (ctx *MessageContext) UserColor() int {
	return ctx.Session.State.UserColor(ctx.Message.Author.ID, ctx.Message.ChannelID)
}

func (ctx *MessageContext) IsModerator() bool {
	for _, r := range ctx.Message.Member.Roles {
		if r == ctx.Env.RoleMod {
			return true
		}
	}

	return false
}

func (ctx *MessageContext) IsHelper() bool {
	for _, r := range ctx.Message.Member.Roles {
		if r == ctx.Env.RoleHelper {
			return true
		}
	}

	return false
}
