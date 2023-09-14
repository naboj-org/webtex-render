package main

import (
	"github.com/naboj-org/webtex-render/webtex_api"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"text/template"
)

type Renderer struct {
	Name         string               `yaml:"name"`
	Key          string               `yaml:"key"`
	TemplatePath string               `yaml:"template"`
	Template     *template.Template   `yaml:"-"`
	CacheDir     string               `yaml:"-"`
	Engine       webtex_api.TexEngine `yaml:"-"`
	EngineString string               `yaml:"engine"`
}

type RendererMap map[string]Renderer

func loadConfig(fileName string) (RendererMap, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var conf []Renderer
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	rendererMap := RendererMap{}
	for _, renderer := range conf {
		renderer.Template, err = template.ParseFiles(renderer.TemplatePath)
		if err != nil {
			return nil, err
		}

		renderer.CacheDir = path.Join("cache", renderer.Name)
		err := os.MkdirAll(renderer.CacheDir, 0755)
		if err != nil {
			return nil, err
		}

		renderer.Engine, err = webtex_api.GetEngine(renderer.EngineString)
		if err != nil {
			return nil, err
		}

		rendererMap[renderer.Name] = renderer
	}

	return rendererMap, nil
}
