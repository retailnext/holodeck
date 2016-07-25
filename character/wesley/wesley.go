package wesley

import (
	"container/list"
	"fmt"

	"github.com/nlopes/slack"
	"github.com/retailnext/holodeck/character"
)

func New(c *slack.Client) character.Character {
	return &wesley{
		c:              c,
		recentMessages: make(map[string]*list.List),
	}
}

type wesley struct {
	c              *slack.Client
	recentMessages map[string]*list.List
}

func (w *wesley) Respond(msgEvent *slack.MessageEvent) string {
	msg := msgEvent.Msg

	if msg.SubType == "message_changed" {
		msg = *msgEvent.SubMessage
	}

	l := w.recentMessages[msg.User]
	if l == nil {
		l = list.New()
		w.recentMessages[msg.User] = l
	}

	var ret string

	if msgEvent.Msg.SubType == "message_changed" {
		for e := l.Front(); e != nil; e = e.Next() {
			prevMsg := e.Value.(slack.Msg)
			if prevMsg.Timestamp == msg.Timestamp && prevMsg.Text != msg.Text {
				ret = fmt.Sprintf("<@%s> edited `%s` to `%s`! I'm reporting this to the Captain!", msg.User, prevMsg.Text, msg.Text)
				break
			}
		}
	}

	l.PushFront(msg)
	if l.Len() > 100 {
		l.Remove(l.Back())
	}

	return ret
}
