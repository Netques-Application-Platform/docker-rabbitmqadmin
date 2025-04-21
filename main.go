package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

func main() {
	var (
		queue    string
		user     string
		password string
		exchange string
		vhost    string
		host     string
		port     string
		debug    bool
	)

	viper.SetDefault("RABBIT_QUEUE", "nap-tasks")
	viper.SetDefault("RABBIT_USER", "guest")
	viper.SetDefault("RABBIT_PASSWORD", "guest")
	viper.SetDefault("RABBIT_EXCHANGE", "amq.default")
	viper.SetDefault("RABBIT_VHOST", "/")
	viper.SetDefault("RABBIT_HOST", "127.0.0.1")
	viper.SetDefault("RABBIT_PORT", "15672")
	viper.SetDefault("DEBUG", false)
	viper.AutomaticEnv()

	queue = viper.GetString("RABBIT_QUEUE")
	user = viper.GetString("RABBIT_USER")
	password = viper.GetString("RABBIT_PASSWORD")
	exchange = viper.GetString("RABBIT_EXCHANGE")
	vhost = viper.GetString("RABBIT_VHOST")
	host = viper.GetString("RABBIT_HOST")
	port = viper.GetString("RABBIT_PORT")
	debug = viper.GetBool("DEBUG")

	fullUrl := fmt.Sprintf("http://%v:%v/api/exchanges/%v/%v/publish", host, port, url.PathEscape(vhost), exchange)
	if debug {
		fmt.Println("DEBUG", user, password, exchange, vhost, queue, host, port)
	}

	fmt.Println("posting", fullUrl, queue, os.Args[1])

	payload := map[string]interface{}{
		"properties":       map[string]interface{}{},
		"routing_key":      queue,
		"payload":          os.Args[1],
		"payload_encoding": "string",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %s\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not read body: %v", err)
		os.Exit(1)
	}

	pf := pubResponse{}
	err = json.Unmarshal(content, &pf)
	if err != nil {
		fmt.Printf("could not unmarshall: %v", err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK && pf.Routed {
		fmt.Println("message published")
	} else if resp.StatusCode != http.StatusOK {
		fmt.Printf("error publishing message: %s\n", resp.Status)
		os.Exit(1)
	} else {
		fmt.Printf("message published but not routed\n")
		os.Exit(1)
	}
}

type pubResponse struct {
	Routed bool `json:"routed"`
}
