package RESTServer

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/eshu0/simplelogger/interfaces"
)

type IRServerConfig interface {
	GetAddress() string
	Save(ConfigFilePath string, Log slinterfaces.ISimpleLogger) bool
	Load(ConfigFilePath string, Log slinterfaces.ISimpleLogger) (IRServerConfig, bool)
	AddHandler(Handler RESTHandler)
	AddDefaultHandler(Handler RESTHandler)
	GetHandlers() []RESTHandler
	GetDefaultHandlers() []RESTHandler
}

type RServerConfig struct {
	Port            string        `json:"port"`
	Handlers        []RESTHandler `json:"handlers"`
	DefaultHandlers []RESTHandler `json:"defaulthandlers"`
}

//server.Config = RServerConfig{}
//server.Config.DefaultHandlers = []RESTHandler{}

func (rsc *RServerConfig) GetHandlers() []RESTHandler {
	return rsc.Handlers
}

func (rsc *RServerConfig) GetDefaultHandlers() []RESTHandler {
	return rsc.DefaultHandlers
}

func (rsc *RServerConfig) GetAddress() string {
	return ":" + rsc.Port
}

func (rsc *RServerConfig) AddDefaultHandler(Handler RESTHandler) {
	rsc.DefaultHandlers = append(rsc.DefaultHandlers, Handler)
}

func (rsc *RServerConfig) AddHandler(Handler RESTHandler) {
	rsc.Handlers = append(rsc.Handlers, Handler)
}

func (rsc *RServerConfig) Save(ConfigFilePath string, Log slinterfaces.ISimpleLogger) bool {
	bytes, err1 := json.MarshalIndent(rsc, "", "\t") //json.Marshal(p)
	if err1 != nil {
		Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", ConfigFilePath, err1.Error())
		return false
	}

	err2 := ioutil.WriteFile(ConfigFilePath, bytes, 0644)
	if err2 != nil {
		Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", ConfigFilePath, err2.Error())
		return false
	}

	return true

}

func (rsc *RServerConfig) Load(ConfigFilePath string, Log slinterfaces.ISimpleLogger) (IRServerConfig, bool) {
	ok, err := rsc.checkFileExists(ConfigFilePath)
	if ok {
		bytes, err1 := ioutil.ReadFile(ConfigFilePath) //ReadAll(jsonFile)
		if err1 != nil {
			Log.LogErrorf("LoadFile()", "Reading '%s' failed with %s ", ConfigFilePath, err1.Error())
			return nil, false
		}

		rserverconfig := RServerConfig{}

		err2 := json.Unmarshal(bytes, &rserverconfig)

		if err2 != nil {
			Log.LogErrorf("LoadFile()", " Loading %s failed with %s ", ConfigFilePath, err2.Error())
			return nil, false
		}

		Log.LogDebugf("LoadFile()", "Read Port %s ", rserverconfig.Port)
		//rs.Log.LogDebugf("LoadFile()", "Port in config %s ", rs.Config.Port)

		return &rserverconfig, true
	} else {

		if err != nil {
			Log.LogErrorf("LoadFile()", "'%s' was not found to load with error: %s", ConfigFilePath, err.Error())
		} else {
			Log.LogErrorf("LoadFile()", "'%s' was not found to load", ConfigFilePath)
		}

		return nil, false
	}
}

func (rsc *RServerConfig) checkFileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}
