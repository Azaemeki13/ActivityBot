package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func loadBot(token string) (*discordgo.Session, string, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, "", err
	}
	// s is server socket ?  ;; identify holds communication infos ;;  Intent is what I want to receive
	// Guild is server in discord, need guild info if i want roles and server basic infos
	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers
	err = session.Open()
	if err != nil {
		return nil, "", fmt.Errorf("error opening socket connection: %w", err)
	}
	fmt.Println("Connection ready !")
	guildID, err := fetchGuildID(session)
	if err != nil {
		return nil, "", err
	}
	return session, guildID, nil
}

func main() {
	token := os.Getenv("ALBION_HELPER")
	config := "./config.json"
	if token == "" {
		fmt.Println(":warning: No bot token found in env.")
		return
	}
	cfg, err := loadConfig(config)
	if err != nil {
		fmt.Println("Error, failed to load config:", err)
		return
	}
	fmt.Println("Loaded config successfully !")
	session, guildID, err := loadBot(token)
	if err != nil {
		log.Fatal(err)
		defer session.Close()
	}
	discordRoles, err := serverConfig(session, (guildID))
	if err != nil {
		fmt.Printf("No roles on the server")
	}
	botMax, err := botMaxPosition(session, guildID, discordRoles)
	if err != nil {
		log.Fatal(err)
	}
	diag := validateRoles(discordRoles, cfg, botMax)
	printDiagnostics(diag, discordRoles)
	fmt.Printf("Loaded %d classes from %s\n", len(cfg.Classes), config) // to get rid of
}
