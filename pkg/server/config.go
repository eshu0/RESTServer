package RESTServer

import (
	"fmt"

	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	appconf "github.com/eshu0/appconfig/pkg"
	appconfint "github.com/eshu0/appconfig/pkg/interfaces"
)

//RServerConfig This struct is the configuration for the REST server
type RServerConfig struct {
	Helper *appconf.AppConfigHelper `json:"-"`
	cache  *ConfigData              `json:"-"`
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
	helper := appconf.NewAppConfigHelperWithDefault(filepath, dc.SetServerDefaultConfig)

	if helper != nil {
		dc.Helper = helper
		// we call this after the helper has been set!
		dc.Helper.LoadedConfig.SetDefaults()
	}

	return dc

}

//Loads the config from disk
func (rsc *RServerConfig) Load() error {
	fmt.Printf("conf before load %v\n", rsc.Helper.LoadedConfig)

	// load the data
	if err := rsc.Helper.Load(); err != nil {
		return err
	}

	fmt.Printf("conf after load %v\n", rsc.Helper.LoadedConfig)

	// reset the cache
	rsc.cache = nil

	// this rebuilds the cache
	rsc.GetConfigData()

	fmt.Printf("conf after data config %v\n", rsc.Helper.LoadedConfig)

	return nil
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
	if rsc.cache == nil {
		fmt.Println("cache is nil")
		data := rsc.Helper.LoadedConfig.GetItem("Data")
		fmt.Printf("data %v\n", data)
		Config, ok := data.(*ConfigData) //(map[string]*ConfigData)
		if ok {
			fmt.Printf("cast ok %v\n", Config)
			rsc.cache = Config
			return Config
		}

		Config1, ok1 := data.(ConfigData) //(map[string]*ConfigData)
		if ok1 {
			fmt.Printf("cast1 ok %v\n", Config1)
			rsc.cache = &Config1
			return &Config1
		}

		Config2, ok2 := data.((map[string]*ConfigData))
		if ok2 {
			fmt.Printf("cast2 ok %v\n", Config2)
		}

		Config3, ok3 := data.((map[string]interface{}))
		if ok3 {
			fmt.Printf("cast3 ok %v\n", Config3)
			for key, element := range Config3 {
				fmt.Println("Key:", key, "=>", "Element:", element)
			}
			Config4, ok4 := Config3["Data"].(*ConfigData)
			if ok4 {
				fmt.Printf("cast4 ok %v\n", Config4)
				rsc.cache = Config4
				return Config4
			}
		}

		fmt.Printf("cast failed %v\n", Config)
		return nil

	}
	return rsc.cache

}

//SetConfigData sets the config data to the store
func (rsc *RServerConfig) SetConfigData(data *ConfigData) {

	// reset the cache
	rsc.cache = nil
	fmt.Printf("before %v\n", rsc.Helper.LoadedConfig)

	// set the data ietm
	rsc.Helper.LoadedConfig.SetItem("Data", data)

	fmt.Printf("after %v\n", rsc.Helper.LoadedConfig)

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
		panic("config data was nil")
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
