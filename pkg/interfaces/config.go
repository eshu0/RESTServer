package rsinterfaces

import (
	Handlers "github.com/eshu0/RESTServer/pkg/handlers"
)

// IRServerConfig  Server configuration interface
// Gets, Setters etc
type IRServerConfig interface {
	GetAddress() string
	HasTemplate() bool
	GetCacheTemplates() bool
	GetTemplatePath() string
	GetTemplateFileTypes() []string

	Save(ConfigFilePath string) error
	Load(ConfigFilePath string) (IRServerConfig, error)
	AddHandler(Handler Handlers.RESTHandler)
	AddDefaultHandler(Handler Handlers.RESTHandler)

	GetHandlers() []Handlers.RESTHandler
	GetHandlersLen() int

	GetDefaultHandlers() []Handlers.RESTHandler
	GetDefaultHandlersLen() int
}
