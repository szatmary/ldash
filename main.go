package main

import (
	"flag"
	"log"
	"net/http"
	"path"
	"time"

	"./handlers"
	"./utils"
	"github.com/gorilla/mux"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 || args[0] == "" {
		utils.GetMainLogger().Errorf("Usage: need base dir\n")
		return
	}

	utils.GetMainLogger().Infof("baseDir %s \n", args[0])

	// clean the segment folder
	utils.RemoveContents(args[0])

	file_downloadHandler := &handlers.FileDownloadHandler{
		StartTime: time.Now(),
		BaseDir:   path.Dir(args[0]),
	}

	file_uploadHandler := &handlers.FileUploadHandler{
		BaseDir: path.Dir(args[0]),
	}

	dash_playHandler := &handlers.DashPlayHandler{}

	// open a thread to clean expired files
	go utils.CheckExpire(args[0])

	r := mux.NewRouter()
	r.Handle("/ldash/upload/{name:[a-zA-Z0-9/._-]+}", file_uploadHandler)
	r.Handle("/ldash/download/{name:[a-zA-Z0-9/._-]+}", file_downloadHandler)
	r.Handle("/ldash/play/{videoid}.html", dash_playHandler)

	utils.GetMainLogger().Infof("start server\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
