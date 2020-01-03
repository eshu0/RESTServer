package RESTServer

type RESTHandler struct {
	URL        string `json:"url"`
	MethodName string `json:"methodname"`
	HTTPMethod string `json:"httpmethod"`
	FunctionalClass   string `json:"functionalclass"`
}


type DefaultRESTHandler struct {
	RESTHandler
	MappedClass interface{} `json:"-"`
}
