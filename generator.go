package redis

import (
	"github.com/eden-framework/eden-framework/pkg/generator"
	"github.com/eden-framework/eden-framework/pkg/generator/files"
)

type GenerationPlugin struct {
}

func (g *GenerationPlugin) NewApplicationGenerationPoint(opt generator.ServiceOption) string {
	panic("implement me")
}

func (g *GenerationPlugin) FileGenerationPoint(opt generator.ServiceOption) *files.GoFile {
	file := files.NewGoFile("global")
	file.WithBlock(``)

	return file
}
