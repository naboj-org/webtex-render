package webtex_api

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type TexEngine string

const (
	TEX_LUALATEX TexEngine = "lualatex"
	TEX_XELATEX  TexEngine = "xelatex"
)

func GetEngine(engine string) (TexEngine, error) {
	switch strings.ToLower(engine) {
	case "xelatex":
		return TEX_XELATEX, nil
	case "lualatex":
		return TEX_LUALATEX, nil
	default:
		return "", fmt.Errorf("unknown engine '%s'", engine)
	}
}

// execCommand runs a given command with args in the workdir
func execCommand(workdir, command string, arg ...string) (string, error) {
	log.Println(command, arg)
	cmd := exec.Command(command, arg...)
	cmd.Dir = workdir
	var b bytes.Buffer
	cmd.Stderr = &b
	cmd.Stdout = &b
	err := cmd.Run()
	return b.String(), err
}

// runTex starts latexmk to generate `render.dvi` from `in.tex` in a given directory
func runTex(directory string, engine TexEngine) error {
	output, err := execCommand(directory, "/usr/bin/latexmk", fmt.Sprintf("-%s", engine), "-norc", "-dvi", "-jobname=render", "in.tex")
	if err != nil {
		log.Println("latexmk failed:", err)
		log.Println(output)
	}
	return err
}

// convertDvi converts `render.dvi` to `rednder.svg` in a given directory
func convertDvi(directory string) error {
	output, err := execCommand(directory, "/usr/bin/dvisvgm", "--exact", "--bbox=papersize", "--no-fonts", "render.dvi")
	if err != nil {
		log.Println("dvisvgm failed:", err)
		log.Println(output)
	}
	return err
}
