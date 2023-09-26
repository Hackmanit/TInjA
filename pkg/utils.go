package pkg

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"html"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	NoColor = 0
	Red     = 1
	Yellow  = 2
	Green   = 3
	Cyan    = 4

	PrecedingMissing  = "PrecedingMissing"
	SubsequentMissing = "SubsequentMissing"
	BothMissing       = "BothMissing"
)

func PrintVerbose(msg string, c int, threshold int) {
	if config.Verbosity >= threshold {
		if c == Red {
			msg = color.RedString("[ERR] ") + msg
		} else if c == Yellow {
			msg = color.YellowString("[!] ") + msg
		} else if c == Green {
			msg = color.GreenString("[+] ") + msg
		} else if c == Cyan {
			msg = color.CyanString("[*] ") + msg
		}

		fmt.Print(msg)
	}
}

func Print(msg string, c int) {
	PrintVerbose(msg, c, 0)
}

func PrintFatal(msg string) {
	Print(msg, Red)
	os.Exit(1)
}

func between(value string, preceding string, subsequent string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, preceding)
	posLast := strings.LastIndex(value, subsequent)

	if posFirst == -1 && posLast == -1 {
		return BothMissing
	}
	if posFirst == -1 {
		return PrecedingMissing
	}
	if posLast == -1 {
		return SubsequentMissing
	}
	posFirstAdjusted := posFirst + len(preceding)
	if posFirstAdjusted >= posLast { // There is nothing between
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func ReadLocalFile(path string, name string) []string {
	path = strings.TrimPrefix(path, "file:")

	if strings.HasPrefix(strings.ToLower(path), "file:") {
		PrintFatal("Please make sure that path: is lowercase: " + path)
	}

	w, err := os.ReadFile(path)
	if err != nil {
		PrintFatal("The specified " + name + " file path " + path + " couldn't be found: " + err.Error() + "\n")
	}

	return strings.Split(string(w), "\n")
}

func SliceReplaceBetween(slice []string, offset int, replace []string) []string {
	// replace offset element in slice with "replace" slice. Order remains
	return append(append(slice[:offset], replace...), slice[offset+1:]...)
}

// taken from https://stackoverflow.com/a/32350135 and modified
func getToken(length int) string {
	if length > 56 {
		Print("length must not be greater than 56", Red)
	}
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

func isEncoded(response string, polyglot string) (bool, string) {
	// html escape + unescape may transform differently
	if response == html.EscapeString(polyglot) || html.UnescapeString(response) == polyglot {
		return true, " (HTML encoded)"
	} else if response == url.PathEscape(polyglot) || response == url.QueryEscape(polyglot) {
		return true, " (URL encoded)"
	}
	return false, ""
}
