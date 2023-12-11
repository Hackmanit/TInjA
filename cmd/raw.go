package cmd

import (
	"example/user/tinja/pkg"

	"github.com/spf13/cobra"
)

var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "Scan using a Raw file\n\nThe file has to have a RAW data with the following structure on each line:\n" + getRawStructure(),
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Scan(config, version, pkg.JSONL)
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)

	rawCmd.PersistentFlags().StringVarP(&rawPath, "raw", "R", "", "Raw file with crawl results")
	rawCmd.PersistentFlags().BoolVarP(&httpP, "http", "", false, "Send HTTP Protoco")
	rawCmd.MarkPersistentFlagRequired("raw")
}

func getRawStructure() string {
	structure := `
GET /listproducts.php?artist=1414 HTTP/1.1
Host: testphp.vulnweb.com
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7
Accept-Encoding: gzip, deflate
Accept-Language: en-US,en;q=0.9
Cache-Control: max-age=0
Connection: close
Sec-Ch-Ua: 
Sec-Ch-Ua-Mobile: ?0
Sec-Ch-Ua-Platform: ""
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.171 Safari/537.36
`
	return structure
}
