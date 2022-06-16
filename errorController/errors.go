package errorController

import (
	"github.com/DisgoOrg/disgohook"
	"github.com/DisgoOrg/disgohook/api"
  "log"
)

const (
	webhookURLWeb = "947157399368769566/ChphRmBehex1A3XiDNrQfQkDfjSHLPb1w9QC3QM0fpH-kNQVrOjAgioHHULSMtlhcyow"
)

var (
	WD Message
)

type Message struct {
	Content string
	Site    string
}

func (m *Message) SendErrorWHWeb() {
	webhook, err := disgohook.NewWebhookClientByToken(nil, nil, webhookURLWeb)
	if err != nil {
		panic(err)
	}

	_, err = webhook.SendEmbeds(api.NewEmbedBuilder().
		SetTitle("New error on the site!").
		SetDescription(m.Site + ": " + m.Content + " @here").
		SetColor(14177041).
		Build(),
	)
	if err != nil {
		log.Println(err)
	}
}
