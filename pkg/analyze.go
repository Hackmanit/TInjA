package pkg

/* Documentation

The `analyze.go` file contains a package called `pkg` that provides functionality for analyzing template injection vulnerabilities. Let's go through each method and describe what it does:

- `analyze(name string, typ int, u string) reportParameter`: This method is the entry point for the template injection vulnerability analysis. It takes in the name, type, and URL as parameters. It performs the analysis by calling other methods, such as `analyzeReflection`, `detectTemplateInjection`, and `identifyTemplateEngine`. It returns a `reportParameter` struct containing the analysis results.

- `resetEverything()`: This method resets all previous analysis results. It sets various variables to their initial values, including `errorShown`, `reflections`, and `polyglotMap`.

- `detectTemplateInjection(name string, typ int, u string) (bool, []reportRequest)`: This method detects template injection vulnerabilities. It takes in the name, type, and URL as parameters. It checks if triggering an error is possible with a universal error polyglot. It then checks if there are reflections or if the input is being reflected. Next, it checks if one of the three universal non-error polyglots is processed, and if so, it sends the polyglot and checks the response code. If the type is a header and the name is "host," it tries more polyglots if needed. It returns a boolean indicating if a vulnerability is detected and a slice of `reportRequest` containing the details of the requests made.

- `identifyTemplateEngine(name string, typ int, u string) (string, []reportRequest, error)`: This method identifies the template engine used in the vulnerability. It takes in the name, type, and URL as parameters. If the engine was already identified during the detection phase, it returns its name. Otherwise, it sends polyglots and checks the response to identify the engine. It returns the identified engine's name, a slice of `reportRequest` containing the details of the requests made, and an error if any.

- `sendPolyglot(name string, typ int, polyglot string, u string) (int, reportRequest)`: This method sends a polyglot request to test for vulnerabilities. It takes in the name, type, polyglot, and URL as parameters. It builds an HTTP request, sets the request parameters based on the type, and performs the request. It returns the response code and a `reportRequest` struct containing the details of the request and response.

- `checkInjectionIndicators(body string, headers http.Header, status int, polyglot string, u string, typ int) (int, string, error)`: This method checks for injection indicators in the response body, headers, and status. It takes in the response body, headers, status code, polyglot, URL, and type as parameters. It checks if the status code is different and if the default status code has changed. If there are reflections, it checks for indicators in the body. It also handles specific cases for different template engines. It returns the indicator code, conclusion message, and an error if any.

- `checkResponses(polyglot string, responses []string) int`: This method checks the responses from different template engines for a polyglot. It takes in the polyglot and a slice of responses as parameters. It matches the responses against the expected engine responses and returns the indicator code indicating if the responses match, and if an error was thrown or not.

- `getIdentifiedEngine() string`: This method returns the identified template engine. It checks the possible engines and returns the engine name if it is the only one identified.

- `hasPolyglotImpact(polyglot string) bool`: This method checks if a polyglot has an impact by comparing its responses across different template engines. It also considers the length limit set in the configuration.

- `runHTMLinHeadless(html string, url string) string`: This method runs HTML in a headless browser using the Rod library. It connects to the URL, sets the document content to the provided HTML, waits for stability, and returns the HTML of the page.

- `checkBodyIndicator(body string, polyglot string, reflection structs.Reflection) (string, string)`: This method checks for indicators in the response body. It takes in the body, polyglot, and a `Reflection` struct representing the reflection point. It handles specific cases for different template engines and returns the response and conclusion based on the indicators found.

- `checkForDistinctTemplateEngineResponse(polyglot string, stringBetween string)`: This method checks if the recieved response was definitely rendered by a template engine. This is the case, if stringBetween matches with an expected response of at least one template engine.

- `setTemplateEngine(names []string)`: This method sets the possible template engines based on the identified engine names. It takes in a slice of engine names as a parameter.

- `getAllPossibleEngines() string`: This method returns a string containing all the possible template engines based on the `possibleEngines` map.

These methods work together to perform the template injection vulnerability analysis, including detection and identification of the template engine used.

*****************/

import (
	"bytes"
	"errors"
	"example/user/tinja/pkg/structs"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"moul.io/http2curl"
)

// Variable declarations
var errorShown bool
var reflected bool
var onlyErrorResponses bool
var recievedModifiedRenderedResponse bool
var reflections []structs.Reflection
var polyglotMap map[string]string
var notTested = "notTested"
var tested = "tested"
var indicatorNone = 0
var indicatorUnmodified = 1
var indicatorModified = 2
var indicatorError = 3
var indicatorIdentified = 4
var indicatorNotValid = 5
var rodBrowser *rod.Browser
var statusCodeChanged bool
var CounterPolyglotsGlobal = 0

