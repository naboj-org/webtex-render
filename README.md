<div align="center">
    <img src="https://user-images.githubusercontent.com/11409143/218767084-7c551090-f706-4363-a500-419b42839eec.png" width="128" height="128" />
    <h3>WebTeX Render</h3>
    <p>Render Pandoc's `--webtex` to local images</p>
</div>

## About the project

WebTeX Render is a simple program used to generate static images from LaTeX equations. It was designed to be used as
an alternative to <https://latex.codecogs.com/>. It also allows you to use arbitrary LaTeX packages in your images.

We also provide a simple web API to generate images on demand.

## Usage

```shell
# 1. Convert pandoc document to HTML
pandoc input.md -s --webtex='eqn://' -o pandoc.html

# 2. Generate WebTeX images
./wr -input pandoc.html -output final.html -template template.tex

# Or using one-liner: 
pandoc input.md -s --webtex='eqn://' -t html | ./wr -output final.html -template template.tex
```

When using Pandoc's standalone mode (`-s`), we recommend using the following CSS to align inline equations:

```css
img.math.inline {
    vertical-align: bottom !important;
}
```

## Using the web API

The `webtex_web` requires a configuration file to define "renderers".

```yaml
- name: render_name
  key: secret
  template: template.tex
```

The TeX image can be then generated at URL `http://...:8080/render/render_name/?eq=EQUATION&sig=SIGNATURE`.
`EQUATION` should be a base64 encoded TeX. `SIGNATURE` should be a HMAC-SHA-256 of the `EQUATION` using the key from config.

## Minimal template

```tex
\documentclass{standalone}

\usepackage{scrextend}
\changefontsizes{16pt}

\begin{document}
\strut{}${{ .Equation }}$
\end{document}
```

We have changed the font size to 16pt as this is the default on web. We have also added `\strut{}` at the start of
the document to make sure that all output files are at least `1em` high.

## Building

```shell
go build -o wr ./webtex_render
go build -o wr_web ./webtex_web
```

## Docker image

We also provide a Docker image. The image contains both `wr` and `wr_web` binaries, but is set up to start the web
API by default.

```shell
docker run --rm \
  -v ./config.yml:/app/config.yml:ro \
  -v ./template.tex:/app/template.tex:ro \
  -p 8080:8080 \
  ghcr.io/naboj-org/webtex_render
```

