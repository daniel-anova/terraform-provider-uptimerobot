package uptimerobotapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func New(apiKey string) UptimeRobotApiClient {
	return UptimeRobotApiClient{apiKey}
}

type UptimeRobotApiClient struct {
	apiKey string
}

func (client UptimeRobotApiClient) MakeCall(
	endpoint string,
	params string,
) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Making request to: %#v", endpoint)

	url := "https://api.uptimerobot.com/v2/" + endpoint

	payload := strings.NewReader(
		fmt.Sprintf("api_key=%s&format=json&%s", client.apiKey, params),
	)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] Got response: %#v", res)
	log.Printf("[DEBUG] Got body: %#v", string(body))

	fmt.Printf("Got response: %#v\n", res)
	fmt.Printf("Got body: %#v\n", string(body))

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	if result["stat"] != "ok" {
		message, _ := json.Marshal(result["error"])
		return nil, errors.New("Got error from UptimeRobot: " + string(message))
	}

	return result, nil
}
