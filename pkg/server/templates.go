package RESTServer

import (
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	Request "github.com/eshu0/RESTServer/pkg/request"
)

// Create Handlers helper functions
func (rs *RServer) CreateTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string) Handlers.RESTHandler {
	return rs.CreateJSONTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, Blob, Filename, false, false)
}

func (rs *RServer) CreateSpecificTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string) Handlers.RESTHandler {
	return rs.CreateJSONSpecificTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, Blob, Path, false, false)
}
func (rs *RServer) CreateJSONTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string, HandleJSONRequest bool, HandleJSONResponse bool) Handlers.RESTHandler {
	drhr := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, HandleJSONRequest, HandleJSONResponse)
	drhr.TemplateBlob = Blob
	drhr.TemplateFileName = Filename
	drhr.TemplateName = TemplateName
	return drhr
}

func (rs *RServer) CreateJSONSpecificTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string, HandleJSONRequest bool, HandleJSONResponse bool) Handlers.RESTHandler {
	drhr := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, HandleJSONRequest, HandleJSONResponse)
	drhr.TemplateBlob = Blob
	drhr.TemplatePath = Path
	drhr.TemplateName = TemplateName
	return drhr
}

// end create handlers

func (rs *RServer) AddJSONTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Filename string) {
	rs.Config.AddHandler(rs.CreateJSONTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, "", Filename, true, false))
}

func (rs *RServer) AddTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Filename string) {
	rs.AddTemplateHandlerWithBlob(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, "", Filename)
}

func (rs *RServer) AddTemplateHandlerWithBlob(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string) {
	rs.Config.AddHandler(rs.CreateTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, Blob, Filename))
}

func (rs *RServer) AddBlobTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string) {
	rs.Config.AddHandler(rs.CreateSpecificTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, Blob, Path))
}

func (rs *RServer) AddSpecificTemplateHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string) {
	rs.Config.AddHandler(rs.CreateSpecificTemplateHandler(URL, MethodName, HTTPMethod, FunctionalClass, TemplateName, Blob, Filename))
}

func (rs *RServer) MakeTemplateHandlerFunction(handler Handlers.RESTHandler, any interface{}) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		rs.LogDebug("MakeTemplateHandlerFunction", "MakeTemplateHandlerFunction called")

		if rs.Config.GetCacheTemplates() {
			rs.LogDebugf("MakeTemplateHandlerFunction", "Looking up template %s for %s ", handler.TemplateName, handler.URL)
			t := rs.Templates.Lookup(handler.TemplateName)
			rs.Invoke(any, handler.MethodName, Request.CreateServerTemplateRequest(w, r, t))
		} else {

			t, err := loadTemplate(rs, handler)

			if err != nil {
				rs.LogErrorf("MakeTemplateHandlerFunction", "Error : %s", err.Error())
				return
			}

			if handler.MethodName == "" {
				rs.LogDebugf("MakeTemplateHandlerFunction", "No method so default to loading the template %s for %s", handler.TemplateBlob, handler.URL)
				request := Request.CreateServerTemplateRequest(w, r, t)
				terr := request.Template.Execute(request.Writer, nil)
				if terr != nil {
					rs.LogErrorf("MakeTemplateHandlerFunction", "Template.Execute Error : %s", terr.Error())
					return
				}
			} else {
				if handler.JSONRequest {
					data, jsonerr := rs.RequestHelper.ReadJSONRequest(r, handler.JSONRequestType)
					if jsonerr != nil {
						rs.Invoke(any, handler.MethodName, Request.CreateServerTemplatedPayloadRequest(w, r, t, data))
					} else {
						rs.LogErrorf("MakeTemplateHandlerFunction", "ReadJSONRequest Error : %s", jsonerr.Error())
						return
					}
				} else {
					rs.Invoke(any, handler.MethodName, Request.CreateServerTemplateRequest(w, r, t))
				}
			}

		}

	}

}

func loadTemplate(rs *RServer, handler Handlers.RESTHandler) (*template.Template, error) {
	if handler.TemplatePath != "" {
		rs.LogDebugf("loadTemplate", "We have a template path %s for %s", handler.TemplatePath, handler.URL)
		t, err := template.ParseFiles(handler.TemplatePath)
		if err != nil {
			rs.LogErrorf("loadTemplate", "Failed to load template path for %s", handler.TemplatePath)
			rs.LogErrorEf("loadTemplate", "Load template path err: %s ", err)
			rs.LogDebugf("loadTemplate", "Failed loading template so trying blob for %s", handler.URL)
			return loadBlobTemplate(rs, handler)
		}
		return t, err
	} else {
		if handler.TemplateFileName != "" {
			tfilepath := rs.Config.GetTemplatePath() + handler.TemplateFileName
			rs.LogDebugf("loadTemplate", "We have a template filename %s for %s", tfilepath, handler.URL)
			t, err := template.ParseFiles(tfilepath)
			if err != nil {
				rs.LogErrorf("loadTemplate", "Failed to load template filename for %s", tfilepath)
				rs.LogErrorEf("loadTemplate", "Load template filename err: %s ", err)
				rs.LogDebugf("loadTemplate", "Failed loading template so trying blob for %s", handler.URL)
				return loadBlobTemplate(rs, handler)
			}
			return t, err
		} else {
			return loadBlobTemplate(rs, handler)
		}
	}
}

func loadBlobTemplate(rs *RServer, handler Handlers.RESTHandler) (*template.Template, error) {

	if handler.TemplateBlob != "" {
		rs.LogDebugf("loadBlobTemplate", "We have a template blob %s for %s", handler.TemplateBlob, handler.URL)
		t, err := template.New(handler.TemplateName).Parse(handler.TemplateBlob)
		return t, err
	} else {
		rs.LogDebugf("loadBlobTemplate", "No template blob for %s", handler.URL)
		return nil, errors.New("No template blob set")
	}
}

func (rs *RServer) LoadTemplates() {

	if rs.Config.HasTemplate() && rs.Config.GetCacheTemplates() {

		var allFiles []string
		TemplatePath := rs.Config.GetTemplatePath()
		files, err := ioutil.ReadDir(TemplatePath)
		if err != nil {
			rs.Templates = nil
			rs.LogErrorf("LoadTemplates", "ReadDir - Error : %s", err.Error())
			return
		}

		filetypes := rs.Config.GetTemplateFileTypes()
		for _, file := range files {
			filename := file.Name()
			for _, filetype := range filetypes {
				if strings.HasSuffix(filename, filetype) {
					allFiles = append(allFiles, TemplatePath+filename)
				}
			}
		}

		templates, terr := template.ParseFiles(allFiles...)
		if terr != nil {
			rs.Templates = nil
			rs.LogErrorEf("LoadTemplates", "ParseFiles - Error : %s", terr)
			return
		}
		rs.LogDebug("LoadTemplates", "Loaded Templates")
		rs.Templates = templates
	} else {
		if !rs.Config.HasTemplate() {
			rs.LogDebug("LoadTemplates", "No Template")
		}

		if !rs.Config.GetCacheTemplates() {
			rs.LogDebug("LoadTemplates", "Not Caching Templates")
		}

		rs.Templates = nil
	}
}
