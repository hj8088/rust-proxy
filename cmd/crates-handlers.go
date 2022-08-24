package cmd

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getCratesDir(filePath string) (string, error) {
	root := DefaultConfig.ProjectRoot
	if root == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Print(err)
			return "", err
		}
		root = cwd
	}
	return path.Join(root, filePath), nil
}

func GetCrates(ctx *gin.Context) {
	var (
		err       error
		crate     = ctx.Param("crate")
		version   = ctx.Param("version")
		cratesDir string
	)

	if cratesDir, err = getCratesDir("crates"); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	var (
		cratePath string
		fileName  = fmt.Sprintf("%s-%s.crate", crate, version)
	)
	switch len(crate) {
	case 1:
		fallthrough
	case 2:
		cratePath = filepath.Join(cratesDir, "1", crate, fileName)
	case 3:
		cratePath = filepath.Join(cratesDir, crate[:1], crate, fileName)
	default:
		cratePath = filepath.Join(cratesDir, crate[:2], crate[2:4], crate, fileName)
	}

	var sInfo os.FileInfo
	if sInfo, err = os.Stat(cratePath); os.IsNotExist(err) {
		if DefaultConfig.RemoteProxyURL != "" {
			urlPath := strings.TrimRight(DefaultConfig.RemoteProxyURL, "/")
			urlPath += "/" + crate + "/" + version + "/download"
			if err = doSyncFromRemote(cratePath, urlPath); err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	} else if !sInfo.Mode().IsRegular() {
		ctx.JSON(http.StatusInternalServerError, errors.New("crate isn't a regular file"))
		return
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(cratePath)

	return
}
