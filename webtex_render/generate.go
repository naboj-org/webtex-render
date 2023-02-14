package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type TemplateContext struct {
	Equation string
}

func generateSvg(config Config, equation, filename string) error {
	dir, err := os.MkdirTemp("", "webtex_render")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	file, err := os.Create(path.Join(dir, "in.tex"))
	if err != nil {
		log.Println("creating in.tex file failed:", err)
		return err
	}

	err = config.Template.Execute(file, TemplateContext{Equation: equation})
	if err != nil {
		log.Println("executing in.tex template failed:", err)
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = runTex(dir)
	if err != nil {
		return err
	}

	err = convertDvi(dir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(filename), 0755)
	if err != nil {
		log.Println("creating output directory failed:", err)
		return err
	}

	outputFile, err := os.Create(filename)
	if err != nil {
		log.Println("opening output file", filename, "failed:", err)
		return err
	}
	defer outputFile.Close()

	dvisvgmResult, err := os.Open(path.Join(dir, "render.svg"))
	if err != nil {
		log.Println("opening render.svg failed:", err)
		return err
	}
	defer dvisvgmResult.Close()

	_, err = io.Copy(outputFile, dvisvgmResult)
	return err
}

func EquationSvg(config Config, equation string) (string, error) {
	h := sha1.New()
	h.Write([]byte(equation))
	sha := hex.EncodeToString(h.Sum(nil))
	filename := path.Join(config.EquationDirectory, sha[0:2], fmt.Sprintf("%s.svg", sha))

	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// Generate
		err := generateSvg(config, equation, filename)
		return filename, err
	} else {
		return "", err
	}
}
