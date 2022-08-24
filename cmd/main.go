package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"net/url"
	"path/filepath"
	"strings"
)

func newApp(name string) *cli.App {
	app := &cli.App{
		Name:  name,
		Usage: "is a tool to proxy the crates.io repositories, for offline usage of cargo",
		Action: func(context *cli.Context) error {

			r, _ := configureServerHandler()

			var (
				serverPort  = context.Int("server-port")
				projectRoot = context.String("project-root")
			)

			DefaultConfig = &Config{
				AuthPassEnvVar:      "",
				AuthUserEnvVar:      "",
				DefaultEnv:          "",
				ProjectRoot:         projectRoot,
				GitBinPath:          "git",
				UploadPack:          true,
				ReceivePack:         true,
				RemoteProxyURL:      context.String("remote-crates-url"),
				RemoteRustStaticURL: strings.TrimSuffix(context.String("remote-rust-static-url"), "/"),
			}

			if _, err := url.Parse(DefaultConfig.RemoteRustStaticURL); err != nil {
				return err
			}
			if _, err := url.Parse(DefaultConfig.RemoteProxyURL); err != nil {
				return err
			}

			return r.Run(fmt.Sprintf(":%d", serverPort))
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "server-port",
				Value: 8080,
			},
			&cli.StringFlag{
				Required: true,
				Name:     "project-root",
			},
			&cli.StringFlag{
				Name:     "remote-crates-url",
				Value:    "https://crates-io.proxy.ustclug.org/api/v1/crates",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "remote-rust-static-url",
				Value:    "https://mirrors.ustc.edu.cn/rust-static",
				Required: false,
			},
		},
	}

	return app
}

func Main(args []string) {
	// Set the minio app name.
	appName := filepath.Base(args[0])

	// Run the app - exit on error.
	if err := newApp(appName).Run(args); err != nil {
		log.Fatal(err)
	}
}
