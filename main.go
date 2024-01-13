package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"ffcu/configuration"
	"ffcu/utils"
	"log"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	config, err := configuration.ReadConfig()
	if err != nil {
		log.Println("can't read config, creating a new one")

		config, err = configuration.CreateConfig()
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer config.SaveConfig()

	app := &cli.App{
		Name:  "ffcu",
		Usage: "Helper CLI tool to auto update your firefox config",

		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "ffcu config set-dir xxx",
				Subcommands: []*cli.Command{
					{
						Name:        "set-dir",
						Description: "The dir where the user.js and chrome dir is",
						Action: func(ctx *cli.Context) error {
							dir := ctx.Args().First()
							if dir == "" {
								return errors.New("please provide a directory")
							}

							if _, err := os.Stat(path.Join(dir, "prefs.js")); os.IsNotExist(err) {
								return errors.New("the directory you provided does not have a prefs.js")
							}

							config.ProfileDir = dir
							config.SaveConfig()

							return nil
						},
					},
					{
						Name:        "set-user.js-url",
						Description: "Download the user.js from this url",
						Action: func(ctx *cli.Context) error {
							u := ctx.Args().First()
							if u == "" {
								return errors.New("provide a url")
							}

							parsedUrl, err := url.Parse(u)
							if err != nil {
								return err
							}

							config.UserJsUrl = parsedUrl.String()

							return nil
						},
					},
					{
						Name:        "set-chrome-zip",
						Description: "Download the zip which has the latest chrome dir from this url",
						Action: func(ctx *cli.Context) error {
							u := ctx.Args().First()
							if u == "" {
								return errors.New("provide a url")
							}

							parsedUrl, err := url.Parse(u)
							if err != nil {
								return err
							}

							config.ZippedChromeUrl = parsedUrl.String()

							return nil
						},
					},
				},
			},
			{
				Name:  "update",
				Usage: "Kills firefox and updates it with the latest downloaded files",
				Action: func(ctx *cli.Context) error {
					log.Println("killing firefox")

					if err := utils.KillProcess("firefox"); err != nil {
						log.Println(err)
						log.Println("proceeding")
					}

					log.Println("waiting 5 seconds")
					time.Sleep(time.Second * 5)

					wg := sync.WaitGroup{}
					wg.Add(2)

					// user.js
					go func() {
						defer wg.Done()

						log.Println("downloading", config.UserJsUrl)

						userJs, err := utils.DownloadBytes(config.UserJsUrl)
						if err != nil {
							log.Println(err)
							return
						}

						log.Println("writing user.js to", config.ProfileDir)

						if err := os.
							WriteFile(path.Join(config.ProfileDir, "user.js"), userJs, os.ModePerm); err != nil {
							log.Println(err)
						}
					}()

					// chrome zipped
					go func() {
						defer wg.Done()

						log.Println("downloading", config.ZippedChromeUrl)

						zipped, err := utils.DownloadBytes(config.ZippedChromeUrl)
						if err != nil {
							log.Println(err)
							return
						}

						z, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
						if err != nil {
							log.Println(err)
							return
						}

						for _, file := range z.File {
							fileInfo := file.FileInfo()

							if fileInfo.IsDir() {
								os.MkdirAll(file.Name, os.ModePerm)
								continue
							}

							opened, err := file.Open()
							if err != nil {
								log.Println(err)
								return
							}

							defer opened.Close()

							local, err := os.Create(file.Name)
							if err != nil {
								log.Println(err)
								return
							}

							writer := bufio.NewWriter(local)
							if _, err := writer.ReadFrom(opened); err != nil {
								log.Println(err)
							}
						}

					}()

					wg.Wait()

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
