package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Bot struct {
	activeTasks   map[string]*time.Ticker
	db            *gorm.DB
	purgeInterval time.Duration
	maxDuration   time.Duration
	minDuration   time.Duration
}

type Task struct {
	ChannelID            string `gorm:"primaryKey"`
	PurgeDurationSeconds int
}

func main() {
	bot := NewBot()
	if err := bot.Run(); err != nil {
		log.Fatal("Error running bot: ", err)
	}
}

func NewBot() *Bot {
	return &Bot{
		activeTasks:   make(map[string]*time.Ticker),
		purgeInterval: 33 * time.Second,
		maxDuration:   3333 * 24 * time.Hour,
		minDuration:   30 * time.Second,
	}
}

func (b *Bot) Run() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	discordKey := os.Getenv("DISCORD_KEY")
	if discordKey == "" {
		return fmt.Errorf("DISCORD_KEY is not set")
	}

	var err error
	b.db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := b.db.AutoMigrate(&Task{}); err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	dg, err := discordgo.New("Bot " + discordKey)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	dg.AddHandler(b.messageCreate)
	dg.AddHandler(b.ready)

	if err := dg.Open(); err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is ready")
	fmt.Printf("Logged in as: %s\n", s.State.User.Username)

	var tasks []Task
	if err := b.db.Find(&tasks).Error; err != nil {
		log.Println("Error querying tasks:", err)
		return
	}

	for _, task := range tasks {
		channelID := task.ChannelID
		duration := time.Duration(task.PurgeDurationSeconds) * time.Second

		// Fetch channel directly from Discord API
		channel, err := s.Channel(channelID)
		if err != nil {
			log.Printf("Error fetching channel %s: %v", channelID, err)
			b.deleteTaskDB(channelID)
			continue
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			log.Printf("Channel %s is not a text channel", channelID)
			b.deleteTaskDB(channelID)
			continue
		}

		b.setPurgeTaskLoop(s, channelID, duration)
	}
}

func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Printf("Received message from %s: %s", m.Author.ID, m.Content)

	if strings.HasPrefix(m.Content, "<@") && strings.Contains(m.Content, s.State.User.ID) {
		log.Println("Bot mentioned in message")

		if !b.isAdminOrOwner(s, m.GuildID, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You need to be an administrator or the owner of the guild to use this bot.")
			return
		}

		args := strings.Fields(m.Content)
		if len(args) < 2 {
			log.Println("Insufficient arguments")
			return
		}

		command := strings.ToLower(args[1])
		channelID := m.ChannelID

		botMention := fmt.Sprintf("<@%s>", s.State.User.ID)

		switch command {
		case "help":
			helpMessage := fmt.Sprintf(`
**AVAILABLE COMMANDS**

purge old messages:
%s 30s
%s 5m
%s 24h
%s 2d
or any custom duration

stop purge task:
%s stop

list purge tasks:
%s list

get help:
%s help`, botMention, botMention, botMention, botMention, botMention, botMention, botMention)
			s.ChannelMessageSend(m.ChannelID, helpMessage)
		case "stop":
			b.stopAndDeleteTask(channelID)
			s.ChannelMessageSend(m.ChannelID, "Purging stopped.")
		case "list":
			b.listPurgeTasks(s, m.GuildID, m.ChannelID)
		default:
			duration, err := b.parseDuration(command)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Invalid duration. Type @bot help for available commands.")
				return
			}
			b.setPurgeTaskLoop(s, channelID, duration)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Messages older than %s will be deleted on a rolling basis in this channel.", formatDuration(duration)))
		}
	}
}

func (b *Bot) isAdminOrOwner(s *discordgo.Session, guildID, userID string) bool {
	// Fetch member directly if state cache is not available
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		log.Println("Error fetching member from API:", err)
		return false
	}

	// Fetch guild directly if state cache is not available
	guild, err := s.Guild(guildID)
	if err != nil {
		log.Println("Error fetching guild from API:", err)
		return false
	}

	// Check if the user is the owner
	if guild.OwnerID == userID {
		return true
	}

	// Check if the user has the Administrator permission
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			log.Println("Error fetching role from API:", err)
			continue
		}
		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			return true
		}
	}

	return false
}

func (b *Bot) parseDuration(input string) (time.Duration, error) {
	re := regexp.MustCompile(`^(\d+)([smhd])$`)
	match := re.FindStringSubmatch(input)
	if len(match) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	num, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("error parsing number: %v", err)
	}

	fmt.Printf("Parsed number: %d, unit: %s\n", num, match[2])

	switch match[2] {
	case "s":
		return time.Duration(num) * time.Second, nil
	case "m":
		return time.Duration(num) * time.Minute, nil
	case "h":
		return time.Duration(num) * time.Hour, nil
	case "d":
		return time.Duration(num) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid duration unit: %s", match[2])
	}
}

