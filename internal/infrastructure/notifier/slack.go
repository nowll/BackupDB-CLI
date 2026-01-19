package notifier

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	client  *slack.Client
	channel string
}

func NewSlackNotifier(token, channel string) *SlackNotifier {
	return &SlackNotifier{
		client:  slack.New(token),
		channel: channel,
	}
}

func (n *SlackNotifier) Notify(ctx context.Context, message string) error {
	_, _, err := n.client.PostMessageContext(
		ctx,
		n.channel,
		slack.MsgOptionText(message, false),
	)

	if err != nil {
		return fmt.Errorf("send slack notification: %w", err)
	}

	return nil
}
