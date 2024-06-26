package server

import (
	"botpull/configs"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type MyBot struct {
	Token string
}

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutting down or not")
)

type IDiscordServer interface {
	Start()
}

type discordServer struct {
	cfg configs.IConfig
	dg  *discordgo.Session
}

func NewDiscordServer(cfg configs.IConfig) IDiscordServer {
	dg, err := discordgo.New("Bot " + cfg.App().GetToken())
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	return &discordServer{
		dg:  dg,
		cfg: cfg,
	}
}

// Start the server
func (s *discordServer) Start() {

	// Register the messageCreate func as a callback for MessageCreate events.
	s.dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	s.dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err := s.dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	s.dg.Close()

	// Create a ticker that ticks every 5 seconds
	// ticker := time.NewTicker(60 * time.Second)
	// defer ticker.Stop()

	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, os.Interrupt)
	// log.Println("Press Ctrl+C to exit")

	// Run a loop to wait for events
	// for {
	// 	select {
	// 	case <-stop:
	// 		if *RemoveCommands {
	// 			log.Println("Removing commands...")
	// 			for _, v := range registeredCommands {
	// 				err := s.dg.ApplicationCommandDelete(s.dg.State.User.ID, *GuildID, v.ID)
	// 				if err != nil {
	// 					log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
	// 				}
	// 			}
	// 		}
	// 		log.Println("Gracefully shutting down.")
	// 		return
	// 	case <-ticker.C:
	// 		// Run your periodic task here
	// 		fmt.Println("run time: ", ticker.C)
	// 		// go s.messageHandler("keng", "1245596903765184592")
	// 	}
	// }
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// fmt.Println("Message received")
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Printf("Message content: %s\n", m.Content) // ตรวจสอบเนื้อหาของข้อความ
	fmt.Printf("Author: %s\n", m.Author.GlobalName)

	// If the message is "!voice", list members in the voice channel of the same guild
	if m.Content == "!voice" {
		guildID := m.GuildID
		fmt.Printf("Guild ID from Message: %s\n", guildID) // พิมพ์ Guild ID ที่ได้รับจากข้อความ

		guild, err := s.State.Guild(guildID)
		if err != nil {
			fmt.Println("Error getting guild:", err)
			return
		}
		fmt.Printf("Guild ID: %s, Guild Name: %s\n", guild.ID, guild.Name) // พิมพ์ข้อมูลของกิลด์

		fmt.Printf("Number of Voice States: %d\n", len(guild.VoiceStates))
		for i, vs := range guild.VoiceStates {
			fmt.Printf("VoiceState %d: UserID=%s, ChannelID=%s\n", i, vs.UserID, vs.ChannelID)
		}

		voiceMembers := make(map[string][]string)
		for _, vs := range guild.VoiceStates {
			user, err := s.User(vs.UserID)
			if err != nil {
				fmt.Println("Error getting user:", err)
				continue
			}
			fmt.Printf("User found: %s (ID: %s)\n", user.Username, user.ID)
			voiceMembers[vs.ChannelID] = append(voiceMembers[vs.ChannelID], user.Username)
		}

		response := "Members in voice channels:\n"
		for channelID, members := range voiceMembers {
			channel, err := s.Channel(channelID)
			if err != nil {
				fmt.Println("Error getting channel:", err)
				continue
			}
			response += fmt.Sprintf("Channel: %s\n", channel.Name)
			for _, member := range members {
				response += fmt.Sprintf(" - %s\n", member)
			}
		}

		fmt.Printf("Response: %s\n", response) // เพิ่มการพิมพ์ response ก่อนส่งข้อความ
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
