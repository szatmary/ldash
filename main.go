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

	utils.GetMainLogger().Infof("baseDir", args[0])

	// clean the segment folder
	utils.RemoveContents(args[0])

	ldash_downloadHandler := &handlers.LDASHDownloadHandler{
		StartTime: time.Now(),
		BaseDir:   path.Dir(args[0]),
	}

	ldash_uploadHandler := &handlers.LDASHUploadHandler{
		BaseDir: path.Dir(args[0]),
	}

	// open a thread to clean expired segments
	//go utils.CheckExpire(args[0])

	r := mux.NewRouter()
	r.Handle("/ldash/upload/{name:[a-zA-Z0-9/._]+}", ldash_uploadHandler)
	r.Handle("/ldash/download/{name:[a-zA-Z0-9/._]+}", ldash_downloadHandler)

	utils.GetMainLogger().Infof("start server\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
