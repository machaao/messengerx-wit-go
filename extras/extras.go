package extras

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type memeResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Nsfw      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
}

// GetJoke uses a 3rd party API which is liscensed under MIT License
// https://github.com/Sv443/JokeAPI
func GetJoke(tag string) string {

	var url string = "https://sv443.net/jokeapi/v2/joke/Any?format=txt&contains=" + tag

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Accept", "text/plain")
	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

// GetMemes uses a 3rd party API which is liscensed under MIT License
// https://github.com/R3l3ntl3ss/Meme_Api
func GetMemes() (string, string, string) {

	var url string = "https://meme-api.herokuapp.com/gimme"

	var jsonBody memeResponse

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &jsonBody)

	return jsonBody.Title, jsonBody.URL, jsonBody.PostLink
}
