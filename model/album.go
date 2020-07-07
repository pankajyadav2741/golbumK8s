package model

//Album : Struct
type Album struct {
	Name  string  `json:"album_name"`
	Image []Image `json:"image"`
}

//Image : Struct
type Image struct {
	Name string `json:"image_name"`
}
