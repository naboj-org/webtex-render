package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/naboj-org/webtex-render/webtex_api"
	"log"
	"net/http"
	"path"
)

var Renderers RendererMap

func getRoot(c *gin.Context) {
	c.String(http.StatusOK, "https://github.com/naboj-org/webtex-render")
}

func getRender(c *gin.Context) {
	rendererName := c.Param("renderer")
	renderer, exists := Renderers[rendererName]
	if !exists {
		c.String(http.StatusForbidden, "wrong renderer credentials")
		return
	}

	mac := hmac.New(sha256.New, []byte(renderer.Key))
	mac.Write([]byte(c.Query("eq")))
	expected := mac.Sum(nil)
	signature, err := hex.DecodeString(c.Query("sig"))
	if err != nil {
		c.String(http.StatusForbidden, "wrong renderer credentials")
		return
	}

	if !hmac.Equal(expected, signature) {
		c.String(http.StatusForbidden, "wrong renderer credentials")
		return
	}

	equation, err := base64.StdEncoding.DecodeString(c.Query("eq"))
	if err != nil {
		c.String(http.StatusBadRequest, "could not base64 decode eq query parameter")
		return
	}

	outputPath, err := webtex_api.EquationSvg(string(equation), renderer.CacheDir, renderer.Template)
	if err != nil {
		c.String(http.StatusInternalServerError, "error generating equation")
		log.Printf("Error while generating: %v\n", err)
		return
	}

	c.File(path.Join(renderer.CacheDir, outputPath))
}

func main() {
	var err error
	Renderers, err = loadConfig("config.yml")
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.GET("/", getRoot)
	r.GET("/render/:renderer/", getRender)
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
