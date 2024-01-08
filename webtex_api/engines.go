package webtex_api

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type TexEngine interface {
	Command(filename string) []string
	MoveOutput(directory, input, output string) error
}

func latexmkCommand(engine string, filename string) []string {
	return []string{"/usr/bin/latexmk", fmt.Sprintf("-%s", engine), "-norc", "-dvi", "-jobname=render", filename}
}

type LuaLaTeXEngine struct{}

func (LuaLaTeXEngine) Command(filename string) []string {
	return latexmkCommand("lualatex", filename)
}

func (LuaLaTeXEngine) MoveOutput(directory, _, output string) error {
	return os.Rename(path.Join(directory, "render.dvi"), path.Join(directory, output))
}

type XeLaTeXEngine struct{}

func (XeLaTeXEngine) Command(filename string) []string {
	return latexmkCommand("xelatex", filename)
}

func (XeLaTeXEngine) MoveOutput(directory, _, output string) error {
	return os.Rename(path.Join(directory, "render.dvi"), path.Join(directory, output))
}

type TectonicEngine struct{}

func (TectonicEngine) Command(filename string) []string {
	return []string{"/usr/bin/tectonic", "-X", "compile", "--untrusted", "--outfmt", "xdv", filename}
}

func (TectonicEngine) MoveOutput(directory, input, output string) error {
	resultName := strings.TrimSuffix(input, ".tex") + ".xdv"
	return os.Rename(path.Join(directory, resultName), path.Join(directory, output))
}
