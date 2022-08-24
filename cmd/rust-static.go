package cmd

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func trySyncRustStaticFile(relFilePath string) error {
	var (
		err error
	)
	filePath := filepath.Join(DefaultConfig.ProjectRoot, "rust-static", relFilePath)
	var sInfo os.FileInfo
	if sInfo, err = os.Stat(filePath); os.IsNotExist(err) {
		u, _ := url.Parse(DefaultConfig.RemoteRustStaticURL)
		u.Path = filepath.Join(u.Path, relFilePath)
		if err = doSyncFromRemote(filePath, u.String()); err != nil {
			return err
		}
	} else if !sInfo.Mode().IsRegular() {
		return errors.New("crate isn't a regular file")
	}
	return nil
}

func registerRustStaticRouter(router *gin.Engine) {

	router.GET("/rust-static/:handle/*fileSuffixPath", func(ctx *gin.Context) {
		var (
			err            error
			handle         = ctx.Param("handle")
			fileSuffixPath = ctx.Param("fileSuffixPath")
		)
		relPath := filepath.Join(handle, fileSuffixPath)
		if err = trySyncRustStaticFile(relPath); err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename="+filepath.Base(fileSuffixPath))
		ctx.File(filepath.Join(DefaultConfig.ProjectRoot, "rust-static", relPath))
	})

}
