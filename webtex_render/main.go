package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

const VERSION = "v23.1"

func main() {
	var config Config
	inputFilename := flag.String("input", "-", "input file or - for stdin")
	outputFilename := flag.String("output", "-", "output file or - for stdout")
	equationDirectory := flag.String("eqdir", "equations", "directory to output equations")
	inputUrl := flag.String("inurl", "eqn://", "input webtex URL prefix")
	outputUrl := flag.String("outurl", "equations", "image src root")
	templateFilename := flag.String("template", "", "TeX template file")
	onlyInner := flag.Bool("innerhtml", false, "export only inner HTML of the result (without <html><body> tags)")
	version := flag.Bool("version", false, "prints current roxy version")
	flag.Parse()

	if *version {
		log.Printf("WebTeX Render version %v\n", VERSION)
		os.Exit(0)
	}

	var err error
	if *inputFilename == "-" {
		config.InputFile = os.Stdin
	} else {
		config.InputFile, err = os.Open(*inputFilename)
		if err != nil {
			panic(err)
		}
		defer config.InputFile.Close()
	}

	if *outputFilename == "-" {
		config.OutputFile = os.Stdout
	} else {
		config.OutputFile, err = os.Create(*outputFilename)
		if err != nil {
			panic(err)
		}
		defer config.OutputFile.Close()
	}

	config.EquationDirectory = *equationDirectory
	config.InputURL = *inputUrl
	config.OutputURL = *outputUrl
	config.OnlyInnerHTML = *onlyInner

	if *templateFilename == "" {
		log.Println("please provide a template file.")
		return
	}

	config.Template, err = template.ParseFiles(*templateFilename)
	if err != nil {
		panic(err)
	}

	err = parseHtml(config)
	if err != nil {
		panic(err)
	}
}
