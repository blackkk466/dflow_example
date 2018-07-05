package main

import (
	"time"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"github.com/jinzhu/gorm"
)

type Product struct{
	gorm.Model
	Name   string `json:"product-type"`
	Amount int    `json:"amount"`
}

type Session struct{
	gorm.Model
	Name string `gorm:"primary_key"`
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

type CountryDateResponse struct{
	total_population TotalPopulation
}

type TotalPopulation struct {
	date time.Time
	population int
}

type PopulationAPIResp CountryDateResponse

func APIAIPopulationEndpoint(c *gin.Context) {

	req := c.Request
	w := c.Writer

	decoder := json.NewDecoder(req.Body)


	var t APIAIRequest
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error in decoding the Request data1", http.StatusInternalServerError)
	}


	var speechResponse string
	var msg APIAIMessage

	var sess Session

	if db.Debug().Where("name = ?",t.SessionID).First(&sess).RecordNotFound() {
		if t.QueryResult.Intent.DisplayName != "email.control"{
			msg = APIAIMessage{Source: "No Email", FulfillText: "Please write your email to start a session."}
			json.NewEncoder(w).Encode(msg)

		}else{
			sess.Name = t.SessionID
			db.Create(&sess)
			fmt.Println("New session started")
		}
		return
	}else{
		if sess.UpdatedAt.IsZero() && time.Now().Sub(sess.UpdatedAt).Seconds() < 20.0{
			sess.UpdatedAt = time.Now()
			db.Save(&sess)
			fmt.Println("Update is NOT zero!")

		}else{
			if t.QueryResult.Intent.DisplayName != "email.control"{
				msg = APIAIMessage{Source: "Backend", FulfillText: "You haven't been writing for more than 20 seconds. Please write your email again."}
				json.NewEncoder(w).Encode(msg)

				db.Unscoped().Delete(&sess)
				return
			}
			sess.UpdatedAt = time.Now()
			db.Save(&sess)
			return
		}
	}

/*	fmt.Println("\n********************************")
	spew.Dump(t)
	fmt.Println("********************************\n")*/


	switch t.QueryResult.Intent.DisplayName {
	case "PopIntent":

		countryV := t.QueryResult.Parameters["geo-country"]
		date := time.Now().Local().Format("2006-01-02")

		country := countryV.(string)  //.GetStringValue()

		//Now make a call to the external Population.io Example
		fmt.Println("http://api.population.io/1.0/population/" + country + "/" + date)
		apiResponse, err := http.Get("http://api.population.io/1.0/population/" + country + "/" + date)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error in decoding the Request data2", http.StatusInternalServerError)
		}
		defer apiResponse.Body.Close()

		var populationAPIResponse PopulationAPIResp

		err = json.NewDecoder(apiResponse.Body).Decode(&populationAPIResponse)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error in decoding the Request data3", http.StatusInternalServerError)
		}
		//fmt.Println("\n-----------------",err,"---------------------\n")
		//spew.Dump(populationAPIResponse)
		//Create Response Message
		speechResponse = fmt.Sprintf("On %s, the Population Statistics shows a total of %d people in %s.", date, populationAPIResponse.total_population.population, country)
		msg = APIAIMessage{Source: "Population.io API", FulfillText: speechResponse }


	case "notes":
		notes := getAllNotes()

		var titles string
		titles = ""

		for _,n := range notes{
			titles += n.Title
			titles += " ,"
		}

		speechResponse = fmt.Sprintf("Here are your note title(s) %s ", titles)
		msg = APIAIMessage{Source: "Note Database", FulfillText: speechResponse}

	case "products":

		productName := t.QueryResult.Parameters["product-type"].(string)
		//amount := int(t.QueryResult.Parameters["amount"].(float64))

		var product Product

		if db.Debug().Where("name = ?", productName).First(&product).RecordNotFound() {
			speechResponse = fmt.Sprintf("There is no product named %s.",productName)

		}else{

			spew.Dump(product)
/*			if ( amount > product.Amount || amount < 1 ){
				speechResponse = fmt.Sprintf("You can't buy %d %s",amount,productName)

			}else{
				db.Debug().Model(&product).UpdateColumn("amount",gorm.Expr("amount - ?",amount))
				speechResponse = fmt.Sprintf("Your order of %d %s has been completed successfuly",amount,productName)
			}*/
		}

		//msg = APIAIMessage{Source: "Product Database", FulfillText: speechResponse}

	case "email.control":
		speechResponse = "Your session is restarted."
		msg = APIAIMessage{Source: "Email Case", FulfillText: speechResponse}

	default:
		http.Error(w,"Intents couldn`t read",http.StatusBadRequest)
		speechResponse = fmt.Sprintf("Sorry, something went wrong please try again.")
		msg = APIAIMessage{Source: "Switch Case", FulfillText: speechResponse}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

//API.AI Response
/*type APIAIMessage struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Source      string `json:"source"`
}*/
/*/
type APIAIRequest struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Result    struct {
		Parameters map[string]string `json:"parameters"`
		Contexts   []interface{}     `json:"contexts"`
		Metadata   struct {
			IntentID                  string `json:"intentId"`
			WebhookUsed               string `json:"webhookUsed"`
			WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
			IntentName                string `json:"intentName"`
		} `json:"metadata"`
		Score float32 `json:"score"`
	} `json:"result"`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
	SessionID       string      `json:"sessionId"`
	OriginalRequest interface{} `json:"originalRequest"`
}*/


/*INSERT INTO database1.products (name,amount) VALUES ("apple", 10)
INSERT INTO database1.products (name,amount) VALUES ("pear", 15)
INSERT INTO database1.products (name,amount) VALUES ("peach", 8)
INSERT INTO database1.products (name,amount) VALUES ("cucumber", 6)
INSERT INTO database1.products (name,amount) VALUES ("tomato", 20)
INSERT INTO database1.products (name,amount) VALUES ("potato", 18)
INSERT INTO database1.products (name,amount) VALUES ("onion", 12)*/