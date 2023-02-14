package main

import (
	"flag"
	"log"
	"os"
	"text/template"
)

func main() {
	var config Config
	var inputFilename = flag.String("input", "-", "input file or - for stdin")
	var outputFilename = flag.String("output", "-", "output file or - for stdout")
	var equationDirectory = flag.String("eqdir", "equations", "directory to output equations")
	var inputUrl = flag.String("inurl", "eqn://", "input webtex URL prefix")
	var outputUrl = flag.String("outurl", "equations", "image src root")
	var templateFilename = flag.String("template", "", "TeX template file")
	var onlyInner = flag.Bool("innerhtml", false, "export only inner HTML of the result (without <html><body> tags)")
	flag.Parse()

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
