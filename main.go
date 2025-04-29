package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"techgentsia-bot/gl"
	"techgentsia-bot/tl"
	"time"
)

const (
	BotToken    string = "BOT_TOKEN"
	TLUserName  string = "TL_USERNAME"
	TLApiToken  string = "TL_API_TOKEN"
	GitlabToken string = "GITLAB_TOKEN"
)

func main() {
	log.Println("Server started...")
	checkContainsValidEnvironment()
	go initServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigChan:
			log.Printf("Received signal: %v. Exiting...", sig)
			return
		}
	}
}

func initServer() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv(BotToken))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized telegram account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	logger := tl.NewTL(bot, os.Getenv(TLUserName), os.Getenv(TLApiToken))
	gitlab := gl.NewGitLab(os.Getenv(GitlabToken), gl.GetProjectConfig())
	authenticated := false
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			authenticated = false
		}
	}()

	for update := range updates {
		if update.Message != nil {
			if keyMatched, _ := regexp.MatchString("^PIN-", update.Message.Text); keyMatched {
				pin := strings.Split(update.Message.Text, "PIN-")[1]
				if os.Getenv("LOGIN_PIN") == pin {
					logger.SendMessage("✅ PIN accepted. Please choose an option", update.Message.Chat.ID, update.Message.MessageID)
					authenticated = true
				} else {
					logger.SendMessage("❌ Incorrect PIN. Please try again", update.Message.Chat.ID, update.Message.MessageID)
				}
				continue
			}
			if !authenticated {
				logger.SendMessage("🔐 Please enter your PIN", update.Message.Chat.ID, update.Message.MessageID)
				continue
			}
			if matched, _ := regexp.MatchString("^B-", update.Message.Text); matched {
				projectID := strings.Split(update.Message.Text, "B-")[1]
				branches := gitlab.Branches(projectID)
				var msg []string
				for _, b := range branches {
					msg = append(msg, b.Name)
				}
				msgString := strings.Join(msg, "\n")
				if msgString == "" {
					logger.SendMessage("No branches found.", update.Message.Chat.ID, update.Message.MessageID)
					continue
				}
				logger.SendMessage(msgString, update.Message.Chat.ID, update.Message.MessageID)
			} else {
				switch update.Message.Text {
				case "P":
					{
						projects := gitlab.MyProjects()
						var msg []string
						for _, project := range projects {
							msg = append(msg, strconv.Itoa(project.Id)+"-"+project.Name)
						}
						msgString := strings.Join(msg, "\n")
						if msgString == "" {
							logger.SendMessage("No projects found.", update.Message.Chat.ID, update.Message.MessageID)
							continue
						}
						logger.SendMessage(msgString, update.Message.Chat.ID, update.Message.MessageID)
					}
				case "C":
					{
						loc, err := time.LoadLocation("Asia/Kolkata")
						if err != nil {
							panic(err)
						}
						now := time.Now().In(loc)
						startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
						endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-1), loc)
						commits, _ := gitlab.Commits(startOfDay, endOfDay)
						var msg []string
						for _, commit := range commits {
							msg = append(msg, commit.Title)
						}

						msgString := strings.Join(msg, "\n")
						if msgString == "" {
							logger.SendMessage("No commits found.", update.Message.Chat.ID, update.Message.MessageID)
							continue
						}
						logger.SendMessage(msgString, update.Message.Chat.ID, update.Message.MessageID)
					}
				case "Log":
					loc, err := time.LoadLocation("Asia/Kolkata")
					if err != nil {
						panic(err)
					}
					now := time.Now().In(loc)
					startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
					endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Nanosecond*time.Second-1), loc)
					commits, err := gitlab.Commits(startOfDay, endOfDay)
					if err != nil {
						logger.SendMessage(err.Error(), update.Message.Chat.ID, update.Message.MessageID)
						continue
					}
					var msg []string
					for _, commit := range commits {
						msg = append(msg, commit.Title)
					}

					msgString := strings.Join(msg, "\n")
					logger.LogToday(update, msgString)
				case "Hello":
					logger.SendMessage("Hello there! Yes, I’m alive and listening. How can I help you today?", update.Message.Chat.ID, update.Message.MessageID)
				case "/help":
					msg := `Available Commands:
					1. P – List all projects with ID.
					2. B-{ProjectID} – List all branches with ProjectID.
					3. C – List all commits from the project you configured.
					4. Log – Log your activity for today.
					6. Hello – Check if I’m alive!`
					logger.SendMessage(msg, update.Message.Chat.ID, update.Message.MessageID)
				default:
					logger.SendMessage("Oops! I didn’t recognize that command. Try /help to see what I can do!", update.Message.Chat.ID, update.Message.MessageID)
				}
			}
		}
	}
}

func checkContainsValidEnvironment() {
	requiredVars := []string{BotToken, TLUserName, TLApiToken, GitlabToken}
	for _, env := range requiredVars {
		if os.Getenv(env) == "" {
			panic(fmt.Sprintf("%s not configured.", env))
		}
	}
}
