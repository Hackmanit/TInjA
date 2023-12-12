package pkg

/* Documentation

The `reflection.go` file contains functions related to analyzing reflections in the response body and headers during template injection scanning. Let's go through each function:

The `analyzeReflection` function takes a name string (`name`), a type integer (`typ`), and a URL string (`u`) as input. It builds a request based on the provided URL and scan configuration. It then sends the request with a token value injected into the specified parameter (query parameter, POST parameter, or header) and analyzes the response for reflections. If a reflection is found in the response headers or body, it adds the reflection information to the `reflections` slice. It also performs additional reflection checks on the URLs specified in the `config.URLsReflection` and adds reflections if found. The function returns a boolean indicating whether any reflections were found and an error if any occurred.

The `doReflectionCheckRequest` function takes a URL string (`u`) as input. It builds a request based on the provided URL and scan configuration and sends the request. It returns the response body as a string.

The `addReflections` function takes a response body string (`body`), a token string (`token`), a count integer (`count`), and a URL reflection string (`urlReflection`) as input. It splits the response body by the token and adds each reflection's preceding and subsequent parts to the `reflections` slice. It also appends the URL reflection if provided. This function is used to process reflections found in the response body and add them to the `reflections` slice.

These functions work together to analyze reflections in the response body and headers during template injection scanning. They identify and collect reflections, which are stored in the `reflections` slice for further processing and reporting.

********/

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Hackmanit/TInjA/pkg/structs"
)

// TODO anstelle von DefaultReq ein neues Request stellen  mit Nonce und überprüfen, dass response bis auf anderen input gleich ist
func analyzeReflection(name string, typ int, u string) (bool, error) {
	req, err := buildRequest(u, config)
	if err != nil {
		msg := "Error: analyzeReflection: buildRequest: " + err.Error()
		Print(msg+"\n", Red)
		return false, err
	}

	tokenlength := 16
	// LengthLimit == 0 => unlimited length
	if config.LengthLimit > 0 && config.LengthLimit < tokenlength {
		tokenlength = config.LengthLimit
	}
	token := getToken(tokenlength)
	typString := ""
	switch typ {
	case QUERY:
		setQuery(req, name, token)
		typString = "query parameter"
	case POST:
		req, _ = setPost(req, name, token)
		typString = "POST parameter"
	case HEADER:
		setHeader(req, name, token)
		typString = "header"
	}
	body, headers, _, _, err := doRequest(req)
	if err != nil {
		PrintVerbose("Error: Analyze Reflection: "+err.Error()+"\n", Yellow, 1)
		return false, err
	}
	for header, headervaluearray := range headers {
		for _, headervalue := range headervaluearray {
			if strings.Contains(headervalue, token) {
				msg := fmt.Sprintln("Value ", token, " of "+typString+" ", name, " is being reflected in the ", header, " header")
				PrintVerbose(msg+"\n", Cyan, 1)
				precedingSubsequent := strings.SplitN(headervalue, token, 2)
				reflections = append(reflections, structs.Reflection{
					ReflectionType: structs.ReflectionHeader,
					HeaderName:     header,
					Preceding:      precedingSubsequent[0],
					Subsequent:     precedingSubsequent[1],
				})
			}
		}
	}
	if count := strings.Count(body, token); count > 0 {
		msg := fmt.Sprintln("Value ", token, " of "+typString+" ", name, " is being reflected "+strconv.Itoa(count)+" time(s) in the response body")
		PrintVerbose(msg+"\n", Cyan, 1)
		addReflections(body, token, "")
	}
	for _, u := range config.URLsReflection {
		body = doReflectionCheckRequest(u)
		if count := strings.Count(body, token); count > 0 {
			msg := fmt.Sprintln("Value ", token, " of "+typString+" ", name, " is being reflected "+strconv.Itoa(count)+" time(s) in the response body of "+u)
			PrintVerbose(msg+"\n", Cyan, 1)
			addReflections(body, token, u)
		}
	}

	return len(reflections) > 0, nil
}

func doReflectionCheckRequest(u string) string {
	config.Data = ""
	config.Parameters = []string{}
	req, _ := buildRequest(u, config)
	body, _, _, _, err := doRequest(req)
	if err != nil {
		PrintVerbose("Error: doReflectionRequest: "+err.Error()+"\n", Yellow, 1)
	}
	return body
}

func addReflections(body string, token string, urlReflection string) {
	for strings.Contains(body, token) { // loop through every occurance of token
		precedingSubsequent := strings.SplitN(body, token, 2)
		preceding := precedingSubsequent[0]
		subsequent := precedingSubsequent[1]
		// reduce length of preceding/subsequent chars according to the config
		if len(preceding) > config.PrecedingLength {
			preceding = preceding[len(preceding)-config.PrecedingLength:]
		}
		if len(subsequent) > config.SubsequentLength {
			subsequent = subsequent[0:config.SubsequentLength]
		}
		reflections = append(reflections, structs.Reflection{
			ReflectionType: structs.ReflectionBody,
			Preceding:      preceding,
			Subsequent:     subsequent,
			ReflectionURL:  urlReflection,
		})
		body = precedingSubsequent[1]
	}
}
