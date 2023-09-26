package pkg

/* Documentation

The `scan.go` file contains functions related to scanning URLs and analyzing parameters for template injection vulnerabilities. Let's go through each function:

The `Scan` function is the main entry point for initiating the scan. It takes a `structs.Config` object (`configParam`), a version string (`version`), and a type integer (`typ`) as input. It sets up the scan configuration, including rate limiting, and proxy settings. It then loops through the URLs or crawls specified in the configuration and performs the scan on each target. It collects the scan results, generates a report if requested, and calculates the scan duration.

The `scanURL` function takes a URL string (`u`), a `structs.Crawl` object (`crawl`), and a type integer (`typ`) as input. It performs the scan on the specified URL or crawl. It sends a default request to the target and analyzes the query parameters, post parameters, and headers for template injection vulnerabilities. It collects the scan results, including vulnerable parameters, and returns a `ReportWebpage` object containing the scan information.

The `getPostParams` function retrieves post parameters from the configuration data. It parses the data string and returns a map containing the key-value pairs of the post parameters.

The `setProxy` function sets up a proxy with the specified `proxyURL` and `proxyCertPath` in the scan configuration. It configures the HTTP client transport to use the proxy for connections. If a proxy certificate path is provided, it loads the certificate and appends it to the certificate pool used for TLS connections.

These functions work together to perform the scan on URLs or crawls, analyze parameters for template injection vulnerabilities, and handle proxy settings if configured.

Additionally, the file includes variable declarations, including global variables used throughout the scanning process.

*******/

import (
	"crypto/tls"
	"crypto/x509"
	"example/user/tinja/pkg/structs"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/time/rate"
)

var (
	currentDate, defaultBody string
	reportPath               string
	defaultHeader            http.Header
	defaultStatus            int
	report                   Report
	config                   structs.Config
	postParams               map[string]string
	limiter                  *rate.Limiter
	caCertPool               *x509.CertPool
	boolReport               bool
	configHeadersBackup      []string
)

var webpageID = 0
var userInputCounter = 0

var (
	QUERY  = 1
	POST   = 2
	HEADER = 3
)

var (
	URL   = 1
	JSONL = 2
)

func Scan(configParam structs.Config, version string, typ int) {
	// set config for every file in package pkg
	config = configParam
	// Get Date for filenames etc
	currentDate = time.Now().Format("2006-01-02_15-04-05")
	/******************************************/

	msg := fmt.Sprintf("TInjA %s started at %s\n", version, currentDate)
	PrintVerbose(msg, NoColor, 1)

	// get start time to calculate scan duration
	start := time.Now()
	// Shall a report be created?
	boolReport = config.ReportPath != ""
	// Making the random generator really random
	rand.Seed(time.Now().UnixNano())

	/**** Fill Report *****/
	if boolReport {
		report.Name = "Template_Injection_Scanner"
		report.Version = version

		report.Config = &config

		report.Date = currentDate
		report.Duration = "Not finished yet"

		report.Command = fmt.Sprint(os.Args)
		report.Command = strings.TrimPrefix(report.Command, "[")
		report.Command = strings.TrimSuffix(report.Command, "]")

		// create the file, so that later file errors are avoided
		reportPath = generateReport(report, currentDate)
	}
	/***********************/

	// Setting up proxy (e.g. burp), if wanted
	if config.ProxyURL != "" {
		setProxy()
	}

	// Setting up rate limiter
	ratelimit := config.Ratelimit
	if ratelimit <= 0 {
		ratelimit = math.MaxFloat64
	}
	limiter = rate.NewLimiter(rate.Limit(ratelimit), 1)

	// Create a header backup if JSONL mode is used
	if typ == JSONL {
		configHeadersBackup = config.Headers
	}

	// loop through every URL
	length := len(config.URLs) + len(config.Crawls)
	for i := 0; i < length; i++ {
		var u string
		var crawl structs.Crawl
		switch typ {
		case JSONL:
			u = config.Crawls[i].Request.Endpoint
			crawl = config.Crawls[i]
		case URL:
			u = config.URLs[i]
			crawl = structs.Crawl{}
		default:
			Print("Scan: Unknown typ: "+strconv.Itoa(typ), Red)
		}
		progress := fmt.Sprintf("(%d/%d)", i+1, length)
		msg := fmt.Sprintf("\nAnalyzing URL%s: %s\n", progress, u)
		Print(msg, NoColor)
		Print("===============================================================\n", NoColor)

		repWebpage := scanURL(u, crawl, typ)

		// add webpage to report
		if boolReport {
			if repWebpage.IsVulnerable {
				report.SuspectedVulnerableURLs++
			}

			addWebpageToReport(repWebpage, reportPath)
		}

		Print("===============================================================\n\n", NoColor)
	}

	/* Scan finished */
	msg = "Successfully finished the scan\n"
	PrintVerbose(msg, NoColor, 1)

	if boolReport {
		PrintVerbose("Suspected vulnerable URLs: "+strconv.Itoa(report.SuspectedInjections)+"\n", Green, 1)
		PrintVerbose("Suspected template injections: "+strconv.Itoa(report.SuspectedInjections)+"\n", Green, 1)
		msg = fmt.Sprintf("%d High, %d Medium, %d Low certainty\n\n", report.High, report.Medium, report.Low)
		PrintVerbose(msg, Green, 1)
	}
	duration := time.Since(start)
	msg = fmt.Sprintf("Duration: %s\n", duration)
	PrintVerbose(msg, NoColor, 1)
	averagePolyglots := float64(CounterPolyglotsGlobal) / float64(userInputCounter)
	Print(fmt.Sprint("Average polyglots sent per user input: ", averagePolyglots, "\n\n"), NoColor)
	/****************/

	if boolReport {
		report.Duration = duration.String()
		updateReportsFirstLine(report, reportPath)
	}

}

