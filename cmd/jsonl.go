package cmd

import (
	"example/user/tinja/pkg"

	"github.com/spf13/cobra"
)

var jsonlCmd = &cobra.Command{
	Use:   "jsonl",
	Short: "Scan using a JSONL file\n\nThe file has to have a JSON object with the following structure on each line:\n" + getStructure(),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Scan(config, version, pkg.JSONL)
	},
}

func init() {
	rootCmd.AddCommand(jsonlCmd)

	jsonlCmd.PersistentFlags().StringVarP(&jsonlPath, "jsonl", "j", "", "JSONL file with crawl results")

	jsonlCmd.MarkPersistentFlagRequired("jsonl")
}

func getStructure() string {
	structure := `
{
	"request":{
		"method":"POST",
		"endpoint":"http://example.com/path",
		"body":"name=kirlia",
		"headers":{
			"Content-Type":"application/x-www-form-urlencoded"
		}
	}
}
`
	return structure
}
