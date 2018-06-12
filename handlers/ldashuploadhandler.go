package handlers

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"../utils"
)

// UploadHandler handles for http upload
type LDASHUploadHandler struct {
	BaseDir string
}

func (u *LDASHUploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	utils.GetUploadLogger().Infof("Received upload request\n")
	curFileURL := req.URL.EscapedPath()[len("/ldash/upload"):]
	curFilePath := path.Join(u.BaseDir, curFileURL)
	u.serveHTTPImpl(curFilePath, w, req)
}

func (u *LDASHUploadHandler) serveHTTPImpl(curFilePath string, w http.ResponseWriter, req *http.Request) {

	// rewrite, mostly for manifest
	if _, err := os.Stat(curFilePath); err == nil {
		utils.GetUploadLogger().Debugf("rewrite file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
		data, _ := ioutil.ReadAll(req.Body)
		err = ioutil.WriteFile(curFilePath, data, 0644)
		if err != nil {
			utils.GetUploadLogger().Errorf("fail to create file %v \n", err)
		}
		return
	}

	// create, mostly for segment
	f, rerr := os.Create(curFilePath)
	if rerr != nil {
		utils.GetUploadLogger().Errorf("fail to create file %v \n", rerr)
	}
	utils.GetUploadLogger().Debugf("create file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
	defer f.Close()
	_, rerr = io.Copy(f, req.Body)
	if rerr != nil {
		utils.GetUploadLogger().Errorf("fail to create file %v \n", rerr)
	}
	utils.GetUploadLogger().Debugf("write file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
}
