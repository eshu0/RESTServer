package RESTServer

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"

	"github.com/eshu0/simplelogger"
	"github.com/eshu0/simplelogger/interfaces"
	"github.com/gorilla/mux"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function

var Server *http.Server

type RServerConfig struct {
	Port            string        `json:"port"`
	Handlers        []RESTHandler `json:"handlers"`
	DefaultHandlers []RESTHandler `json:"defaulthandlers"`
}

type RServer struct {
	Config         RServerConfig              `json:"-"`
	Log            slinterfaces.ISimpleLogger `json:"-"`
	FunctionalMap  map[string]interface{}     `json:"-"`
	ConfigFilePath string                     `json:"-"`
	//Server *http.Server	`json:"-"`
}

func NewRServer() (*RServer, *os.File) {

	server := RServer{}
	server.Config = RServerConfig{}
	server.Config.DefaultHandlers = []RESTHandler{}
	server.FunctionalMap = make(map[string]interface{})

	// this is the dummy logger object
	logger := &simplelogger.SimpleLogger{}

	// lets open a flie log using the session
	f1 := logger.OpenSessionFileLog("restserver.log", "123")

	server.Log = logger

	return &server, f1
}

func (server *RServer) AddDefaults() {

	// Default commands for server
	// These should be removed if not required
	rsc := RServerCommand{Server: server}

	server.FunctionalMap["RServerCommand"] = rsc

	server.Config.DefaultHandlers = append(server.Config.DefaultHandlers, rsc.CreateShutDownHandler())
	server.Config.DefaultHandlers = append(server.Config.DefaultHandlers, rsc.CreateListHandler())
	server.Config.DefaultHandlers = append(server.Config.DefaultHandlers, rsc.CreateLoadConfigHandler())
	server.Config.DefaultHandlers = append(server.Config.DefaultHandlers, rsc.CreateSaveConfigHandler())

	for _, handl := range server.Config.DefaultHandlers {
		server.Log.LogDebugf("NewRServerWithDefaults", "Default Handler: Added %s", handl.MethodName)
	}
}

func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) {

	rs.Log.LogDebugf("Invoke", "Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	val := reflect.ValueOf(any)

	meth := val.MethodByName(name)
	if !meth.IsZero() && !meth.IsNil() {
		meth.Call(inputs)
	} else {
		rs.Log.LogDebugf("Invoke", "Method: %s could not be found ", name)
	}

}

func (rs *RServer) MakeHandler(MethodName string, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs.Invoke(any, MethodName, w, r)
	}
}

func (rs *RServer) MapFunctionsToHandlers() *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Config.Handlers {

		funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

		if ok {
			rs.Log.LogDebugf("MapFunctionsToHandlers", "Handlers: Adding %s", handl.MethodName)
			r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
		} else {
			rs.Log.LogError("MapFunctionsToHandlers", "Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)
		}
	}

	for _, handl := range rs.Config.DefaultHandlers {

		funcclass, ok := rs.FunctionalMap[handl.FunctionalClass]

		if ok {
			rs.Log.LogDebugf("MapFunctionsToHandlers", "Default Handlers: Adding %s", handl.MethodName)
			r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
		} else {
			rs.Log.LogError("MapFunctionsToHandlers", "Default Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)
		}

	}

	return r
}

func (rs *RServer) ShutDown() {
	if Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.Log.LogDebugf("Shutdown", "HTTP server Shutdown: %v", err)
			}
		} else {
			rs.Log.LogError("Shutdown", "Called but context.Background() was nil")
		}

	} else {
		rs.Log.LogError("Shutdown", "Called but server was nil")
	}

}

func (rs *RServer) ListenAndServe() {
	r := rs.MapFunctionsToHandlers()

	Server = &http.Server{Addr: ":" + rs.Config.Port, Handler: r}

	rs.FunctionalMap["RServerCommand"] = RServerCommand{Server: rs}

	rs.Log.LogDebugf("ListenAndServe", "Listening on: %s", rs.Config.Port)

	if err := Server.ListenAndServe(); err != http.ErrServerClosed {
		rs.Log.LogErrorf("HTTP server ListenAndServe", "%v", err)
	}
}

func (rs *RServer) SaveJSONFile() bool {
	filepath := rs.ConfigFilePath + ".json"
	bytes, err1 := json.MarshalIndent(rs.Config, "", "\t") //json.Marshal(p)
	if err1 != nil {
		rs.Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", filepath, err1.Error())
		return false
	}

	err2 := ioutil.WriteFile(filepath, bytes, 0644)
	if err2 != nil {
		rs.Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", filepath, err2.Error())
		return false
	}

	return true

}

func (rs *RServer) LoadJSONFile() bool {
	filepath := rs.ConfigFilePath + ".json"
	ok, err := rs.CheckFileExists(filepath)
	if ok {
		bytes, err1 := ioutil.ReadFile(filepath) //ReadAll(jsonFile)
		if err1 != nil {
			rs.Log.LogErrorf("LoadFile()", "Reading '%s' failed with %s ", filepath, err1.Error())
			return false
		}

		rserverconfig := RServerConfig{}

		err2 := json.Unmarshal(bytes, &rserverconfig)

		if err2 != nil {
			rs.Log.LogErrorf("LoadFile()", " Loading %s failed with %s ", filepath, err2.Error())
			return false
		}

		rs.Config = rserverconfig
		rs.Log.LogDebugf("LoadFile()", "Read Port %s ", rserverconfig.Port)
		rs.Log.LogDebugf("LoadFile()", "Port in config %s ", rs.Config.Port)

		return true
	} else {

		if err != nil {
			rs.Log.LogErrorf("LoadFile()", "'%s' was not found to load with error: %s", filepath, err.Error())
		} else {
			rs.Log.LogErrorf("LoadFile()", "'%s' was not found to load", filepath)
		}

		return false
	}
}

func (rs *RServer) CheckFileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}
