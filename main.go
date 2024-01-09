package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/bwmarrin/discordgo"
)

var (
	discordToken string
	giteaToken   string
	giteaURL     string
)

func main() {

	giteaToken = os.Getenv("GITEA_TOKEN")
	giteaURL = os.Getenv("GITEA_URL")
	discordToken = os.Getenv("DISCORD_TOKEN")

	if giteaToken == "" || giteaURL == "" || discordToken == "" {
		fmt.Println("Environment variables GITEA_TOKEN, GITEA_URL, and DISCORD_TOKEN must be set.")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}

// validateUsername checks if the username complies with Gitea requirements
func validateUsername(username string) (bool, error) {
	// Regular expression for valid username
	var validUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([_.-]?[a-zA-Z0-9]+)*$`)

	// Check length constraints (1 to 39 characters)
	if len(username) < 1 || len(username) > 39 {
		return false, nil
	}

	return validUsernameRegex.MatchString(username), nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the command is !register
	if strings.HasPrefix(m.Content, "!register") {
		args := strings.Fields(m.Content)

		var username, email string

		// Loop through arguments to find username and email
		for _, arg := range args[1:] { // Skip the command itself
			if strings.HasPrefix(arg, "--username=") {
				username = strings.TrimPrefix(arg, "--username=")
			} else if strings.HasPrefix(arg, "--email=") {
				email = strings.TrimPrefix(arg, "--email=")
			}
		}

		// If no username parameter is provided, use the user's Discord nickname or username
		if username == "" {
			member, err := s.GuildMember(m.GuildID, m.Author.ID)
			if err != nil {
				username = m.Author.Username
			} else if member.Nick != "" {
				username = member.Nick
			} else {
				username = m.Author.Username
			}
		}

		// Remove trailing period from username if present
		username = strings.TrimRight(username, ".")

		// If no email parameter is provided, use the default email format
		if email == "" {
			email = fmt.Sprintf("%s@git.beisel.it", username)
		}

		// Validate the username
		isValid, err := validateUsername(username)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Invalid username. Username must consist of alphanumeric characters, dashes, underscores, and dots. It cannot begin or end with non-alphanumeric characters, and consecutive non-alphanumeric characters are forbidden.")
			return
		}
		if !isValid {
			s.ChannelMessageSend(m.ChannelID, "Invalid username. Username must consist of alphanumeric characters, dashes, underscores, and dots. It cannot begin or end with non-alphanumeric characters, and consecutive non-alphanumeric characters are forbidden.")
			return
		}

		fmt.Println("Creating user:", username)

		// Auto-generate a password
		password, err := generatePassword(12) // You can choose the length of the password
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to generate a secure password.")
			return
		}

		// Create user in Gitea
		err = createUserInGitea(username, email, password)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to create user in Gitea.")
			return
		}

		// Send DM with login details and instructions
		dmChannel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to send DM.")
			return
		}

		dmMessage := fmt.Sprintf("You can now log in to %s using:\n\nUsername: %s\nPassword: %s\n\nYour E-Mail has temporarily been set to \"%s\". Please make sure to set this to your correct e-mail address and update your password immediately.", giteaURL, username, password, email)
		_, err = s.ChannelMessageSend(dmChannel.ID, dmMessage)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to send login details via DM.")
			return
		}
	}
}

// generatePassword creates a secure, random password of a given length.
func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var password strings.Builder
	for i := 0; i < length; i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password.WriteByte(charset[charIndex.Int64()])
	}
	return password.String(), nil
}

func createUserInGitea(username, email, password string) error {
	client, _ := gitea.NewClient(giteaURL, gitea.SetToken(giteaToken))

	_, _, err := client.AdminCreateUser(gitea.CreateUserOption{
		LoginName: username,
		Username:  username,
		Email:     email,
		Password:  password,
	})

	return err
}
