package pkg

/* Documentation

The `ReadJSONL` function accepts a `jsonlPath` parameter, which is the path to the JSONL file to be read. It opens the file, reads its content line by line, and parses each line as a JSON object into a `structs.Crawl` object. The parsed crawls are then stored in a slice and returned by the function.

**********/

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Hackmanit/TInjA/pkg/structs"
)

// ReadJSONL reads a JSONL file and returns a slice of structs.Crawl.
// Each line in the JSONL file is expected to contain a valid JSON object representing a crawl.
// The function returns the parsed crawls or an error if there was an issue reading or parsing the file.
func ReadJSONL(jsonlPath string) []structs.Crawl {
	fileJSONL, err := os.Open(jsonlPath)
	if err != nil {
		PrintFatal("Error: ReadJSONL: " + err.Error() + "\n")
	}
	defer fileJSONL.Close()

	fileScanner := bufio.NewScanner(fileJSONL)
	fileScanner.Split(bufio.ScanLines)
	var crawls []structs.Crawl

	for fileScanner.Scan() {
		var crawl structs.Crawl
		err = json.Unmarshal(fileScanner.Bytes(), &crawl)
		if err != nil {
			PrintFatal("Error: ReadJSONL: " + err.Error() + "\n")
		}
		crawls = append(crawls, crawl)
	}

	return crawls
	/*
		 FORMAT

			{
		    "request":{
		        "method":"POST","endpoint":"http://example.com/path","body":"name=kirlia","headers":{
		            "Content-Type":"application/x-www-form-urlencoded"
		        }
		    }
	*/
}
