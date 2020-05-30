package RESTHelpers

import (
	"net/http"
	"strconv"

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

func (rs *RequestHelper) GetRequestId(r *http.Request, name string) *int {
	vars := mux.Vars(r)
	rs.Log.LogInfof("GetRequestId","Got the following %s for %v ",name, vars[name])
	Id, err := strconv.Atoi(vars[name])
	if err != nil {
		rs.Log.LogErrorf("GetRequestId","Got the following error parsing %s for %s",name, err.Error())
		return nil
	}else{
		return &Id
	}
}

func (rs *RequestHelper) GetRequestIds(r *http.Request, names []string) map[string]*int{
	vars := mux.Vars(r)
	results := make(map[string]*int)
	for _, name := range names {
		rs.Log.LogInfof("GetRequestIds","Got the following %s for %v",name, vars[name])
		id, err := strconv.Atoi(vars[name])
		if err != nil {
			rs.Log.LogErrorf("GetRequestIds","Got the following error parsing %s for %s",name, err.Error())
			results[name] = nil
		}else{
			results[name] = &id
		}
	}
	return results
}