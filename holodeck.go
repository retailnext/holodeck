package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
	"github.com/retailnext/holodeck/character"
	"github.com/retailnext/holodeck/character/wesley"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Must specify API key")
		os.Exit(1)
	}

	api := slack.New(os.Args[1])

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	fmt.Println("starting loop")

	characters := []character.Character{
		wesley.New(api),
	}

	c, err := api.JoinChannel("eng-cloud")
	fmt.Println(c, err)

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				fmt.Printf("slack hello: %v\n", ev)
				// Ignore hello

			case *slack.MessageEvent:
				for _, char := range characters {
					if msg := char.Respond(ev); msg != "" {
						rtm.SendMessage(rtm.NewOutgoingMessage(msg, ev.Channel))
					}
				}

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:

				// Ignore other events..
				fmt.Printf("Unexpected: %+v\n", msg.Data)
			}
		}
	}

}
