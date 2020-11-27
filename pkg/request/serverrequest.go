package RESTRequest

import (
	"html/template"
	"net/http"
)

type ServerRequest struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Template *template.Template
	Payload  interface{}
}

func CreateServerRawRequest(w http.ResponseWriter, r *http.Request) ServerRequest {
	sr := ServerRequest{}
	sr.Writer = w
	sr.Request = r
	sr.Template = nil
	sr.Payload = nil
	return sr
}

func CreateServerTemplateRequest(w http.ResponseWriter, r *http.Request, t *template.Template) ServerRequest {
	sr := ServerRequest{}
	sr.Writer = w
	sr.Request = r
	sr.Template = t
	sr.Payload = nil
	return sr
}

func CreateServerPayloadRequest(w http.ResponseWriter, r *http.Request, data interface{}) ServerRequest {
	sr := ServerRequest{}
	sr.Writer = w
	sr.Request = r
	sr.Template = nil
	sr.Payload = data
	return sr
}

func CreateServerTemplatedPayloadRequest(w http.ResponseWriter, r *http.Request, t *template.Template, data interface{}) ServerRequest {
	sr := ServerRequest{}
	sr.Writer = w
	sr.Request = r
	sr.Template = t
	sr.Payload = data
	return sr
}
