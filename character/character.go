package character

import "github.com/nlopes/slack"

type Character interface {
	Respond(*slack.MessageEvent) string
}
