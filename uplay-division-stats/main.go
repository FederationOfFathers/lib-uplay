package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/FederationOfFathers/lib-uplay"
	"github.com/koding/multiconfig"
)

type Uplay struct {
	Username string
	Password string
	Platform string
	Gamer    string
	UplayID  string
}

var config = &Uplay{}

func init() {
	var m multiconfig.Loader
	var l multiconfig.Loader
	var found = true
	if _, err := os.Stat(".uplay.toml"); err == nil {
		l = &multiconfig.TOMLLoader{Path: ".uplay.toml"}
	} else if _, err := os.Stat(".uplay.json"); err == nil {
		l = &multiconfig.JSONLoader{Path: ".uplay.json"}
	} else if _, err := os.Stat(fmt.Sprintf("%s/.uplay.toml", os.Getenv("HOME"))); err == nil {
		l = &multiconfig.TOMLLoader{Path: fmt.Sprintf("%s/.uplay.toml", os.Getenv("HOME"))}
	} else if _, err := os.Stat(fmt.Sprintf("%s/.uplay.json", os.Getenv("HOME"))); err == nil {
		l = &multiconfig.JSONLoader{Path: fmt.Sprintf("%s/.uplay.json", os.Getenv("HOME"))}
	} else {
		found = false
	}
	if found {
		m = multiconfig.MultiLoader(l, &multiconfig.FlagLoader{}, &multiconfig.EnvironmentLoader{})
	} else {
		m = multiconfig.MultiLoader(&multiconfig.FlagLoader{}, &multiconfig.EnvironmentLoader{})
	}
	m.Load(config)
}

func main() {
	client := uplay.New(config.Username, config.Password)
	if err := client.Authenticate(); err != nil {
		log.Fatalf("Error authenticating to ubisoft: %s", err.Error())
	}
	if config.UplayID == "" && config.Gamer != "" {
		if profiles, err := client.UserSearch(uplay.GuessPlatform(config.Platform), config.Gamer); err != nil {
			log.Fatalf("Error searching for user profile: %s", err.Error())
		} else {
			if len(profiles) < 1 {
				log.Fatal("No matching profiles found")
			}
			config.UplayID = profiles[0].UplayID
		}
	}
	if config.UplayID == "" {
		log.Fatal("No UplayID Speficied")
	}
	if stats, err := client.DivisionStats(uplay.GuessPlatform(config.Platform), config.UplayID); err != nil {
		log.Fatalf("Error fetching stats: %s", err.Error())
	} else {
		y, _ := yaml.Marshal(stats)
		fmt.Println(string(y))
	}
}
