package pkg

/* Documentation

The `engines.go` file contains a package called `pkg` that defines possible Template Injection Polyglots and the expected answers of various template engines. Here's a breakdown of the class and its contents:

- The file imports the `structs` package from the `example/user/tinja/pkg/structs` directory.

- It declares a slice of `structs.Engine` named `engines` and an empty `map[string]bool` named `possibleEngines`.

- The file defines multiple variables, each representing a template injection polyglot or an expected error response for a specific template engine. Examples of these variables include `err1`, `err2`, `errPython`, `nonerr1`, `nonerr2`, and more.

- The file initializes the `possibleEngines` map, but its initialization is incomplete and not shown in the provided code snippet.

- The file defines the initialization function `init()` which populates the `engines` slice with `structs.Engine` instances representing various template engines along with their respective polyglots and expected responses.

- For each template engine, the `init()` function creates a new `structs.Engine` instance and adds it to the `engines` slice.

- Each `structs.Engine` instance has properties such as `Name` (name of the template engine), `Language` (programming language associated with the engine), `Version` (version of the engine), and `Polyglots` (a map of polyglots and their expected responses).

- The `Polyglots` map contains entries for each polyglot, where the key is the polyglot string and the value is the expected response from the associated template engine. Examples of polyglot entries are `err1: respError`, `err2: respError`, `errRuby: respError`, `nonerr1: respUnmodified`, `nonerr2: "#{1}1@*"`, and more.

This `engines.go` file provides a collection of template engines, their respective polyglots, and the expected responses from those engines. These can be used for comparison and analysis when detecting and identifying template injection vulnerabilities.

********/

import "example/user/tinja/pkg/structs"

var engines []structs.Engine

var err1 = "<%'${{/#{@}}%>{{"
var err2 = "<%'${{#{@}}%>"
var err3 = "${{<%[%'\"}}%\\"
var err4 = "<#set($x<%={{={@{#{${xux}}%>)"
var err5 = "<%={{={@{#{${xu}}%>"

var errPython = "${{/#}}"
var errRuby = "<%{{#{%>}"
var errDotnet = "{{@"
var errJava = "<%'#{@}"
var errJava2 = "<th:t=\"${xu}#foreach."
var errPHP = "{{/}}"
var errJavascript = "<%${{#{%>}}"
var errGolang = "{{"
var errElixir = "<%"

var nonerr1 = "p \">[[${{1}}]]"
var nonerr2 = "<%=1%>@*#{1}"
var nonerr3 = "{##}/*{{.}}*/"

var nonerrPython = "{#${{1}}#}}"
var nonerrRuby = "<%=1%>#{2}{{a}}"
var nonerrDotnet = "{{1}}@*"
var nonerrJava = "a\">##[[${1}]]"
var nonerrPHP = "{{7}}}"
var nonerrJavascript = "//*<!--{##<%=1%>{{!--{{1}}--}}-->*/#}"
var nonerrGolang = "{{.}}"
var nonerrElixir = "<%%a%>"
var nonerrIdent1 = "{{1in[1]}}"
var nonerrIdent2 = "${\"<%-1-%>\"}"
var nonerrIdent3 = "#evaluate(\"a\")"

var respError = "error"
var respEmpty = "empty"
var respUnmodified = "unmodified"
var respIdentified = "identified"
var respContinue = "continue"

var possibleEngines map[string]bool

