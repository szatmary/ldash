package handlers

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"../utils"
	"github.com/fujiwara/shapeio"
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

	utils.GetDownloadLogger().Debugf("file %s was downloaded @ %v \n", curFileURL, time.Now().Format(time.RFC3339))

	w.Header().Set("Content-Type", "video/MP2T")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Transfer-Encoding", "chunked")
	//w.Header().Set("Connection", "Keep-Alive")

	w.WriteHeader(http.StatusOK)
	writer := shapeio.NewWriter(w)
	bufferSize := 2048

	buffer := make([]byte, bufferSize)

	for {
		for {
			bytesread, err := file.Read(buffer)

			if bytesread > 0 {
				_, errpr := writer.Write(buffer[:bufferSize])
				if errpr != nil {
					panic(err)
				}
			}

			if bytesread != bufferSize {
				break
			}
		}

		if err != nil {
			if err == io.EOF { // break out if reach to the file end
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}