func init() {
	rodBrowser = rod.New()
}

// analyze performs the template injection vulnerability analysis
func analyze(name string, typ int, u string) (reportParameter, bool) {
	// Initialize report parameter and request slices
	var repParam reportParameter
	var requestsDetect, requestsIdentify []reportRequest
	var detected bool
	var err error

	// Reset all previous analysis results
	resetEverything()

	// Analyze reflection
	reflected, err = analyzeReflection(name, typ, u)
	if boolReport && err != nil {
		repParam.ErrorMessages = append(repParam.ErrorMessages, err.Error())
	}

	// Detect template injection vulnerabilities
	detected, requestsDetect = detectTemplateInjection(name, typ, u)
	if boolReport && err != nil {
		repParam.ErrorMessages = append(repParam.ErrorMessages, err.Error())
	}

	// If template injection is detected, identify the template engine
	if detected {
		repParam.TemplateEngine, requestsIdentify, err = identifyTemplateEngine(name, typ, u)
		if boolReport && err != nil {
			repParam.ErrorMessages = append(repParam.ErrorMessages, err.Error())
		}
	}

	certainty := "None"
	if repParam.TemplateEngine == "unknown" && !recievedModifiedRenderedResponse {
		certainty = certaintyLow
	} else if repParam.TemplateEngine == "unknown" && recievedModifiedRenderedResponse {
		certainty = certaintyMedium
	} else if repParam.TemplateEngine != "" && !recievedModifiedRenderedResponse {
		// The host header has great limitations which characters are allowed and which not. This often leads to false positives for template engines which often responds with unmodified or error, like dot
		if typ == HEADER && strings.EqualFold(name, "host") {
			certainty = certaintyLow
		} else {
			certainty = certaintyMedium
		}
	} else if repParam.TemplateEngine != "" && recievedModifiedRenderedResponse {
		certainty = certaintyHigh
	}

	if boolReport {
		// Set report parameter values
		repParam.Reflections = reflections
		repParam.AreErrorsThrown = errorShown
		repParam.Requests = append(repParam.Requests, requestsDetect...)
		repParam.Requests = append(repParam.Requests, requestsIdentify...)
		repParam.Certainty = certainty
		repParam.IsVulnerable = repParam.TemplateEngine != ""
	}

	switch repParam.TemplateEngine {
	case "":
		PrintVerbose("No template engine could be detected\n\n", NoColor, 1)
	case "unknown":
		Print("A template engine was detected, but could not be identified (certainty: "+certainty+")\n\n", Green)
	default:
		Print(repParam.TemplateEngine+" was identified (certainty: "+certainty+")\n\n", Green)
	}

	return repParam, statusCodeChanged
}

// resetEverything resets all previous analysis results
func resetEverything() {
	errorShown = true
	// set onlyErrorResponses = false if a response is no error response
	onlyErrorResponses = true
	recievedModifiedRenderedResponse = false
	reflected = false
	statusCodeChanged = false

	reflections = []structs.Reflection{}

	possibleEngines = map[string]bool{}
	for _, engine := range engines {
		possibleEngines[engine.Name] = true
	}

	polyglotMap = map[string]string{
		err1:             notTested,
		err2:             notTested,
		err3:             notTested,
		err4:             notTested,
		err5:             notTested,
		errPython:        notTested,
		errRuby:          notTested,
		errDotnet:        notTested,
		errJava:          notTested,
		errJava2:         notTested,
		errPHP:           notTested,
		errJavascript:    notTested,
		errGolang:        notTested,
		errElixir:        notTested,
		nonerr1:          notTested,
		nonerr2:          notTested,
		nonerr3:          notTested,
		nonerrPython:     notTested,
		nonerrRuby:       notTested,
		nonerrDotnet:     notTested,
		nonerrJava:       notTested,
		nonerrPHP:        notTested,
		nonerrJavascript: notTested,
		nonerrGolang:     notTested,
		nonerrElixir:     notTested,
		nonerrIdent1:     notTested,
		nonerrIdent2:     notTested,
		nonerrIdent3:     notTested,
	}
}

