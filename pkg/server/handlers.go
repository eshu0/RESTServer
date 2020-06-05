package RESTServer

import (
	"net/http"
	"strings"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"

	mux "github.com/gorilla/mux"
)

// Create Handlers helper functions
func (rs *RServer) CreateStaticHandler(URL string, StaticDir string) Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}
	drhr.URL = URL
	drhr.StaticDir = StaticDir
	drhr.JSONRequest = false
	drhr.JSONResponse = false
	return drhr
}

func (rs *RServer) CreateFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, HandleJSONRequest bool, HandleJSONResponse bool) Handlers.RESTHandler {
	drhr := Handlers.RESTHandler{}
	drhr.URL = URL
	drhr.MethodName = MethodName
	drhr.HTTPMethod = HTTPMethod
	drhr.FunctionalClass = FunctionalClass
	drhr.JSONRequest = HandleJSONRequest
	drhr.JSONResponse = HandleJSONResponse
	return drhr
}
// end 

func (rs *RServer) MakeHandlerFunction(MethodName string, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if handler.HandleJSONResponse {
			if handler.HandleJSONRequest {
				data, jsonerr := rs.RequestHelper.ReadJSONRequest(r)
				if err != nil {
					resp := rs.Invoke(any, handler.MethodName,data)
				}else{
					rs.Log.LogErrorf("MakeHandlerFunction", "ReadJSONRequest Error : %s", jsonerr.Error())
					return			
				}
			} else{ 
				resp := rs.Invoke(any, handler.MethodName, r)
			}

			rs.RequestHelper.WriteJSON(w,resp)
			
		} else {
			if handler.HandleJSONRequest {
				data, jsonerr := rs.RequestHelper.ReadJSONRequest(r)
				if err != nil {
					rs.Invoke(any, handler.MethodName, w, r, data)
				}else{
					rs.Log.LogErrorf("MakeHandlerFunction", "ReadJSONRequest Error : %s", jsonerr.Error())
					return			
				}
			} else{ 
				rs.Invoke(any, MethodName, w, r)
			}
		}

	}
}

func (rs *RServer) AddStaticHandler(URL string, StaticDir string)  {
	rs.Config.AddHandler(rs.CreateStaticHandler(URL,StaticDir))
}

func (rs *RServer) AddFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string)  {
	rs.Config.AddHandler(rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass))
}

func (rs *RServer) addHandlerToRouter(r *mux.Router, handl Handlers.RESTHandler){

	funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

	if ok {
		if handl.TemplatePath != "" || handl.TemplateFileName != "" || handl.TemplateBlob != ""  {
			rs.Log.LogDebugf("addHandlertoRouter", "Handlers: Adding Template function %s - %s",handl.HTTPMethod, handl.MethodName)
			
			if handl.HTTPMethod != ""{
				if strings.Contains(handl.HTTPMethod,",") {
					rs.Log.LogDebugf("addHandlertoRouter", "Method is multiple")	
					r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(strings.Split(handl.HTTPMethod,",")...)
				}else{
					r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(handl.HTTPMethod)
				}
			}else{
				r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass))
			}
		} else {
			rs.Log.LogDebugf("addHandlertoRouter", "Handlers: Adding %s - %s",handl.HTTPMethod, handl.MethodName)
			if handl.HTTPMethod != ""{
				if strings.Contains(handl.HTTPMethod,",") {
					rs.Log.LogDebugf("addHandlertoRouter", "Method is multiple")
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass)).Methods(strings.Split(handl.HTTPMethod,",")...)
				}else{
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
				}
			}else{
				r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl.MethodName, funcclass))
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
	r.HandleFunc(handl.URL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
}

func (rs *RServer) MapFunctionsToHandlers() *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Config.GetHandlers() {
		rs.addHandlerToRouter(r,handl)
	}

	for _, handl := range rs.Config.GetDefaultHandlers() {
		rs.addHandlerToRouter(r,handl)
	}
	
	r.Use(mux.CORSMethodMiddleware(r))

	return r
}
