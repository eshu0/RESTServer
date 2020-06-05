package RESTHelpers

import (
	"net/http"
	"strconv"
	"encoding/json"
	"io/ioutil"

	sli "github.com/eshu0/simplelogger/interfaces"
	mux "github.com/gorilla/mux"
)

type RequestHelper struct {
	Log sli.ISimpleLogger      `json:"-"`
}

func NewRequestHelper(logger sli.ISimpleLogger) *RequestHelper {

	helper := RequestHelper{}
	helper.Log = logger

	return &helper
}

func (rh *RequestHelper) GetRequestId(r *http.Request, name string) *int {
	vars := mux.Vars(r)
	rh.Log.LogDebugf("GetRequestId","Got the following %v for %s ",vars[name],name)
	Id, err := strconv.Atoi(vars[name])
	if err != nil {
		rh.Log.LogErrorf("GetRequestId","Got the following error parsing %s for %s",name, err.Error())
		return nil
	}else{
		return &Id
	}
}

func (rh *RequestHelper) GetRequestIds(r *http.Request, names []string) map[string]*int{
	vars := mux.Vars(r)
	results := make(map[string]*int)
	for _, name := range names {
		rh.Log.LogDebugf("GetRequestIds","Got the following %v for %s",vars[name], name)
		id, err := strconv.Atoi(vars[name])
		if err != nil {
			rh.Log.LogErrorf("GetRequestIds","Got the following error parsing %s for %s",name, err.Error())
			results[name] = nil
		}else{
			results[name] = &id
		}
	}
	return results
}

func (rh *RequestHelper) ParseForm(r *http.Request, names []string) map[string]string{
	results := make(map[string]string)

	if err := r.ParseForm(); err != nil {
		rh.Log.LogErrorf("ParseForm","Got the following error parsing form %s",err.Error())
		return results
	}
	for _, name := range names {
		v := r.FormValue(name)
		rh.Log.LogDebugf("ParseForm","Got the following %s for %s",v,name)
		results[name] = v
	}
	return results
}

func (rh *RequestHelper) ReadBody(r *http.Request) ([]byte, error) {
	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		rh.Log.LogErrorf("ReadBody","Got the following error while reading body %s",err1.Error())
		return []byte{},err1
	}
	rh.Log.LogDebugf("ReadBody","Got the following request body %s",string(body))
	return body,err1
}

func (rh *RequestHelper) ReadJSONRequest(r *http.Request,Data interface{}) (interface{}, error) {
	body, err := rh.ReadBody(r)
	if err != nil {
		rh.Log.LogErrorf("ReadJSONRequest","Got the following error while reading body %s",err.Error())
		return nil, err
	}
	rh.Log.LogDebugf("ReadJSONRequest","Got the following request body %s",string(body))

	//err := json.NewDecoder(string(body)).Decode(&Data)
	err = json.Unmarshal(body, &Data)
	if err != nil {
		rh.Log.LogErrorf("ReadJSONRequest","Got the following error while unmarchsalling JSON %s",err.Error())
		return nil, err
	}

	return Data, nil
}