// detectTemplateInjection detects template injection vulnerabilities
func detectTemplateInjection(name string, typ int, u string) (bool, []reportRequest) {
	var repRequests []reportRequest
	var repRequest reportRequest
	var respCode int
	success := false

	// Check if triggering error is possible with universal error polyglot
	respCode, repRequest = sendPolyglot(name, typ, err1, u, false)
	errorShown = respCode == indicatorError || respCode == indicatorIdentified
	polyglotMap[err1] = tested
	success = errorShown
	if boolReport {
		repRequests = append(repRequests, repRequest)
	}

	// Check if no errors are thrown and input is not being reflected
	if !errorShown && !reflected {
		PrintVerbose("No errors are thrown and input is not being reflected.\n", NoColor, 1)
		return false, repRequests
	}

	// Check if one of the three universal non error polyglots is processed - if input is reflected
	for _, polyglot := range []string{nonerr1, nonerr2, nonerr3} {
		// return if engine was already identified. But continue if input is being reflected and no definitive rendered response was received yet (to diminish false positives)
		if getIdentifiedEngine() != "" && getIdentifiedEngine() != "unkown" && (!reflected || (reflected && recievedModifiedRenderedResponse)) {
			return true, repRequests
		}
		if respCode, repRequest = sendPolyglot(name, typ, polyglot, u, false); respCode != indicatorNone && respCode != indicatorUnmodified && respCode != indicatorNotValid {
			success = true
		}
		polyglotMap[polyglot] = tested
		if boolReport {
			repRequests = append(repRequests, repRequest)
		}
	}

	// 3 of the 4 Universal Detection Polyglots might not work due to character limitations of the host header.
	// Hence try more polyglots if needed.
	if typ == HEADER && strings.EqualFold(name, "host") && !success {
		// return if engine was already identified
		if getIdentifiedEngine() != "" && getIdentifiedEngine() != "unkown" && (!reflected || (reflected && recievedModifiedRenderedResponse)) {
			return true, repRequests
		}
		for polyglot := range polyglotMap {
			if respCode, repRequest = sendPolyglot(name, typ, polyglot, u, false); respCode != indicatorNone && respCode != indicatorUnmodified && respCode != indicatorNotValid {
				success = true
			}
			polyglotMap[polyglot] = tested
			if boolReport {
				repRequests = append(repRequests, repRequest)
			}
		}
	}

	if statusCodeChanged {
		return false, repRequests
	}

	return success, repRequests
}

func identifyTemplateEngine(name string, typ int, u string) (string, []reportRequest, error) {
	var repRequests []reportRequest
	var repRequest reportRequest
	var err error
	// If Engine was already identified during detection phase, return its name
	if getIdentifiedEngine() != "" && getIdentifiedEngine() != "unkown" && (!reflected || (reflected && recievedModifiedRenderedResponse)) && !onlyErrorResponses {
		return getIdentifiedEngine(), nil, err
	}
	Print("\nA template engine was successfully detected and is now being identified.\n", NoColor)
	for polyglot := range polyglotMap {
		if statusCodeChanged {
			break
		}
		if hasPolyglotImpact(polyglot) {
			_, repRequest = sendPolyglot(name, typ, polyglot, u, false)
			polyglotMap[polyglot] = tested
			if boolReport {
				repRequests = append(repRequests, repRequest)
			}
		}
		switch engine := getIdentifiedEngine(); engine {
		case "":
			continue
		case "unkown":
			return "unknown", repRequests, err
		default:
			if !onlyErrorResponses && (!reflected || (reflected && recievedModifiedRenderedResponse)) {
				return engine, repRequests, err
			} else {
				continue
			}
		}
	}
	return getAllPossibleEngines(), repRequests, err
}

