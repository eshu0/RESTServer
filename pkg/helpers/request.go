package RESTHelpers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	sli "github.com/eshu0/simplelogger/pkg/interfaces"
	mux "github.com/gorilla/mux"
)

type RequestHelper struct {
	Log sli.ISimpleLogger `json:"-"`
}

func NewRequestHelper(logger sli.ISimpleLogger) *RequestHelper {

	helper := RequestHelper{}
	helper.Log = logger

	return &helper
}

func (rh *RequestHelper) GetRequestId(r *http.Request, name string) *int {
	vars := mux.Vars(r)
	rh.Log.LogDebugf("GetRequestId", "Got the following %v for %s ", vars[name], name)
	Id, err := strconv.Atoi(vars[name])
	if err != nil {
		rh.Log.LogErrorf("GetRequestId", "Got the following error parsing %s for %s", name, err.Error())
		return nil
	} else {
		return &Id
	}
}

func (rh *RequestHelper) GetRequestIds(r *http.Request, names []string) map[string]*int {
	vars := mux.Vars(r)
	results := make(map[string]*int)
	for _, name := range names {
		rh.Log.LogDebugf("GetRequestIds", "Got the following %v for %s", vars[name], name)
		id, err := strconv.Atoi(vars[name])
		if err != nil {
			rh.Log.LogErrorf("GetRequestIds", "Got the following error parsing %s for %s", name, err.Error())
			results[name] = nil
		} else {
			results[name] = &id
		}
	}
	return results
}

func (rh *RequestHelper) ParseForm(r *http.Request, names []string) map[string]string {
	results := make(map[string]string)

	if err := r.ParseForm(); err != nil {
		rh.Log.LogErrorf("ParseForm", "Got the following error parsing form %s", err.Error())
		return results
	}
	for _, name := range names {
		v := r.FormValue(name)
		rh.Log.LogDebugf("ParseForm", "Got the following %s for %s", v, name)
		results[name] = v
	}
	return results
}

func (rh *RequestHelper) ReadBody(r *http.Request) ([]byte, error) {
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		rh.Log.LogErrorf("ReadBody", "Got the following error while reading body %s", err1.Error())
		return []byte{}, err1
	}
	rh.Log.LogDebugf("ReadBody", "Got the following request body %s", string(body))
	return body, err1
}

func (rh *RequestHelper) ReadJSONRequest(r *http.Request, Data interface{}) (interface{}, error) {
	body, err := rh.ReadBody(r)

	if err != nil {
		rh.Log.LogErrorf("ReadJSONRequest", "Got the following error while reading body %s", err.Error())
		return nil, err
	}

	rh.Log.LogDebugf("ReadJSONRequest", "Got the following request body %s", string(body))

	d := map[string]interface{}{}
	json.Unmarshal(body, &d)

	firstArg := reflect.TypeOf(Data)
	s := reflect.New(firstArg).Elem()

	structPtr := reflect.New(firstArg).Elem()
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		for j, f := range d {
			rh.Log.LogDebugf("ReadJSONRequest", "j :%+v\n", j)
			rh.Log.LogDebugf("ReadJSONRequest", "%v - %v - %v - %v\n", typeOfT, typeOfT.Field(i), typeOfT.Field(i).Tag, typeOfT.Field(i).Tag.Get("json"))
			rh.Log.LogDebugf("ReadJSONRequest", "%v - %v - %v - %v\n", typeOfT, typeOfT.Field(i), typeOfT.Field(i).Tag, typeOfT.Field(i).Tag.Get("json"))

			withoutomit := typeOfT.Field(i).Tag.Get("json")
			withoutomit = strings.Replace(withoutomit, ",omitempty", "", -1)
			if withoutomit == j {
				rh.Log.LogDebugf("ReadJSONRequest", "Name :%+v\n", typeOfT.Field(i).Name)

				fl := structPtr.FieldByName(typeOfT.Field(i).Name)
				rh.Log.LogDebugf("ReadJSONRequest", "Kind :%+v\n", fl.Kind())

				switch fl.Kind() {
				case reflect.Bool:
					fl.SetBool(f.(bool))
				case reflect.Int, reflect.Int64:
					c, _ := f.(float64)
					rh.Log.LogDebugf("ReadJSONRequest", "c :%+v\n", c)

					fl.SetInt(int64(c))
				case reflect.String:
					rh.Log.LogDebugf("ReadJSONRequest", "f :%+v\n", f)
					fl.SetString(f.(string))
				}
			}
		}
	}
	rh.Log.LogDebugf("ReadJSONRequest", "%+v\n", structPtr)

	return structPtr.Interface(), nil
}
