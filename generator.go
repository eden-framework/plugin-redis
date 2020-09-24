package plugin_redis

import (
	"fmt"
	"github.com/eden-framework/eden-framework/pkg/generator"
	"github.com/eden-framework/eden-framework/pkg/generator/files"
	"path"
)

type GenerationPlugin struct {
}

func (g *GenerationPlugin) NewApplicationGenerationPoint(opt generator.ServiceOption, cwd string) string {
	globalPkgPath := path.Join(opt.PackageName, "internal/global")
	globalFilePath := path.Join(cwd, "internal/global")
	tpl := fmt.Sprintf(`,
		{{ .UseWithoutAlias "github.com/eden-framework/eden-framework/pkg/application" "" }}.WithConfig(&{{ .UseWithoutAlias "%s" "%s" }}.RedisConfig)`, globalPkgPath, globalFilePath)
	return tpl
}

func (g *GenerationPlugin) FileGenerationPoint(opt generator.ServiceOption, cwd string) *files.GoFile {
	file := files.NewGoFile("global")
	file.WithBlock(`
var RedisConfig = struct {
	Redis *{{ .UseWithoutAlias "github.com/eden-framework/plugin-redis" "" }}.Redis
}{
	Redis: &{{ .UseWithoutAlias "github.com/eden-framework/plugin-redis" "" }}.Redis{
		Host: "localhost",
		Port: 6379,
	},
}
`)

	return file
}
