package cmd

/* Documentation

The `root.go` file contains the main command and configuration setup for the TInjA (Template INJection Analyzer) CLI tool. Here's an overview of the file:

- The `cmd` package is imported to define the command-line interface.
- Several packages, including `fmt`, `os`, `strings`, `example/user/tinja/pkg`, `example/user/tinja/pkg/structs`, `github.com/fatih/color`, `github.com/spf13/cobra`, `github.com/spf13/pflag`, and `github.com/spf13/viper`, are imported to support command-line handling, configuration, and output formatting.
- Constants `version`, `envPrefix`, `useragent`, and `uaChrome` are defined to store version information, environment prefix, user agent string, and Chrome user agent string respectively.
- Variables are declared for command-line flags and configuration settings.
- The `rootCmd` variable represents the main command for the CLI tool. It is of type `*cobra.Command` and contains information such as the command name, version, short and long descriptions, and the main function to execute when the command is run.
- The `Execute` function is defined to execute the root command. If an error occurs, it prints an error message and exits the program.
- The `init` function is called to initialize the CLI tool. It sets up command-line flags, initializes the configuration, and performs necessary bindings between flags and configuration settings.
- The `initConfig` function is responsible for initializing the configuration based on various sources such as command-line flags, environment variables, and config files.
- The `bindFlags` function binds each command-line flag to its associated configuration setting using Viper.
- The `getLogo` function returns a formatted logo string for the TInjA tool.

Overall, `root.go` sets up the main command and its configuration, handles command-line flags, and initializes the necessary components for the TInjA CLI tool to run.

*******/

import (
	"fmt"
	"os"
	"strings"

	"github.com/Hackmanit/TInjA/pkg"
	"github.com/Hackmanit/TInjA/pkg/structs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	version   = "v1.2.0"
	envPrefix = "TINJA"
	useragent = "TInjA " + version
	uaChrome  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
)

var ( // commandline flags
	timeout, verbosity, precedingLength, subsequentLength, lengthLimit     int
	ratelimit                                                              float64
	cfgFile, data, reportPath, proxyCertPath, proxyURL, jsonlPath, rawPath string
	cookies, headers, parameters, urls, urlsReflection, testheaders        []string
	config                                                                 structs.Config
	uac, csti, escapeJSON, httpP                                           bool
)

