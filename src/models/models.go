package models

type LightState struct {
	On        bool
	Reachable bool
	Bri       int
	Hue       int
	Sat       int
}

type Light struct {
	State LightState
	Name  string
}
