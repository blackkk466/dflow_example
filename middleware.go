package main

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

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

	if t.QueryResult.Intent.DisplayName == "menu.show"{
		speechResponse = fmt.Sprintf("**MENU** PRODUCT - AMOUNT - PRICE %s",getProducts() )
		msg = APIAIMessage{Source: "Menu", FulfillText: speechResponse}
		json.NewEncoder(w).Encode(msg)
		return
	}

	if db.Where("name = ?",t.SessionID).First(&sess).RecordNotFound() {
		if t.QueryResult.Intent.DisplayName != "email.control"{
			msg = APIAIMessage{Source: "No Email", FulfillText: "Please write your email to start a session."}

		}else{
			sess.Name = t.SessionID
			sess.Email = t.QueryResult.Parameters["email"].(string)
			db.Create(&sess)
			speechResponse = fmt.Sprintf("Your session for %s is started. The products you can buy and" +
				" their amounts are: %s", sess.Email, getProducts())
			msg = APIAIMessage{Source: "Email Case", FulfillText: speechResponse}
		}

		json.NewEncoder(w).Encode(msg)
		return
	}else{
		if time.Now().Sub(sess.UpdatedAt).Seconds() < 25.0{
			sess.UpdatedAt = time.Now()
			db.Save(&sess)
			//fmt.Println("Update is NOT zero!")

		}else{
			if t.QueryResult.Intent.DisplayName != "email.control"{
				msg = APIAIMessage{Source: "Backend", FulfillText: "You haven't been writing for more than 25 seconds. Please write your email again."}
				json.NewEncoder(w).Encode(msg)

				db.Unscoped().Delete(&sess)
				return
			}
			sess.UpdatedAt = time.Now()
			sess.Email = t.QueryResult.Parameters["email"].(string)
			db.Save(&sess)
			return
		}
	}

	/*	fmt.Println("\n********************************")
		spew.Dump(t)
		fmt.Println("********************************\n")*/

	switch t.QueryResult.Intent.DisplayName {
	case "products":

		productName := t.QueryResult.Parameters["product-type"].(string)
		amount := int(t.QueryResult.Parameters["amount"].(float64))
		theTime := t.QueryResult.Parameters["time"].(string)

		var product Product

		if db.Where("name = ?", productName).First(&product).RecordNotFound() {
			speechResponse = fmt.Sprintf("There is no product named %s.",productName)

		}else{

			//spew.Dump(product)
			if ( amount > product.Amount || amount < 1 ){
				speechResponse = fmt.Sprintf("You can't buy %d %s",amount,productName)

			}else{
				total := float64(product.Price) * t.QueryResult.Parameters["amount"].(float64)
				db.Model(&product).UpdateColumn("amount",gorm.Expr("amount - ?",amount))
				speechResponse = fmt.Sprintf("Your order of %d %s for %s has been completed successfuly. Total PRICE is %.2f" +
					"And it will be delivered on %s",amount,productName,sess.Email,total,theTime)
			}
		}

		msg = APIAIMessage{Source: "Product Database", FulfillText: speechResponse}

	case "notes":
		notes := getAllNotes()

		var titles string
		titles = ""

		for _,n := range notes{
			titles += n.Title
			titles += " ,"
		}

		speechResponse = fmt.Sprintf("Here is your note title(s) %s ", titles)
		msg = APIAIMessage{Source: "Note Database", FulfillText: speechResponse}

	case "email.control":
		speechResponse = fmt.Sprintf("Your session for %s is restarted.",sess.Email)
		msg = APIAIMessage{Source: "Email Case", FulfillText: speechResponse}

	default:
		http.Error(w,"Intents couldn`t read",http.StatusBadRequest)
		speechResponse = fmt.Sprintf("Sorry, something went wrong please try again.")
		msg = APIAIMessage{Source: "Switch Case", FulfillText: speechResponse}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}