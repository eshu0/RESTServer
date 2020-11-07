package RESTHelpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	sli "github.com/eshu0/simplelogger/interfaces"
)

type ResponseHelper struct {
	Log sli.ISimpleLogger `json:"-"`
}

func NewResponseHelper(logger sli.ISimpleLogger) *ResponseHelper {

	helper := ResponseHelper{}
	helper.Log = logger

	return &helper
}

func (rh *ResponseHelper) WriteIndentJSON(w http.ResponseWriter, Data interface{}) (bool, error) {
	bytes, err := json.MarshalIndent(Data, "", "\t")
	if err != nil {
		rh.Log.LogErrorf("WriteIndentJSON()", "MarshalIndent json failed with %s ", err.Error())
		return false, err
	}
	fmt.Fprint(w, string(bytes))
	return true, nil
}

func (rh *ResponseHelper) WriteJSON(w http.ResponseWriter, Data interface{}) (bool, error) {
	bytes, err := json.Marshal(Data)
	if err != nil {
		rh.Log.LogErrorf("WriteJSON()", "Marshal json failed with %s", err.Error())
		return false, err
	}
	fmt.Fprint(w, string(bytes))
	return true, nil

}
