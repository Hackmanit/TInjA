package structs

type Config struct {
	Timeout          int
	Verbosity        int
	PrecedingLength  int
	SubsequentLength int
	LengthLimit      int

	Ratelimit float64

	Data          string
	ProxyCertPath string
	ProxyURL      string
	ReportPath    string
	UserAgent     string

	Cookies        []string
	Headers        []string
	Parameters     []string
	URLs           []string
	URLsReflection []string
	Crawls         []Crawl

	UserAgentChrome bool
	EscapeJSON      bool
	CSTI            bool
}
