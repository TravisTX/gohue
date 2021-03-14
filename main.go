package main

import (
	"encoding/json"
	"fmt"
	"gohue/models"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// todo: configuration
const (
	hueUsername    = "kfZqr9nqnhINjSKYlPXQ4R6TacR9nPE5Q9UOOC14"
	hueBridgeIp    = "10.0.0.112"
	hueLightNumber = "6"
)

func main() {
	printLights()
}

func printLights() {
	lights := getLights()

	keys := make([]string, 0, len(lights))
	for key := range lights {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	longestName := 0
	for _, key := range keys {
		length := len(fmt.Sprintf("%s %s", key, lights[key].Name))
		if length > longestName {
			longestName = length
		}
	}

	for _, key := range keys {
		light := lights[key]
		lightDisplay := getLightDisplay(key, light, longestName)
		fmt.Println(lightDisplay)
	}
}

func getLights() map[string]models.Light {
	url := fmt.Sprintf("http://%s/api/%s/lights", hueBridgeIp, hueUsername)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	var lights map[string]models.Light
	json.Unmarshal([]byte(sb), &lights)

	return lights
}

func getLightDisplay(lightId string, light models.Light, longestName int) string {
	onDisplay := ""
	v := 0.5
	if light.State.On {
		onDisplay = ""
		v = 1
	}

	c := colorful.Hsv(float64(light.State.Hue)/182, float64(light.State.Sat)/256, v)
	briDisplay := getAsciiProgressBar(int(float32(light.State.Bri) / 254 * 100))
	briDisplay = fmt.Sprintf("%s \x1b[38;2;%d;%d;%dm%s\x1b[0m", onDisplay, int32(c.R*255), int32(c.G*255), int32(c.B*255), briDisplay)

	out := fmt.Sprintf("%s %s", lightId, light.Name)

	out = out + strings.Repeat(" ", longestName-len(out)+2) + briDisplay
	return out
}

func getAsciiProgressBar(percentage int) string {
	maxWidth := 10
	unitsFilled := percentage * maxWidth / 100
	unitsRemaining := maxWidth - unitsFilled
	padding := ""
	if percentage < 100 {
		padding = " "
	}
	return fmt.Sprintf("%s%s %s%v%%", strings.Repeat("█", unitsFilled), strings.Repeat("░", unitsRemaining), padding, percentage)
}