func scanURL(u string, crawl structs.Crawl, typ int) ReportWebpage {
	var repWebpage ReportWebpage
	var err error
	var responseDump string
	repWebpage.URL = u
	repWebpage.ID = webpageID
	webpageID++

	/****** Setting up client *****/
	msg := "Setting up client\n"
	PrintVerbose(msg, NoColor, 2)

	http.DefaultClient.Jar, _ = cookiejar.New(nil)
	http.DefaultClient.Timeout = time.Duration(config.Timeout) * time.Second
	// Disable certificate verification for websites with invalid certificates
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if typ == JSONL {
		config.Headers = configHeadersBackup
		var headerSlice []string
		// Setting cookies, user-agent and other headers
		for k, v := range crawl.Request.Headers {
			switch strings.ToLower(k) {
			case "cookie":
				cookieSlice := []*http.Cookie{}
				for _, cookie := range strings.Split(v, ";") {
					newCookie := strings.SplitN(cookie, "=", 2)
					if len(newCookie) == 2 {
						cookie := http.Cookie{
							Name:  strings.TrimSpace(newCookie[0]),
							Value: strings.TrimSpace(newCookie[1]),
						}
						cookieSlice = append(cookieSlice, &cookie)
					}
				}
				// Convert string u to URL in order to use it to set the specified Cookies
				var uURL *url.URL
				if uURL, err = url.Parse(crawl.Request.Endpoint); err != nil {
					msg := "Couldn't convert " + crawl.Request.Endpoint + " to URL:" + err.Error() + "\n"
					Print(msg, Red)
				}
				http.DefaultClient.Jar.SetCookies(uURL, cookieSlice)
			default:
				headerSlice = append(config.Headers, k+":"+v)
			}
		}
		// add headers from JSONL before the command line argument headers, so that the command line argument headers have priority
		config.Headers = append(headerSlice, config.Headers...)
		// Set Body
		config.Data = crawl.Request.Body
	}

	// Setting cookies, specified by setcookies
	if len(config.Cookies) > 0 {
		cookieSlice := []*http.Cookie{}
		for _, c := range config.Cookies {
			c = strings.TrimSuffix(c, "\r")
			c = strings.TrimSpace(c)
			if c == "" {
				continue
			} else if !strings.Contains(c, "=") {
				msg = "Specified cookie %s doesn't contain a = and will be skipped\n"
				PrintVerbose(msg, NoColor, 2)
				continue
			} else {
				cSlice := strings.SplitAfterN(c, "=", 2)
				cSlice[0] = strings.TrimSuffix(cSlice[0], "=")

				cookie := http.Cookie{
					Name:  cSlice[0],
					Value: cSlice[1],
				}
				cookieSlice = append(cookieSlice, &cookie)
			}
		}

		// Convert string u to URL in order to use it to set the specified Cookies
		var uURL *url.URL
		if uURL, err = url.Parse(u); err != nil {
			msg := "Couldn't convert " + u + " to URL:" + err.Error() + "\n"
			Print(msg, Red)
		}

		http.DefaultClient.Jar.SetCookies(uURL, cookieSlice)
	}

	/*******************************/

	msg = "Sending default request\n"
	PrintVerbose(msg, NoColor, 2)
	reqDefault, _ := buildRequest(u, config)
	if err != nil {
		msg := "Error: ScanURL: buildRequest: " + err.Error()
		Print(msg+"\n", Red)
		repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, msg)
		return repWebpage
	}
	defaultBody, defaultHeader, defaultStatus, responseDump, err = doRequest(reqDefault)
	if err != nil {
		msg := fmt.Sprintf("Couldn't connect to URL: %s", err.Error())
		Print(msg+"\n", Red)
		repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, msg)
		return repWebpage
	} else if strings.HasPrefix(defaultBody, "<html><head><title>Burp Suite") {
		msg := fmt.Sprintf("Couldn't connect to URL: \n%s", defaultBody)
		Print(msg+"\n", Red)
		repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, msg)
		return repWebpage
	}
	msg = fmt.Sprintf("Status code %d\n", defaultStatus)
	PrintVerbose(msg, NoColor, 1)

	var repParam reportParameter
	var statusCodeChanged bool

	values := reqDefault.URL.Query()
	msg = fmt.Sprintf("Found %d query parameters\n", len(values))
	PrintVerbose(msg, NoColor, 2)
	for k, v := range values {
		msg = fmt.Sprintln("Analyzing query parameter ", k, " => ", v)
		PrintVerbose(msg, NoColor, 1)
		userInputCounter++
		repParam, statusCodeChanged = analyze(k, QUERY, u)
		if boolReport {
			repParam.Type = "Query"
			repParam.Name = k
			repParam.DefaultValues = v
			repWebpage.IsVulnerable = repWebpage.IsVulnerable || repParam.IsVulnerable
			requestError := false
			for _, rr := range repParam.Requests {
				requestError = rr.Error != ""
				if requestError {
					break
				}
			}
			// only add repParam to the report, if the param is vulnerable or if there are error messages or if there are reflections
			if repParam.IsVulnerable || len(repParam.ErrorMessages) > 0 || len(repParam.Reflections) > 0 || requestError {
				repWebpage.Parameters = append(repWebpage.Parameters, repParam)
				if repWebpage.Certainty != certaintyHigh && repWebpage.Certainty != certaintyMedium && repParam.Certainty == certaintyLow {
					repWebpage.Certainty = certaintyLow
					report.Low += 1
					report.SuspectedInjections += 1
				} else if repWebpage.Certainty != certaintyHigh && repParam.Certainty == certaintyMedium {
					repWebpage.Certainty = certaintyMedium
					report.Medium += 1
					report.SuspectedInjections += 1
				} else if repParam.Certainty == certaintyHigh {
					repWebpage.Certainty = certaintyHigh
					report.High += 1
					report.SuspectedInjections += 1
				}
			}
			if statusCodeChanged {
				repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, "Status code changed. Skipping this URL.")
				break
			}
		}
	}

	if !statusCodeChanged {
		postParams = getPostParams()
		msg = fmt.Sprintf("Found %d post parameters\n", len(postParams))
		PrintVerbose(msg, NoColor, 2)
		for k, v := range postParams {
			msg = fmt.Sprintln("Analyzing post parameter ", k, " => ", v)
			PrintVerbose(msg, NoColor, 1)
			userInputCounter++
			repParam, statusCodeChanged := analyze(k, POST, u)
			if boolReport {
				repParam.Type = "POST"
				repParam.Name = k
				repParam.DefaultValues = []string{v}
				repWebpage.IsVulnerable = repWebpage.IsVulnerable || repParam.IsVulnerable
				requestError := false
				for _, rr := range repParam.Requests {
					requestError = rr.Error != ""
					if requestError {
						break
					}
				}
				// only add repParam to the report, if the param is vulnerable or if there are error messages or if there are reflections
				if repParam.IsVulnerable || len(repParam.ErrorMessages) > 0 || len(repParam.Reflections) > 0 || requestError {
					repWebpage.Parameters = append(repWebpage.Parameters, repParam)
					if repWebpage.Certainty != certaintyHigh && repWebpage.Certainty != certaintyMedium && repParam.Certainty == certaintyLow {
						repWebpage.Certainty = certaintyLow
						report.Low += 1
						report.SuspectedInjections += 1
					} else if repWebpage.Certainty != certaintyHigh && repParam.Certainty == certaintyMedium {
						repWebpage.Certainty = certaintyMedium
						report.Medium += 1
						report.SuspectedInjections += 1
					} else if repParam.Certainty == certaintyHigh {
						repWebpage.Certainty = certaintyHigh
						report.High += 1
						report.SuspectedInjections += 1
					}
				}
				if statusCodeChanged {
					repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, "Status code changed. Skipping this URL.")
					break
				}
			}
		}
	}

	if !statusCodeChanged {
		headerMap := reqDefault.Header
		// Remove headers which are *very* unlikely embedded into a template
		headerMap.Del("Content-Type")
		headerMap.Del("User-Agent")

		// Add headers which are likely embedded into a template
		headerMap["Host"] = []string{reqDefault.Host}
		if headerMap["X-Forwarded-For"] == nil {
			headerMap["X-Forwarded-For"] = []string{"added-by-TInjA"}
		}
		if headerMap["Origin"] == nil {
			headerMap["Origin"] = []string{"added-by-TInjA"}
		}
		msg = fmt.Sprintf("Found %d headers\n", len(headerMap))
		PrintVerbose(msg, NoColor, 2)
		for k, v := range headerMap {
			msg = fmt.Sprintln("Analyzing header ", k, " => ", v)
			PrintVerbose(msg, NoColor, 1)
			userInputCounter++
			repParam, statusCodeChanged := analyze(k, HEADER, u)
			if boolReport {
				repParam.Type = "Header"
				repParam.Name = k
				repParam.DefaultValues = v
				repWebpage.IsVulnerable = repWebpage.IsVulnerable || repParam.IsVulnerable
				requestError := false
				for _, rr := range repParam.Requests {
					requestError = rr.Error != ""
					if requestError {
						break
					}
				}
				// only add repParam to the report, if the param is vulnerable or if there are error messages or if there are reflections
				if repParam.IsVulnerable || len(repParam.ErrorMessages) > 0 || len(repParam.Reflections) > 0 || requestError {
					repWebpage.Parameters = append(repWebpage.Parameters, repParam)
					if repWebpage.Certainty != certaintyHigh && repWebpage.Certainty != certaintyMedium && repParam.Certainty == certaintyLow {
						repWebpage.Certainty = certaintyLow
						report.Low += 1
						report.SuspectedInjections += 1
					} else if repWebpage.Certainty != certaintyHigh && repParam.Certainty == certaintyMedium {
						repWebpage.Certainty = certaintyMedium
						report.Medium += 1
						report.SuspectedInjections += 1
					} else if repParam.Certainty == certaintyHigh {
						repWebpage.Certainty = certaintyHigh
						report.High += 1
						report.SuspectedInjections += 1
					}
				}
				if statusCodeChanged {
					repWebpage.ErrorMessages = append(repWebpage.ErrorMessages, "Status code changed. Skipping this URL.")
					break
				}
			}
		}
	}

	if boolReport {
		// only add default request and response to report if there are error messages or if a parameter was added
		if len(repWebpage.Parameters) > 0 || len(repWebpage.ErrorMessages) > 0 {
			reqDefault, _ := buildRequest(u, config)
			dumpReqBytes, _ := httputil.DumpRequest(reqDefault, true)
			repWebpage.ReportDefault.Request = string(dumpReqBytes)
			repWebpage.ReportDefault.Response = responseDump
		}
		repWebpage.ReportDefault.StatusCode = defaultStatus
	}

	return repWebpage
}

