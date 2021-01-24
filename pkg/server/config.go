package RESTServer

import (
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	appconfint "github.com/eshu0/appconfig/pkg/interfaces"
	appconf "github.com/eshu0/appconfig/pkg"

)

//DefaultFilePath is the default path for the server config
const DefaultFilePath = "./config.json"

//RServerConfig This struct is the configuration for the REST server
type RServerConfig struct {

	Parent            *appconf.AppConfig
	Port              string                 `json:"port,omitempty"`
	Handlers          []Handlers.RESTHandler `json:"handlers,omitempty"`
	DefaultHandlers   []Handlers.RESTHandler `json:"defaulthandlers,omitempty"`
	TemplateFilepath  string                 `json:"templatefilepath,omitempty"`
	TemplateFileTypes []string               `json:"templatefiletypes,omitempty"`
	CacheTemplates    bool                   `json:"cachetemplates,omitempty"`
}

func NewRServerConfig() *RServerConfig {
	conf := NewAppConfig()
	dc := &DummyConfig{}
	Config, ok := conf.(*RServerConfig)
	if ok {
		dc.Parent = Config
		dc.Parent.SetDefaultFunc(SetServerDefaultConfig)
		dc.Parent.SetDefaults()
		return dc
	}

	return nil

}

//NewRServerConfig creates a new server configuation with default settings
func SetServerDefaultConfig(Config appconfint.IAppConfig) {
	Config.SetItem("DefaultHandlers",[]Handlers.RESTHandler{})
	Config.SetItem("Handlers",[]Handlers.RESTHandler{})
	Config.SetItem("Port", "7777")
	Config.SetItem("TemplateFileTypes",[]string{".tmpl", ".html"})
	Config.SetItem("CacheTemplates", = false)
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
