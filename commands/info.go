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
	embed.SetColor(0x5ff5d2)
	embed.SetTitle("Bot Information")
	embed.SetDescription("Below you'll find some information about the current bot instance:")
	embed.AddField("Processor", cpu)
	embed.AddField("Isolation", isolationType)
	embed.AddField("Operating System", fmt.Sprintf("`%s`", runtime.GOOS))
	embed.AddField("Architecture", fmt.Sprintf("`%s`", runtime.GOARCH))
	embed.AddField("Server Count", fmt.Sprintf("%d server(s)", len(s.State.Guilds)))
	embed.AddField("Git Revision", fmt.Sprintf("`%s`", utils.ReadGitRevision()))
	embed.SetFooter("Nayuki", s.State.User.AvatarURL("128"))
	embed.Timestamp = time.Now().Format(time.RFC3339)
	embed.InlineAllFields()

	return utils.ReplyEmbed(s, i, embed.MessageEmbed)
}
