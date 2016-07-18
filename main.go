package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var (
	configFile  = flag.String("config", "", "Path to config.yml")
	httpAddress = flag.String("httpAddress", ":8080", "HTTP port")

	configs []RepoConfig

	scripTemplate = template.New("script")
)

type tplData struct {
	ENV map[string]string
	Hub Payload
}

func main() {
	flag.Parse()
	err := parseConfig()
	if err != nil {
		log.Fatalf("Can't load config: %+v", err)
	}

	router := httprouter.New()
	router.POST("/docker/:apikey", hook)

	log.Printf("Now listening on %s", *httpAddress)
	log.Fatal(http.ListenAndServe(*httpAddress, router))
}

func parseConfig() error {
	configBytes, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configBytes, &configs)
	return nil
}

func getConfigsForApiKey(apikey string) ([]RepoConfig, error) {
	result := make([]RepoConfig, 0, 10)
	for _, config := range configs {
		if apikey == config.ApiKey {
			result = append(result, config)
		}
	}
	var err error
	if len(result) == 0 {
		err = fmt.Errorf("Can't find configs for api key %s", apikey)
	}
	return result, err
}

func hook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	apiKey := ps.ByName("apikey")
	configs, err := getConfigsForApiKey(apiKey)
	if err != nil {
		log.Printf("Api key %s does not exist", apiKey)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var payload Payload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("Can't parse payload from Docker Hub: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dockerTag := payload.PushData.Tag
	for _, repoConfig := range configs {
		if repoConfig.Tag == dockerTag {
			log.Printf("Received valid call for key %s", apiKey)
			go executeScript(repoConfig, payload)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	log.Printf("Tag %s is not configured", dockerTag)
	w.WriteHeader(http.StatusBadRequest)
}

func executeScript(config RepoConfig, payload Payload) {
	tpl, err := scripTemplate.Parse(config.Script)
	if err != nil {
		log.Printf("Can't parse script template: %+v", err)
		return
	}
	tplVars := tplData{
		ENV: make(map[string]string),
		Hub: payload,
	}
	for _, envPair := range os.Environ() {
		parts := strings.Split(envPair, "=")
		if len(parts) == 2 {
			tplVars.ENV[parts[0]] = parts[1]
		} else {
			log.Printf("Unusual environment value %s", envPair)
		}
	}
	var scriptBuffer bytes.Buffer
	err = tpl.Execute(&scriptBuffer, tplVars)
	if err != nil {
		log.Printf("Can't execute script template: %+v", err)
		return
	}
	script := scriptBuffer.String()
	log.Printf("Executing %s", script)

	args := strings.Split(script, " ")
	scriptCommand := exec.Command(args[0], args[1:]...)
	scriptCommand.Env = os.Environ()
	// TODO wrap this in a nicer writer
	scriptCommand.Stdout = os.Stdout
	scriptCommand.Stderr = os.Stderr
	err = scriptCommand.Run()
	if err != nil {
		log.Printf("Error running script: %+v", err)
		return
	}
	_, err = http.Get(payload.CallbackUrl)
	if err != nil {
		log.Printf("Failed to call callback URL: %+v", err)
	}
}
