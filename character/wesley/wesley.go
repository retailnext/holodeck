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
		userNames:      make(map[string]string),
	}
}

type wesley struct {
	c              *slack.Client
	recentMessages map[string]*list.List
	userNames      map[string]string
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
				name := w.userNames[msg.User]
				if name == "" {
					user, err := w.c.GetUserInfo(msg.User)
					if err != nil {
						fmt.Printf("error fetching user info for %s: %s", msg.User, err)
					}

					if user != nil {
						name = user.Name
						w.userNames[msg.User] = name
					} else {
						name = msg.User
					}
				}
				ret = fmt.Sprintf("@%s edited %q to %q! I'm reporting this to the Captain!", name, prevMsg.Text, msg.Text)
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
