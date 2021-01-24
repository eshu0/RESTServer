package RESTServer

import (
	"context"
	"html/template"
	"net/http"
	"reflect"
	"strings"

	Helpers "github.com/eshu0/RESTServer/pkg/helpers"
	appconf "github.com/eshu0/appconfig/pkg"

	sl "github.com/eshu0/simplelogger/pkg"
	sli "github.com/eshu0/simplelogger/pkg/interfaces"
)

// This is the http Server that will host the HTTP requests
var Server *http.Server

type RServer struct {
	sl.AppLogger
	Config *RServerConfig `json:"-"`
	// This map is designed for the functions were there is no types
	// w http.ResponseWriter, r *http.Request
	RawFunctions map[string]interface{} `json:"-"`
	// This accepts Request.ServerRequest
	TypedMap        map[string]interface{}  `json:"-"`
	Templates       *template.Template      `json:"-"`
	RequestHelper   *Helpers.RequestHelper  `json:"-"`
	ResponseHelper  *Helpers.ResponseHelper `json:"-"`
	NotFoundHandler func(w http.ResponseWriter, r *http.Request)
}

//NewRServer a config is needed
func NewRServer(config *RServerConfig) *RServer {
	return NewRServerCustomLog(config, sl.NewApplicationNowLogger())
}

//NewRServerCustomLog if a different logger is wanted instead of the default one
func NewRServerCustomLog(config *RServerConfig, logger sli.ISimpleLogger) *RServer {
	server := RServer{}
	if config == nil {
		// this creates a new server config with defaults
		config = NewRServerConfig()
	}

	server.Config = config
	server.Log = logger
	server.RawFunctions = make(map[string]interface{})
	server.TypedMap = make(map[string]interface{})
	server.RequestHelper = Helpers.NewRequestHelper(server.Log)
	server.ResponseHelper = Helpers.NewResponseHelper(server.Log)
	return &server
}

//Invoke this invokes the func based on the config
func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) []reflect.Value {

	rs.LogDebugf("Invoke", "Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i := range args {
		val := reflect.ValueOf(args[i])
		rs.LogDebugf("Invoke", "ValueOf of arg at [%d] ", i, val)
		rs.LogDebugf("Invoke", "ValueOf of arg at  %v ", val)
		inputs[i] = val
	}

	val := reflect.ValueOf(any)
	rs.LogDebugf("Invoke", "ValueOf %v ", val)
	rs.LogDebugf("Invoke", "Looking up method by %s", name)
	meth := val.MethodByName(name)
	rs.LogDebugf("Invoke", "MethodByName %v", meth)
	/*
		numIn := meth.NumIn() //Count inbound parameters

		for i := 0; i < numIn; i++ {
			inV := meth.In(i)
			in_Kind := inV.Kind() //func
			rs.LogDebugf("Parameter IN: "+strconv.Itoa(i)+"\nKind: %v\nName: %v\n-----------", in_Kind, inV.Name())
		}
	*/
	if meth.IsValid() && !meth.IsNil() {
		return meth.Call(inputs)
	}

	rs.LogErrorf("Invoke", "Method: %s could not be found", name)
	return nil
}

//RegisterFunction Registers a func xyz (w http.ResponseWriter, r *http.Request) for processing
func (rs *RServer) RegisterFunction(FunctionClass string, data interface{}) {
	rs.RawFunctions[FunctionClass] = data
}

//Register Registers a func xyz (request Request.ServerRequest) for processing
func (rs *RServer) Register(FunctionClass string, data interface{}) {
	rs.TypedMap[FunctionClass] = data
}

/// GENERAL OPERATIONS

//ShutDown this shutsdown the server
func (rs *RServer) ShutDown() {
	if Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.LogErrorf("Shutdown", "HTTP server Shutdown: %v", err)
			}

		} else {
			rs.LogError("Shutdown", "Called but context.Background() was nil")
		}

	} else {
		rs.LogError("Shutdown", "Called but server was nil")
	}
}

//ListenAndServe This starts the Server, runs the mapping, loads templates and hosts on the port and arddress set in the config
func (rs *RServer) ListenAndServe() {

	httphandler := rs.MapFunctionsToHandlers()

	if rs.NotFoundHandler != nil {
		rs.LogInfo("ListenAndServe", "Not Found (404) Handler is set")
		httphandler.NotFoundHandler = http.HandlerFunc(rs.NotFoundHandler)
	}

	rs.LogDebug("ListenAndServe", "LoadTemplates started")
	rs.LoadTemplates()
	rs.LogDebug("ListenAndServe", "LoadTemplates finished")
	address := rs.Config.GetAddress()

	Server = &http.Server{Addr: address, Handler: httphandler}

	rs.PrintDetails()
	rs.LogInfo("ListenAndServe", address+" Starting.....")

	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		rs.LogErrorEf("ListenAndServe", "HTTP server ListenAndServe %v", err)
	}
}

//SaveConfig saves server config to disk
func (rs *RServer) SaveConfig() bool {

	if rs.Config == nil {
		rs.LogError("SaveConfig", "Config was nil")
		return false
	}

	if err := rs.Config.Save(); err != nil {
		rs.LogErrorEf("SaveConfig", "SaveConfig - %v", err)
		return false
	}
	return true
}

//LoadConfig loads server config from disk
func (rs *RServer) LoadConfig() bool {

	if rs.Config == nil {
		rs.LogError("LoadConfig", "Config was nil")
		return false
	}

	if err := rs.Config.Load(); err != nil {
		rs.LogErrorEf("LoadConfig", "LoadConfig - %v", err)
		return false
	}

	return true
}

//DefaultServer Creates a default server
func DefaultServer(ConfigFilePath *string) *RServer {
	// Create a new REST Server
	server := NewRServer(nil)

	// has a conifg file been provided?
	if ConfigFilePath != nil && len(*ConfigFilePath) > 0 {
		server.LogDebugf("DefaultServer", "Custom config file path is %s", *ConfigFilePath)

		// load this first
		server.Config.ConfigFilePath = *ConfigFilePath

	} else {
		server.LogDebugf("DefaultServer", "Default config file path is %s", appconf.DefaultFilePath)
		// load this first
		server.Config.ConfigFilePath = appconf.DefaultFilePath
	}
	ok := server.LoadConfig()

	// we failed to load the configuration file
	if !ok {
		server.LogError("DefaultServer ", "failed to load configuration file")
	}

	server.LogDebugf("DefaultServer", "Loaded config from %s", server.Config.ConfigFilePath)

	return server
}

func (rs *RServer) PrintDetails() {

	rs.LogInfof("PrintDetails", "Address: %s", rs.Config.GetAddress())
	rs.LogInfof("PrintDetails", "Template Filepath: %s", rs.Config.GetTemplatePath())
	rs.LogInfof("PrintDetails", "Template FileTypes: %s", strings.Join(rs.Config.GetTemplateFileTypes(), ","))
	rs.LogInfof("PrintDetails", "Cache Templates: %t", rs.Config.GetCacheTemplates())
	rs.LogInfo("PrintDetails", "Handlers: ")
	for _, handl := range rs.Config.GetHandlers() {
		rs.LogInfof("PrintDetails", "Handler: %s", handl.MethodName)
	}

	rs.LogInfo("PrintDetails", "DefaultHandlers: ")

	for _, handl := range rs.Config.GetDefaultHandlers() {
		rs.LogInfof("PrintDetails", "Default Handler: %s", handl.MethodName)
	}

}
