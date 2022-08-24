package cmd

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type router struct {
	reg     *regexp.Regexp
	method  string
	handler func(HandlerReq)
	Rpc     string
}

var routes = []router{
	{regexp.MustCompile("(.*?)/git-upload-pack$"), http.MethodPost, serviceRpc, "upload-pack"},
	{regexp.MustCompile("(.*?)/git-receive-pack$"), http.MethodPost, serviceRpc, "receive-pack"},

	{regexp.MustCompile("(.*?)/info/refs$"), http.MethodGet, getInfoRefs, ""},

	{regexp.MustCompile("(.*?)/HEAD$"), http.MethodGet, getTextFile, ""},
	{regexp.MustCompile("(.*?)/objects/info/alternates$"), http.MethodGet, getTextFile, ""},
	{regexp.MustCompile("(.*?)/objects/info/http-alternates$"), http.MethodGet, getTextFile, ""},
	{regexp.MustCompile("(.*?)/objects/info/[^/]*$"), http.MethodGet, getTextFile, ""},

	{regexp.MustCompile("(.*?)/objects/info/packs$"), http.MethodGet, getInfoPacks, ""},
	{regexp.MustCompile("(.*?)/objects/[0-9a-f]{2}/[0-9a-f]{38}$"), http.MethodGet, getLooseObject, ""},
	{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.pack$"), http.MethodGet, getPackFile, ""},
	{regexp.MustCompile("(.*?)/objects/pack/pack-[0-9a-f]{40}\\.idx$"), http.MethodGet, getIdxFile, ""},
}

func configureServerHandler() (*gin.Engine, error) {
	//gin.SetMode(gin.ReleaseMode)
	e := gin.Default()

	e.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, nil)
	})

	e.GET("/api/v1/crates/:crate/:version/download", GetCrates)

	e.Any("/crates.io-index/*action", func(ctx *gin.Context) {
		url := strings.ToLower(ctx.Request.URL.Path)
		method := ctx.Request.Method

		for _, _route := range routes {
			if method != _route.method {
				continue
			}
			if m := _route.reg.FindStringSubmatch(url); m != nil {
				file := strings.Replace(url, m[1]+"/", "", 1)
				dir, err := getGitDir(m[1])

				if err != nil {
					log.Print(err)
					renderNotFound(ctx.Writer)
					return
				}

				hr := HandlerReq{ctx.Writer, ctx.Request, _route.Rpc, dir, file}
				_route.handler(hr)
				return
			}
		}

		renderMethodNotAllowed(ctx.Writer, ctx.Request)
	})

	registerRustStaticRouter(e)

	return e, nil
}
