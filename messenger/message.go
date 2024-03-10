package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Jeff-Rdg/webhook_discord/color"
	"github.com/Jeff-Rdg/webhook_discord/emoji"
	"io"
	"log"

	"net/http"
	"time"
)

// message JSON sending to Discord
type message struct {
	Content     string   `json:"content"`
	Embeds      []embeds `json:"embeds"`
	Attachments []string `json:"attachments"`
}

type embeds struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Url         string             `json:"url"`
	Color       color.DiscordColor `json:"color"`
	Author      author
	Timestamp   time.Time `json:"timestamp"`
}

type author struct {
	Name string `json:"name"`
}

func createMessage(logger SenderLog) *message {

	req, _ := json.Marshal(logger.Response)

	embeds := []embeds{{
		Title:       logger.Uri,
		Description: string(req),
		Url:         logger.UriUrl,
		Color:       color.RandomColor(),
		Author: author{
			Name: logger.ApplicationName,
		},
		Timestamp: time.Now(),
	}}

	return &message{
		Content:     fmt.Sprint(emoji.Warning, " ", logger.ApplicationName, " - ", logger.Enviroment),
		Embeds:      embeds,
		Attachments: nil,
	}
}

func (m *message) sendToDiscord(url string) {
	body, err := json.Marshal(m)
	if err != nil {
		log.Fatalln("Erro ao serializar a mensagem:", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Erro ao enviar a solicitação HTTP:", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Erro ao ler a resposta do servidor:", err)
	}

	if resp.StatusCode >= 300 {
		log.Println("Houve um erro ao conectar com servidor")
	}

	log.Printf("Status Code: %d", resp.StatusCode)
}

// SenderLog Json to prepare to sending
type SenderLog struct {
	Enviroment      string      `json:"enviroment"`
	Response        interface{} `json:"response"`
	Uri             string      `json:"uri"`
	UriUrl          string      `json:"uri_url"`
	ApplicationName string      `json:"application_name"`
	WebhookUrl      string      `json:"webhook_url"`
}

// SendLog method to sending to Discord Webhook
func SendLog(logger SenderLog) {
	messenger := createMessage(logger)
	messenger.sendToDiscord(logger.WebhookUrl)
}
