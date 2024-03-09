package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jeff-rdg/discord_webhook/color"
	"github.com/jeff-rdg/discord_webhook/emoji"
	"io"
	"log"

	"net/http"
	"os"
	"time"
)

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

	req, _ := json.Marshal(logger.Request)

	embeds := []embeds{{
		Title:       logger.Uri,
		Description: string(req),
		Url:         logger.UriUrl,
		Color:       color.RandomColor(),
		Author: author{
			Name: logger.Application,
		},
		Timestamp: time.Now(),
	}}

	return &message{
		Content:     fmt.Sprint(emoji.Warning, " - ", logger.Enviroment),
		Embeds:      embeds,
		Attachments: nil,
	}
}

func (m *message) sendToDiscord() {
	body, err := json.Marshal(m)
	if err != nil {
		log.Fatalln("Erro ao serializar a mensagem:", err)
	}

	url := os.Getenv("WEBHOOK_URL")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln("Erro ao enviar a solicitação HTTP:", err)
	}

	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Erro ao ler a resposta do servidor:", err)
	}

	log.Printf("Status Code: %d", resp.StatusCode)
	log.Printf("Resposta do servidor: %s", result)
}

type SenderLog struct {
	Enviroment  string      `json:"enviroment"`
	Request     interface{} `json:"request"`
	Uri         string      `json:"uri"`
	UriUrl      string      `json:"uri_url"`
	Application string      `json:"application"`
}

func newLogger(env, uri, uriUrl, app string, req interface{}) SenderLog {
	return SenderLog{
		Enviroment:  env,
		Request:     req,
		Uri:         uri,
		UriUrl:      uriUrl,
		Application: app,
	}
}

func SendLog(env, uri, uriUrl, app string, req interface{}) {
	logger := newLogger(env, uri, uriUrl, app, req)

	messenger := createMessage(logger)
	messenger.sendToDiscord()
}
