package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var discord_token string
var hex_cache [32]byte
var we_are_happy = false
var we_have_printed = false

const guthib_url = "https://api.github.com/users/sukus21/repos"
const pattern = "(?i)spin"

type GitHubRepo struct {
	Name string `json:"name"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	discord_token = os.Getenv("DISCORD_TOKEN")

	discord, err := discordgo.New("Bot " + discord_token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	discord.AddHandler(discordEventHandler)

	err = discord.Open()
	defer discord.Close()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}

	c := cron.New()
	_, err = c.AddFunc("@every 5s", cronJob)
	if err != nil {
		fmt.Println("Error scheduling cron job:", err)
		return
	}
	_, err = c.AddFunc("@every 5s", func() {
		if we_are_happy && !we_have_printed {
			discord.ChannelMessageSend("1445386351225602150", "SUKUS RELASED SPIN!!!")
			we_have_printed = true
		}
	})
	if err != nil {
		fmt.Println("Error scheduling cron job:", err)
		return
	}
	c.Start()

	// Keep until terminated by Ctrl+C
	fmt.Println("Bot running....")
	inter := make(chan os.Signal, 1)
	signal.Notify(inter, os.Interrupt)
	<-inter
	fmt.Println("Shutting down...")
}

func cronJob() {
	repos, err := fetchGitHubRepos()
	if err != nil {
		fmt.Println("Error fetching GitHub repos in cron job:", err)
		return
	}

	reposStr := strings.Join(repos, ",")
	hash := sha256.Sum256([]byte(reposStr))
	if hash != hex_cache {
		fmt.Println("Change detected in GitHub repos")
		re, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
			return
		}

		for _, repo := range repos {
			if re.MatchString(repo) {
				we_are_happy = true
			}
		}

		hex_cache = hash
	}
}

func fetchGitHubRepos() ([]string, error) {
	// Fetch GitHub repos
	resp, err := http.Get(guthib_url)
	if err != nil {
		fmt.Println("Error fetching GitHub repos:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var repos []GitHubRepo
	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		fmt.Println("Error decoding GitHub response:", err)
		return nil, err
	}

	var repoNames []string
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Name)
	}
	return repoNames, nil
}

func discordEventHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {}
