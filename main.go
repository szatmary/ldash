package main

import (
	"flag"
	"fmt"
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
	if len(args) != 1 {
		fmt.Println("Usage: need base dir")
		return
	}

	fmt.Println("baseDir", args[0])

	if args[0] == "" {
		fmt.Println("Usage: need base dir")
		return
	}

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
	
	fmt.Println("start server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
