package main

import (
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"github.com/jinzhu/gorm"
)

type Note struct{
	ID int `gorm:"primary_key";"AUTO_INCREMENT" json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type Product struct{
	gorm.Model
	Name   string `json:"product-type"`
	Amount int    `json:"amount"`
	Price float32 `json:"price"`
}

type Session struct{
	gorm.Model
	Name string `gorm:"primary_key"`
	Email string
}

type APIAIRequest struct {
	ID        string    `json:"responseId"`
	SessionID       string      `json:"session"`
	QueryResult struct{
		QueryText string `json:"queryText"`
		Parameters map[string]interface{} `json:"parameters"`
		AllRequiredPresent bool `json:"allRequiredParamsPresent"`
		FulfillText string `json:"fulfillmentText"`
		OutputContexts []dialogflow.Context `json:"outputContexts"`
		Intent struct{
			Name string `json:"name"`
			DisplayName string `json:"displayName"`
		}`json:"intent"`

	} `json:"queryResult"`
}

type APIAIMessage struct {
	FulfillText string `json:"fulfillmentText"`
	Source      string `json:"source"`
	OutputContexts []dialogflow.Context `json:"outputContexts"`
}




/*INSERT INTO database1.product (name,amount,price) VALUES ("apple", 10, 3.2)
INSERT INTO database1.product (name,amount,price) VALUES ("pear", 15, 4.5)
INSERT INTO database1.product (name,amount,price) VALUES ("peach", 8, 1.0)
INSERT INTO database1.product (name,amount,price) VALUES ("cucumber", 6, 2.1)
INSERT INTO database1.product (name,amount,price) VALUES ("tomato", 20, 3.3)
INSERT INTO database1.product (name,amount,price) VALUES ("potato", 18, 4.7)
INSERT INTO database1.product (name,amount,price) VALUES ("onion", 12, 5.5)*/