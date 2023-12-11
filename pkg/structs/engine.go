package structs

type Engine struct {
	Name            string
	Language        string
	Version         string
	Polyglots       map[string]string
	VerifyReflected string
	VerifyError     string
}
