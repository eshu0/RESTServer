package RESTServer

import (
	"net/http"
	"html/template"
	"errors"
    "io/ioutil"
	"strings"

	//Helpers "github.com/eshu0/RESTServer/pkg/helpers"
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	//Config "github.com/eshu0/RESTServer/pkg/config"

	//sl "github.com/eshu0/simplelogger"
	//sli "github.com/eshu0/simplelogger/interfaces"
	//mux "github.com/gorilla/mux"
)


func (rs *RServer) MakeTemplateHandlerFunction(handler Handlers.RESTHandler, any interface{}) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		rs.Log.LogDebug("MakeTemplateHandlerFunction", "MakeTemplateHandlerFunction called")
		
		if rs.Config.GetCacheTemplates(){
			rs.Log.LogDebugf("MakeTemplateHandlerFunction", "Looking up template %s for %s ",handler.TemplateName,handler.URL)
			t := rs.Templates.Lookup(handler.TemplateName) 
			rs.Invoke(any, handler.MethodName, w, r, t)
		} else {
			t := template.New(handler.TemplateName) 
			err := errors.New("should not see this error")
			
			if handler.TemplatePath != "" {
				rs.Log.LogDebugf("MakeTemplateHandlerFunction", "We have a template path %s for %s", handler.TemplatePath,handler.URL)
				t, err = template.ParseFiles(handler.TemplatePath)
			} else {
				if handler.TemplateFileName != "" {
					tfilepath := rs.Config.GetTemplatePath() + handler.TemplateFileName 
					rs.Log.LogDebugf("MakeTemplateHandlerFunction", "We have a template filename %s for %s", tfilepath,handler.URL)
					t, err = template.ParseFiles(tfilepath)
				} else {
					if handler.TemplateBlob != "" {
						rs.Log.LogDebugf("MakeTemplateHandlerFunction", "We have a template blob %s for %s", handler.TemplateBlob,handler.URL)
						t, err = template.New(handler.TemplateName).Parse(handler.TemplateBlob)
					} else {
						err = errors.New("No template set")
					}
				}
			}
		
			if err != nil {
				rs.Log.LogErrorf("MakeTemplateHandlerFunction", "Error : %s", err.Error())
				return
			}
			rs.Invoke(any, handler.MethodName, w, r, t)
		}
	
	}

}

func (rs *RServer) AddTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Filename string)  {
	rs.AddTemplateHandlerWithBlob(URL,MethodName,HTTPMethod,FunctionalClass,TemplateName,"",Filename)
}

func (rs *RServer) AddTemplateHandlerWithBlob(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string)  {
	rs.Config.AddHandler(rs.CreateTemplateHandler(URL,MethodName,HTTPMethod,FunctionalClass,TemplateName,Blob,Filename))
}

func (rs *RServer) AddBlobTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string)  {
	rs.Config.AddHandler(rs.CreateBlobTemplateHandler(URL,MethodName,HTTPMethod,FunctionalClass,TemplateName,Blob,Path))
}

func (rs *RServer) AddSpecificTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string)  {
	rs.Config.AddHandler(rs.CreateSpecificTemplateHandler(URL,MethodName,HTTPMethod,FunctionalClass,TemplateName,Blob,Filename))
}


func (rs *RServer) CreateTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Filename string) Handlers.RESTHandler {
	drhr := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass)
	drhr.TemplateBlob = Blob
	drhr.TemplateFileName = Filename
	drhr.TemplateName = TemplateName		
	return drhr
}

func (rs *RServer) CreateBlobTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string) Handlers.RESTHandler {
	drhr := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass)
	drhr.TemplateBlob = Blob
	drhr.TemplatePath = Path
	drhr.TemplateName = TemplateName		
	return drhr
}

func (rs *RServer) CreateSpecificTemplateHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, TemplateName string, Blob string, Path string) Handlers.RESTHandler {
	drhr := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass)
	drhr.TemplateBlob = Blob
	drhr.TemplatePath = Path
	drhr.TemplateName = TemplateName		
	return drhr
}

func (rs *RServer) LoadTemplates(){

	if rs.Config.HasTemplate() && rs.Config.GetCacheTemplates() {

		var allFiles []string
		TemplatePath := rs.Config.GetTemplatePath() 
		files, err := ioutil.ReadDir(TemplatePath)
		if err != nil {
			rs.Templates = nil
			rs.Log.LogErrorf("LoadTemplates", "ReadDir - Error : %s", err.Error())
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
			rs.Log.LogErrorf("LoadTemplates", "ParseFiles - Error : %s", terr.Error())
			return 
		}
		rs.Log.LogDebug("LoadTemplates", "Loaded Templates")
		rs.Templates = templates
	} else {
		if !rs.Config.HasTemplate() {
			rs.Log.LogDebug("LoadTemplates", "No Template")
		}

		if !rs.Config.GetCacheTemplates() {
			rs.Log.LogDebug("LoadTemplates", "Not Caching Templates")
		}

		rs.Templates = nil
	}
}
