package structs

type Crawl struct {
	Request struct {
		Method   string            `json:"method"`
		Endpoint string            `json:"endpoint"`
		Body     string            `json:"body"`
		Headers  map[string]string `json:"headers"`
	} `json:"request"`
	Response struct {
		StatusCode   int      `json:"status_code"`
		Technologies []string `json:"technologies"`
	} `json:"response"`
}

/*
	 FORMAT

		{
	    "request":{
	        "method":"POST","endpoint":"http://127.0.0.1:13370/ruby/Erb","body":"name=katana","headers":{
	            "Content-Type":"application/x-www-form-urlencoded"
	        }
	    },
	    "response":{
	        "status_code":200,"technologies":["Nginx:1.23.3"]
	    }
	}
*/
