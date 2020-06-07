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

func (rs *RServer) AddStaticHandler(URL string, StaticDir string)  {
	rs.Config.AddHandler(rs.CreateStaticHandler(URL,StaticDir))
}

func (rs *RServer) AddFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string)  {
	rs.Config.AddHandler(rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass, false, false))
}

func (rs *RServer) AddJSONRequestFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, DataType interface{})  {
	fc := rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass, true, false)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) AddJSONResponseFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, DataType interface{})  {
	fc := rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass, false, true)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) AddJSONFunctionHandler(URL string, MethodName string,HTTPMethod string, FunctionalClass string, DataType interface{})  {
	fc := rs.CreateFunctionHandler(URL,MethodName,HTTPMethod,FunctionalClass, true, true)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) MakeHandlerFunction(handler Handlers.RESTHandler, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		// JSON request expected
		if handler.JSONRequest {
			// parse the JSON
			data, jsonerr := rs.RequestHelper.ReadJSONRequest(r,handler.JSONRequestType)
			
			// no error
			if jsonerr == nil {
				// are we returning JSON?
				if handler.JSONResponse {
					// we are invoking a JSON method this should do the writing
					rs.Log.LogDebugf("MakeHandlerFunction", "The response is JSON response for %s and %s(data)",handl.HTTPMethod, handl.MethodName)
					resp := rs.Invoke(any,handler.MethodName,data)
					if len(resp) > 0 { 
						rs.ResponseHelper.WriteJSON(w,resp[0].Interface())
					}
				} else{ 
					rs.Log.LogDebugf("MakeHandlerFunction", "The response is not a JSON response for %s and %s(w,r,data)",handl.HTTPMethod, handl.MethodName)
					// we are invoking the method that will do the writing out etc
					rs.Invoke(any, handler.MethodName, w, r, data)
				}
			}else{
				rs.Log.LogErrorf("MakeHandlerFunction", "ReadJSONRequest Error : %s", jsonerr.Error())
				return			
			}
	
		} else {
			if handler.JSONResponse {
				rs.Log.LogDebugf("MakeHandlerFunction", "The response is JSON response for %s and %s(r)",handl.HTTPMethod, handl.MethodName)
				// we are just letting the request do the work and then the data will be returned
				resp := rs.Invoke(any,handler.MethodName,r)
				if len(resp) > 0{ 
					rs.ResponseHelper.WriteJSON(w,resp[0].Interface())
				}
			} else{ 
				rs.Log.LogDebugf("MakeHandlerFunction", "The response is not a JSON response for %s and %s(w,r)",handl.HTTPMethod, handl.MethodName)
				// not json request or response -> raw read/write
				rs.Invoke(any,handler.MethodName, w, r)
			}
		}

	}
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
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass)).Methods(strings.Split(handl.HTTPMethod,",")...)
				}else{
					r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass)).Methods(handl.HTTPMethod)
				}
			}else{
				r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass))
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
