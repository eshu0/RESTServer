package RESTConfig

import (
	"encoding/json"
	"io/ioutil"
	"os"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	slinterfaces "github.com/eshu0/simplelogger/pkg/interfaces"
)

const string DefaultFilePath = "./config.json"

//
// Server configuration interface
// Gets, Setters etc
//
type IRServerConfig interface {
	GetAddress() string
	HasTemplate() bool
	GetCacheTemplates() bool
	GetTemplatePath() string
	GetTemplateFileTypes() []string

	Save(ConfigFilePath string, Log slinterfaces.ISimpleLogger) bool
	Load(ConfigFilePath string, Log slinterfaces.ISimpleLogger) (IRServerConfig, bool)
	AddHandler(Handler Handlers.RESTHandler)
	AddDefaultHandler(Handler Handlers.RESTHandler)

	GetHandlers() []Handlers.RESTHandler
	GetHandlersLen() int

	GetDefaultHandlers() []Handlers.RESTHandler
	GetDefaultHandlersLen() int
}

type RServerConfig struct {
	Port              string                 `json:"port"`
	Handlers          []Handlers.RESTHandler `json:"handlers"`
	DefaultHandlers   []Handlers.RESTHandler `json:"defaulthandlers"`
	TemplateFilepath  string                 `json:"templatefilepath"`
	TemplateFileTypes []string               `json:"templatefiletypes"`
	CacheTemplates    bool                   `json:"cachetemplates"`
}

func NewRServerConfig() IRServerConfig {
	Config := RServerConfig{}
	Config.DefaultHandlers = []Handlers.RESTHandler{}
	Config.Handlers = []Handlers.RESTHandler{}
	Config.Port = "7777"
	Config.TemplateFileTypes = []string{".tmpl", ".html"}
	Config.CacheTemplates = false
	return &Config
}

func (rsc *RServerConfig) HasTemplate() bool {

	if &rsc.TemplateFilepath == nil {
		return false
	}

	if rsc.TemplateFilepath == "" {
		return false
	}

	return true
}

func (rsc *RServerConfig) GetTemplatePath() string {
	return rsc.TemplateFilepath
}

func (rsc *RServerConfig) GetCacheTemplates() bool {
	return rsc.CacheTemplates
}

func (rsc *RServerConfig) GetTemplateFileTypes() []string {
	return rsc.TemplateFileTypes
}

func (rsc *RServerConfig) GetHandlersLen() int {
	return len(rsc.Handlers)
}

func (rsc *RServerConfig) GetHandlers() []Handlers.RESTHandler {
	return rsc.Handlers
}

func (rsc *RServerConfig) GetDefaultHandlers() []Handlers.RESTHandler {
	return rsc.DefaultHandlers
}

func (rsc *RServerConfig) GetDefaultHandlersLen() int {
	return len(rsc.DefaultHandlers)
}

func (rsc *RServerConfig) GetAddress() string {
	return ":" + rsc.Port
}

func (rsc *RServerConfig) AddDefaultHandler(Handler Handlers.RESTHandler) {
	rsc.DefaultHandlers = append(rsc.DefaultHandlers, Handler)
}

func (rsc *RServerConfig) AddHandler(Handler Handlers.RESTHandler) {
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