func (b *Bot) setPurgeTaskLoop(s *discordgo.Session, channelID string, duration time.Duration) {
	if duration < b.minDuration {
		duration = b.minDuration
	} else if duration > b.maxDuration {
		duration = b.maxDuration
	}

	fmt.Printf("Setting purge task for channel %s with duration %v\n", channelID, duration)

	b.stopTask(channelID)
	ticker := time.NewTicker(b.purgeInterval)
	b.activeTasks[channelID] = ticker

	go func() {
		for range ticker.C {
			b.purgeChannel(s, channelID, duration)
		}
	}()

	b.updateTaskDB(channelID, int(duration.Seconds()))
}

func (b *Bot) stopTask(channelID string) {
	if ticker, ok := b.activeTasks[channelID]; ok {
		ticker.Stop()
		delete(b.activeTasks, channelID)
	}
}

func (b *Bot) stopAndDeleteTask(channelID string) {
	b.stopTask(channelID)
	b.deleteTaskDB(channelID)
}

func (b *Bot) purgeChannel(s *discordgo.Session, channelID string, duration time.Duration) {
	var lastMessageID string
	threshold := time.Now().Add(-duration)

	for {
		messages, err := s.ChannelMessages(channelID, 100, lastMessageID, "", "")
		if err != nil {
			log.Println("Error fetching messages:", err)
			return
		}

		if len(messages) == 0 {
			break
		}

		for _, msg := range messages {
			// log.Printf("Checking message %s from %s, timestamp: %s", msg.ID, msg.Author.ID, msg.Timestamp)

			if msg.Timestamp.Before(threshold) {
				err = s.ChannelMessageDelete(channelID, msg.ID)
				if err != nil {
					log.Printf("Error deleting message %s: %v", msg.ID, err)
				} else {
					log.Printf("Deleted message %s", msg.ID)
				}
			}
		}

		lastMessageID = messages[len(messages)-1].ID
	}
}

func (b *Bot) updateTaskDB(channelID string, durationSeconds int) {
	task := Task{ChannelID: channelID, PurgeDurationSeconds: durationSeconds}
	if err := b.db.Save(&task).Error; err != nil {
		log.Println("Error updating database:", err)
	}
}

func (b *Bot) deleteTaskDB(channelID string) {
	if err := b.db.Delete(&Task{}, "channel_id = ?", channelID).Error; err != nil {
		log.Println("Error deleting from database:", err)
	}
}

func (b *Bot) listPurgeTasks(s *discordgo.Session, guildID, channelID string) {
	var tasks []Task
	if err := b.db.Find(&tasks).Error; err != nil {
		log.Println("Error querying tasks:", err)
		s.ChannelMessageSend(channelID, "Error retrieving tasks.")
		return
	}

	var sb strings.Builder
	for _, task := range tasks {
		ch, err := s.State.Channel(task.ChannelID)
		if err != nil {
			log.Println("Error fetching channel:", err)
			continue
		}
		if ch.GuildID == guildID {
			fmt.Println(task.PurgeDurationSeconds)
			duration := time.Duration(task.PurgeDurationSeconds) * time.Second
			sb.WriteString(fmt.Sprintf("<#%s>, duration: %s\n", ch.ID, formatDuration(duration)))
		}
	}

	if sb.Len() == 0 {
		s.ChannelMessageSend(channelID, "No purge tasks found for this guild.")
	} else {
		s.ChannelMessageSend(channelID, sb.String())
	}
}

func formatDuration(duration time.Duration) string {
	if duration%(24*time.Hour) == 0 {
		days := duration / (24 * time.Hour)
		if days > 1 {
			return fmt.Sprintf("%d days", days)
		}
		return "1 day"
	}
	if duration%time.Hour == 0 {
		hours := duration / time.Hour
		if hours > 1 {
			return fmt.Sprintf("%d hours", hours)
		}
		return "1 hour"
	}
	if duration%time.Minute == 0 {
		minutes := duration / time.Minute
		if minutes > 1 {
			return fmt.Sprintf("%d minutes", minutes)
		}
		return "1 minute"
	}
	seconds := duration / time.Second
	if seconds > 1 {
		return fmt.Sprintf("%d seconds", seconds)
	}
	return "1 second"
}
