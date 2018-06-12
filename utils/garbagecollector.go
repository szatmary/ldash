package utils

import (
	"os"
	"path/filepath"
	"time"
)

func deleteFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		GetGCloadLogger().Errorf("%v\n", err)
		return nil
	}
	if info.Mode().IsRegular() {
		if info.ModTime().Before(time.Now().Add(time.Minute * -1)) {
			os.Remove(path)
			GetGCloadLogger().Debugf("Deleted file %s \n", path)
		}
	}
	return nil
}

// RemoveContents remove the all files in the work folder at beginning
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		GetGCloadLogger().Errorf("%v\n", err)
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			GetGCloadLogger().Errorf("%v\n", err)
			return err
		}
	}
	return nil
}

func removeExpiredContents(dir string) {
	GetGCloadLogger().Infof("start delete: %s \n", dir)
	err := filepath.Walk(dir, deleteFile)
	if err != nil {
		GetGCloadLogger().Errorf("%v\n", err)
	}
}

// CheckExpire garbage collect the expired segment and manifest
func CheckExpire(dir string) {
	for {
		removeExpiredContents(dir)
		time.Sleep(5 * 1000 * time.Millisecond)
	}
}
