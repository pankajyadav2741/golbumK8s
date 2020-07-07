package model

import (
	"fmt"
	"net/http"
	"github.com/gocql/gocql"
	"github.com/kelseyhightower/envconfig"
	"github.com/pankajyadav2741/golbumK8s/utils"
)

type Conf struct {
	DbHost string `envconfig:"DB_HOST"`
}

var cluster *gocql.ClusterConfig

//Albums : Slice of struct
var albums []Albums

//init : Establish Database Connections
func init() {
	db := &Conf{}
	err := envconfig.Process("", db)
	if err != nil {
		fmt.Println("Error in envconfig")
		log.Fatal(err.Error())
	}
	fmt.Println("DB HOST : %v", db.DbHost)
	cluster = gocql.NewCluster(db.DbHost)
	cluster.Keyspace = "system"
	cluster.Timeout  = time.Second * 20
	cluster.ConnectTimeout  = time.Second * 20
	
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		fmt.Println("Create session failed")
		panic(err)
	}
	fmt.Println("Cassandra init done")

	//Create Keyspace
	if err := Session.Query(`CREATE KEYSPACE IF NOT EXISTS albumspace WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`).Exec(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Keyspace created")
	}
	
	cluster.Keyspace = "albumspace"
	session, err := cluster.CreateSession()
	defer session.Close()
	if err != nil {
		log.Fatal("createSession:", err)
	}
	
	//Create Table albumtable
	if err := session.Query(`CREATE TABLE IF NOT EXISTS albumtable(albname TEXT PRIMARY KEY, imagelist LIST<TEXT>);`).Exec(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Table albumtable created")
	}
}

//AlbumExists : Check if an album exists
func AlbumExists(albName string) (bool, *utils.ApplicationError) {
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		fmt.Println("Create session failed")
		return err
	}
	
	row = Session.Execute("SELECT * FROM albumtable WHERE albname = ? limit 1", albName)
    if len(row.current_rows) == 0 { 
		return false, &utils.ApplicationError{
			Message:    fmt.Sprintf("Image %v not found in album %v", imgName, albName),
			StatusCode: http.StatusNotFound,
		}
	}
	
	return true, &utils.ApplicationError{
		Message:    fmt.Sprintf("Image %v found in album %v", imgName, albName),
		StatusCode: http.StatusConflict,
	}
}

//ImageExists : Check if an image exists in an album
func ImageExists(albName, imgName string) (bool, *utils.ApplicationError) {
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		fmt.Println("Create session failed")
		return err
	}
	
	//r = session.execute(f"SELECT * FROM test_table WHERE id = {some_id} limit 1")
    
	row = Session.Execute("SELECT * FROM albumtable WHERE albname = ? CONTAINS ?", albName, imgName)
    if len(row.current_rows) == 0 { 
		return false, &utils.ApplicationError{
			Message:    fmt.Sprintf("Image %v not found in album %v", imgName, albName),
			StatusCode: http.StatusNotFound,
		}
	}
	
	return true, &utils.ApplicationError{
		Message:    fmt.Sprintf("Image %v found in album %v", imgName, albName),
		StatusCode: http.StatusConflict,
	}
}

//ShowAlbum : Show all albums
func ShowAlbum() ([]Album, *utils.ApplicationError) {
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	
	iter := Session.Query("SELECT albname FROM albumtable;").Iter()
	var data string
	for iter.Scan(&data) {
		return data, nil
		//json.NewEncoder(w).Encode(data)
	}
	if err := iter.Close(); err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("DB ERROR: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
}


//AddAlbum : Create a new album
func AddAlbum(albName string) *utils.ApplicationError {
	if ok, err := AlbumExists(albName); ok != false {
		return err
	}
	
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	if err := Session.Query(`INSERT INTO albumtable (albname) VALUES (?) IF NOT EXISTS;`, albName).Exec(); err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("INSERT album %v operation failed", albName),
			StatusCode: http.StatusInternalServerError,
		}
	} 
	return nil
}

//DeleteAlbum : Delete an existing album
func DeleteAlbum(albName string) *utils.ApplicationError {
	if ok, err := AlbumExists(albName); ok != true {
		return err
	}
	
	Session, err := cluster.CreateSession()
	defer Session.Close()	
	if err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	
	if err := Session.Query(`DELETE FROM albumtable WHERE albname=? IF EXISTS;`, albName).Exec(); err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("DELETE album %v operation failed", albName),
			StatusCode: http.StatusInternalServerError,
		}
	} 	
	return nil
}

//ShowImagesInAlbum : Show all images in an album
func ShowImagesInAlbum(albName string) ([]Image, *utils.ApplicationError) {
	if ok, err := AlbumExists(albName); ok != true {
		return nil, err
	}

	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	
	iter := Session.Query("SELECT imagelist FROM albumtable WHERE albname=?;", albName).Iter()
	var data []string
	for iter.Scan(&data) {
		return data, nil
		//json.NewEncoder(w).Encode(data)
	}
	if err := iter.Close(); err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("DB ERROR: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
}


//ShowImage : Show a particular image inside an album
func ShowImage(albName, imgName string) (*Image, *utils.ApplicationError) {
	if ok, err := AlbumExists(albName); ok != true {
		return nil, err
	}

	if ok, err := ImageExists(albName, imgName); ok != true {
		return nil, err
	}

	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}
	
	iter := Session.Query("SELECT imagelist FROM albumtable WHERE albname='?';", albName).Iter()
	var data []string
	for iter.Scan(&data) {
		for _, img := range data {
			if img == imgName {
				return img, nil
				//json.NewEncoder(w).Encode(img)
			}
		}
	}
	if err := iter.Close(); err != nil {
		return nil, &utils.ApplicationError{
			Message:    fmt.Sprintf("DB ERROR: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	}
}


//AddImage : Create an image in an album
func AddImage(albName, imgName string) *utils.ApplicationError {
	if ok, err := AlbumExists(albName); ok != true {
		return err
	}

	if ok, err := ImageExists(albName, imgName); ok != true {
		return err
	}
	
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := Session.Query(`UPDATE albumtable SET imagelist=imagelist+ ? WHERE albname=?;`, []string{imgName}, albName).Exec(); err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("DB ERROR: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	} 
	return nil
}

//DeleteImage : Delete an image in an album
func DeleteImage(albName, imgName string) *utils.ApplicationError {
	if ok, err := AlbumExists(albName); ok != true {
		return err
	}

	if ok, err := ImageExists(albName, imgName); ok != false {
		return err
	}
	
	Session, err := cluster.CreateSession()
	defer Session.Close()
	if err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("Create session failed"),
			StatusCode: http.StatusInternalServerError,
		}
	}

	if err := Session.Query(`UPDATE albumtable SET imagelist=imagelist-? WHERE albname=?;`, []string{imgName}, albName).Exec(); err != nil {
		return &utils.ApplicationError{
			Message:    fmt.Sprintf("DB ERROR: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
	} 
	return nil
}
