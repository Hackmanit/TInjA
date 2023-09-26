package structs

var ReflectionBody = "Body"
var ReflectionHeader = "Header"

type Reflection struct {
	ReflectionType string
	HeaderName     string `json:",omitempty"`
	ReflectionURL  string `json:",omitempty"`
	Preceding      string
	Subsequent     string
}
