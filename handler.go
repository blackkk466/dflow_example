package main

import (
	"encoding/json"
	"os"
	"reflect"
	. "github.com/mlabouardy/dialogflow-watchnow-messenger/models"
	"github.com/gin-gonic/gin"
)

type Note struct{
	ID int `gorm:"primary_key";"AUTO_INCREMENT" json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func verificationEndPoint(c *gin.Context) {

	challenge := c.Request.URL.Query().Get("hub.challenge")
	mode := c.Request.URL.Query().Get("hub.mode")
	token := c.Request.URL.Query().Get("hub.verify_token")
	if mode != "" && token == os.Getenv("VERIFY_TOKEN") {
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte(challenge))
	} else {
		c.Writer.WriteHeader(404)
		c.Writer.Write([]byte("Error, wrong validation token"))
	}
}

func messagesEndPoint(c *gin.Context) {
	var callback Callback
	json.NewDecoder(c.Request.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				if !reflect.DeepEqual(event.Message, Message{}) && event.Message.Text != "" {
					processMessage(event)
				}
			}
		}
		c.Writer.WriteHeader(200)
		c.Writer.Write([]byte("Got your message"))
	} else {
		c.Writer.WriteHeader(404)
		c.Writer.Write([]byte("Message not supported"))
	}
}

func getAllNotes() []Note{
	var noteList []Note

	db.Debug().Find(&noteList)

	return noteList
}