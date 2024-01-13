package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Tronikelis/ffcu/configuration"
	"github.com/Tronikelis/ffcu/utils"

	"github.com/urfave/cli/v2"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.MkdirAll(path.Join(homeDir, ".ffcu"), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	configDir := path.Join(homeDir, ".ffcu/config.json")

	config, err := configuration.ReadConfig(configDir)
	if err != nil {
		log.Println("can't read config, creating a new one")

		config, err = configuration.CreateConfig(configDir)
		if err != nil {
			log.Fatalln(err)
		}
	}

	defer config.SaveConfig(configDir)

	app := &cli.App{
		Name:  "ffcu",
		Usage: "Helper CLI tool to auto update your firefox config",

		Commands: []*cli.Command{
			{
				Name: "config",
				Subcommands: []*cli.Command{
					{
						Name:        "set-profile-dir",
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
					if !config.IsFilledOut() {
						return errors.New("config is not filled out")
					}

					log.Println("killing firefox")

					if err := utils.KillProcess("firefox"); err != nil {
						log.Println(err)
						log.Println("proceeding")
					}

					log.Println("waiting for 1 second")
					time.Sleep(time.Second)

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

						chromeDir := path.Join(config.ProfileDir, "chrome")

						log.Println("removing", chromeDir)

						if err := os.RemoveAll(chromeDir); err != nil {
							log.Println(err)
							return
						}

						log.Println("creating", chromeDir)

						if err := os.Mkdir(chromeDir, os.ModePerm); err != nil {
							log.Println(err)
							return
						}

						for _, file := range z.File {
							fileInfo := file.FileInfo()

							split := strings.Split(file.Name, "chrome/")
							if len(split) == 1 {
								continue
							}

							underChrome := split[len(split)-1]
							joinedPath := path.Join(chromeDir, underChrome)

							if fileInfo.IsDir() {
								if err := os.MkdirAll(joinedPath, os.ModePerm); err != nil {
									log.Println(err)
									return
								}

								continue
							}

							// joinedPath is a file from here on out

							opened, err := file.Open()
							if err != nil {
								log.Println(err)
								return
							}

							defer opened.Close()

							local, err := os.Create(joinedPath)
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

					log.Println("Done, open firefox to enjoy your new theme 🥳")

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
