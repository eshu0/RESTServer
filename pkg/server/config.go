package RESTServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	rsinterfaces "github.com/eshu0/RESTServer/pkg/interfaces"
)

//DefaultFilePath is the default path for the server config
const DefaultFilePath = "./config.json"

//RServerConfig This struct is the configuration for the REST server
type RServerConfig struct {
	rsinterfaces.IRServerConfig
	Port              string                 `json:"port,omitempty"`
	Handlers          []Handlers.RESTHandler `json:"handlers,omitempty"`
	DefaultHandlers   []Handlers.RESTHandler `json:"defaulthandlers,omitempty"`
	TemplateFilepath  string                 `json:"templatefilepath,omitempty"`
	TemplateFileTypes []string               `json:"templatefiletypes,omitempty"`
	CacheTemplates    bool                   `json:"cachetemplates,omitempty"`
}

//NewRServerConfig creates a new server configuation with default settings
func NewRServerConfig() rsinterfaces.IRServerConfig {
	Config := RServerConfig{}
	Config.DefaultHandlers = []Handlers.RESTHandler{}
	Config.Handlers = []Handlers.RESTHandler{}
	Config.Port = "7777"
	Config.TemplateFileTypes = []string{".tmpl", ".html"}
	Config.CacheTemplates = false
	return &Config
}

//HasTemplate returns if a teplate path has been set
func (rsc *RServerConfig) HasTemplate() bool {
	return &(rsc.TemplateFilepath) != nil && len(rsc.TemplateFilepath) > 0
}

//GetTemplatePath returns the template path
func (rsc *RServerConfig) GetTemplatePath() string {
	return rsc.TemplateFilepath
}

//GetCacheTemplates returns the cached template paths
func (rsc *RServerConfig) GetCacheTemplates() bool {
	return rsc.CacheTemplates
}

//GetTemplateFileTypes returns the file types for the templates, such as .tmpl, .html
func (rsc *RServerConfig) GetTemplateFileTypes() []string {
	return rsc.TemplateFileTypes
}

//GetHandlersLen this gets length handlers
func (rsc *RServerConfig) GetHandlersLen() int {
	return len(rsc.Handlers)
}

//GetHandlers this gets the handlers from the config
func (rsc *RServerConfig) GetHandlers() []Handlers.RESTHandler {
	return rsc.Handlers
}

//GetDefaultHandlers this gets the default handlers
func (rsc *RServerConfig) GetDefaultHandlers() []Handlers.RESTHandler {
	return rsc.DefaultHandlers
}

//GetDefaultHandlersLen this gets length default handlers
func (rsc *RServerConfig) GetDefaultHandlersLen() int {
	return len(rsc.DefaultHandlers)
}

//GetAddress this gets the server address
func (rsc *RServerConfig) GetAddress() string {
	return ":" + rsc.Port
}

//AddDefaultHandler this adds a default handler to the configuration
func (rsc *RServerConfig) AddDefaultHandler(Handler Handlers.RESTHandler) {
	rsc.DefaultHandlers = append(rsc.DefaultHandlers, Handler)
}

//AddHandler this adds a handler to the configuration
func (rsc *RServerConfig) AddHandler(Handler Handlers.RESTHandler) {
	rsc.Handlers = append(rsc.Handlers, Handler)
}

//Save This saves the configuration from a file path
func (rsc *RServerConfig) Save(ConfigFilePath string) error {
	bytes, err1 := json.MarshalIndent(rsc, "", "\t") //json.Marshal(p)
	if err1 != nil {
		//Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", ConfigFilePath, err1.Error())
		return err1
	}

	err2 := ioutil.WriteFile(ConfigFilePath, bytes, 0644)
	if err2 != nil {
		//Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", ConfigFilePath, err2.Error())
		return err2
	}

	return nil

}

//Load This loads the configuration from a file path
func (rsc *RServerConfig) Load(ConfigFilePath string) (rsinterfaces.IRServerConfig, error) {
	ok, err := rsc.checkFileExists(ConfigFilePath)
	if ok {
		bytes, err1 := ioutil.ReadFile(ConfigFilePath) //ReadAll(jsonFile)
		if err1 != nil {
			return nil, fmt.Errorf("Reading '%s' failed with %s ", ConfigFilePath, err1.Error())
		}

		rserverconfig := RServerConfig{}

		err2 := json.Unmarshal(bytes, &rserverconfig)

		if err2 != nil {
			return nil, fmt.Errorf("Loading %s failed with %s ", ConfigFilePath, err2.Error())
		}

		//Log.LogDebugf("LoadFile()", "Read Port %s ", rserverconfig.Port)
		//rs.Log.LogDebugf("LoadFile()", "Port in config %s ", rs.Config.Port)
		return &rserverconfig, nil
	}

	if err != nil {
		return nil, fmt.Errorf("'%s' was not found to load with error: %s", ConfigFilePath, err.Error())
	}

	return nil, fmt.Errorf("'%s' was not found to load", ConfigFilePath)
}

func (rsc *RServerConfig) checkFileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}
