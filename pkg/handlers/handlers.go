package RESTHandlers

type RESTHandler struct {
	URL             	string `json:"url"`
	MethodName     	 	string `json:"methodname"`
	HTTPMethod      	string `json:"httpmethod"`
	FunctionalClass	 	string `json:"functionalclass"`
	StaticDir 			string `json:"staticdir"`
	TemplatePath  		string `json:"templatepath"`
	TemplateFileName  	string `json:"templatefilename"`	
	TemplateBlob  		string `json:"templateblob"`	
	TemplateName  		string `json:"templatename"`		
}