func sendPolyglot(name string, typ int, polyglot string, u string, backslashed bool) (int, reportRequest) {
	// Initialize report request and response indicator variables
	var repRequest reportRequest
	var respIndicator int

	// Print possible engines
	PrintVerbose("Possible Engines: "+getAllPossibleEngines()+"\n", Yellow, 2)

	// Build the HTTP request
	req, err := buildRequest(u, config)
	if err != nil {
		msg := "Error: sendPolyglot: buildRequest: " + err.Error()
		Print(msg+"\n", Red)
		if boolReport {
			repRequest.Conclusion = msg
			repRequest.Error = err.Error()
		}
		return indicatorNotValid, repRequest
	}

	// Set request parameters based on the request type
	switch typ {
	case QUERY:
		setQuery(req, name, polyglot)
	case POST:
		req, _ = setPost(req, name, polyglot)
	case HEADER:
		setHeader(req, name, polyglot)
	}

	var dumpReqBytes []byte
	var bodyBackup []byte

	if req.Body != nil {
		// Backup the request body, because it can only be read once
		bodyBackup, err = io.ReadAll(req.Body)
		if err != nil {
			msg := "Error: sendPolyglot: " + err.Error()
			Print(msg+"\n", Red)
			repRequest.Error = err.Error()
			return indicatorNotValid, repRequest
		}
		if boolReport {
			// Restore the request body for the first use
			req.Body = io.NopCloser(bytes.NewReader(bodyBackup))
			// Dump the request including the body
			dumpReqBytes, _ = httputil.DumpRequest(req, true)
		}
		// Restore the request body for the second use
		req.Body = io.NopCloser(bytes.NewReader(bodyBackup))
	} else {
		// Dump the request without the body
		dumpReqBytes, _ = httputil.DumpRequest(req, false)
	}

	if !backslashed {
		CounterPolyglotsGlobal++
	}

	// Perform the HTTP request and retrieve the response body, headers, status, and dumped response
	body, headers, status, dumpRes, err := doRequest(req)
	if err != nil {
		PrintVerbose("Error: sendPolyglot: "+err.Error()+"\n", Yellow, 1)
		if boolReport {
			repRequest.Error = err.Error()
		}
		return indicatorNotValid, repRequest
	}

	// Populate the report request with the response details and dumped request
	if boolReport {
		repRequest.Response = dumpRes
		repRequest.Request = string(dumpReqBytes)
	}

	// Add the request as curl command to the report
	command, err := http2curl.GetCurlCommand(req)
	if err != nil {
		PrintVerbose("Error: sendPolyglot: "+err.Error()+"\n", Yellow, 1)
		repRequest.Error = err.Error()
	}
	commandFixed := strings.Replace(command.String(), "-d ''", "-d '"+string(bodyBackup)+"'", 1)
	if boolReport {
		repRequest.CurlCommand = commandFixed
	}
	PrintVerbose("Curl command: "+commandFixed+"\n", NoColor, 2)

	if boolReport {
		repRequest.Polyglot = polyglot
	}

	// Check injection indicators in the response body, headers, and status
	respIndicator, repRequest.Conclusion, err = checkInjectionIndicators(body, headers, status, polyglot, u, typ, backslashed, name)
	if err != nil {
		PrintVerbose("Error: sendPolyglot: "+err.Error()+"\n", Yellow, 1)
		if boolReport {
			repRequest.Error = err.Error()
		}
	}
	if respIndicator == indicatorIdentified {
		recievedModifiedRenderedResponse = true
	}
	if respIndicator != indicatorError && respIndicator != indicatorNotValid && !backslashed {
		onlyErrorResponses = false
	}
	return respIndicator, repRequest
}

