package webtex_api

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GetEngine(engine string) (TexEngine, error) {
	switch strings.ToLower(engine) {
	case "xelatex":
		return XeLaTeXEngine{}, nil
	case "lualatex":
		return LuaLaTeXEngine{}, nil
	case "tectonic":
		return TectonicEngine{}, nil
	default:
		return nil, fmt.Errorf("unknown engine '%s'", engine)
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
	command := engine.Command("in.tex")
	output, err := execCommand(directory, command[0], command[1:]...)
	if err != nil {
		log.Println("latex failed:", err)
		log.Println(output)
		return err
	}

	return engine.MoveOutput(directory, "in.tex", "render.dvi")
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
