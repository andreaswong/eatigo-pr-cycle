package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/andreaswong/eatigo-pr-cycle/permutation"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

func main() {
	// .env for local, for heroku, set it in app dashboard
	// see .env.example
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read and decode config to Config struct
	config := &Config{}
	var configMap = map[string]string{} //key=>val
	for _, configLine := range os.Environ() {
		keyToVal := strings.SplitN(configLine, "=", 2)
		if len(keyToVal) == 2 {
			configMap[keyToVal[0]] = keyToVal[1]
		}
	}

	err = mapstructure.Decode(configMap, config)
	if err != nil {
		log.Fatal("Unable to decode from .env to %v, error=%v", config, err)
		panic(nil)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8515"
	}
	http.ListenAndServe(":"+port, NewDefaultHandler(config))
}

type DefaultHandler struct {
	config *Config
}

func NewDefaultHandler(config *Config) *DefaultHandler {
	return &DefaultHandler{
		config: config,
	}
}

func (h *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	weekday := int(now.Weekday())

	objJson := h.config.ObjectsJson
	engineers := []string{}
	if err := json.Unmarshal([]byte(objJson), &engineers); err != nil {
		log.Fatal("error unmarshaling objects, error=%v", err)
		panic("Unable to unmarshal OBJECTS_JSON from env (is it a json array string?)")
	}

	// Last guy reviews first guy
	permutation := permutation.Permute(engineers)
	engineers = permutation[weekday%len(permutation)]
	engineers = append(engineers, engineers[0])

	for _, cycle := range permutation {
		fmt.Printf("%v\n", cycle)
	}

	slackURL := h.config.SlackURL
	text := fmt.Sprintf(`curl -d '{"text": "%s"}' -X POST %s`,
		strings.Join(engineers, " `reviews PR of` "),
		slackURL,
	)
	w.Write([]byte(text))
}

type Config struct {
	ObjectsJson string `mapstructure:"OBJECTS_JSON"`
	SlackURL    string `mapstructure:"SLACK_URL"`
}
