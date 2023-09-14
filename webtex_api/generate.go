package webtex_api

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"
)

type TemplateContext struct {
	Equation string
}

func generateSvg(template *template.Template, equation, filename string, engine TexEngine) error {
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

	err = template.Execute(file, TemplateContext{Equation: equation})
	if err != nil {
		log.Println("executing in.tex template failed:", err)
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = runTex(dir, engine)
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

func EquationSvg(equation string, cacheDir string, template *template.Template, engine TexEngine) (string, error) {
	h := sha1.New()
	h.Write([]byte(equation))
	sha := hex.EncodeToString(h.Sum(nil))
	filenameEnd := path.Join(sha[0:2], fmt.Sprintf("%s.svg", sha))
	filename := path.Join(cacheDir, filenameEnd)

	if _, err := os.Stat(filename); err == nil {
		return filenameEnd, nil
	} else if errors.Is(err, os.ErrNotExist) {
		// Generate
		err := generateSvg(template, equation, filename, engine)
		return filenameEnd, err
	} else {
		return "", err
	}
}
