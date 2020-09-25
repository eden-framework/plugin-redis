package main

import (
	"fmt"
	"github.com/eden-framework/plugins"
	"path"
)

var Plugin GenerationPlugin

type GenerationPlugin struct {
}

func (g *GenerationPlugin) GenerateEntryPoint(opt plugins.Option, cwd string) string {
	globalPkgPath := path.Join(opt.PackageName, "internal/global")
	globalFilePath := path.Join(cwd, "internal/global")
	tpl := fmt.Sprintf(`,
		{{ .UseWithoutAlias "github.com/eden-framework/eden-framework/pkg/application" "" }}.WithConfig(&{{ .UseWithoutAlias "%s" "%s" }}.RedisConfig)`, globalPkgPath, globalFilePath)
	return tpl
}

func (g *GenerationPlugin) GenerateFilePoint(opt plugins.Option, cwd string) []*plugins.FileTemplate {
	file := plugins.NewFileTemplate("global", path.Join(cwd, "internal/global/redis.go"))
	file.WithBlock(`
var RedisConfig = struct {
	Redis *{{ .UseWithoutAlias "github.com/eden-framework/plugin-redis/redis" "" }}.Redis
}{
	Redis: &{{ .UseWithoutAlias "github.com/eden-framework/plugin-redis/redis" "" }}.Redis{
		Host: "localhost",
		Port: 6379,
	},
}
`)

	return []*plugins.FileTemplate{file}
}