func checkInjectionIndicators(body string, headers http.Header, status int, polyglot string, u string, typ int, backslashed bool, name string) (int, string, error) {
	var conclusion, response string
	if typ == HEADER && strings.HasPrefix(strconv.Itoa(status), "4") {
		msg := "The polyglot " + polyglot + " was rejected with a " + strconv.Itoa(status)
		PrintVerbose(msg+"\n", Cyan, 1)
		return indicatorNotValid, msg, nil
	}
	// check1: is the status code different and a 5xx status code. If true, check if the default status code has changed!
	if defaultStatus != status && strings.HasPrefix(strconv.Itoa(status), "5") {
		req, err := buildRequest(u, config)
		if err != nil {
			msg := "Error: checkInjectionIndicators: buildRequest: " + err.Error()
			Print(msg+"\n", Red)
			return indicatorNotValid, msg, err
		}
		_, _, checkStatus, _, err := doRequest(req)
		if err != nil {
			checkResponses(polyglot, []string{respError}, backslashed)
			conclusion = fmt.Sprintf("Couldn't connect to URL: %s", err.Error())
			Print(conclusion+"\n", Red)
			return indicatorNotValid, conclusion, err
		}
		if checkStatus != defaultStatus {
			conclusion = "The default status code changed to " + strconv.Itoa(status) + ". Skipping this URL."
			Print(conclusion+"\n", Red)
			statusCodeChanged = true
			return indicatorNotValid, conclusion, errors.New(conclusion)
		} else {
			// check if backshlashed also triggers error
			if !backslashed {
				if indicator, _ := sendPolyglot(name, typ, backslashPolyglot(polyglot), u, true); indicator == indicatorError {
					conclusion = "The backshlashed polyglot also throws an error; Therefore the error is most likely not thrown by a template engine."
					return indicatorNotValid, conclusion, nil
				}
			}

			checkResponses(polyglot, []string{respError}, backslashed)
			conclusion = "The polyglot " + polyglot + " triggered an error: Status Code " + strconv.Itoa(status)
			if !backslashed {
				Print(conclusion+"\n", Yellow)
			} else {
				PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
			}
			return indicatorError, conclusion, nil
		}
	} else if reflected {
		var err error
		responses := []string{} // there might be multiple reflection points. Hence, save all reflections in an array
		reflectionBodies := make(map[string]string)
		backslashedErrorsToo := false
		backslashedSend := false
		for _, reflection := range reflections {
			switch reflection.ReflectionType {
			case structs.ReflectionBody:
				urlHeadless := u
				bodyToCheck := body
				if reflection.ReflectionURL != "" {
					if reflectionBodies[reflection.ReflectionURL] == "" { // check if body was already requested. This saves multiple requests to the same URL, if the input is reflected multiple times
						reflectionBodies[reflection.ReflectionURL] = doReflectionCheckRequest(reflection.ReflectionURL)
					}
					bodyToCheck = reflectionBodies[reflection.ReflectionURL]
					urlHeadless = reflection.ReflectionURL
				}
				response, conclusion = checkBodyIndicator(bodyToCheck, polyglot, reflection, backslashed, name, typ, u)
				if response == respContinue {
					continue
				} else if response == respIdentified {
					return indicatorIdentified, conclusion, nil
				}
				if !backslashed && response == respUnmodified && config.CSTI && strings.Contains(bodyToCheck, "</script>") {
					var bodyNew string
					bodyNew, err = runHTMLinHeadless(bodyToCheck, urlHeadless)
					response, conclusion = checkBodyIndicator(bodyNew, polyglot, reflection, backslashed, name, typ, u)
				}
				if !backslashed && response == respError && !backslashedSend {
					// check if backshlashed also triggers error
					if indicator, _ := sendPolyglot(name, typ, backslashPolyglot(polyglot), u, true); indicator == indicatorError {
						backslashedErrorsToo = true
					}
					backslashedSend = true
				}
				// continue if both polyglot and backslashed polyglot throw an error
				if response == respContinue || response == respError && backslashedErrorsToo {
					continue
				} else if response == respIdentified {
					return indicatorIdentified, conclusion, err
				}
				responses = append(responses, response)
			case structs.ReflectionHeader:
				conclusion = ""
				for _, headerValue := range headers.Values(reflection.HeaderName) {
					var conclusionAppend string
					response, conclusionAppend = checkBodyIndicator(headerValue, polyglot, reflection, backslashed, name, typ, u)
					conclusion += conclusionAppend
					responses = append(responses, response)
				}
			default:
				Print("Unknown ReflectionType: "+reflection.ReflectionType+"\n", Red)
			}
		}
		var printResponses []string
		for _, v := range responses {
			if len(v) > 100 {
				printResponses = append(printResponses, v[:100]+"(longer than 100 characters)")
			} else {
				printResponses = append(printResponses, v)
			}
		}
		msg := "The polyglot " + polyglot + " returned the response(s) " + fmt.Sprint(printResponses)
		conclusion = conclusion + msg
		if !backslashed {
			Print(msg+"\n", Cyan)
		} else {
			PrintVerbose("Backslashed: "+msg+"\n", Cyan, 2)
		}
		return checkResponses(polyglot, responses, backslashed), conclusion, err
	} else {
		PrintVerbose("The polyglot "+polyglot+" did not trigger an error and input is being not reflected\n", NoColor, 2)
	}
	conclusion = "No indicator could be identified"
	return indicatorNone, conclusion, nil
}

func backslashPolyglot(polyglot string) string {
	var result strings.Builder
	for _, char := range polyglot {
		result.WriteString("\\")
		result.WriteRune(char)
	}
	return result.String()
}

