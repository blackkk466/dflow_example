package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	apiai "github.com/mlabouardy/dialogflow-go-client/models"
	. "github.com/mlabouardy/dialogflow-watchnow-messenger/models"
)

const (
	FACEBOOK_API   = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
	//IMAGE          = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
	//VISIT_SHOW_URL = "http://www.blog.labouardy.com/bot-in-messenger-with-dialogflow-golang/"
)
/*
func buildCarousel(shows []Show) []Element {
	elements := make([]Element, 0)
	for _, show := range shows {
		element := Element{
			Title:    show.Title,
			ImageURL: show.Cover,
			DefaultAction: DefaultAction{
				Type: "web_url",
				URL:  VISIT_SHOW_URL,
			},
			Buttons: []Button{
				Button{
					Type:  "web_url",
					Title: "Watch now",
					URL:   VISIT_SHOW_URL,
				},
			},
		}
		elements = append(elements, element)
	}
	return elements
}*/

func processMessage(event Messaging) {
	var userQuery = event.Message.Text
	var dialogFlowResponse = getResponse(userQuery)
	client := &http.Client{}
	var response Response

	if !reflect.DeepEqual(dialogFlowResponse.Metadata, apiai.Metadata{}) && dialogFlowResponse.Metadata.IntentName == "notes" {
		//var showType = dialogFlowResponse.Parameters["note-type"]

		var notes []Note
		var titles []string



		notes = getAllNotes()

		for _, note := range notes{
			titles = append(titles,note.Title)
		}

		response = Response{
			Recipient: User{
				ID: event.Sender.ID,
			},
			Message: Message{
				Text: titles[1],
			},
		}
	} else {
		response = Response{
			Recipient: User{
				ID: event.Sender.ID,
			},
			Message: Message{
				Text: dialogFlowResponse.Fulfillment.Speech,
			},
		}
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)

	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
