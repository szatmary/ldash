package handlers

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"
        "strings"
	"fmt"
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
	
        writer := shapeio.NewWriter(w)

	if strings.HasSuffix(curFilePath, ".mpd") { // we don't do chunkced download for manifest
		utils.GetDownloadLogger().Debugf("download manifest file \n")
		io.Copy(writer, file)
		return
	}

	bufferSize := 20480

	buffer := make([]byte, bufferSize)
        var read_err error
        bytesread := 0
	for {
		for {
			bytesread, read_err = file.Read(buffer)
			if read_err != nil {
                                 utils.GetDownloadLogger().Errorf("check: %v \n", read_err)
				if read_err != io.EOF { // print out if read error
					utils.GetDownloadLogger().Errorf("Failed to read file: %v \n", err)
				}
			}

                        utils.GetDownloadLogger().Debugf("read %d bytes \n", bytesread)

			if bytesread > 0 {
				prefix := fmt.Sprintf("%d\r\n", bytesread)
				init := []byte(prefix)
				end := []byte("\r\n")
				new_data := append(init, buffer...)
				msg := append(new_data, end...)
				_, errpr := writer.Write(msg)
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
			if read_err == io.EOF { // break out if reach to the file end
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}

	// write the end of chunk
        msg := []byte("0\r\n\r\n")
        _, errpr := writer.Write(msg)
        if errpr != nil {
          panic(errpr)
        }
}