func checkResponses(polyglot string, responses []string, backslashed bool) int {
	matchGlobal := false
	onlyUnmodified := true
	errorThrown := false
	for _, engine := range engines {
		match := false
		for _, response := range responses {
			encodedUnmodified, _ := isEncoded(response, polyglot)
			if response != respUnmodified && !encodedUnmodified {
				onlyUnmodified = false
			}
			if response == respError {
				errorThrown = true
			}
			encoded, _ := isEncoded(response, engine.Polyglots[polyglot])
			// Skip if no errors are shown and the engine's anticipated response would be an error AND the response is either empty or unmodified
			if engine.Polyglots[polyglot] == respError && !errorShown && (response == respUnmodified || response == respEmpty || encoded) {
				match = true
				matchGlobal = true
				break
			}
			if response == engine.Polyglots[polyglot] || encoded || (engine.Polyglots[polyglot] == respUnmodified && encodedUnmodified) {
				match = true
				matchGlobal = true
				break
			}
			// if an engines polyglot contains Arbitrary chars, we need to do a regex check
			if strings.Contains(engine.Polyglots[polyglot], "ARBITRARY") {
				for i, poly := range []string{engine.Polyglots[polyglot], html.EscapeString(engine.Polyglots[polyglot]), url.QueryEscape(engine.Polyglots[polyglot]), url.PathEscape(engine.Polyglots[polyglot])} {
					splitted := strings.Split(poly, "ARBITRARY")
					arbitrary := ".*"
					if len(splitted) == 3 {
						arbitrary = `\w{` + splitted[1] + `}`
					}
					pattern := `(?s)^` + regexp.QuoteMeta(splitted[0]) + arbitrary + regexp.QuoteMeta(splitted[len(splitted)-1]) + `$`
					if matchExpr, _ := regexp.MatchString(pattern, response); matchExpr {
						match = true
						matchGlobal = true
						break
					} else if i == 0 {
						// both check html.UnescapeString(response) as well as html.EscapeString(engine.Polyglots[polyglot]), because escape and unescape might behave differently
						if matchExpr, _ := regexp.MatchString(pattern, html.UnescapeString(response)); matchExpr {
							match = true
							matchGlobal = true
							break
						}
					}
				}
			}
		}
		if !match && !backslashed {
			// Skip for first universal error polyglot, as it might mean, that errors are catched
			if polyglot == err1 && engine.Polyglots[polyglot] == respError {
				continue
			}
			possibleEngines[engine.Name] = false
		}
	}
	if !matchGlobal {
		return indicatorNone
	}
	if errorThrown {
		return indicatorError
	} else if onlyUnmodified {
		return indicatorUnmodified
	} else {
		return indicatorModified
	}
}

// returns empty string, if more than 1 engine is possible, returns unknown if no known engine is possible
func getIdentifiedEngine() string {
	identifiedEngine := "unkown"
	for engine, possible := range possibleEngines {
		if possible {
			if identifiedEngine == "unkown" {
				identifiedEngine = engine
			} else {
				return ""
			}
		}
	}

	return identifiedEngine
}

func hasPolyglotImpact(polyglot string) bool {
	if statusCodeChanged {
		return false
	}

	if config.LengthLimit > 0 && len(polyglot) > config.LengthLimit {
		PrintVerbose("Polyglot "+polyglot+" ("+strconv.Itoa(len(polyglot))+") is longer than lengthlimit ("+strconv.Itoa(config.LengthLimit)+") and is skipped.\n", NoColor, 2)
		return false
	}

	// Keep checking if only one engine is left, but only error repsonses have been returned
	if engine := getIdentifiedEngine(); engine != "" && onlyErrorResponses {
		for _, e := range engines {
			if e.Name == engine {
				if e.Polyglots[polyglot] == respError {
					return false
				} else {
					return true
				}
			}
		}
	}

	for polyglot2, tested2 := range polyglotMap {
		if polyglot == polyglot2 && tested2 == tested {
			return false
		}
	}

	response := ""

	// check if there are minimum 2 different responses for the polyglot
	for _, engine := range engines {
		if !possibleEngines[engine.Name] {
			continue
		}
		if response == "" {
			response = engine.Polyglots[polyglot]
		} else {
			if reflected && response != engine.Polyglots[polyglot] {
				return true
			} else {
				// if input is not reflected, one possible polyglot response must be an error, and another one not an error
				if (response == respError && engine.Polyglots[polyglot] != respError) || (response != respError && engine.Polyglots[polyglot] == respError) {
					return true
				}
			}
		}
	}

	return false
}

func runHTMLinHeadless(html string, url string) (string, error) {
	returnVal := html
	err := rod.Try(func() {
		returnVal = rodBrowser.MustConnect().MustPage(url).MustSetDocumentContent(html).Timeout(15 * time.Second).MustWaitStable().MustHTML()
	})
	if err != nil {
		Print("runHTMLinHeadless: Error opening page: "+err.Error()+"\n", Red)
	}
	err = rod.Try(func() {
		for _, p := range rodBrowser.MustPages() {
			p.MustClose()
		}
	})
	if err != nil {
		Print("runHTMLinHeadless: Error closing page: "+err.Error()+"\n", Red)
	}

	return returnVal, err
}

