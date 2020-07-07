package service

import (
	"github.com/pankajyadav2741/golbumK8s/model"
	"github.com/pankajyadav2741/golbumK8s/utils"
)

//ShowAlbum : Show all albums
func ShowAlbum() ([]string, *utils.ApplicationError) {
	return model.ShowAlbum()
}

//AddAlbum : Create a new album
func AddAlbum(albName string) *utils.ApplicationError {
	return model.AddAlbum(albName)
}

//DeleteAlbum : Delete an existing album
func DeleteAlbum(albName string) *utils.ApplicationError {
	return model.DeleteAlbum(albName)
}

//ShowImagesInAlbum : Show all images in an album
func ShowImagesInAlbum(albName string) ([]string, *utils.ApplicationError) {
	return model.ShowImagesInAlbum(albName)
}

//ShowImage : Show a particular image inside an album
func ShowImage(albName, imgName string) (string, *utils.ApplicationError) {
	return model.ShowImage(albName, imgName)
}

//AddImage : Create an image in an album
func AddImage(albName, imgName string) *utils.ApplicationError {
	return model.AddImage(albName, imgName)
}

//DeleteImage : Delete an image in an album
func DeleteImage(albName, imgName string) *utils.ApplicationError {
	return model.DeleteImage(albName, imgName)
}
