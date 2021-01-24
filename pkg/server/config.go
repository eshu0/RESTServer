package RESTServer

import (
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	appconf "github.com/eshu0/appconfig/pkg"
	appconfint "github.com/eshu0/appconfig/pkg/interfaces"
)

//RServerConfig This struct is the configuration for the REST server
type RServerConfig struct {
	Helper *appconf.AppConfigHelper `json:"-"`
	data   *ConfigData              `json:"-"`
}

//ConfigData the data to be stored
type ConfigData struct {
	Port              string                 `json:"port,omitempty"`
	Handlers          []Handlers.RESTHandler `json:"handlers,omitempty"`
	DefaultHandlers   []Handlers.RESTHandler `json:"defaulthandlers,omitempty"`
	TemplateFilepath  string                 `json:"templatefilepath,omitempty"`
	TemplateFileTypes []string               `json:"templatefiletypes,omitempty"`
	CacheTemplates    bool                   `json:"cachetemplates,omitempty"`
}

//NewRServerConfig creates new server config
func NewRServerConfig(filepath string) *RServerConfig {
	dc := &RServerConfig{}
	dc.Helper = appconf.NewAppConfigHelperWithDefault(filepath, dc.SetServerDefaultConfig)
	return dc

}

//SetServerDefaultConfig ets the defult items
func (rsc *RServerConfig) SetServerDefaultConfig(Config appconfint.IAppConfig) {

	Data := &ConfigData{}
	Data.DefaultHandlers = []Handlers.RESTHandler{}
	Data.Handlers = []Handlers.RESTHandler{}
	Data.Port = "7777"
	Data.TemplateFileTypes = []string{".tmpl", ".html"}
	Data.CacheTemplates = false

	rsc.SetConfigData(Data)
}

//GetConfigData returns the config data from the store
func (rsc *RServerConfig) GetConfigData() *ConfigData {
	if rsc.data == nil {
		data := rsc.Helper.Config.GetItem("Data")
		Config, ok := data.(*ConfigData)
		if ok {
			rsc.data = Config
			return Config
		}
		return nil

	}
	return rsc.data

}

//SetConfigData sets the config data to the store
func (rsc *RServerConfig) SetConfigData(data *ConfigData) {
	rsc.Helper.Config.SetItem("Data", data)
	rsc.data = nil
}

//HasTemplate returns if a teplate path has been set
func (rsc *RServerConfig) HasTemplate() bool {
	d := rsc.GetConfigData()
	if d == nil {
		return false
	}

	return &(d.TemplateFilepath) != nil && len(d.TemplateFilepath) > 0
}

//GetTemplatePath returns the template path
func (rsc *RServerConfig) GetTemplatePath() string {
	d := rsc.GetConfigData()
	if d == nil {
		return ""
	}
	return d.TemplateFilepath
}

//GetCacheTemplates returns the cached template paths
func (rsc *RServerConfig) GetCacheTemplates() bool {
	d := rsc.GetConfigData()
	if d == nil {
		return false
	}
	return d.CacheTemplates
}

//GetTemplateFileTypes returns the file types for the templates, such as .tmpl, .html
func (rsc *RServerConfig) GetTemplateFileTypes() []string {
	d := rsc.GetConfigData()
	if d == nil {
		return []string{}
	}
	return d.TemplateFileTypes
}

//GetHandlersLen this gets length handlers
func (rsc *RServerConfig) GetHandlersLen() int {
	d := rsc.GetConfigData()
	if d == nil {
		return -1
	}
	return len(d.Handlers)
}

//GetHandlers this gets the handlers from the config
func (rsc *RServerConfig) GetHandlers() []Handlers.RESTHandler {
	d := rsc.GetConfigData()
	if d == nil {
		return []Handlers.RESTHandler{}
	}
	return d.Handlers
}

//GetDefaultHandlers this gets the default handlers
func (rsc *RServerConfig) GetDefaultHandlers() []Handlers.RESTHandler {
	d := rsc.GetConfigData()
	if d == nil {
		return []Handlers.RESTHandler{}
	}
	return d.DefaultHandlers
}

//GetDefaultHandlersLen this gets length default handlers
func (rsc *RServerConfig) GetDefaultHandlersLen() int {
	d := rsc.GetConfigData()
	if d == nil {
		return -1
	}
	return len(d.DefaultHandlers)
}

//GetAddress this gets the server address
func (rsc *RServerConfig) GetAddress() string {
	d := rsc.GetConfigData()
	if d == nil {
		return ":7777"
	}
	return ":" + d.Port
}

//AddDefaultHandler this adds a default handler to the configuration
func (rsc *RServerConfig) AddDefaultHandler(Handler Handlers.RESTHandler) {
	d := rsc.GetConfigData()
	if d != nil {
		handlers := d.DefaultHandlers
		handlers = append(handlers, Handler)
		d.DefaultHandlers = handlers
		rsc.SetConfigData(d)
	}
}

//AddHandler this adds a handler to the configuration
func (rsc *RServerConfig) AddHandler(Handler Handlers.RESTHandler) {
	d := rsc.GetConfigData()
	if d != nil {
		handlers := d.Handlers
		handlers = append(handlers, Handler)
		d.Handlers = handlers
		rsc.SetConfigData(d)
	}
}
