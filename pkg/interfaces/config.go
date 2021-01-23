package rsinterfaces

import (
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
	slinterfaces "github.com/eshu0/simplelogger/pkg/interfaces"
)

// IRServerConfig  Server configuration interface
// Gets, Setters etc
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
