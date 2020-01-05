package RESTServer

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/eshu0/RESTServer/pkg/handlers"
	"github.com/eshu0/simplelogger/interfaces"
)

type IRServerConfig interface {
	GetAddress() string
	Save(ConfigFilePath string, Log slinterfaces.ISimpleLogger) bool
	Load(ConfigFilePath string, Log slinterfaces.ISimpleLogger) (IRServerConfig, bool)
	AddHandler(Handler RESTServer.RESTHandler)
	AddDefaultHandler(Handler RESTServer.RESTHandler)
	GetHandlers() []RESTServer.RESTHandler
	GetDefaultHandlers() []RESTServer.RESTHandler
}

type RServerConfig struct {
	Port            string                   `json:"port"`
	Handlers        []RESTServer.RESTHandler `json:"handlers"`
	DefaultHandlers []RESTServer.RESTHandler `json:"defaulthandlers"`
}

func NewRServerConfig() IRServerConfig {
	Config := RServerConfig{}
	Config.DefaultHandlers = []RESTServer.RESTHandler{}
	Config.Handlers = []RESTServer.RESTHandler{}
	Config.Port = "7777"
	return &Config
}

func (rsc *RServerConfig) GetHandlers() []RESTServer.RESTHandler {
	return rsc.Handlers
}

func (rsc *RServerConfig) GetDefaultHandlers() []RESTServer.RESTHandler {
	return rsc.DefaultHandlers
}

func (rsc *RServerConfig) GetAddress() string {
	return ":" + rsc.Port
}

func (rsc *RServerConfig) AddDefaultHandler(Handler RESTServer.RESTHandler) {
	rsc.DefaultHandlers = append(rsc.DefaultHandlers, Handler)
}

func (rsc *RServerConfig) AddHandler(Handler RESTServer.RESTHandler) {
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