func getPostParams() map[string]string {
	m := make(map[string]string)
	for _, p := range strings.Split(config.Data, "&") {
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		} else {
			Print(p+" cannot be split by =\n", Yellow)
		}
	}
	return m
}

/* Setting proxy with specified proxyURL and proxyCertPath */
func setProxy() {
	proxyURL, err := url.Parse(config.ProxyURL)
	if err != nil {
		msg := "setProxy: url.Parse: " + err.Error() + "\n"
		PrintFatal(msg)
	}
	if config.ProxyCertPath != "" {
		caCert, err := os.ReadFile(config.ProxyCertPath)
		if err != nil {
			msg := "setProxy: os.ReadFile: " + err.Error() + "\n"
			PrintFatal(msg)
		}
		//caCertPool,err := x509.SystemCertPool() // throws error: "crypto/x509: system root pool is not available on Windows"
		caCertPool = x509.NewCertPool()

		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			msg := "setProxy: could not append cert\n"
			PrintFatal(msg)
		}

		tlsConfig := &tls.Config{RootCAs: caCertPool,
			InsecureSkipVerify: true}
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: tlsConfig}

		err = http2.ConfigureTransport(tr)
		if err != nil {
			msg := fmt.Sprintf("setProxy: Cannot switch to HTTP2: %s\n", err.Error())
			PrintFatal(msg)
		}

		http.DefaultTransport = tr
	} else {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: tlsConfig}

		http.DefaultTransport = tr
	}
}
