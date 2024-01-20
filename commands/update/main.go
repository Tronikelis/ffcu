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

func insertChrome(fConfig *ffcu.Config, log *log.Logger) error {
	log.Println("downloading", fConfig.ZippedChromeUrl)

	zipped, err := utils.DownloadBytes(fConfig.ZippedChromeUrl)
	if err != nil {
		return err
	}

	z, err := zip.NewReader(bytes.NewReader(zipped), int64(len(zipped)))
	if err != nil {
		return err
	}

	chromeDir := path.Join(fConfig.ProfileDir, "chrome")

	log.Println("removing", chromeDir)

	if err := os.RemoveAll(chromeDir); err != nil {
		return err
	}

	log.Println("creating", chromeDir)

	if err := os.Mkdir(chromeDir, os.ModePerm); err != nil {
		return err
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
				return err
			}

			continue
		}

		// joinedPath is a file from here on out

		opened, err := file.Open()
		if err != nil {
			return err
		}

		defer opened.Close()

		local, err := os.Create(joinedPath)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(local)
		if _, err := writer.ReadFrom(opened); err != nil {
			log.Println(err)
		}
	}

	return nil
}

func insertUserJs(fConfig *ffcu.Config, log *log.Logger) error {
	USER_OVERRIDES_JS_NAME := "user.overrides.js"
	USER_JS_NAME := "user.js"

	log.Println("downloading", fConfig.UserJsUrl)

	userJs, err := utils.DownloadBytes(fConfig.UserJsUrl)
	if err != nil {
		return err
	}

	userJsDir := path.Join(fConfig.ProfileDir, USER_JS_NAME)
	userJsOverrideDir := path.Join(fConfig.ProfileDir, USER_OVERRIDES_JS_NAME)

	log.Println("removing", USER_JS_NAME)
	if err := os.RemoveAll(userJsDir); err != nil {
		return err
	}

	log.Println("creating", USER_JS_NAME)
	userJsOpened, err := os.Create(userJsDir)
	if err != nil {
		return err
	}

	defer userJsOpened.Close()

	log.Println("writing", USER_JS_NAME, "to", fConfig.ProfileDir)
	if _, err := userJsOpened.Write(append(userJs, []byte("\n")...)); err != nil {
		return err
	}

	_, err = os.Stat(userJsOverrideDir)
	if os.IsNotExist(err) {
		log.Println(USER_OVERRIDES_JS_NAME, "not found, exiting")
		return nil
	}
	if err != nil {
		return err
	}

	overrideBytes, err := os.ReadFile(userJsOverrideDir)
	if err != nil {
		return err
	}

	log.Println("appending", USER_OVERRIDES_JS_NAME, "to", USER_JS_NAME)
	if _, err := userJsOpened.Write(overrideBytes); err != nil {
		return err
	}

	return nil
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
				defer wg.Done()

				logger := log.New(os.Stdout, "[insertChrome] ", log.Flags())

				if err := insertChrome(fConfig, logger); err != nil {
					logger.Println(err)
				}
			}()
		}

		if fConfig.UserJsUrl != "" {
			wg.Add(1)

			go func() {
				defer wg.Done()

				logger := log.New(os.Stdout, "[insertUserJs] ", log.Flags())

				if err := insertUserJs(fConfig, logger); err != nil {
					logger.Println(err)
				}
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