func init() {

	/* Begin Ruby */
	// Erb/Erubi/Erubis
	engines = append(engines, structs.Engine{
		Name:     "Erb/Erubi/Erubis",
		Language: "Ruby",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: "<%a%>", nonerrIdent1: respUnmodified, nonerrIdent2: "${\"\"}", nonerrIdent3: respUnmodified,
		},
	})
	// Haml
	engines = append(engines, structs.Engine{
		Name:     "Haml",
		Language: "Ruby",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "<%=1%>@*1", nonerr3: respUnmodified,
			nonerrRuby: "<%=1%>2{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respError, nonerrJavascript: "<!-- /*<!--{##<%=1%>{{!--{{1}}--}}-->*/#} -->", nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respError,
		},
	})
	// Liquid
	engines = append(engines, structs.Engine{
		Name:     "Liquid",
		Language: "Ruby",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: "<%'$%>", err3: "$%\\", err4: "<#set($x<%=%>)", err5: "<%=%>", errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respEmpty, errPython: "$", errJavascript: "<%$", errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Slim
	engines = append(engines, structs.Engine{
		Name:     "Slim",
		Language: "Ruby",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "<%=1%>@*1", nonerr3: respUnmodified,
			nonerrRuby: "<%=1%>2{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Mustache
	engines = append(engines, structs.Engine{
		Name:     "Mustache",
		Language: "Ruby",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: "<#set($x<%=%>)", err5: "<%=%>", errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$]]", nonerr2: respUnmodified, nonerr3: "{##}/*#&lt;Mustache:0xARBITRARY16ARBITRARY&gt;*/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "@*", nonerrJava: respUnmodified, nonerrPHP: "}", nonerrPython: "{#$#}}", nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: "#&lt;Mustache:0xARBITRARY16ARBITRARY&gt;", nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End Ruby */

	/* Begin Golang */
	// html/template
	engines = append(engines, structs.Engine{
		Name:     "html/template",
		Language: "Golang",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respError, errJava: "&lt;%'#{@}", errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: "&lt;%",
			nonerr1: "p \">[[$1]]", nonerr2: "&lt;%=1%>@*#{1}", nonerr3: "{##}/*ARBITRARY*/",
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: "ARBITRARY", nonerrElixir: "&lt;%%a%>", nonerrIdent1: respError, nonerrIdent2: "${\"&lt;%-1-%>\"}", nonerrIdent3: respUnmodified,
		},
	})
	// text/template
	engines = append(engines, structs.Engine{
		Name:     "text/template",
		Language: "Golang",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: "{##}/*ARBITRARY*/",
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: "ARBITRARY", nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End Golang */

	/* Begin Dotnet */
	// RazorEngine.NetCore
	engines = append(engines, structs.Engine{
		Name:     "RazorEngine.NetCore",
		Language: "Dotnet",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respError, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respUnmodified, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "<%=1%>", nonerr3: respUnmodified,
			nonerrRuby: respUnmodified, nonerrDotnet: "{{1}}", nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// DotLiquid
	engines = append(engines, structs.Engine{
		Name:     "DotLiquid",
		Language: "Dotnet",
		Version:  "",
		Polyglots: map[string]string{err1: respError, err2: "<%'$%>", err3: respError, err4: "<#set($x<%=%>)", err5: "<%=%>", errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respEmpty, errPython: "$", errJavascript: "<%$", errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: "//*<!--{##<%=1%>-}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Scriban
	engines = append(engines, structs.Engine{
		Name:     "Scriban",
		Language: "Dotnet",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: "<%'$%>", err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: "<%", errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: "<%$", errGolang: respEmpty, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Scriban Liquid
	engines = append(engines, structs.Engine{
		Name:     "Scriban (Liquid mode)",
		Language: "Dotnet",
		Version:  "",
		Polyglots: map[string]string{err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respEmpty, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Fluid
	engines = append(engines, structs.Engine{
		Name:     "Fluid",
		Language: "Dotnet",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End Dotnet

	/* Begin Elixir */
	// EEx
	engines = append(engines, structs.Engine{
		Name:     "EEx",
		Language: "Elixir",
		Version:  "",
		Polyglots: map[string]string{err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: respUnmodified, nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: "<%a%>", nonerrIdent1: respUnmodified, nonerrIdent2: respError, nonerrIdent3: respUnmodified,
		},
	})
	/* End Elixir */

	/* Begin Java */
	// Groovy
	engines = append(engines, structs.Engine{
		Name:     "Groovy",
		Language: "Java",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: "p \">[[1]]", nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: "a\">##[[1]]", nonerrPHP: respUnmodified, nonerrPython: "{#1#}}", nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: respError, nonerrIdent1: respUnmodified, nonerrIdent2: "<%-1-%>", nonerrIdent3: respUnmodified,
		},
	})
	// Freemarker
	engines = append(engines, structs.Engine{
		Name:     "Freemarker",
		Language: "Java",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respError, nonerr2: "<%=1%>@*1", nonerr3: respUnmodified,
			nonerrRuby: "<%=1%>2{{a}}", nonerrDotnet: respUnmodified, nonerrJava: "a\">##[[1]]", nonerrPHP: respUnmodified, nonerrPython: respError, nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: "<%-1-%>", nonerrIdent3: respUnmodified,
		},
	})
	// Velocity
	engines = append(engines, structs.Engine{
		Name:     "Velocity",
		Language: "Java",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respUnmodified, errJava2: respError,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respUnmodified, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: respUnmodified, nonerr3: "{",
			nonerrRuby: respUnmodified, nonerrDotnet: respUnmodified, nonerrJava: "a\">", nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{", nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: "a",
		},
	})
	// Thymeleaf
	engines = append(engines, structs.Engine{
		Name:     "Thymeleaf",
		Language: "Java",
		Version:  "",
		Polyglots: map[string]string{err1: respError, err2: respError, err3: respError, err4: respUnmodified, err5: respUnmodified, errJava2: respError,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respUnmodified, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p \">1", nonerr2: respUnmodified, nonerr3: respUnmodified,
			nonerrRuby: respUnmodified, nonerrDotnet: respUnmodified, nonerrJava: "a\">##1", nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Thymeleaf Inline
	engines = append(engines, structs.Engine{
		Name:     "Thymeleaf (Inline)",
		Language: "Java",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respError, errJava: respError, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respError,
			nonerr1: "p", nonerr2: respError, nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: respError, nonerrJava: "a", nonerrPHP: respError, nonerrPython: respError, nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respError, nonerrIdent1: respError, nonerrIdent2: respError, nonerrIdent3: respError,
		},
	})
	/* End Java */

	/* Begin PHP */
	// Blade
	engines = append(engines, structs.Engine{
		Name:     "Blade",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Twig
	engines = append(engines, structs.Engine{
		Name:     "Twig",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: "1", nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Twig Sandbox
	engines = append(engines, structs.Engine{
		Name:     "Twig (Sandbox)",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Mustache PHP
	engines = append(engines, structs.Engine{
		Name:     "Mustache.PHP",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "@*", nonerrJava: respUnmodified, nonerrPHP: "}", nonerrPython: "{#$#}}", nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Smarty/Smarty (Security)
	engines = append(engines, structs.Engine{
		Name:     "Smarty/Smarty (Security)",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respError, errJava: respError, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: "<%=1%>@*#1", nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: "a\">##[[$1]]", nonerrPHP: "7}", nonerrPython: respError, nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: "$<%-1-%>", nonerrIdent3: respUnmodified,
		},
	})
	// Latte/Latte (Sandbox)
	engines = append(engines, structs.Engine{
		Name:     "Latte/Latte (Sandbox)",
		Language: "PHP",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respError, errJava: respError, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[${1}]]", nonerr2: "<%=1%>@*#1", nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "{1}@*", nonerrJava: "a\">##[[$1]]", nonerrPHP: "{7}}", nonerrPython: respError, nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: "{1}", nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End PHP */

	/* Begin Python */
	// Jinja2/Jinja2 (Sandbox)
	engines = append(engines, structs.Engine{
		Name:     "Jinja2/Jinja2 (Sandbox)",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: "True", nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Tornado
	engines = append(engines, structs.Engine{
		Name:     "Tornado",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: "True", nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Mako
	engines = append(engines, structs.Engine{
		Name:     "Mako",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: "p \">[[{1}]]", nonerr2: respError, nonerr3: respUnmodified,
			nonerrRuby: respError, nonerrDotnet: respUnmodified, nonerrJava: "a\">##[[1]]", nonerrPHP: respUnmodified, nonerrPython: "{#{1}#}}", nonerrJavascript: respError, nonerrGolang: respUnmodified, nonerrElixir: respError, nonerrIdent1: respUnmodified, nonerrIdent2: "<%-1-%>", nonerrIdent3: respUnmodified,
		},
	})
	// Django
	engines = append(engines, structs.Engine{
		Name:     "Django",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respError, nonerr2: respUnmodified, nonerr3: "/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: respError, nonerrJava: respUnmodified, nonerrPHP: respError, nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// SimpleTemplateEngine
	engines = append(engines, structs.Engine{
		Name:     "SimpleTemplateEngine",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: respError, nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respError, nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: "True", nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Pystache
	engines = append(engines, structs.Engine{
		Name:     "Pystache",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: "%>", err3: "$%\\", err4: "<#set($x<%=%>)", err5: "<%=%>", errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respEmpty, errPython: respError, errJavascript: "", errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p \">[[$]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "@*", nonerrJava: respUnmodified, nonerrPHP: "}", nonerrPython: "{#$#}}", nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Cheetah3
	engines = append(engines, structs.Engine{
		Name:     "Cheetah3",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: "p \">[[{1}]]", nonerr2: "1@*#{1}", nonerr3: "{",
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: "a\">", nonerrPHP: respUnmodified, nonerrPython: "{#{1}#}}", nonerrJavascript: "//*<!--{", nonerrGolang: respUnmodified, nonerrElixir: respError, nonerrIdent1: respUnmodified, nonerrIdent2: "<%-1-%>", nonerrIdent3: respUnmodified,
		},
	})
	// Chameleon
	engines = append(engines, structs.Engine{
		Name:     "Chameleon",
		Language: "Python",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respError,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p \">[[{1}]]", nonerr2: respUnmodified, nonerr3: respUnmodified,
			nonerrRuby: respUnmodified, nonerrDotnet: respUnmodified, nonerrJava: "a\">##[[1]]", nonerrPHP: respUnmodified, nonerrPython: "{#{1}#}}", nonerrJavascript: respError, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End Python */

	/* Begin Javascript */
	// Handlebars
	engines = append(engines, structs.Engine{
		Name:     "Handlebars",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "@*", nonerrJava: respUnmodified, nonerrPHP: respError, nonerrPython: "{#$#}}", nonerrJavascript: "//*<!--{##<%=1%>-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// EJS
	engines = append(engines, structs.Engine{
		Name:     "EJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: respUnmodified, nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: "<%a%>", nonerrIdent1: respUnmodified, nonerrIdent2: "${\"1\"}", nonerrIdent3: respUnmodified,
		},
	})
	// Underscore
	engines = append(engines, structs.Engine{
		Name:     "Underscore",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: respError, nonerrIdent1: respUnmodified, nonerrIdent2: respError, nonerrIdent3: respUnmodified,
		},
	})
	// VueJS
	engines = append(engines, structs.Engine{
		Name:     "VueJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respError,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: "p &quot;&gt;[[$1]]", nonerr2: "&lt;%=1%&gt;@*#{1}", nonerr3: respError,
			nonerrRuby: "&lt;%=1%&gt;#{2}", nonerrDotnet: "1@*", nonerrJava: "a&quot;&gt;##[[${1}]]", nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: "<!--[-->//*<!--{##<%=1%>{{!--{{1}}--}}-->*/#}<!--]-->", nonerrGolang: respError, nonerrElixir: "&lt;%%a%&gt;", nonerrIdent1: respError, nonerrIdent2: "${&quot;&lt;%-1-%&gt;&quot;}", nonerrIdent3: respUnmodified,
		},
	})
	// MustacheJS
	engines = append(engines, structs.Engine{
		Name:     "MustacheJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: respError, nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: respError, nonerrDotnet: respError, nonerrJava: respUnmodified, nonerrPHP: respError, nonerrPython: respError, nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Pug
	engines = append(engines, structs.Engine{
		Name:     "Pug",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respError, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respError,
			nonerr1: "<p>\">[[${{1}}]]</p>", nonerr2: "<%=1%>@*1", nonerr3: respError,
			nonerrRuby: "<%=1%>2{{a}}", nonerrDotnet: respError, nonerrJava: respError, nonerrPHP: respError, nonerrPython: respError, nonerrJavascript: "<!--*<!--{##<%=1%>{{!--{{1}}--}}-->*/#}-->", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respError, nonerrIdent3: "<div id=\"evaluate\" a=\"a\"></div>",
		},
	})
	// Pug (Inline)
	engines = append(engines, structs.Engine{
		Name:     "Pug (Inline)",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respUnmodified, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respUnmodified, nonerr2: "<%=1%>@*1", nonerr3: respUnmodified,
			nonerrRuby: "<%=1%>2{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respError, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// AngularJS
	engines = append(engines, structs.Engine{
		Name:     "AngularJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: "<th:t=\\\"${xu}#foreach.< p=\"\"></th:t=\\\"${xu}#foreach.<>",
			errRuby: "&lt;%{{#{%&gt;}", errDotnet: respUnmodified, errJava: "&lt;%'#{@}", errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: "&lt;%",
			nonerr1: "class=\"ng-binding\">p \"&gt;[[$1]]", nonerr2: "&lt;%=1%&gt;@*#{1}", nonerr3: respError,
			nonerrRuby: "class=\"ng-binding\">&lt;%=1%&gt;#{2}", nonerrDotnet: "1@*", nonerrJava: "a\"&gt;##[[${1}]]", nonerrPHP: "7}", nonerrPython: "{#$1#}}", nonerrJavascript: respUnmodified, nonerrGolang: respUnmodified, nonerrElixir: "&lt;%%a%&gt;", nonerrIdent1: respError, nonerrIdent2: "${\"&lt;%-1-%&gt;\"}", nonerrIdent3: respUnmodified,
		},
	})
	// HoganJS
	engines = append(engines, structs.Engine{
		Name:     "HoganJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: "<#set($x<%=ARBITRARY<#set($x<%={{={@{#{${xux}}%>)", err5: "<%=ARBITRARY<%={{={@{#{${xu}}%>", errJava2: respUnmodified,
			errRuby: "<%{%>}", errDotnet: "@", errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respEmpty, errElixir: respUnmodified,
			nonerr1: "p \">[[$]]", nonerr2: respUnmodified, nonerr3: "{##}/**/",
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "@*", nonerrJava: respUnmodified, nonerrPHP: "}", nonerrPython: "{#$#}}", nonerrJavascript: "//*<!--{##<%=1%>--}}-->*/#}", nonerrGolang: respEmpty, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Nunjucks
	engines = append(engines, structs.Engine{
		Name:     "Nunjucks",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respEmpty, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// Dot
	engines = append(engines, structs.Engine{
		Name:     "Dot",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respError, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respError, nonerr2: respUnmodified, nonerr3: "{##}",
			nonerrRuby: respError, nonerrDotnet: respError, nonerrJava: respUnmodified, nonerrPHP: respError, nonerrPython: respError, nonerrJavascript: "/#}", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	// VelocityJS
	engines = append(engines, structs.Engine{
		Name:     "VelocityJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respUnmodified, errJava2: respError,
			errRuby: respUnmodified, errDotnet: respUnmodified, errJava: respUnmodified, errPHP: respUnmodified, errPython: respError, errJavascript: respError, errGolang: respUnmodified, errElixir: respUnmodified,
			nonerr1: respError, nonerr2: respUnmodified, nonerr3: "{",
			nonerrRuby: respUnmodified, nonerrDotnet: respUnmodified, nonerrJava: "a\">", nonerrPHP: respUnmodified, nonerrPython: respError, nonerrJavascript: "//*<!--{", nonerrGolang: respUnmodified, nonerrElixir: respUnmodified, nonerrIdent1: respUnmodified, nonerrIdent2: respError, nonerrIdent3: respEmpty,
		},
	})
	// Eta
	engines = append(engines, structs.Engine{
		Name:     "Eta",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respUnmodified, errJava: respError, errPHP: respUnmodified, errPython: respUnmodified, errJavascript: respError, errGolang: respUnmodified, errElixir: respError,
			nonerr1: respUnmodified, nonerr2: "1@*#{1}", nonerr3: respUnmodified,
			nonerrRuby: "1#{2}{{a}}", nonerrDotnet: respUnmodified, nonerrJava: respUnmodified, nonerrPHP: respUnmodified, nonerrPython: respUnmodified, nonerrJavascript: "//*<!--{##1{{!--{{1}}--}}-->*/#}", nonerrGolang: respUnmodified, nonerrElixir: respError, nonerrIdent1: respUnmodified, nonerrIdent2: "${\"\"}", nonerrIdent3: respUnmodified,
		},
	})
	//TwigJS
	engines = append(engines, structs.Engine{
		Name:     "TwigJS",
		Language: "Javascript",
		Version:  "",
		Polyglots: map[string]string{
			err1: respError, err2: respError, err3: respError, err4: respError, err5: respError, errJava2: respUnmodified,
			errRuby: respError, errDotnet: respError, errJava: respUnmodified, errPHP: "NaN", errPython: respError, errJavascript: respError, errGolang: respError, errElixir: respUnmodified,
			nonerr1: "p \">[[$1]]", nonerr2: respUnmodified, nonerr3: respError,
			nonerrRuby: "<%=1%>#{2}", nonerrDotnet: "1@*", nonerrJava: respUnmodified, nonerrPHP: "7}", nonerrPython: "}", nonerrJavascript: "//*<!--", nonerrGolang: respError, nonerrElixir: respUnmodified, nonerrIdent1: respError, nonerrIdent2: respUnmodified, nonerrIdent3: respUnmodified,
		},
	})
	/* End Javascript */
}
