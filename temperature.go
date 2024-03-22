package main

import (
	"encoding/json"
	"errors"
	"os/exec"
	"reflect"
	"runtime"
	"slices"
	"strings"
)

type JsonMap = map[string]any

// Return the temp (if found) along with whether the temp was present/found or not.
func extractTemp(jsonMap JsonMap) (float64, bool) {
	hasFound := false
	extractedTemp := float64(0)

	// Recurse until you found the temp
	for key, value := range jsonMap {
		if strings.Contains(key, "temp") && strings.Contains(key, "input") {
			if temp, ok := value.(float64); ok {
				extractedTemp += temp
				hasFound = true
			}
		} else if reflect.TypeOf(value).Kind() == reflect.Map {
			if valueAsJson, ok := value.(JsonMap); ok {
				extractedTemp2, hasFound2 := extractTemp(valueAsJson)

				if hasFound2 {
					extractedTemp += extractedTemp2
					hasFound = true
				}
			}
		}
	}

	return extractedTemp, hasFound
}

// Use a platform dependant API like `lm-sensors` in linux to get the temperature
// of hardwares whose names match the names of elements in  `sensors`
// and then parse the information and return it.
func getAvgTempOfSensors(sensors []string) (float64, error) {
	avgCpuTemp := float64(0)

	switch runtime.GOOS {
	case "darwin": // AKA mac
		// Nice joke
	case "linux":
		// Use sensors (AKA lm-sensors) to get the temps in JSON
		// and then take average temp from the parsed JSON and return.
		// Return an error if something goes wrong.
		command := exec.Command("sensors", "-j")
		outputAsBytesArray, commandErr := command.Output()

		if commandErr != nil {
			return 0.0, commandErr
		}

		tempsJson := map[string]JsonMap{}

		if jsonParsingError := json.Unmarshal(outputAsBytesArray, &tempsJson); jsonParsingError != nil {
			return 0.0, jsonParsingError
		}

		hardwareTempSum := float64(0)
		hardwareCount := float64(0)

		for hardwareName, value := range tempsJson {
			if slices.Contains(sensors, hardwareName) {
				tempOfHardware, hadTemperatureReading := extractTemp(value)

				if hadTemperatureReading {
					hardwareTempSum += tempOfHardware
					hardwareCount++
				}
			}
		}

		avgCpuTemp = hardwareTempSum / hardwareCount
		return avgCpuTemp, nil
	case "windows":
		// Maybe TODO if I start using windows again
	}

	// If it has reached this line, then the OS was not supported.
	return 0.0, errors.New("OS not supported")
}
