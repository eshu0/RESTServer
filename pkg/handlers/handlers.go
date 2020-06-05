package RESTHandlers

type RESTHandler struct {
	URL             	string `json:"url"`
	MethodName     	 	string `json:"methodname"`
	HTTPMethod      	string `json:"httpmethod"`
	FunctionalClass	 	string `json:"functionalclass"`
	// Static file handling 
	StaticDir 			string `json:"staticdir"`
	// Template details
	TemplatePath  		string `json:"templatepath"`
	TemplateFileName  	string `json:"templatefilename"`	
	TemplateBlob  		string `json:"templateblob"`	
	TemplateName  		string `json:"templatename"`		
	//JSON Handling
	JSONRequest			bool `json:"jsonrequest"` 
	JSONResponse		bool `json:"jsonresponse"` 
}
