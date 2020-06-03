package RESTServer

import (
	"net/http"
	"strings"

	//Helpers "github.com/eshu0/RESTServer/pkg/helpers"
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	//Config "github.com/eshu0/RESTServer/pkg/config"

	//sl "github.com/eshu0/simplelogger"
	//sli "github.com/eshu0/simplelogger/interfaces"
	mux "github.com/gorilla/mux"
)

func (rs *RServer) MakeHandlerFunction(MethodName string, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs.Invoke(any, MethodName, w, r)
	}
}

func (rs *RServer) AddStaticHandler(URL string, StaticDir string)  {
	rs.Config.AddHandler(rs.CreateStaticHandler(URL,StaticDir))
}

func (rs *RServer) AddFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string)  {
	rs.Config.AddHandler(rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass))
}

func (rs *RServer) CreateStaticHandler(URL string, StaticDir string) Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}
	drhr.URL = URL
	drhr.StaticDir = StaticDir
	return drhr
}

func (rs *RServer) CreateFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string) Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}
	drhr.URL = URL
	drhr.MethodName = MethodName
	drhr.HTTPMethod = HTTPMethod
	drhr.FunctionalClass = FunctionalClass
	return drhr
}

func (rs *RServer) addHandlerToRouter(r *mux.Router, handl Handlers.RESTHandler){

	funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

	if ok {
		if handl.TemplatePath != "" || handl.TemplateFileName != "" || handl.TemplateBlob != ""  {
			rs.Log.LogDebugf("addHandlertoRouter", "Handlers: Adding Template function %s", handl.MethodName)
			
			if handl.HTTPMethod != ""{
				if strings.Contains(handl.HTTPMethod,",") {
					r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(http.MethodOptions,strings.Split(handl.HTTPMethod,",")...)
					r.Use(mux.CORSMethodMiddleware(r))

				}else{
					r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(handl.HTTPMethod,http.MethodOptions)
					r.Use(mux.CORSMethodMiddleware(r))
				}
			}else{
				r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass))
				r.Use(mux.CORSMethodMiddleware(r))

			}
		} else {
			rs.Log.LogDebugf("addHandlertoRouter", "Handlers: Adding %s", handl.MethodName)
			if handl.HTTPMethod != ""{
				if strings.Contains(handl.HTTPMethod,",") {
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass)).Methods(http.MethodOptions,strings.Split(handl.HTTPMethod,",")...)
					r.Use(mux.CORSMethodMiddleware(r))
				}else{
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass)).Methods(handl.HTTPMethod, http.MethodOptions)
					r.Use(mux.CORSMethodMiddleware(r))
				}
			}else{
				r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass))
				r.Use(mux.CORSMethodMiddleware(r))
			}
		}
	} else {
		if handl.StaticDir != "" {
			rs.Log.LogDebugf("addHandlertoRouter", "Handlers: Adding route %s for  static directory %s", handl.URL, handl.StaticDir)		
			r.PathPrefix(handl.URL).Handler(http.StripPrefix(handl.URL, http.FileServer(http.Dir(handl.StaticDir))))

		} else {
			rs.Log.LogError("addHandlertoRouter", "Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)		
		}
	}
}

func (rs *RServer) MapFunctionsToHandlers() *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Config.GetHandlers() {
		rs.addHandlerToRouter(r,handl)
	}

	for _, handl := range rs.Config.GetDefaultHandlers() {
		rs.addHandlerToRouter(r,handl)
	}

	return r
}
