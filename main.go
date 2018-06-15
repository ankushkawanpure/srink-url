package main

import (
	"log"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/speps/go-hashids"
	"time"
	"os"
)

//Main Function
func main() {
	dbPath := "./sample.db"
	setDomainEnv()
	initializeSqllite(dbPath)
	startServer()
}

// Function to set Domain Environment element
func setDomainEnv() {
	domainURL := os.Getenv("DOMAIN")
	if domainURL == "" {
		os.Setenv("DOMAIN", "http://localhost:8080")
	}
}

func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/shrink", createHandler). Methods("POST")
	router.HandleFunc("/{shortURL}", expandHandler). Methods("GET")
	router.PathPrefix(	"/").Handler(http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", router))
}

type MyURL struct {
	LongURL  string `json:"longUrl,omitempty"`
	ShortURL string `json:"shortUrl,omitempty"`
}


func initializeSqllite(dbPath string) {

	os.Remove(dbPath)
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement := "CREATE TABLE if not exists `urls` ( `shorturl` TEXT NOT NULL, `longurl` INTEGER NOT NULL, PRIMARY KEY(`shorturl`) )"
	_, err = db.Exec(statement)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database initialize")
}

func addURL (shortURL string, longURL string, dbPath string) string {

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	statement, err := tx.Prepare("INSERT INTO `urls`(`shorturl`,`longurl`) VALUES (?,?)")

	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	result, err := statement.Exec(shortURL, longURL)

	if err != nil {
		return ""
	}

	_, err = result.LastInsertId()
	//log.Printf("done")
	tx.Commit()
	return ""
}

func getShortURL(longURL string, dbPath string) string {

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, err := db.Prepare("select `shorturl` from `urls` where longurl = ?")

	if err != nil {
		log.Fatal(err)
		return ""
	}

	defer statement.Close()

	var shortURL string

	err = statement.QueryRow(longURL).Scan(&shortURL)

	if err != nil {
		return ""
	}

	//log.Printf(shortURL)
	return shortURL
}

func getLongURL(shortURL string, dbPath string) string {
	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	statement, err := db.Prepare("select `longurl` from `urls` where shorturl = ?")

	if err != nil {
		return ""
	}

	defer statement.Close()

	var longURL string

	err = statement.QueryRow(shortURL).Scan(&longURL)

	if err != nil {
		return ""
	}
	//log.Printf(longURL)
	return longURL
}


func createHandler(w http.ResponseWriter, request *http.Request)  {

	dbPath := "./sample.db"

	if (request.Method == "POST") {
		var url MyURL

		_ = json.NewDecoder(request.Body).Decode(&url)


		longURL := url.LongURL
		log.Printf(longURL)

		var shortURL string

		shortURL = getShortURL(longURL, dbPath)


		if shortURL != "" {
			//log.Printf("found shourt url")
			log.Printf("Already generated shorturl" + shortURL)
			// return short url
		} else {
			hd := hashids.NewData()
			hd.Salt = "Admiration we surrounded possession frequently he"
			hd.MinLength = 6
			h, _ := hashids.NewWithData(hd)
			now := time.Now()

			shortURL, _ = h.Encode([]int{int(now.Unix())})
			addURL(shortURL, longURL, dbPath)
			//log.Printf("short url added")
			log.Printf("Newly generated shorturl" + shortURL)
		}

		resp := map[string]string{"shortURL": os.Getenv("DOMAIN") + "/" + shortURL, "longURL": "" + longURL}
		//resp := map[string]string{"shortURL": "http://localhost:8080/" + shortURL, "longURL": "" + longURL}
		jsonData, _ := json.Marshal(resp)
		w.Write(jsonData)

	}
}


func expandHandler(w http.ResponseWriter, request *http.Request) {

	dbPath := "./sample.db"

	if (request.Method == "GET") {
		shortURL := mux.Vars(request)["shortURL"]
		log.Printf(shortURL);

		longURL := getLongURL(shortURL, dbPath)

		if longURL == "" {
			log.Printf("in here")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)

			http.ServeFile(w, request, "./public/404.html")
			return
		}

		http.Redirect(w, request, longURL, 301)
	}
}


