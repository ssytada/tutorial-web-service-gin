package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

func getAlbums(c *gin.Context) {
	var albums []album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	for rows.Next() {
		var alb album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message: getAlbums:": err})
			return
		}
		albums = append(albums, alb)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message: postAlbums:": err})
		return
	}
	alb, err := addAlbum(newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message: postAlbums:": err})
		return
	}
	newAlbum.ID = alb
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func addAlbum(alb album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?,?,?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func getAlbumByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	alb, err := albumByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, alb)

}

func albumByID(id int64) (album, error) {
	var alb album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "docker",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "album",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connected!")

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.Run("localhost:8080")
}
