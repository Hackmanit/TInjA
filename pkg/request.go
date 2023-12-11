package pkg

/* Documentation

The `request.go` file contains several functions related to making HTTP requests and handling connections. Let's go through each function:

The `buildRequest` function takes a URL string (`u`) and a `structs.Config` object (`conf`) as input and constructs an `http.Request` object based on the provided information. It determines whether to issue a GET or POST request based on the presence of data in the configuration. It sets headers and query parameters according to the configuration values, and returns the constructed request.

The `setQuery` function takes an `http.Request` object (`req`), a key string (`key`), and a payload string (`payload`) as input. It sets the specified key-value pair as a query parameter in the request's URL.

The `setPost` function takes an `http.Request` object (`req`), a key string (`key`), and a payload string (`payload`) as input. It sets the specified key-value pair as a POST parameter in the request's body. It returns a new `http.Request` object with the updated body.

The `setHeader` function takes an `http.Request` object (`req`), a key string (`key`), and a payload string (`payload`) as input. It sets the specified key-value pair as a header in the request.

The `dialConnection` function takes a scheme string (`scheme`) and a host string (`host`) as input. It establishes a network connection using the appropriate protocol (HTTP or HTTPS) based on the scheme. It returns a `net.Conn` object representing the connection.

The `doRequest` function takes an `http.Request` object (`req`) as input and performs an HTTP request using the provided request object. It handles scenarios where there are host header manipulations or the request needs to be sent through a proxy. It returns the response body, headers, status code, dump (string representation of the response), and any error that occurred during the request.

Overall, the file contains functions for building requests, setting query parameters and headers, establishing network connections, and performing HTTP requests with various configurations.

*********/

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"example/user/tinja/pkg/structs"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func buildRequest(u string, conf structs.Config) (req *http.Request, err error) {
	// Shall a GET or POST request be issued?
	if conf.Data != "" {
		req, err = http.NewRequest("POST", u, bytes.NewBufferString(conf.Data))
		if err != nil {
			Print("buildRequest: "+err.Error()+"\n", Red)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			Print("buildRequest: "+err.Error()+"\n", Red)
			return
		}
	}

	req.Header.Set("User-Agent", config.UserAgent)

	/* add headers */
	for _, h := range conf.Headers {
		// hv is a Slice where hv[0] = header name and hv[1] = header value
		hv := strings.SplitN(h, ":", 2)
		if len(hv) != 2 {
			msg := fmt.Sprintf("Could not split %s into header name and value\n", h)
			PrintVerbose(msg, Yellow, 1)
			continue
		}
		if strings.EqualFold(hv[0], "Host") {
			msg := fmt.Sprintf("Overwriting Host:%s with Host:%s\n", req.URL.Host, hv[1])
			PrintVerbose(msg, NoColor, 2)
			req.Host = hv[1]
		} else {
			// check if header already exists
			if val := req.Header.Get(hv[0]); val != "" {
				msg := fmt.Sprintf("Overwriting %s:%s with %s\n", hv[0], val, hv[1])
				PrintVerbose(msg, Red, 2)
			}
			// set header
			req.Header.Set(hv[0], hv[1])
		}
	}
	/***************/

	/* add query parameters */
	q := req.URL.Query()
	for _, p := range conf.Parameters {
		// pv is a Slice where pv[0] = parameter name and pv[1] = parameter value
		pv := strings.SplitN(p, "=", 2)
		if len(pv) != 2 {
			msg := fmt.Sprintf("Could not split %s into parameter name and value\n", p)
			PrintVerbose(msg, Yellow, 1)
			continue
		}
		// check if parameter already exists
		if val := req.URL.Query().Get(pv[0]); val != "" {
			msg := fmt.Sprintf("Overwriting %s=%s with %s\n", pv[0], val, pv[1])
			PrintVerbose(msg, NoColor, 2)
		}
		// set parameter
		q.Set(pv[0], pv[1])
	}
	req.URL.RawQuery = q.Encode()
	/***********************/

	return
}

func setQuery(req *http.Request, key string, payload string) {
	q := req.URL.Query()
	q.Set(key, payload)
	req.URL.RawQuery = q.Encode()
}

func setPost(req *http.Request, key string, payload string) (*http.Request, error) {
	conf := config

	bodyString := ""
	for k, v := range postParams {
		bodyString += k + "="
		if k == key {
			bodyString += url.QueryEscape(payload)
		} else {
			bodyString += v + "?"
		}
	}
	bodyString = strings.TrimSuffix(bodyString, "?")
	conf.Data = bodyString
	return buildRequest(req.URL.String(), conf)
}

func setHeader(req *http.Request, key string, payload string) {
	if key == "Host" { // Host header tests need a custom implementation as the net/http library doesn't allow host header manipulations
		req.Host = payload + ".com"
	} else {
		req.Header.Set(key, payload)
	}
}

// dialConnection establishes the appropriate connection based on the scheme (http or https)
func dialConnection(scheme, host string) (net.Conn, error) {
	if scheme == "https" {
		return tls.Dial("tcp", host, &tls.Config{RootCAs: caCertPool,
			InsecureSkipVerify: true})
	}
	return net.Dial("tcp", host)
}

