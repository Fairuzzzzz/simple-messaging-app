package ws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Fairuzzzzz/fiber-boostrap/app/models"
	"github.com/Fairuzzzzz/fiber-boostrap/app/repository"
	"github.com/Fairuzzzzz/fiber-boostrap/pkg/env"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServeWSMessaging(app *fiber.App) {
	clients := make(map[*websocket.Conn]bool)

	broadcast := make(chan models.MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true
		for {
			var msg models.MessagePayload
			if err := c.ReadJSON(&msg); err != nil {
				log.Println("error payload")
				break
			}
			msg.Date = time.Now()
			err := repository.InsertNewMessage(context.Background(), msg)
			if err != nil {
				log.Println(err)
			}
			broadcast <- msg
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Println("failed to write json: ", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf(
		"%s:%s",
		env.GetEnv("APP_HOST", "localhost"),
		env.GetEnv("APP_PORT_SOCKET", "8080"),
	)))
}