func checkBodyIndicator(body string, polyglot string, reflection structs.Reflection, backslashed bool, name string, typ int, u string) (string, string) {
	var conclusion string
	// Thymeleaf / ThymeleafInline specific
	if strings.Contains(body, "org.thymeleaf.exceptions") && strings.Contains(body, polyglot) {
		possible := []string{}
		for _, engine := range engines {
			if engine.Name == "Thymeleaf" || engine.Name == "Thymeleaf (Inline)" {
				if engine.Polyglots[polyglot] == respError {
					possible = append(possible, engine.Name)
				}
			}
		}
		conclusion = fmt.Sprint("The polyglot "+polyglot+" triggered a ", possible, " error message")
		if !backslashed {
			Print(conclusion+"\n", Green)
		} else {
			PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
		}
		setTemplateEngine(possible)
		return respIdentified, conclusion
	}
	switch stringBetween := between(body, reflection.Preceding, reflection.Subsequent); stringBetween {
	case BothMissing:
		conclusion = "The polyglot " + polyglot + " triggered an error: " + BothMissing
		if !backslashed {
			Print(conclusion+"\n", Yellow)
		} else {
			PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
		}
		return respError, conclusion
	case PrecedingMissing:
		conclusion = "The polyglot " + polyglot + " triggered a PrecedingMissing"
		if !backslashed {
			Print(conclusion+"\n", Yellow)
		} else {
			PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
		}
		// Angular.js specific
		if strings.Contains(body, "class=\"ng-binding\"") {
			conclusion = "Preceding was replaced with class=\"ng-binding\""
			Print(conclusion+"\n", Green)
			setTemplateEngine([]string{"AngularJS"})
			return respIdentified, conclusion
		}
		// Pystache specific
		if polyglot == err2 && strings.Contains(body, "%>"+reflection.Subsequent) {
			conclusion = "Preceding was removed and %> rendered"
			setTemplateEngine([]string{"Pystache"})
			return respIdentified, conclusion
		} else if polyglot == errJavascript && strings.Contains(body, reflection.Subsequent) {
			conclusion = "Preceding was removed and empty string rendered"
			Print(conclusion+"\n", Green)
			setTemplateEngine([]string{"Pystache"})
			return respIdentified, conclusion
		}
		// Mustache.php specific
		if strings.Contains(body, "Unclosed tag: ") && strings.Contains(body, " on line ") {
			conclusion = "The polyglot " + polyglot + " triggered a Mustache.PHP error message"
			Print(conclusion+"\n", Green)
			setTemplateEngine([]string{"Mustache.PHP"})
			return respIdentified, conclusion
		}
		return respContinue, conclusion
	case SubsequentMissing:
		conclusion = "The polyglot " + polyglot + " triggered a SubsequentMissing"
		Print(conclusion+"\n", Yellow)
		/* Velocity / Velocityjs / Cheetah3 + HoganJS + Pug specific */
		possible := []string{}
		for _, engine := range engines {
			switch engine.Name {
			// all 3 might remove everything after the polyglot in the same line
			case "Cheetah3", "Velocity", "VelocityJS":
				if !possibleEngines[engine.Name] {
					continue
				}
				if strings.Contains(body, reflection.Preceding+engine.Polyglots[polyglot]) {
					possible = append(possible, engine.Name)
					conclusion = "Subsequent was removed and " + engine.Polyglots[polyglot] + " rendered"
				}
			// HoganJS might remove the first subsequent character
			case "HoganJS":
				betweenHogan := between(body, reflection.Preceding, reflection.Subsequent[1:])
				if betweenHogan == "" {
					betweenHogan = respEmpty
				}
				if betweenHogan == engine.Polyglots[polyglot] {
					conclusion = "The first subsequent character was removed and " + engine.Polyglots[polyglot] + " rendered"
					Print(conclusion+"\n", Green)
					setTemplateEngine([]string{"HoganJS"})
					return respIdentified, conclusion
				}
			// Pug messes with the subsequent string. Because nonce => <nonce></nonce>
			case "Pug":
				pugsExpectedAnswer := engine.Polyglots[polyglot]
				// if len(engine.Polyglots[polyglot]) > 10, then it isn't "error", "unmodified" or so on
				if len(pugsExpectedAnswer) > 10 && pugsExpectedAnswer != respError && pugsExpectedAnswer != respEmpty && pugsExpectedAnswer != respUnmodified && strings.Contains(body, pugsExpectedAnswer) {
					conclusion = "Subsequent was modified and " + pugsExpectedAnswer + " rendered"
					Print(conclusion+"\n", Green)
					setTemplateEngine([]string{"Pug"})
					return respIdentified, conclusion
				}
			}
		}
		if len(possible) > 0 {
			Print(conclusion+"\n", Green)
			setTemplateEngine(possible)
			return respIdentified, conclusion
		}
		// Generic Error Detection
		lowerbody := strings.ToLower(body)
		if strings.Contains(lowerbody, "error") { // maybe use parseerror and syntaxerror instead
			conclusion = "The polyglot " + polyglot + " triggered an error message, because it contained the word error"
			if !backslashed {
				Print(conclusion+"\n", Green)
			} else {
				PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
			}
			return respError, conclusion
		}
		/***********/
		return respContinue, conclusion
	default:
		isEncoded, addition := isEncoded(stringBetween, polyglot)
		if !isEncoded {
			// DotLiquid specific
			if strings.Contains(stringBetween, "DotLiquid.Exceptions") {
				conclusion = "The polyglot " + polyglot + " triggered a DotLiquid error message"
				Print(conclusion+"\n", Green)
				setTemplateEngine([]string{"DotLiquid"})
				return respIdentified, conclusion
			}
			// Scriban / Scriban (liquid mode) specific
			if strings.Contains(stringBetween, "Scriban.Template") {
				possible := []string{}
				for _, engine := range engines {
					if engine.Name == "Scriban" || engine.Name == "Scriban (Liquid mode)" {
						if engine.Polyglots[polyglot] == respError {
							possible = append(possible, engine.Name)
						}
					}
				}
				conclusion = fmt.Sprint("The polyglot "+polyglot+" triggered a ", possible, " error message")
				Print(conclusion+"\n", Green)
				setTemplateEngine(possible)
				return respIdentified, conclusion
			}
			// Fluid specific
			if strings.Contains(stringBetween, "Fluid.ParseException") {
				conclusion = "The polyglot " + polyglot + " triggered a Fluid error message"
				Print(conclusion+"\n", Green)
				setTemplateEngine([]string{"Fluid"})
				return respIdentified, conclusion
			}
			// Generic Error Detection
			if strings.Contains(strings.ToLower(stringBetween), "error") || strings.Contains(strings.ToLower(stringBetween), "exception") || strings.Contains(strings.ToLower(stringBetween), "unexpected") {
				conclusion = "The polyglot " + polyglot + " triggered an error message, because the rendered response contained the word error, exception or unexpected"
				if !backslashed {
					Print(conclusion+"\n", Yellow)
				} else {
					PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
				}
				return respError, conclusion
			}
		}

		switch stringBetween {
		case polyglot:
			stringBetween = respUnmodified
			conclusion = "The polyglot " + polyglot + " was rendered unmodified"
			return stringBetween, conclusion
		case "":
			stringBetween = respEmpty
			conclusion = "The polyglot " + polyglot + " was rendered as empty string"
			return stringBetween, conclusion
		default:
			// check if the response was definitely rendered by a template engine. If it is simply the encoded polyglot, there's no need to check
			if !isEncoded {
				checkForDistinctTemplateEngineResponse(polyglot, stringBetween)
			}
		}

		var printBetween string
		if len(stringBetween) > 100 {
			printBetween = stringBetween[:100] + "(longer than 100 characters)"
		} else {
			printBetween = stringBetween
		}
		conclusion = "The polyglot " + polyglot + " was rendered in a modified way: [" + printBetween + "]" + addition
		if !backslashed {
			Print(conclusion+"\n", Yellow)
		} else {
			PrintVerbose("Backslashed: "+conclusion+"\n", Cyan, 2)
		}
		return stringBetween, conclusion
	}
}

func checkForDistinctTemplateEngineResponse(polyglot string, stringBetween string) {
	if !recievedModifiedRenderedResponse && reflected && stringBetween != respError && stringBetween != respEmpty && stringBetween != respUnmodified && stringBetween != respContinue {
		for _, engine := range engines {
			if engine.Polyglots[polyglot] == stringBetween {
				recievedModifiedRenderedResponse = true
			}
		}
	}
}

func setTemplateEngine(names []string) {
	for engine := range possibleEngines {
		match := false
		for _, name := range names {
			if engine == name {
				match = true
			}
		}
		possibleEngines[engine] = match
	}
}

func getAllPossibleEngines() string {
	engines := ""
	for engine, possible := range possibleEngines {
		if possible {
			if engines != "" {
				engines += ", "
			}
			engines += engine
		}
	}
	return engines
}
