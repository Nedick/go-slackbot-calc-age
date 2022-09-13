package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func printCommandEvent(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	currentTime := time.Now()
	SLACK_BOT_TOKEN := goDotEnvVariable("SLACK_BOT_TOKEN")
	SLACK_APP_TOKEN := goDotEnvVariable("SLACK_APP_TOKEN")

	bot := slacker.NewClient(SLACK_BOT_TOKEN, SLACK_APP_TOKEN)

	go printCommandEvent(bot.CommandEvents())

	bot.Command("My year of birth is <year>", &slacker.CommandDefinition{
		Description: "Year of birth calculator",
		// Example:     "My year of birth is 1990",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yearOfBirth, err := strconv.Atoi(year)
			if err != nil {
				println("error")
			}
			age := currentTime.Year() - yearOfBirth
			r := fmt.Sprintf("You are %d years old", age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
