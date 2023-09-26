package cmd

import (
	"example/user/tinja/pkg"

	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Scan a single or multiple URLs",
	Run: func(cmd *cobra.Command, args []string) {
		pkg.Scan(config, version, pkg.URL)
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)

	urlCmd.PersistentFlags().StringVarP(&data, "data", "d", "", "add data to the body and use POST method")
	urlCmd.PersistentFlags().StringSliceVarP(&parameters, "parameter", "p", []string{}, "add custom parameter(s)")
	urlCmd.PersistentFlags().StringSliceVarP(&urls, "url", "u", []string{}, "URL(s) to scan")
	urlCmd.PersistentFlags().StringSliceVar(&urlsReflection, "reflection", []string{}, "URL(s) to check for reflection")
	urlCmd.PersistentFlags().IntVar(&lengthLimit, "lengthlimit", 0, "limit the polyglot length. 0 is unlimited (default 0)")

	urlCmd.MarkPersistentFlagRequired("url")
}
