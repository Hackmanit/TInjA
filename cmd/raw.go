package cmd

import (
	"github.com/Hackmanit/TInjA/pkg"

	"github.com/spf13/cobra"
)

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Scan using a Raw file",
	Long:  getLogo() + "\n\nthe Template INJection Analyzer. (" + version + ")\n" + getCopyright() + "\n\nScan using a Raw file. The file has to have a RAW data with the following structure on each line:\n" + getRawStructure(),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Scan(config, version, pkg.JSONL)
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)

	rawCmd.PersistentFlags().StringVarP(&rawPath, "raw", "R", "", "Raw file with crawl results")
	rawCmd.PersistentFlags().BoolVar(&httpP, "http", false, "Send HTTP Protocol")
	rawCmd.MarkPersistentFlagRequired("raw")
}

func getRawStructure() string {
	structure := `
POST /path?query=param HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.171 Safari/537.36
Content-Length: 11

name=kirlia
`
	return structure
}
