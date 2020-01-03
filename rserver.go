package RESTServer

import (
	"fmt"
	"net/http"
	"reflect"
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/gorilla/mux"
	"github.com/eshu0/simplelogger/interfaces"
)

// The new router function creates the router and
// returns it to us. We can now use this function
// to instantiate and test the router outside of the main function
type RServer struct {
		Port string `json:"port"`
		Handlers []RESTHandler `json:"handlers"`
		Log slinterfaces.ISimpleLogger	`json:"_"`
}

func (rs *RServer) Invoke(any interface{}, name string, args ...interface{}) {

	fmt.Println(fmt.Sprintf("Invoke - Method: Looking up %s ", name))

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
		fmt.Println(fmt.Sprintf("Invoke - Method: %s could not be found ", name))
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
		}
		/*
			switch handl.FunctionalClass {
			case "type1":
				r.HandleFunc(handl.URL, makeHandler(handl.MethodName, handl.Data)).Methods(handl.HTTPMethod)
				break
			case "type2":
				r.HandleFunc(handl.URL, makeHandler(handl.MethodName, handl.Data)).Methods(handl.HTTPMethod)
				break
			}
		*/
	}
	return r
}

func (rs *RServer) ListenAndServe(FunctionalMap map[string]interface{}) {
	r := rs.MapFunctionsToHandlers(FunctionalMap)
	http.ListenAndServe(":"+rs.Port, r)
}

func (rs *RServer) SaveToFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		bytes, err := json.MarshalIndent(rs, "", "\t") //json.Marshal(p)
		if err != nil {
			rs.Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", path, err.Error())
			return
		}

		err = ioutil.WriteFile(path+".json", bytes, 0644)
		if err != nil {
			rs.Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", path, err.Error())
		}
	} else {
		rs.Log.LogErrorf("SaveToFile()", "'%s' was not found to save", path)
	}
}

func (rs *RServer) LoadFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		filepath := path + ".json"
		bytes, err := ioutil.ReadFile(filepath) //ReadAll(jsonFile)
		if err != nil {
			rs.Log.LogErrorf("LoadFile()", "Reading '%s' failed with %s ", filepath, err.Error())
			return
		}

		var rserver RServer

		err = json.Unmarshal(bytes, rserver)

		if err != nil {
			rs.Log.LogErrorf("LoadFile()", " Loading %s failed with %s ", filepath, err.Error())
			return
		}

		rs.Port = rserver.Port
		rs.Handlers = rserver.Handlers
	} else {
		rs.Log.LogErrorf("LoadFile()", "'%s' was not found to load", path)
	}
}
