package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/azurejelly/nayuki/utils"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/klauspost/cpuid/v2"
)

type InfoCommand struct{}

func (c *InfoCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "info",
		Description: "Provides potentially useful information.",
	}
}

func (c *InfoCommand) Run(s *discordgo.Session, event *discordgo.InteractionCreate) error {
	i := event.Interaction

	cpu := cpuid.CPU.BrandName
	if cpu == "" {
		cpu = "Unknown"
	}

	var isolationType string

	if utils.IsDockerContainer() {
		isolationType = "Running on Docker"
	} else if !cpuid.CPU.VM() {
		isolationType = "None"
	} else {
		isolationType = cpuid.CPU.HypervisorVendorString

		if isolationType == "" {
			isolationType = "Unknown Virtual Machine"
		}
	}

	embed := embed.NewEmbed()
	embed.SetColor(utils.DEFAULT_EMBED_COLOR)
	embed.SetTitle("Information")
	embed.SetDescription("Below you'll find some information about the current bot instance:")
	embed.AddField(":computer: Processor", cpu)
	embed.AddField(":package: Isolation", isolationType)
	embed.AddField(":gear: Operating System", fmt.Sprintf("`%s` on `%s`", runtime.GOOS, runtime.GOARCH))
	embed.AddField(":beginner: Server count", fmt.Sprintf("%d server(s)", len(s.State.Guilds)))
	embed.AddField(":link: Source code", fmt.Sprintf("Available on [`GitHub`](%s)", utils.GITHUB_REPOSITORY))
	embed.AddField(":label: Revision", fmt.Sprintf("`%s`", utils.ReadGitRevision()))
	embed.SetFooter("Nayuki", s.State.User.AvatarURL("128"))
	embed.Timestamp = time.Now().Format(time.RFC3339)
	embed.InlineAllFields()

	return utils.ReplyEmbed(s, i, embed.MessageEmbed)
}
