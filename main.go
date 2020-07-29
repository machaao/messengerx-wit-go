package main

import (
	"strings"

	// "bufio"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"machaao-go/extras"

	"github.com/dgrijalva/jwt-go"
	witai "github.com/wit-ai/wit-go"
	// "os"
	// witai "github.com/wit-ai/wit-go"
)

var machaaoAPIToken string = os.Getenv("MachaaoApiToken")

//WitAiToken Get it from https://wit.ai
var witAiToken string = os.Getenv("WitApiToken")

func getJokeTagUsingWitAI(message string) string {
	client := witai.NewClient(witAiToken)
	// Use client.SetHTTPClient() to set custom http.Client

	msg, _ := client.Parse(&witai.MessageRequest{
		Query: message,
	})

	return msg.Entities["local_search_query"].([]interface{})[0].(map[string]interface{})["value"].(string)
}

func simpleReply(userID []string, message string, apiToken string) {

	if strings.ToLower(message) == "😜 Random Jokes" {
		message = extras.GetJoke("%20")
	} else if message == "🙃 Random Memes" {
		title, url, postlink := extras.GetMemes()

		_ = title

		body := map[string]interface{}{
			"users": userID,
			"message": map[string]interface{}{
				"attachment": map[string]interface{}{
					"type": "template",
					"payload": map[string]interface{}{
						"template_type": "generic",
						"elements": []map[string]interface{}{
							{
								"image_url": url,
								"buttons": []map[string]string{
									{
										"type":  "web_url",
										"url":   postlink,
										"title": "ℹ️ Source",
									},
								},
							},
						},
					},
				},
				"quick_replies": []map[string]string{
					{
						"content_type": "text",
						"payload":      "😜 Random Jokes",
						"title":        "😜 Random Jokes",
					},
					{
						"content_type": "text",
						"payload":      "🙃 Random Memes",
						"title":        "🙃 Random Memes",
					},
				},
			},
		}
		log.Println("Sending Message to user")

		var urlMachaao string = "https://ganglia-dev.machaao.com/v1/messages/send"

		jsonValue, _ := json.Marshal(body)

		// fmt.Println(jsonValue)

		req, err := http.NewRequest("POST", urlMachaao, bytes.NewBuffer(jsonValue))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api_token", apiToken)

		fmt.Println(req)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		bodyf, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(bodyf))

		return

	} else {
		var tag string = getJokeTagUsingWitAI(message)
		message = extras.GetJoke(tag)

		if message[:9] == "Error 106" {
			message = "Sorry, no jokes found"
		}
	}

	log.Println("Sending Message to user")

	var url string = "https://ganglia-dev.machaao.com/v1/messages/send"
	// var url string = "http://127.0.0.1:5000/upload"

	body := map[string]interface{}{
		"users": userID,
		"message": map[string]interface{}{
			"text": message,
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "😜 Random Jokes",
					"title":        "😜 Random Jokes",
				},
				{
					"content_type": "text",
					"payload":      "🙃 Random Memes",
					"title":        "🙃 Random Memes",
				},
			},
		},
	}

	jsonValue, _ := json.Marshal(body)

	// fmt.Println(jsonValue)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_token", apiToken)

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	bodyf, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(bodyf))
}

func messageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	var bodyData string = string(body)
	var tokenString string = bodyData[8:(len(bodyData) - 2)]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(machaaoAPIToken), nil
	})

	_ = token

	if err != nil {
		fmt.Println(err)
	}

	messageData := claims["sub"].(map[string]interface{})["messaging"].([]interface{})[0].(map[string]interface{})["message_data"]
	messageText := messageData.(map[string]interface{})["text"].(string)

	log.Println(messageData)
	log.Println(messageText)

	log.Println(r.Header["User_id"])

	if messageText == "hi" {
		quickReply(r.Header["User_id"], messageText, machaaoAPIToken)
	} else {
		simpleReply(r.Header["User_id"], messageText, machaaoAPIToken)
	}
}

func main() {
	port := getPort()
	http.HandleFunc("/machaao_hook", messageHandler)

	log.Println("[-] Listening on...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(machaaoAPIToken)
	// fmt.Println(witAiToken)

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4747"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}

func quickReply(userID []string, message string, apiToken string) {

	log.Println("Sending QR to user")

	var url string = "https://ganglia-dev.machaao.com/v1/messages/send"
	// var url string = "http://127.0.0.1:5000/upload"

	body := map[string]interface{}{
		"users": userID,
		"message": map[string]interface{}{
			"text": "Hello, My name is Witty - Your funny friend ;)",
			"quick_replies": []map[string]string{
				{
					"content_type": "text",
					"payload":      "😜 Random Jokes",
					"title":        "😜 Random Jokes",
				},
				{
					"content_type": "text",
					"payload":      "🙃 Random Memes",
					"title":        "🙃 Random Memes",
				},
			},
		},
	}

	jsonValue, _ := json.Marshal(body)

	// fmt.Println(jsonValue)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_token", apiToken)

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
}
