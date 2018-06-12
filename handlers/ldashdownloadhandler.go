package handlers

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"../utils"
	"github.com/fujiwara/shapeio"
)

// LDASHDownloadHandler handles the lhls fetching
type LDASHDownloadHandler struct {
	StartTime time.Time
	BaseDir   string
}

func (l *LDASHDownloadHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	utils.GetDownloadLogger().Infof("Received download request\n")
	if strings.HasSuffix(req.URL.EscapedPath(), "html") {
		l.servePlayer(w, req)
	} else {
		l.serveDownload(w, req)
	}

}

func (l *LDASHDownloadHandler) getSourcePath(req *http.Request) string {
	return l.BaseDir
}

func (l *LDASHDownloadHandler) serveDownload(w http.ResponseWriter, req *http.Request) {
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
	w.WriteHeader(http.StatusOK)

	writer := shapeio.NewWriter(w)
	io.Copy(writer, file)
	return
}

func (l *LDASHDownloadHandler) servePlayer(w http.ResponseWriter, req *http.Request) {
	//vars := mux.Vars(req)
	//videoid := vars["videoid"]

	html := `
	<html>
		<head>
			<meta charset="utf-8">
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
			<script type="text/javascript" src="https://bitmovin-a.akamaihd.net/bitmovin-player/stable/7/bitmovinplayer.js"></script>
		</head>
	
		<body>
			<div id="my-player"></div>
			<script type="text/javascript">
				var conf = {
					key: "b050df4b-6966-412d-96bf-a6103e9df1d9",
					source: {
					hls: "playlist.m3u8"
					},
				"tweaks": {
					"max_buffer_level": 4,
				}
				};

				var player = bitmovin.player('my-player');
				
				player.setup(conf).then(function(playerInstance) {
						console.log('Successfully created Bitmovin Player instance');
				}, function(reason) {
						console.log('Error while creating Bitmovin Player instance');
				});
			</script>
		</body>
	</html>`
	w.Write([]byte(html))
}
