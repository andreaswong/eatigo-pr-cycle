package main

import (
	"encoding/json"
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
	envs, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &Config{}
	err = mapstructure.Decode(envs, config)
	if err != nil {
		log.Fatal("Unable to decode from .env to %v, error=%v", config, err)
		panic(nil)
	}

	port := os.Getenv("PORT")
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
	if err := json.Unmarshal([]byte(objJson), engineers); err != nil {
		panic("Unable to unmarshal OBJECTS_JSON from env (is it a json array string?)")
	}

	index := weekday % len(engineers)
	w.Write([]byte(strings.Join(permutation.Permute(engineers)[index], " reviews PR of ")))
}

type Config struct {
	ObjectsJson string `mapstructure:"OBJECTS_JSON"`
}
