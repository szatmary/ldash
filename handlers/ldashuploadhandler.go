package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/juju/loggo"
)

var upload_logger = loggo.GetLogger("ldash_upload")

// UploadHandler handles for http upload
type LDASHUploadHandler struct {
	BaseDir string
}

func (u *LDASHUploadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("here")
	curFileURL := req.URL.EscapedPath()[len("/ldash/upload"):]
	curFilePath := path.Join(u.BaseDir, curFileURL)
	u.serveHTTPImpl(curFilePath, w, req)
}

func (u *LDASHUploadHandler) serveHTTPImpl(curFilePath string, w http.ResponseWriter, req *http.Request) {

	// rewrite, mostly for manifest
	if _, err := os.Stat(curFilePath); err == nil {
		upload_logger.Tracef("rewrite file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
		data, _ := ioutil.ReadAll(req.Body)
		err = ioutil.WriteFile(curFilePath, data, 0644)
		if err != nil {
			upload_logger.Errorf("fail to create file %v \n", err)
		}
		return
	}

	// create, mostly for segment
	f, rerr := os.Create(curFilePath)
	if rerr != nil {
		upload_logger.Errorf("fail to create file %v \n", rerr)
	}
	upload_logger.Debugf("create file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
	defer f.Close()
	_, rerr = io.Copy(f, req.Body)
	if rerr != nil {
		upload_logger.Errorf("fail to create file %v \n", rerr)
	}
	upload_logger.Debugf("write file %s @ %v \n", curFilePath, time.Now().Format(time.RFC3339))
}
