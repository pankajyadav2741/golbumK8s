package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/pankajyadav2741/golbumK8s/utils"

	"github.com/gorilla/mux"
	"github.com/pankajyadav2741/golbumK8s/service"
)

//checkType : Check if parameter passed is string or not
func checkType(param string) bool {
	paramType := reflect.TypeOf(param).Kind()
	if paramType == reflect.String {
		return true
	}
	return false
}

//checkTypeAlbum : Check if album name passed is string or not
func checkTypeAlbum(albName string) (bool, *utils.ApplicationError) {
	if ok := checkType(albName); ok != true {
		return false, &utils.ApplicationError{
			Message:    fmt.Sprintf("Album name %v is not string", albName),
			StatusCode: http.StatusBadRequest,
		}
	}
	return true, nil
}

//checkTypeImage : Check if image name passed is string or not
func checkTypeImage(imgName string) (bool, *utils.ApplicationError) {
	if ok := checkType(imgName); ok != true {
		return false, &utils.ApplicationError{
			Message:    fmt.Sprintf("Image name %v is not string", imgName),
			StatusCode: http.StatusBadRequest,
		}
	}
	return true, nil
}

//ShowAlbum : Show all albums
func ShowAlbum(w http.ResponseWriter, r *http.Request) {
	alb, albErr := service.ShowAlbum()
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	jsonValue, _ := json.Marshal(alb)
	w.Write([]byte("Albums are: \n"))
	w.Write(jsonValue)
}

//AddAlbum : Create a new album
func AddAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	albErr := service.AddAlbum(param["album"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	w.Write([]byte(fmt.Sprintf("Album %v has been added\n", param["album"])))
}

//DeleteAlbum : Delete an existing album
func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	albErr := service.DeleteAlbum(param["album"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	w.Write([]byte(fmt.Sprintf("Album %v has been deleted\n", param["album"])))
}

//ShowImagesInAlbum : Show all images in an album
func ShowImagesInAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	img, albErr := service.ShowImagesInAlbum(param["album"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	jsonValue, _ := json.Marshal(img)
	w.Write([]byte(fmt.Sprintf("Images in album %v are: \n", param["album"])))
	w.Write(jsonValue)
}

//ShowImage : Show a particular image inside an album
func ShowImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	if ok, err := checkTypeImage(param["image"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	img, albErr := service.ShowImage(param["album"], param["image"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	jsonValue, _ := json.Marshal(img)
	w.Write([]byte(fmt.Sprintf("Image %v found in album %v : \n", param["image"], param["album"])))
	w.Write(jsonValue)
}

//AddImage : Add an image to an album
func AddImage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Add Image")
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	if ok, err := checkTypeImage(param["image"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	albErr := service.AddImage(param["album"], param["image"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	w.Write([]byte(fmt.Sprintf("Image %v has been added to album %v\n", param["image"], param["album"])))
}

//DeleteImage : Delete an image from an album
func DeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := mux.Vars(r)
	if ok, err := checkTypeAlbum(param["album"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	if ok, err := checkTypeImage(param["image"]); ok != true {
		//Bad Request
		jsonValue, _ := json.Marshal(err)
		w.WriteHeader(err.StatusCode)
		w.Write(jsonValue)
		return
	}
	albErr := service.DeleteImage(param["album"], param["image"])
	if albErr != nil {
		//return error
		jsonValue, _ := json.Marshal(albErr)
		w.WriteHeader(albErr.StatusCode)
		w.Write([]byte(jsonValue))
		return
	}
	w.Write([]byte(fmt.Sprintf("Image %v has been deleted from album %v\n", param["image"], param["album"])))
}
