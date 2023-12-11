package pkg

/* Documentation

The `Report` structure represents an overall report and includes information about the name, version, number of template injections found, error messages, date, duration, command, configuration, and webpages.

The `GenerateReport` function takes a `Report` object and a `currentDate` string as parameters, and it generates a report file in JSONL format based on the provided report. The report file is saved to the specified path, and the function handles indentation and JSON encoding options based on the configuration settings.

The `addWebpageToReport` takes a ReportWebpage object marshalls it to JSON and appends it to the JSONL report.

************/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"example/user/tinja/pkg/structs"
	"fmt"
	"os"
)

var certaintyVeryHigh = "Very High"
var certaintyHigh = "High"
var certaintyMedium = "Medium"
var certaintyLow = "Low"
var certaintyVeryLow = "Very Low"

type (
	// reportRequest represents a request in the report.
	reportRequest struct {
		Conclusion  string `json:"conclusion"`
		Polyglot    string `json:"polyglot"`
		Error       string `json:"error"`
		CurlCommand string `json:"curlCommand"`
		Request     string `json:"request"`
		Response    string `json:"response"`
	}

	// reportParameter represents a parameter in the report.
	reportParameter struct {
		Name            string               `json:"name"`
		Type            string               `json:"type"`
		DefaultValues   []string             `json:"defaultValues"`
		IsVulnerable    bool                 `json:"isParameterVulnerable"`
		Certainty       string               `json:"certainty"`
		TemplateEngine  string               `json:"identifiedEngine"`
		ErrorMessages   []string             `json:"errorMessages"`
		Reflections     []structs.Reflection `json:"reflections"`
		AreErrorsThrown bool                 `json:"areErrorsShown"`
		Requests        []reportRequest      `json:"requests"`
	}

	reportDefault struct {
		StatusCode int    `json:"statusCode"`
		Request    string `json:"request"`
		Response   string `json:"response"`
	}

	// ReportWebpage represents a webpage in the report.
	ReportWebpage struct {
		ID            int               `json:"id"`
		URL           string            `json:"url"`
		IsVulnerable  bool              `json:"isWebpageVulnerable"`
		Certainty     string            `json:"certainty"`
		ErrorMessages []string          `json:"errorMessages"`
		ReportDefault reportDefault     `json:"default"`
		Parameters    []reportParameter `json:"parameters"`
	}

	// Report represents the overall report structure.
	Report struct {
		Name                    string   `json:"name"`
		Version                 string   `json:"version"`
		SuspectedVulnerableURLs int      `json:"suspectedVulnerableURLs"`
		SuspectedInjections     int      `json:"suspectedTemplateInjections"`
		VeryHigh                int      `json:"veryHighCertainty"`
		High                    int      `json:"highCertainty"`
		Medium                  int      `json:"mediumCertainty"`
		Low                     int      `json:"lowCertainty"`
		VeryLow                 int      `json:"veryLowCertainty"`
		ErrorMessages           []string `json:"errorMessages"`
		Date                    string   `json:"date"`
		Duration                string   `json:"duration"`
		Command                 string   `json:"command"`

		Config *structs.Config `json:"config,omitempty"`
	}
)

// generateReport generates a report and saves it to a file.
func generateReport(report Report, currentDate string) string {
	reportPath := config.ReportPath + currentDate + "_Report.jsonl"

	var file *os.File
	defer file.Close()

	file, err := os.OpenFile(reportPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		msg := fmt.Sprintf("GenerateReport: os.OpenFile: %s\n", err.Error())
		PrintFatal(msg)
	}

	jsonEncoder := json.NewEncoder(file)
	jsonEncoder.SetEscapeHTML(config.EscapeJSON)
	jsonEncoder.Encode(report)

	msg := fmt.Sprintf("Exported report %s\n", reportPath)
	PrintVerbose(msg, NoColor, 1)

	return reportPath
}

// add a webpage to the report
func addWebpageToReport(reportWebpage ReportWebpage, path string) {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Print("Error while opening report: "+err.Error()+"\n", Red)
	}
	defer f.Close()

	jsonEncoder := json.NewEncoder(f)
	jsonEncoder.SetEscapeHTML(config.EscapeJSON)
	err = jsonEncoder.Encode(reportWebpage)

	if err != nil {
		msg := fmt.Sprintf("Error appending webpage to report: %s\n", err.Error())
		Print(msg, Red)
	}

	msg := fmt.Sprintf("Successfully appended to %s\n", path)
	PrintVerbose(msg, NoColor, 1)
}

func replaceFirstLine(filePath, newFirstLine string) error {
	// Open the original file in read mode
	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// Create a temporary file for writing the modified content
	tempFilePath := filePath + ".temp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Write the new first line to the temporary file
	_, err = fmt.Fprint(tempFile, newFirstLine)
	if err != nil {
		return err
	}

	// Create a buffered reader for the original file
	reader := bufio.NewReader(originalFile)

	// Skip the first line in the original file
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Copy the rest of the original file content to the temporary file
	_, err = reader.WriteTo(tempFile)
	if err != nil {
		return err
	}

	// Close the original file and remove it
	originalFile.Close()
	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	// Rename the temporary file to the original file name
	err = os.Rename(tempFilePath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func updateReportsFirstLine(report Report, path string) {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(config.EscapeJSON)
	jsonEncoder.Encode(report)
	err := replaceFirstLine(path, bf.String())
	if err != nil {
		Print("Error while updating reports first line: "+err.Error()+"\n", Red)
	}

	msg := fmt.Sprintf("Finished report %s\n", reportPath)
	PrintVerbose(msg, NoColor, 1)
}
