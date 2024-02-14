[![Release](https://img.shields.io/github/release/Hackmanit/TInjA.svg?color=brightgreen)](https://github.com/Hackmanit/TInjA/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/Hackmanit/TInjA)](https://goreportcard.com/report/github.com/Hackmanit/TInjA)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Hackmanit/TInjA)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
# TInjA – the Template INJection Analyzer

TInjA is a CLI tool for testing web pages for **template injection** vulnerabilities.

It supports **44 of the most relevant template engines** (as of September 2023) for eight different programming languages.

TInjA was developed by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x).

- [Features](#features)
- [Supported Template Engines](#supported-template-engines)
- [Installation](#installation)
    - [Option 1: Prebuilt Binary](#option-1-prebuilt-binary)
    - [Option 2: Install Using Go](#option-2-install-using-go)
- [Usage](#usage)
    - [Specify Headers, Cookies, and POST Body](#specify-headers-cookies-and-post-body)
    - [Scan CSTI in Addition to SSTI](#scan-csti-in-addition-to-ssti)
    - [Generate a JSONL Report](#generate-a-jsonl-report)
    - [Use a Proxy](#use-a-proxy)
    - [Set a Ratelimit](#set-a-ratelimit)
- [Troubleshooting](#troubleshooting)
- [TODOs](#todos)
- [Background Information](#background-information)
- [License](#license)

## Features
- Automatic detection of template injection possibilities and identification of the template engine in use.
    - 44 of the most relevant template engines supported (see [Supported Template Engines](#supported-template-engines)).
    - Both **SSTI** and **CSTI** vulnerabilities are detected.
        - SSTI = server-side template injection
        - CSTI = client-side template injection
- Efficient scanning thanks to the usage of polyglots:
    - On average only five polyglots are sent to the web page until the template injection possibility is detected and the template engine identified.
- Pass crawled URLs to TInjA in JSONL format.
- Pass a raw HTTP request to TInjA.
- Set custom headers, cookies, POST parameters, and query parameters.
- Route the traffic through a proxy (e.g., Burp Suite).
- Configure Ratelimiting.

## Supported Template Engines
### .NET
- DotLiquid
- Fluid
- Razor Engine
- Scriban
### Elixir
- EEx
### Go
- html/template
- text/template
### Java
- Freemarker
- Groovy
- Thymeleaf
- Velocity
### JavaScript
- Angular.js
- Dot
- EJS
- Eta
- Handlebars
- Hogan.js
- Mustache
- Nunjucks
- Pug
- Twig.js
- Underscore
- Velocity.js
- Vue.js
### PHP
- Blade
- Latte
- Mustache.php
- Smarty
- Twig
### Python
- Chameleon
- Cheetah3
- Django
- Jinja2
- Mako
- Pystache
- SimpleTemplate Engine
- Tornado
### Ruby
- ERB
- Erubi
- Erubis
- Haml
- Liquid
- Mustache
- Slim

## Installation
### Option 1: Prebuilt Binary
Prebuilt binaries of TInjA are provided on the [releases page](https://github.com/Hackmanit/TInjA/releases).
### Option 2: Install Using Go
Requirements: go1.21 or higher
```bash
go install -v github.com/Hackmanit/TInjA@latest
```

## Usage
- Scan a single URL: `tinja url -u "http://example.com/"`
- Scan multiple URLs: `tinja url -u "http://example.com/" -u "http://example.com/path2"`
- Scan URLs provided in a file: `tinja url -u "file:/path/to/file"`
- Scan a single URL by passing a file with a raw HTTP request: `tinja raw -R "/path/to/file"`
- Scan URLs with additional information provided in a JSONL file: `tinja jsonl -j "/path/to/file"`
    - Each line of the JSONL file must contain a single JSON object. The whole JSON object must be in one line. Each object must have the following structure *(extra line breaks and indentation are for display purposes only)*:
```json
{
"request":{
    "method":"POST",
    "endpoint":"http://example.com/path",
    "body":"name=Kirlia",
    "headers":{
        "Content-Type":"application/x-www-form-urlencoded"
    }
}
```

### Specify Headers, Cookies, and POST Body
- `--header`/`-H` specifies headers which shall be added to the request.
    - Example: `tinja url -u "http://example.com/" -H "Authentication: Bearer ey..."`
- `--cookie`/`-c` specifies cookies which shall be added to the request.
    - Example: `tinja url -u "http://example.com/" -c "PHPSESSID=ABC123..."`
- `--data`/`-d` specifies the POST body which shall be added to the request.
    - Example: `tinja url -u "http://example.com/" -d "username=Kirlia&password=notguessable"`

### Scan CSTI in Addition to SSTI
- `--csti` enables the scanning for CSTI.
    - Example: `tinja url -u "http://example.com/" --csti`

By default TInjA only scans for SSTI. A headless browser is utilized for scanning for CSTI, which may increase RAM and CPU usage.

### Generate a JSONL Report
- `--reportpath` enables generating a report in JSONL format. The report will be updated after each scanned URL and will be stored at the provided path.
    - Example: `tinja url -u "http://example.com/" --reportpath "/home/user/Documents"`

### Use a Proxy
- `--proxyurl` specifies the URL and port of a proxy to be used for scanning.
    - Example: `tinja url -u "http://example.com/" --proxyurl "http://127.0.0.1:8080"`
- `--proxycertpath` specifies the CA certificate of the proxy in PEM format (needed when scanning HTTPS URLs).
    - Example `tinja url -u "http://example.com/" --proxyurl "http://127.0.0.1:8080" --proxycertpath "/home/user/Documents/cacert.pem"`

To scan HTTPS URLs using a proxy a CA certificate of the proxy in PEM format is needed. Burp Suite CA certificates are provided in DER format, for example. To convert them, the following command can be used:

`openssl x509 -inform DER -outform PEM -text -in cacert.der -out cacert.pem`

### Set a Ratelimit
- `--ratelimit`/`-r` specifies the number of maximum requests per second allowed. By default, this number is unrestricted.
    - Example: `tinja url -u "http://example.com/" --ratelimit 10`

## Troubleshooting
- `[ERR] Couldn't connect to URL: remote error: tls: user canceled`
    - When using a proxy and connecting via HTTPS, the proxy's CA certificate (.pem) needs to be specified with `--proxycertpath` (see [Use a Proxy](#use-a-proxy)).
- `[ERR] Error reading response from target server via proxy: malformed HTTP response "HTTP/1.1"`
    - Disable `Default to HTTP/2 if the server supports it` in Burp Suite's settings (`Network` > `HTTP`).

## TODOs
- `TINJA` marker to mark where the polyglots shall be placed.
- Support for multipart bodies.
- Optional: Blind SSTI Payloads (e.g., sleep payloads).
- Feedback, whether CSTI or SSTI was detected.
- Check headless browser's console for template engine error messages (see https://github.com/go-rod/rod/issues/330).
- Improve Error Detection, when input is not reflected

## Background Information
A blog post providing more information about template injection and [TInjA – the Template INJection Analyzer](https://github.com/Hackmanit/TInjA) can be found here:

[Template Injection Vulnerabilities – Understand, Detect, Identify](https://hackmanit.de/en/blog-en/178-template-injection-vulnerabilities-understand-detect-identify)

TInjA was developed as a part of a master's thesis by Maximilian Hildebrand.
You can find results of the master's thesis publicly available here:
- [Template Injection Table](https://github.com/Hackmanit/template-injection-table)
- [Template Injection Playground](https://github.com/Hackmanit/template-injection-playground)
- [TInjA – the Template INJection Analyzer](https://github.com/Hackmanit/TInjA)
- [Master's Thesis (PDF)](https://www.hackmanit.de/images/download/thesis/Improving-the-Detection-and-Identification-of-Template-Engines-for-Large-Scale-Template-Injection-Scanning-Maximilian-Hildebrand-Master-Thesis-Hackmanit.pdf)

## License
TInjA – the Template INJection Analyzer was developed by [Hackmanit](https://hackmanit.de) and [Maximilian Hildebrand](https://www.github.com/m10x) as a part of his master's thesis. TInjA – the Template INJection Analyzer is licensed under the [Apache License, Version 2.0](license.txt).

<a href="https://hackmanit.de"><img src="https://www.hackmanit.de/templates/hackmanit-v2/img/wbm_hackmanit.png" width="30%"></a>
