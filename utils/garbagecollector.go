package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("garbagecollector")

// RemoveContents remove the all files in the work folder at beginning
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func removeExpiredContents(dir string) {
	logger.Infof("start delete: %s \n", dir)
	d, err := os.Open(dir)
	if err != nil {
		logger.Errorf("dir %v not exists: %v \n", dir, err)
		os.Exit(1)
	}
	defer d.Close()

	videos, err := d.Readdir(-1)
	if err != nil {
		logger.Errorf("cannot read dir:%v \n", err)
		os.Exit(1)
	}

	logger.Infof("scan %s \n" + dir)

	for _, video := range videos {
		videoPath := filepath.Join(dir, video.Name())
		if video.Mode().IsRegular() {
			if video.ModTime().Before(time.Now().Add(time.Minute * -1)) {
				os.Remove(videoPath)
				logger.Debugf("Deleted video %s \n", videoPath)
			}
			continue
		}

		logger.Infof("populate path: %s \n", videoPath)
		d, _ = os.Open(videoPath)
		resolutions, err := d.Readdir(-1)
		if err != nil {
			logger.Errorf("cannot read video dir:%v \n", err)
			os.Exit(1)
		}

		for _, resolution := range resolutions {
			resolutionPath := filepath.Join(videoPath, resolution.Name())
			if resolution.Mode().IsRegular() {
				if resolution.ModTime().Before(time.Now().Add(time.Minute * -1)) {
					os.Remove(resolutionPath)
					logger.Debugf("Deleted resolution %s \n", resolutionPath)
				}
				continue
			}

			logger.Debugf("populate resolutionPath: %s \n", resolutionPath)
			d, _ = os.Open(resolutionPath)
			files, err := d.Readdir(-1)
			if err != nil {
				logger.Errorf("cannot read resolution dir:%v \n", err)
				os.Exit(1)
			}
			for _, file := range files {
				filePath := filepath.Join(resolutionPath, file.Name())
				if file.Mode().IsRegular() {
					if file.ModTime().Before(time.Now().Add(time.Minute * -1)) {
						os.Remove(filePath)
						logger.Debugf("Deleted segment %s \n", filePath)
					}
				}
			}
		}
	}
}

// CheckExpire garbage collect the expired segment and manifest
func CheckExpire(dir string) {
	for {
		removeExpiredContents(dir)
		time.Sleep(5 * 1000 * time.Millisecond)
	}
}
