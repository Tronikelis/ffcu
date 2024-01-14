package update

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Tronikelis/ffcu/ffcu"
	"github.com/Tronikelis/ffcu/utils"
	"github.com/urfave/cli/v2"
)

func insertChrome(fConfig *ffcu.Config, log *log.Logger) {
	log.Println("downloading", fConfig.ZippedChromeUrl)

	zipped, err := utils.DownloadBytes(fConfig.ZippedChromeUrl)
	if err != nil {
		log.Println(err)
		return
	}

	z, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		log.Println(err)
		return
	}

	chromeDir := path.Join(fConfig.ProfileDir, "chrome")

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
}

func insertUserJs(fConfig *ffcu.Config, log *log.Logger) {
	log.Println("downloading", fConfig.UserJsUrl)

	userJs, err := utils.DownloadBytes(fConfig.UserJsUrl)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("writing user.js to", fConfig.ProfileDir)

	if err := os.
		WriteFile(path.Join(fConfig.ProfileDir, "user.js"), userJs, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func action(fConfig *ffcu.Config) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if !fConfig.IsFilledOut() {
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

		if fConfig.ZippedChromeUrl != "" {
			wg.Add(1)

			go func() {
				insertChrome(fConfig, log.New(os.Stdout, "[insertChrome] ", log.Flags()))
				wg.Done()
			}()
		}

		if fConfig.UserJsUrl != "" {
			wg.Add(1)

			go func() {
				insertUserJs(fConfig, log.New(os.Stdout, "[insertUserJs] ", log.Flags()))
				wg.Done()
			}()
		}

		wg.Wait()

		log.Println("Done, open firefox to enjoy your new theme 🥳")

		return nil
	}

}

func Command(fConfig *ffcu.Config) *cli.Command {
	return &cli.Command{
		Name:   "update",
		Usage:  "Kills firefox and updates it with the latest downloaded files",
		Action: action(fConfig),
	}
}
