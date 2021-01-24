package RESTServer

import (
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	appconf "github.com/eshu0/appconfig/pkg"
	appconfint "github.com/eshu0/appconfig/pkg/interfaces"
)

//DefaultFilePath is the default path for the server config
const DefaultFilePath = "./config.json"

//RServerConfig This struct is the configuration for the REST server
type RServerConfig struct {
	Parent *appconf.AppConfig
	Data   *ConfigData
}

//ConfigData the data to be stored
type ConfigData struct {
	Port              string                  `json:"port,omitempty"`
	Handlers          []*Handlers.RESTHandler `json:"handlers,omitempty"`
	DefaultHandlers   []*Handlers.RESTHandler `json:"defaulthandlers,omitempty"`
	TemplateFilepath  string                  `json:"templatefilepath,omitempty"`
	TemplateFileTypes []string                `json:"templatefiletypes,omitempty"`
	CacheTemplates    bool                    `json:"cachetemplates,omitempty"`
}

//NewRServerConfig creates new server config
func NewRServerConfig() *RServerConfig {
	conf := appconf.NewAppConfig()
	dc := &RServerConfig{}
	Config, ok := conf.(*appconf.AppConfig)
	if ok {
		dc.Parent = Config
		dc.Parent.SetDefaultFunc(SetServerDefaultConfig)
		dc.Parent.SetDefaults()
		return dc
	}

	return nil

}

//SetServerDefaultConfig ets the defult items
func SetServerDefaultConfig(Config appconfint.IAppConfig) {
	//Config.SetItem("DefaultHandlers", []Handlers.RESTHandler{})
	//Config.SetItem("Handlers", []Handlers.RESTHandler{})
	//Config.SetItem("Port", "7777")
	//Config.SetItem("TemplateFileTypes", []string{".tmpl", ".html"})
	//Config.SetItem("CacheTemplates", false)

	Data := &ConfigData{}
	Data.DefaultHandlers = []*Handlers.RESTHandler{}
	Data.Handlers = []*Handlers.RESTHandler{}
	Data.Port = "7777"
	Data.TemplateFileTypes = []string{".tmpl", ".html"}
	Data.CacheTemplates = false

	Config.SetItem("Data", Data)
}
func (rsc *RServerConfig) getConfigData() *ConfigData {

	data := rsc.Parent.GetItem("Data")
	Config, ok := conf.(*ConfigData)
	if ok {
		return Config
	}
	return nil
}

//HasTemplate returns if a teplate path has been set
func (rsc *RServerConfig) HasTemplate() bool {
	d := rsc.getConfigData()
	if d == nil {
		return false
	}

	return &(d.TemplateFilepath) != nil && len(d.TemplateFilepath) > 0
}

//GetTemplatePath returns the template path
func (rsc *RServerConfig) GetTemplatePath() string {
	d := rsc.getConfigData()
	if d == nil {
		return ""
	}
	return d.TemplateFilepath
}

//GetCacheTemplates returns the cached template paths
func (rsc *RServerConfig) GetCacheTemplates() bool {
	d := rsc.getConfigData()
	if d == nil {
		return false
	}
	return d.CacheTemplates
}

//GetTemplateFileTypes returns the file types for the templates, such as .tmpl, .html
func (rsc *RServerConfig) GetTemplateFileTypes() []string {
	d := rsc.getConfigData()
	if d == nil {
		return []string{}
	}
	return d.TemplateFileTypes
}

//GetHandlersLen this gets length handlers
func (rsc *RServerConfig) GetHandlersLen() int {
	d := rsc.getConfigData()
	if d == nil {
		return -1
	}
	return len(d.Handlers)
}

//GetHandlers this gets the handlers from the config
func (rsc *RServerConfig) GetHandlers() []*Handlers.RESTHandler {
	d := rsc.getConfigData()
	if d == nil {
		return nil
	}
	return d.Handlers
}

//GetDefaultHandlers this gets the default handlers
func (rsc *RServerConfig) GetDefaultHandlers() []*Handlers.RESTHandler {
	d := rsc.getConfigData()
	if d == nil {
		return nil
	}
	return d.DefaultHandlers
}

//GetDefaultHandlersLen this gets length default handlers
func (rsc *RServerConfig) GetDefaultHandlersLen() int {
	d := rsc.getConfigData()
	if d == nil {
		return -1
	}
	return len(d.DefaultHandlers)
}

//GetAddress this gets the server address
func (rsc *RServerConfig) GetAddress() string {
	d := rsc.getConfigData()
	if d == nil {
		return ":7777"
	}
	return ":" + d.Port
}

//AddDefaultHandler this adds a default handler to the configuration
func (rsc *RServerConfig) AddDefaultHandler(Handler Handlers.RESTHandler) {
	d := rsc.getConfigData()
	if d == nil {
		d.DefaultHandlers = append(d.DefaultHandlers, &Handler)
	}
}

//AddHandler this adds a handler to the configuration
func (rsc *RServerConfig) AddHandler(Handler Handlers.RESTHandler) {
	d := rsc.getConfigData()
	if d == nil {
		d.Handlers = append(d.Handlers, &Handler)
	}
}
