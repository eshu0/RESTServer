package RESTHandlers

//RESTHandler is the struct that represents a REST call
type RESTHandler struct {
	URL             string `json:"url,omitempty"`
	MethodName      string `json:"methodname,omitempty""`
	HTTPMethod      string `json:"httpmethod,omitempty""`
	FunctionalClass string `json:"functionalclass,omitempty""`
	// Static file handling
	StaticDir string `json:"staticdir,omitempty""`
	// Template details
	TemplatePath     string `json:"templatepath,omitempty""`
	TemplateFileName string `json:"templatefilename,omitempty""`
	TemplateBlob     string `json:"templateblob,omitempty""`
	TemplateName     string `json:"templatename,omitempty""`
	//JSON Handling
	JSONRequest      bool        `json:"jsonrequest,omitempty""`
	JSONResponse     bool        `json:"jsonresponse,omitempty""`
	JSONRequestType  interface{} `json:"-"`
	JSONRequestData  string      `json:"jsonin,omitempty""`
	JSONResponseData string      `json:"jsonout,omitempty""`
}
