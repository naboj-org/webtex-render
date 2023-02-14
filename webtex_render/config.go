package main

import (
	"os"
	"text/template"
)

type Config struct {
	InputFile         *os.File
	OutputFile        *os.File
	EquationDirectory string
	InputURL          string
	OutputURL         string
	Template          *template.Template
	OnlyInnerHTML     bool
}
