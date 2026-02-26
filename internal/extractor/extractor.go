package extractor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func GetDependencies(entryPoint string) ([]string, error) {
	var localFiles []string

	trackerPlugin := api.Plugin{
		Name: "dependency-tracker",
		Setup: func(build api.PluginBuild) {
			// 1. FILTRO DE ACTIVOS: Marcamos imágenes y otros archivos no deseados como externos
			// Esto evita que esbuild intente abrirlos o procesarlos.
			build.OnResolve(api.OnResolveOptions{Filter: `\.(png|jpg|jpeg|gif|svg|ico|webp|avif|woff|woff2|ttf|eot)$`},
				func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					return api.OnResolveResult{External: true}, nil
				})

			// 2. RASTREO DE ARCHIVOS DE CÓDIGO:
			build.OnLoad(api.OnLoadOptions{Filter: `.*`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				// Solo agregamos a la lista si es un archivo de código (js, ts, jsx, tsx, etc.)
				ext := strings.ToLower(filepath.Ext(args.Path))
				if isCodeExtension(ext) {
					localFiles = append(localFiles, args.Path)
				}
				return api.OnLoadResult{}, nil
			})
		},
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryPoint},
		Bundle:      true,
		Write:       false,
		Outdir:      "out",
		Packages:    api.PackagesExternal,
		Plugins:     []api.Plugin{trackerPlugin},

		Loader: map[string]api.Loader{
			".js":   api.LoaderJSX,
			".jsx":  api.LoaderJSX,
			".ts":   api.LoaderTS,
			".tsx":  api.LoaderTSX,
			".json": api.LoaderJSON,
			".css":  api.LoaderCSS,
			".scss": api.LoaderCSS,
			".less": api.LoaderCSS,
		},

		LogOverride: map[string]api.LogLevel{
			"commonjs-variable-in-esm": api.LogLevelSilent,
		},
	})

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("error de esbuild al procesar %s: %s", entryPoint, result.Errors[0].Text)
	}

	return localFiles, nil
}

// isCodeExtension ayuda a filtrar solo los archivos que queremos en nuestro Markdown
func isCodeExtension(ext string) bool {
	switch ext {
	case ".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".json", ".css", ".scss", ".less":
		return true
	default:
		return false
	}
}