var rootCmd = &cobra.Command{
	Use:     "tinja",
	Version: version,
	Short:   "TInjA - the Template INJection Analyzer",
	Long:    getLogo() + "\n\nthe Template INJection Analyzer. (" + version + ")\n" + getCopyright(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getLogo() + "\n\nthe Template INJection Analyzer. (" + version + ")\n" + getCopyright() + "\n")
		fmt.Println("Use 'tinja url' to scan a single or multiple URLs.")
		fmt.Println("Use 'tinja raw' to scan a single URL using a raw file.")
		fmt.Println("Use 'tinja jsonl' to scan a single or multiple URLs using crawl results from a jsonl file.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 15, "seconds until timeout")
	rootCmd.PersistentFlags().Float64VarP(&ratelimit, "ratelimit", "r", 0, "number of requests per seconds. 0 is infinite (default 0)")
	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 1, "verbosity of the output. 0 = quiet, 1 = default, 2 = verbose")
	rootCmd.PersistentFlags().IntVar(&precedingLength, "precedinglength", 30, "how many chars shall be memorized, when getting the preceding chars of a body reflection point")
	rootCmd.PersistentFlags().IntVar(&subsequentLength, "subsequentlength", 30, "how many chars shall be memorized, when getting the subsequent chars of a body reflection point")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "set the path for a config file to be read")
	rootCmd.PersistentFlags().StringVar(&reportPath, "reportpath", "", "set the path for a report to be generated")
	rootCmd.PersistentFlags().StringVar(&proxyCertPath, "proxycertpath", "", "set the path for the certificate of the proxy")
	rootCmd.PersistentFlags().StringVar(&proxyURL, "proxyurl", "", "set the URL of the proxy")

	rootCmd.PersistentFlags().StringSliceVarP(&headers, "header", "H", []string{}, "add custom header(s)")
	rootCmd.PersistentFlags().StringSliceVarP(&cookies, "cookie", "c", []string{}, "add custom cookie(s)")

	rootCmd.PersistentFlags().BoolVar(&uac, "useragentchrome", false, "set chrome as user-agent. Default user-agent is '"+useragent+"'")
	rootCmd.PersistentFlags().BoolVar(&csti, "csti", false, "enable scanning for Client-Side Template Injections using a headless browser")
	rootCmd.PersistentFlags().BoolVar(&escapeJSON, "escapereport", false, "escape HTML special chars in the JSON report")
	rootCmd.PersistentFlags().StringSliceVar(&testheaders, "testheaders", []string{}, "headers to test. E.g. --testheaders Host,Origin,X-Forwarded-For")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			fmt.Println("Error reading config file:", err.Error())
		}
	}

	/*** the following code is based on https://github.com/carolynvs/stingoftheviper ***/

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --threads
	// binds to an environment variable SCANNER_THREADS. This helps
	// avoid conflicts.
	viper.SetEnvPrefix(envPrefix)

	// Environment variables can't have dashes in them, so bind them to their equivalent
	// keys with underscores, e.g. --foo-bar to SCANNER_FOO_BAR
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	viper.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(rootCmd)

	/***********************************************************************************/

	/* read cookie file and append its content to slice */
	cookiesReplaced := cookies
	for i := len(cookies) - 1; i >= 0; i-- {
		if strings.HasPrefix(cookies[i], "file:") {
			cookiesReplaced = pkg.SliceReplaceBetween(cookiesReplaced, i, pkg.ReadLocalFile(cookies[i], "cookie"))
		}
	}
	cookies = cookiesReplaced

	/* read parameter file and append its content to slice */
	parametersReplaced := parameters
	for i := len(parameters) - 1; i >= 0; i-- {
		if strings.HasPrefix(parameters[i], "file:") {
			parametersReplaced = pkg.SliceReplaceBetween(parametersReplaced, i, pkg.ReadLocalFile(parameters[i], "parameter"))
		}
	}
	parameters = parametersReplaced

	/* read header file and append its content to slice */
	headersReplaced := headers
	for i := len(headers) - 1; i >= 0; i-- {
		if strings.HasPrefix(headers[i], "file:") {
			headersReplaced = pkg.SliceReplaceBetween(headersReplaced, i, pkg.ReadLocalFile(headers[i], "header"))
		}
	}
	headers = headersReplaced

	/* read url file and append its content to slice */
	urlsReplaced := urls
	for i := len(urls) - 1; i >= 0; i-- {
		if strings.HasPrefix(urls[i], "file:") {
			urlsReplaced = pkg.SliceReplaceBetween(urlsReplaced, i, pkg.ReadLocalFile(urls[i], "url"))
		}
	}
	urls = urlsReplaced
	/*************************************************/

	//set user agent
	if uac {
		useragent = uaChrome
	}

	// get crawls
	var crawls []structs.Crawl
	if rawPath != "" {
		crawls = pkg.ReadRaw(rawPath, httpP)
	}

	if jsonlPath != "" {
		crawls = pkg.ReadJSONL(jsonlPath)
	}

	config = structs.Config{
		// Root
		Timeout:          timeout,
		Ratelimit:        ratelimit,
		Verbosity:        verbosity,
		PrecedingLength:  precedingLength,
		SubsequentLength: subsequentLength,
		ReportPath:       reportPath,
		UserAgentChrome:  uac,
		ProxyCertPath:    proxyCertPath,
		ProxyURL:         proxyURL,
		CSTI:             csti,
		EscapeJSON:       escapeJSON,
		UserAgent:        useragent,
		TestHeaders:      testheaders,
		// URL Command
		Data:           data,
		Cookies:        cookies,
		Headers:        headers,
		Parameters:     parameters,
		URLs:           urls,
		URLsReflection: urlsReflection,
		LengthLimit:    lengthLimit,
		// JSONL + RAW Command
		Crawls: crawls,
		// RAW Command
		HTTP:    httpP,
		RawPath: rawPath,
		// JSONL Command
		JSONLPath: jsonlPath,
	}
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(configName) {
			val := viper.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func getLogo() (logo string) {
	// source: https://patorjk.com/software/taag/#p=display&f=Slant%20Relief&t=TInjA
	logo = `
__/\\\\\\\\\\\\\\\__/\\\\\\\\\\\______________________________/\\\\\\\\\____        
 _\///////\\\/////__\/////\\\///______________________/\\\___/\\\\\\\\\\\\\__       
  _______\/\\\___________\/\\\________________________\///___/\\\/////////\\\_      
   _______\/\\\___________\/\\\______/\\/\\\\\\_________/\\\_\/\\\_______\/\\\_     
    _______\/\\\___________\/\\\_____\/\\\////\\\_______\/\\\_\/\\\\\\\\\\\\\\\_    
     _______\/\\\___________\/\\\_____\/\\\__\//\\\______\/\\\_\/\\\/////////\\\_   
      _______\/\\\___________\/\\\_____\/\\\___\/\\\__/\\_\/\\\_\/\\\_______\/\\\_  
       _______\/\\\________/\\\\\\\\\\\_\/\\\___\/\\\_\//\\\\\\__\/\\\_______\/\\\_ 
        _______\///________\///////////__\///____\///___\//////___\///________\///__`

	logo = strings.ReplaceAll(logo, "_", color.HiRedString("_"))

	return
}

func getCopyright() (copyright string) {
	copyright = `
Published by Hackmanit under http://www.apache.org/licenses/LICENSE-2.0
Author: Maximilian Hildebrand
Repository: https://github.com/Hackmanit/TInjA`

	return
}
