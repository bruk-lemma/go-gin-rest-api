package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var Db *sql.DB

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// get Albums respond with the list of all album as json
func getAlbums(c *gin.Context) {
	//c.IndentedJSON(http.StatusOK, albums)
	rows, err := Db.Query("select * from albums")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching albums"})
		return
	}
	defer rows.Close()

	albums := make([]album, 0)
	for rows.Next() {
		a := album{}
		err := rows.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching albums"})
			return
		}
		albums = append(albums, a)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching albums"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": albums})
}

func postAlbums(c *gin.Context) {

	body := album{}
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "Album not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}
	_, err = Db.Exec("Insert into albums (id,title,artist,price) values ($1,$2,$3,$4)", body.ID, body.Title, body.Artist, body.Price)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, "Could not else{create the new album")
	} else {
		c.JSON(http.StatusOK, "album created successfully")
	}

}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	row := Db.QueryRow("select * from albums where id=$1", id)
	a := album{}
	err := row.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error fetching albums : " + err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": a})

}

func updateAlbumByID(c *gin.Context) {

	body := album{}
	id := c.Param("id")
	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(400, "Album not defined")
		return
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.AbortWithStatusJSON(400, "Bad request")
		return
	}

	_, err = Db.Exec("update albums set title=$1,artist=$2,price=$3 where id=$4", body.Title, body.Artist, body.Price, id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, "Could not update the album")
	} else {
		c.IndentedJSON(http.StatusOK, "album updated successfully")
	}

}

func deleteAlbum(c *gin.Context) {
	id := c.Param("id")
	row := Db.QueryRow("select * from albums where id=$1", id)
	//check if album exists frist and then delete
	a := album{}
	err := row.Scan(&a.ID, &a.Title, &a.Artist, &a.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error fetching albums : " + err.Error()})
		return
	}
	_, err = Db.Exec("delete from albums where id=$1", id)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, "Could not delete the album")
	} else {
		c.IndentedJSON(http.StatusOK, "album deleted successfully")
	}

}

func Conectdatabase() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file please check")
	}

	//we need to read the .env file
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Println(host, port, user, pass, dbname)

	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)

	//fmt.Println("psqlSetup", connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Error connecting to database", err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		panic(err)
	}
	Db = db
	fmt.Println("Successfully connected to database")

}

func main() {

	router := gin.Default()
	Conectdatabase()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.PUT("/albums/:id", updateAlbumByID)
	router.DELETE("/albums/:id", deleteAlbum)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8000")

}
