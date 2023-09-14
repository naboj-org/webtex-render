package main

import (
	"github.com/naboj-org/webtex-render/webtex_api"
	"os"
	"text/template"
)

type Config struct {
	InputFile         *os.File
	OutputFile        *os.File
	EquationDirectory string
	InputURL          string
	OutputURL         string
	Engine            webtex_api.TexEngine
	Template          *template.Template
	OnlyInnerHTML     bool
}
