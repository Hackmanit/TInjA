package pkg

import (
	"bufio"
	"os"
	"strings"

	"github.com/Hackmanit/TInjA/pkg/structs"
)

func ReadRaw(rawPath string, httpP bool) []structs.Crawl {
	fileRaw, err := os.Open(rawPath)
	if err != nil {
		PrintFatal("Error: ReadRaw: " + err.Error() + "\n")
	}
	defer fileRaw.Close()

	fileScanner := bufio.NewScanner(fileRaw)
	index := 0
	crawl := structs.Crawl{}
	crawl.Request.Headers = make(map[string]string)
	var path, body, host string
	bodyswitch := false
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if index == 0 {
			parse := strings.Split(line, " ")
			if len(parse) > 1 {
				crawl.Request.Method = parse[0]
				path = parse[1]
			} else {
				PrintFatal("Error: ReadRaw: The first line of the Raw file is malformed\n")
			}
		} else {
			if strings.HasPrefix(line, "Host: ") {
				host = line[6:]
			} else {
				parse := strings.SplitN(line, ":", 2)
				if len(parse) > 1 {
					crawl.Request.Headers[parse[0]] = parse[1]
				}
			}
			if bodyswitch {
				body = body + line
			}
			if len(line) == 0 {
				bodyswitch = true
			}
		}
		index++
	}
	if config.HTTP {
		crawl.Request.Endpoint = "http://" + host + path
	} else {
		crawl.Request.Endpoint = "https://" + host + path
	}
	crawl.Request.Body = body
	return []structs.Crawl{crawl}
}