func doRequest(req *http.Request) (body string, headers http.Header, status int, dump string, err error) {
	var resp *http.Response
	var conn net.Conn

	err = limiter.Wait(context.Background())
	if err != nil {
		msg := "doRequest rate limiter error: " + err.Error() + "\n"
		Print(msg, Red)
	}

	// if there are no host header manipulations, do a standard client.do
	if req.Host == req.URL.Host || req.Host == req.URL.Hostname() {
		// Do Request
		resp, err = http.DefaultClient.Do(req)

		if err != nil {
			return
		}
		// if there is a host header manipulation, do the following in order to circumvent net/https limitations
	} else {
		// Determine if HTTP or HTTPS is used
		scheme := "http"
		if req.URL.Scheme == "https" {
			scheme = "https"
			if !strings.Contains(req.URL.Host, ":") {
				req.URL.Host += ":443"
			}
		} else {
			if !strings.Contains(req.URL.Host, ":") {
				req.URL.Host += ":80"
			}
		}

		// Construct the request line
		requestLine := []byte(fmt.Sprintf("%s %s %s\r\n", req.Method, req.URL.String(), req.Proto))

		// Read the request body, if present
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, err = io.ReadAll(req.Body)
			if err != nil {
				err = errors.New("Error reading request body: " + err.Error())
				return
			}
		}

		/* Add headers which golang doesn't add before sending the actual request */
		if len(bodyBytes) > 0 {
			req.Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
		}
		req.Header.Set("Host", req.Host)
		req.Header.Set("User-Agent", req.UserAgent())
		/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

		// Construct the headers
		var reqHeaders bytes.Buffer
		req.Header.Write(&reqHeaders)

		// Combine the request line, headers, and request body (if applicable)
		requestData := append(append(append(requestLine, reqHeaders.Bytes()...), []byte("\r\n")...), bodyBytes...)

		timeoutDuration := time.Duration(config.Timeout) * time.Second

		if config.ProxyURL != "" {
			var proxyURL *url.URL
			var proxyConn net.Conn
			// Connect to the proxy
			proxyURL, err = url.Parse(config.ProxyURL) // Replace with your proxy URL
			if err != nil {
				err = errors.New("Error parsing proxy URL:" + err.Error())
				return
			}

			proxyConn, err = dialConnection(proxyURL.Scheme, proxyURL.Host)
			if err != nil {
				err = errors.New("Error connecting to proxy:" + err.Error())
				return
			}
			defer proxyConn.Close()

			// Set a deadline for proxyConn
			err = proxyConn.SetDeadline(time.Now().Add(timeoutDuration))
			if err != nil {
				err = errors.New("Error setting proxyConn deadline:" + err.Error())
				return
			}

			// Send the request through the proxy
			proxyRequest := fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", req.URL.Host, req.URL.Host)
			_, err = proxyConn.Write([]byte(proxyRequest))
			if err != nil {
				err = errors.New("Error sending CONNECT request to proxy: " + err.Error())
				return
			}

			// Read the response from the proxy
			var n int
			proxyResponse := make([]byte, 4096)
			n, err = proxyConn.Read(proxyResponse)
			if err != nil {
				err = errors.New("Error reading response from proxy: " + err.Error())
				return
			}

			proxyResponseStr := string(proxyResponse[:n])
			PrintVerbose("Proxy Response: "+proxyResponseStr+"\n", NoColor, 2)

			// If its a http URL, the proxyconn can be used as is. Otherwise a TLS handshake needs to be established
			if scheme == "http" {
				conn = proxyConn
			} else {
				// Establish a TLS connection to the target server via the proxy
				targetConn := tls.Client(proxyConn, &tls.Config{
					RootCAs:            caCertPool,
					InsecureSkipVerify: true, // Skip certificate verification (for demonstration purposes only)
					ServerName:         req.URL.Hostname(),
				})

				err = targetConn.Handshake()
				if err != nil {
					err = errors.New("Error establishing TLS connection to target server: " + err.Error())
					return
				}

				conn = targetConn
			}
		} else {
			// Send the request using net.Dial or tls.Dial
			conn, err = dialConnection(scheme, req.URL.Host)
			if err != nil {
				err = errors.New("Error establishing connection: " + err.Error())
				return
			}
			defer conn.Close()
		}
		err = conn.SetDeadline(time.Now().Add(timeoutDuration))
		if err != nil {
			err = errors.New("Error setting conn deadline:" + err.Error())
			return
		}

		_, err = conn.Write(requestData)
		if err != nil {
			err = errors.New("Error sending request: " + err.Error())
			return
		}

		// Read the response
		responseData := make([]byte, 4096)
		var n int
		n, err = conn.Read(responseData)
		if err != nil {
			err = errors.New("Error reading response: " + err.Error())
			return
		}

		response := string(responseData[:n])

		// Parse the response
		resp, err = http.ReadResponse(bufio.NewReader(bytes.NewReader([]byte(response))), req)
		if err != nil {
			err = errors.New("Error reading response from target server: " + err.Error())
			return
		}
	}

	if boolReport {
		var dumpBytes []byte
		dumpBytes, err = httputil.DumpResponse(resp, true)
		if err != nil {
			// Sometimes DumpResponse throws an error like "http: ContentLength=54 with Body length 0". However the Content-Length Header and the body length are correct...
			if !strings.Contains(err.Error(), "unexpected EOF") {
				PrintVerbose("Error dumping response: "+err.Error()+"\n", Yellow, 1)
			}
			err = nil
			dumpBytes, _ = httputil.DumpResponse(resp, false)
		}
		dump = string(dumpBytes)
		dumpBytes = nil
	}

	// Read Response
	var resBody []byte
	resBody, err = io.ReadAll(resp.Body)
	if err != nil {
		if strings.Contains(err.Error(), "unexpected EOF") {
			err = nil
		} else {
			err = errors.New("Error reading response body: " + err.Error())
			return
		}
	}

	body = string(resBody)
	resBody = nil
	headers = resp.Header
	status = resp.StatusCode

	return
}
