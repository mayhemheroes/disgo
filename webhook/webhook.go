package webhook

import (
	"context"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
)

type Webhook struct {
	discord.Webhook
	WebhookClient Client
}

func (h *Webhook) Update(webhookUpdate discord.WebhookUpdate) (*Webhook, rest.Error) {
	return h.WebhookClient.UpdateWebhook(webhookUpdate)
}

func (h *Webhook) Delete(opts ...rest.RequestOpt) rest.Error {
	return h.WebhookClient.DeleteWebhook(ctx)
}
