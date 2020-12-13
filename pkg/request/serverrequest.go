package RESTRequest

import (
	"html/template"
	"net/http"
)

//ServerRequest - Represents a request to the server, has Writer for returning response
type ServerRequest struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Template *template.Template
	Payload  interface{}
}

//CreateServerRawRequest - Creates a simple writer/reader server request
func CreateServerRawRequest(w http.ResponseWriter, r *http.Request) ServerRequest {
	sr := ServerRequest{}
	sr.Writer = w
	sr.Request = r
	sr.Template = nil
	sr.Payload = nil
	return sr
}

//CreateServerTemplateRequest - Creates a template request with simple writer/reader server request
func CreateServerTemplateRequest(w http.ResponseWriter, r *http.Request, t *template.Template) ServerRequest {
	sr := CreateServerRawRequest(w, r)
	sr.Template = t
	return sr
}

//CreateServerPayloadRequest - Creates a payload request with simple writer/reader server request
func CreateServerPayloadRequest(w http.ResponseWriter, r *http.Request, data interface{}) ServerRequest {
	sr := CreateServerRawRequest(w, r)
	sr.Payload = data
	return sr
}

//CreateServerTemplatedPayloadRequest - Creates a templated/payload request with simple writer/reader server request
func CreateServerTemplatedPayloadRequest(w http.ResponseWriter, r *http.Request, t *template.Template, data interface{}) ServerRequest {
	sr := CreateServerRawRequest(w, r)
	sr.Template = t
	sr.Payload = data
	return sr
}
