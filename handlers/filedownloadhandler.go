package handlers

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"../utils"
)

type FileDownloadHandler struct {
	StartTime time.Time
	BaseDir   string
}

func (l *FileDownloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	utils.GetDownloadLogger().Infof("Received download request\n")
	l.serveDownload(w, req)
}

func (l *FileDownloadHandler) getSourcePath(req *http.Request) string {
	return l.BaseDir
}

func (l *FileDownloadHandler) isFileUploadingDone(file string) bool {
	symlink := file + ".symlink"
	if _, err := os.Stat(symlink); err == nil {
		// exist
		return true
	}
	// not exist
	return false
}

func (l *FileDownloadHandler) serveDownload(w http.ResponseWriter, req *http.Request) {
	curFileURL := req.URL.EscapedPath()[len("/ldash/download"):]
	curFilePath := path.Join(l.getSourcePath(req), curFileURL)
	file, err := os.Open(curFilePath) // For read access.
	if err != nil {
		utils.GetDownloadLogger().Errorf("Failed to open file: %v \n", err)
		http.NotFound(w, req)
		return
	}
	defer file.Close()

	utils.GetDownloadLogger().Debugf("file %s was requested @ %v \n", curFileURL, time.Now().Format(time.RFC3339))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "Keep-Alive")

	if strings.HasSuffix(curFilePath, ".mpd") {
		w.Header().Set("Content-Type", "application/dash+xml")
	} else {
		w.Header().Set("Content-Type", "video/MP4")
	}

	w.WriteHeader(http.StatusOK)

	bufferSize := 20480
	buffer := make([]byte, bufferSize)
	var read_err error
	bytesread := 0
	for {
		for {
			bytesread, read_err = file.Read(buffer)
			if read_err != nil {
				if read_err != io.EOF { // print out if read error
					utils.GetDownloadLogger().Errorf("Failed to read file: %v \n", err)
				}
			}

			utils.GetDownloadLogger().Debugf("read %d bytes \n", bytesread)

			if bytesread > 0 {
				_, errpr := w.Write(buffer[:bytesread])
				if errpr != nil {
					panic(errpr)
				}
			}

			if bytesread != bufferSize {
				utils.GetDownloadLogger().Debugf("read all existing data \n")
				break
			}
		}

		if read_err != nil {
			if read_err == io.EOF && // if file upload is done and read to eof, the chunk download should be done too.
				l.isFileUploadingDone(curFilePath) {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}
