package RESTServer

import (
	"net/http"
	"strings"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	Request "github.com/eshu0/RESTServer/pkg/request"

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

func (rs *RServer) CreateFunctionHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, HandleJSONRequest bool, HandleJSONResponse bool) Handlers.RESTHandler {
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

func (rs *RServer) AddStaticHandler(URL string, StaticDir string) {
	rs.Config.AddHandler(rs.CreateStaticHandler(URL, StaticDir))
}

func (rs *RServer) AddFunctionHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string) {
	rs.Config.AddHandler(rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, false, false))
}

func (rs *RServer) AddJSONRequestFunctionHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, DataType interface{}) {
	fc := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, true, false)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) AddJSONResponseFunctionHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, DataType interface{}) {
	fc := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, false, true)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) AddJSONFunctionHandler(URL string, MethodName string, HTTPMethod string, FunctionalClass string, DataType interface{}) {
	fc := rs.CreateFunctionHandler(URL, MethodName, HTTPMethod, FunctionalClass, true, true)
	fc.JSONRequestType = DataType
	rs.Config.AddHandler(fc)
}

func (rs *RServer) MakeHandlerFunction(handler Handlers.RESTHandler, any interface{}, istyped bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !istyped {
			rs.LogDebug("MakeHandlerFunction", "Untyped called")
			sr := Request.CreateServerRawRequest(w, r)
			rs.Invoke(any, handler.MethodName, sr.Writer, sr.Request)
		} else {
			rs.LogDebug("MakeHandlerFunction", "Typed called")

			// JSON request expected
			if handler.JSONRequest {
				// parse the JSON
				data, jsonerr := rs.RequestHelper.ReadJSONRequest(r, handler.JSONRequestType)
				// no error
				if jsonerr == nil {
					rs.LogDebugf("MakeHandlerFunction", "Data: %v ", data)

					// are we returning JSON?
					if handler.JSONResponse {
						// we are invoking a JSON method this should do the writing
						rs.LogDebugf("MakeHandlerFunction", "The response is JSON response for %s and %s(data)", handler.HTTPMethod, handler.MethodName)
						resp := rs.Invoke(any, handler.MethodName, Request.CreateServerPayloadRequest(w, r, data))
						if len(resp) > 0 {
							rs.ResponseHelper.WriteJSON(w, resp[0].Interface())
						}
					} else {
						rs.LogDebugf("MakeHandlerFunction", "The response is not a JSON response for %s and %s(w,r,data)", handler.HTTPMethod, handler.MethodName)
						// we are invoking the method that will do the writing out etc
						rs.Invoke(any, handler.MethodName, Request.CreateServerPayloadRequest(w, r, data))
					}
				} else {
					rs.LogErrorEf("MakeHandlerFunction", "ReadJSONRequest Error : %s", jsonerr)
					return
				}

			} else {
				if handler.JSONResponse {
					rs.LogDebugf("MakeHandlerFunction", "The response is JSON response for %s and %s(r)", handler.HTTPMethod, handler.MethodName)
					// we are just letting the request do the work and then the data will be returned
					resp := rs.Invoke(any, handler.MethodName, Request.CreateServerRawRequest(w, r))
					if len(resp) > 0 {
						rs.ResponseHelper.WriteJSON(w, resp[0].Interface())
					}
				} else {
					rs.LogDebugf("MakeHandlerFunction", "The response is not a JSON response for %s and %s(w,r)", handler.HTTPMethod, handler.MethodName)
					// not json request or response -> raw read/write
					rs.Invoke(any, handler.MethodName, Request.CreateServerRawRequest(w, r))
				}
			}

		}

	}
}

func (rs *RServer) regHandlerToRouter(r *mux.Router, handl Handlers.RESTHandler, funcclass interface{}, istyped bool) {
	if len(handl.TemplatePath) > 0 || len(handl.TemplateFileName) > 0 || len(handl.TemplateBlob) > 0 {
		rs.LogDebugf("regHandlerToRouter", "Handlers: Adding Template function %s - %s", handl.HTTPMethod, handl.MethodName)

		if handl.HTTPMethod != "" {
			if strings.Contains(handl.HTTPMethod, ",") {
				rs.LogDebugf("regHandlerToRouter", "Template Method is multiple")
				r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(strings.Split(handl.HTTPMethod, ",")...)
			} else {
				rs.LogDebug("regHandlerToRouter", "Template Method is a single HTTP Verb")
				r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass)).Methods(handl.HTTPMethod)
			}
		} else {
			rs.LogDebug("regHandlerToRouter", "Template Method does not have a HTTP Verb")
			r.HandleFunc(handl.URL, rs.MakeTemplateHandlerFunction(handl, funcclass))
		}
	} else {
		rs.LogDebugf("regHandlerToRouter", "Handlers: Adding %s - %s", handl.HTTPMethod, handl.MethodName)
		if handl.HTTPMethod != "" {
			if strings.Contains(handl.HTTPMethod, ",") {
				rs.LogDebug("regHandlerToRouter", "Method is multiple HTTP Verbs")
				r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass, istyped)).Methods(strings.Split(handl.HTTPMethod, ",")...)
			} else {
				rs.LogDebug("regHandlerToRouter", "Method is a single HTTP Verb")
				r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass, istyped)).Methods(handl.HTTPMethod)
			}
		} else {
			rs.LogDebug("regHandlerToRouter", "Method does not have a HTTP Verb")
			r.HandleFunc(handl.URL, rs.MakeHandlerFunction(handl, funcclass, istyped))
		}
	}
}

func (rs *RServer) addHandlerToRouter(r *mux.Router, handl Handlers.RESTHandler) {

	funcclass, ok := rs.RawFunctions[handl.FunctionalClass]

	if ok {
		rs.regHandlerToRouter(r, handl, funcclass, false)
	} else {

		funcclass, ok = rs.TypedMap[handl.FunctionalClass]
		if ok {
			rs.regHandlerToRouter(r, handl, funcclass, true)
		} else {
			if handl.StaticDir != "" {
				rs.LogDebugf("addHandlertoRouter", "Handlers: Adding route %s for  static directory %s", handl.URL, handl.StaticDir)
				r.PathPrefix(handl.URL).Handler(http.StripPrefix(handl.URL, http.FileServer(http.Dir(handl.StaticDir))))
			} else {
				rs.LogErrorf("addHandlertoRouter", "Handlers Error FunctionalClass (%s) doesn't have a function %s mapped", handl.FunctionalClass, handl.MethodName)
			}
		}
	}

	// this handles the func
	r.HandleFunc(handl.URL, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
	}).Methods(http.MethodOptions)
}

//MapFunctionsToHandlers Maps functions to the handlers
func (rs *RServer) MapFunctionsToHandlers() *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Config.GetHandlers() {
		rs.addHandlerToRouter(r, handl)
	}

	for _, handl := range rs.Config.GetDefaultHandlers() {
		rs.addHandlerToRouter(r, handl)
	}

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}
