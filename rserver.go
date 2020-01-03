package RESTServer

import (
	"context"
	"net/http"
	"reflect"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
	"github.com/eshu0/simplelogger/interfaces"
	"github.com/eshu0/simplelogger"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
type RServer struct {
		Port string `json:"port"`
		Handlers []RESTHandler `json:"handlers"`
		DefaultHandlers []DefaultRESTHandler `json:"defaulthandlers"`
		Log slinterfaces.ISimpleLogger	`json:"-"`
		Server *http.Server	`json:"-"`
}

func NewRServer() (RServer, *os.File){

	server := RServer{}
	server.DefaultHandlers = []DefaultRESTHandler{}

	drh := DefaultRESTHandler{}

	drh.URL = "/Admin/Shutdown"
	drh.MethodName = "ShutDown"
	drh.HTTPMethod = "GET"
	drh.FunctionalClass = "RServerCommand"
	drh.MappedClass = RServerCommand{ Server : &server }

	server.DefaultHandlers = append(server.DefaultHandlers, drh)


	// this is the dummy logger object
	logger := &simplelogger.SimpleLogger{}

	// lets open a flie log using the session
	f1 := logger.OpenSessionFileLog("restserver.log", "123")

	server.Log = logger

	return server,f1
}


func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) {

	rs.Log.LogDebugf("Invoke","Method: Looking up %s ", name)

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	val := reflect.ValueOf(any)

	//if !val.IsNil() {
	meth := val.MethodByName(name)
	if !meth.IsZero() && !meth.IsNil() {
		meth.Call(inputs)
	} else {
		rs.Log.LogDebugf("Invoke","Method: %s could not be found ", name)
	}
	//} else {
	//	fmt.Println(fmt.Sprintf("Invoke - Value: %s could not be found ", any))
	//}

}

func (rs *RServer) MakeHandler(MethodName string, any interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs.Invoke(any, MethodName, w, r)
	}
}

func (rs *RServer) MapFunctionsToHandlers(FunctionalMap map[string]interface{}) *mux.Router {

	r := mux.NewRouter()

	for _, handl := range rs.Handlers {

		funcclass, ok := FunctionalMap[handl.FunctionalClass]

		if ok {
			r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, funcclass)).Methods(handl.HTTPMethod)
		}else{
			rs.Log.LogError("MapFunctionsToHandlers","Handlers Error FunctionalClass (%s) doesn't have a function mapped", handl.FunctionalClass)
		}
	}

	for _, handl := range rs.DefaultHandlers {
		r.HandleFunc(handl.URL, rs.MakeHandler(handl.MethodName, handl.MappedClass)).Methods(handl.HTTPMethod)
	}

	return r
}

func (rs *RServer) ShutDown() {
	if rs.Server != nil {
		backg := context.Background()

		if backg != nil {

			if err := rs.Server.Shutdown(backg); err != nil {
				// Error from closing listeners, or context timeout:
				rs.Log.LogDebugf("Shutdown","HTTP server Shutdown: %v", err)
			}
		}else{
			rs.Log.LogError("Shutdown","Called but context.Background() was nil")
		}

	}else{
		rs.Log.LogError("Shutdown","Called but server was nil")
	}

}

func (rs *RServer) ListenAndServe(FunctionalMap map[string]interface{}) {
	r := rs.MapFunctionsToHandlers(FunctionalMap)

	rs.Server =  &http.Server{Addr: ":"+rs.Port, Handler: r}
	//http.ListenAndServe(":"+rs.Port, r)
/*
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := rs.Server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			rs.Log.LogDebugf("Shutdown","HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()
*/


	if err := rs.Server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		rs.Log.LogErrorf("HTTP server ListenAndServe", "%v", err)
	}
}

func (rs *RServer) SaveJSONFile(path string) bool {
	filepath := path + ".json"
	//ok, err := rs.CheckFileExists(filepath)
	//if  ok {
		bytes, err1 := json.MarshalIndent(rs, "", "\t") //json.Marshal(p)
		if err1 != nil {
			rs.Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", path, err1.Error())
			return false
		}

		err2 := ioutil.WriteFile(filepath, bytes, 0644)
		if err2 != nil {
			rs.Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", path, err2.Error())
			return false
		}

		return true
/*
	} else {

		if(err != nil){
			rs.Log.LogErrorf("SaveToFile()", "'%s' was not found to save with error: %s", filepath, err.Error())
		}else{
			rs.Log.LogErrorf("SaveToFile()", "'%s' was not found to save", filepath)
		}

		return false
	}
	*/
}

func (rs *RServer) LoadJSONFile(path string) bool {
	filepath := path + ".json"
	ok, err := rs.CheckFileExists(filepath)
	if  ok {
		bytes, err1 := ioutil.ReadFile(filepath) //ReadAll(jsonFile)
		if err1 != nil {
			rs.Log.LogErrorf("LoadFile()", "Reading '%s' failed with %s ", filepath, err1.Error())
			return false
		}

		rserver := RServer{}

		err2 := json.Unmarshal(bytes, &rserver)

		if err2 != nil {
			rs.Log.LogErrorf("LoadFile()", " Loading %s failed with %s ", filepath, err2.Error())
			return false
		}

		rs.Port = rserver.Port
		rs.Log.LogDebugf("LoadFile()", "Read Port %s ", rserver.Port)

		rs.Handlers = rserver.Handlers

		return true
	} else {

		if(err != nil){
			rs.Log.LogErrorf("LoadFile()", "'%s' was not found to load with error: %s", filepath, err.Error())
		}else{
			rs.Log.LogErrorf("LoadFile()", "'%s' was not found to load", filepath)
		}

		return false
	}
}

func (rs *RServer)  CheckFileExists(filename string) (bool, error) {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false, err
    }
    return !info.IsDir(), nil
}
