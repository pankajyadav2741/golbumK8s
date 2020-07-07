package app

import (
	"log"
	"net/http"
	"os/signal"
	"time"
	"context"
	"syscall"
	"github.com/gorilla/mux"
	"github.com/pankajyadav2741/golbumK8s/controller"
)

//StartApp : Start Application
func StartApp() {
	//Initialize Router
	myRouter := mux.NewRouter().StrictSlash(true)

	//Show all albums
	myRouter.HandleFunc("/", controller.showAlbum).Methods(http.MethodGet)
	//Create a new album
	myRouter.HandleFunc("/{album}", controller.addAlbum).Methods(http.MethodPost)
	//Delete an existing album
	myRouter.HandleFunc("/{album}", controller.deleteAlbum).Methods(http.MethodDelete)

	//Show all images in an album
	myRouter.HandleFunc("/{album}", controller.showImagesInAlbum).Methods(http.MethodGet)
	//Show a particular image inside an album
	myRouter.HandleFunc("/{album}/{image}", controller.showImage).Methods(http.MethodGet)
	//Create an image in an album
	myRouter.HandleFunc("/{album}/{image}", controller.addImage).Methods(http.MethodPost)
	//Delete an image in an album
	myRouter.HandleFunc("/{album}/{image}", controller.deleteImage).Methods(http.MethodDelete)
	
	srv := &http.Server{
		Handler:      myRouter,
		Addr:         ":5000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
