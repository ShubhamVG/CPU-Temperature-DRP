package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hugolgst/rich-go/client"
)

const (
	CRITICAL_TEMPERATURE = 55.0
	DURATION_COOLDOWN    = time.Second * 5
	EMOJI_HEART_ON_FIRE  = "heart-on-fire"
	EMOJI_EXPLOSION      = "explosion"
)

// These are also consts btw
var (
	// check your sensors and include it here
	SENSORS_TO_USE = []string{
		"INSERT SENSOR NAME" // like "k10temp-pci-24c3",
	}
)

func waitForCooldown() {
	time.Sleep(DURATION_COOLDOWN)
}

func main() {
	clientId := os.Getenv("DISCORD_RICH_PRESENCE_CLIENT_ID")

	if err := client.Login(clientId); err != nil {
		fmt.Println("The Discord app is probably not open.")
		os.Exit(1)
	}

	// First status
	now := time.Now()
	initialActivityErr := client.SetActivity(client.Activity{
		State:      "Time to burn it down.",
		Details:    "CPUwU",
		LargeImage: EMOJI_HEART_ON_FIRE,
		Timestamps: &client.Timestamps{
			Start: &now,
		},
	})

	if initialActivityErr != nil {
		fmt.Println("Failed to set the initial activity.")
		os.Exit(1)
	}

	waitForCooldown()

	for {
		var activity client.Activity
		var err error
		cpuTemp, err := getAvgTempOfSensors(SENSORS_TO_USE)

		if err != nil {
			activity = client.Activity{
				State:      "The chicken has either been overcooked OR IT'S FKING RAW.",
				Details:    "The CPU temperature is not available.",
				LargeImage: EMOJI_EXPLOSION,
				Timestamps: &client.Timestamps{
					Start: &now,
				},
			}
		} else {
			var state string
			var largeImageName string

			if cpuTemp > CRITICAL_TEMPERATURE {
				largeImageName = EMOJI_EXPLOSION
				state = "THE HOLY FAN HAS KICKED IN. WE ARE DOOMED AFHKFADHFJDAFD"
			} else {
				largeImageName = EMOJI_HEART_ON_FIRE
				state = "It's being slow cooked :)"
			}

			cpuTempString := fmt.Sprintf("%.0f°C.", cpuTemp)
			activity = client.Activity{
				State:      state,
				Details:    "The CPU temperature is " + cpuTempString,
				LargeImage: largeImageName,
				Timestamps: &client.Timestamps{
					Start: &now,
				},
			}
		}

		err = client.SetActivity(activity)

		if err != nil {
			fmt.Println("Failed to update CPU temperature.")
			os.Exit(1)
		}

		waitForCooldown()
	}
}
