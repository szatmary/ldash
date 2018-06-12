package handlers

import (
	"net/http"

	"../utils"
	"github.com/gorilla/mux"
)

type DashPlayHandler struct {
}

func (l *DashPlayHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	utils.GetDownloadLogger().Infof("Received play request\n")
	l.servePlayer(w, req)
}

func (l *DashPlayHandler) servePlayer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	videoid := vars["videoid"]

	html := `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="utf-8"/>
			<title>Auto-player instantiation example, single videoElement</title>
			<style>
				body {
				background-color : black;
				margin : 0;
				}
				video {
				left: 50%;
				position: absolute;
				top: 50%;
				transform: translate(-50%, -50%);
				width: 100%;
				max-height: 100%;
				}
			</style>
		</head>
		<body>
			<div>
				<video id="videoPlayer" controls></video>
			</div>
			<script src="https://reference.dashif.org/dash.js/nightly/dist/dash.all$
			<script>
				var player;
				(function(){
					var url = "` + videoid + `.mpd";
					player = dashjs.MediaPlayer().create();
					player.initialize(document.querySelector("#videoPlayer"), url, $
					player.setLowLatencyEnabled(true);
				})();
			</script>
		</body>
	</html>`
	w.Write([]byte(html))
}
