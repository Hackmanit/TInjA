[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

# TInjA - the Template INJection Analyzer

TInjA is a CLI tool for testing web pages for **template injection** vulnerabilities.

- [Features](#features)
- [Supported Template Engines](#supported-template-engines)
- [Installation](#installation)
- [Usage](#usage)
- [Troubleshooting](#troubleshooting)

## Features
- Template injection detection and template engine identification
    - 44 template engines supported (see [Supported Template Engines](#supported-template-engines))
    - Both **SSTI** and **CSTI** are detected
        - SSTI = server-side template injection
        - CSTI = client-side template injection
- Efficient through the usage of polyglots
    - On average 5 polyglots are sent until the injection is detected and the engine identified
- Pass crawled URLs to TInjA in JSONL format
- Set custom headers, cookies, POST parameters, query parameters
- Route the traffic through a proxy (e. g. Burp Suite)
- Ratelimiting

## Supported Template Engines
### Javascript
- Handlebars
- EJS
- Underscore
- Vue.js
- Mustache
- Pug
- Angular.js
- Hogan.js
- Nunjucks
- Dot
- Velocity.js
- Eta
- Twig.js
### Python
- Jinja2
- Tornado
- Mako
- Django
- SimpleTemplate Engine
- Pystache
- Cheetah3
- Chameleon
### Java
- Groovy
- Freemarker
- Velocity
- Thymeleaf
### PHP
- Blade
- Twig
- Mustache.php
- Smarty
- Latte
### Ruby
- ERB
- Erubi
- Erubis
- Haml
- Liquid
- Slim
- Mustache
### Dotnet
- Razor Engine
- DotLiquid
- Scriban
- Fluid
### Golang
- html/template
- text/template
### Elixir
- EEx
## Installation
### Option 1: Pre-built Binary
Prebuilt binaries of TInjA are provided on the releases page.
### Option 2: Install Using Go
go1.18 or higher is required.
```bash
go install -v github.com/Hackmanit/Web-Cache-Vulnerability-Scanner@latest
```
## Usage
```tinja url -u "http://example.com/"``` scan a single URL  
```tinja url -u "http://example.com/" -u "http://example.com/path2``` scan multiple URLs  
```tinja url -u "file:/path/to/file"``` scan by importing URLs from a file  
```tinja jsonl -j "/path/to/file"``` scan by importing URLs and additional information from a JSONL file. The file must contain a JSON object with the following structure on each line:
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

### Specify Headers, Cookies and POST Body
```
tinja url -u "http://example.com/" -H "Authentication: Bearer ej..."
tinja url -u "http://example.com/" -c "PHPSESSID=ABC123..."
tinja url -u "http://example.com/" -d "username=Kirlia&password=notguessable"
```

### Scan CSTI, too
By default TInjA only scans for SSTI. In order to also scan for CSTI a headless browser needs to be utilized, which may increase RAM and CPU usage.
```tinja url -u "http://example.com/" --csti```
### Generate a JSONL Report
A JSONL report is generated and updated after each scanned URL if the flag --reportpath is set.   
```tinja url -u "http://example.com/" --reportpath "/home/user/Documents"```
### Use a Proxy
To scan HTTPS URLs using a proxy, a CA certificate of the proxy in PEM format is needed. Burp Suite certificates are provided in DER format, for example. To convert them, the following command can be used: openssl x509 -inform DER -outform PEM -text -in cacert.der -out cacert.pem.  
```tinja url -u "http://example.com/" --proxyurl "http://127.0.0.1:8080"``` scan HTTP URLs using a proxy  
```tinja url -u "http://example.com/" --proxyurl "http://127.0.0.1:8080" --proxycertpath "/home/user/Documents/cacert.pem"``` scan both HTTP and HTTPS URLs using a proxy  
### Set a Ratelimit
The number of maximum allowed requests per second can be set with --ratelimit/-r. By default, this number is unrestricted.  
```tinja url -u "http://example.com/" --ratelimit 10```

## Troubleshooting
`[ERR] Couldn't connect to URL: remote error: tls: user canceled`: When using a proxy and connecting via https, the proxy's certificate (.pem) needs to be specified with --proxycertpath

`[ERR] Error reading response from target server via proxy: malformed HTTP response "HTTP/1.1"`: Disable "Default to HTTP/2 if the server supports it" in Burp's Settings under Network > HTTP

## TODO
- "TINJA" marker to mark where the polyglots shall be placed
- support for multipart bodies
- optional: Blind SSTI Payloads (e.g. sleep payloads)
- feedback, whether CSTI or SSTI was detected
- check headless browser's console for template engine error messages https://github.com/go-rod/rod/issues/330