// ffetch.go - fetch and store content to SQL
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"time"

	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/SlyMarbo/rss"
)

/* RSS feed's */
var RSS = map[string]string{
	"miau":    "http://www.lets-hack.it/feed/",
	"marmaro": "http://marmaro.de/lue/feed.rss",
	"kuchen":  "https://kuchen.io/feed",
	"g0tmi1k": "https://blog.g0tmi1k.com/atom.xml",
}

// Default SQL Database
const (
	SQLDatabase = "feedme.db"
	FilePath    = "/tmp/feedme/"
)

var cursor *sqlite3.Conn

func die(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkIfExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func initSQL() {
	cursor, _ = sqlite3.Open(SQLDatabase)

	query := "CREATE TABLE IF NOT EXISTS feed(" +
		"id        INTEGER PRIMARY KEY," +
		"site      TEXT," +
		"title     TEXT UNIQUE," +
		"link      TEXT," +
		"date      INTEGER," +
		"read      INTEGER);"

	cursor.Exec(query)
}

func insertSQL(site string, title string, link string, date time.Time, read bool) {
	query := "INSERT INTO feed (site, title, link, date, read) " +
		"VALUES ($site, $title, $link, $date, $read);"

	sql := sqlite3.NamedArgs{
		"$site":  site,
		"$title": title,
		"$link":  link,
		"$date":  date,
		"$read":  read,
	}

	cursor.Exec(query, sql)
}

func writeToFile(filename string, content string) {
	if checkIfExist(FilePath) {
		file, err := os.Create(filename)
		checkErr(err)

		n, err := io.WriteString(file, content)
		if err != nil {
			die("%s", n, err)
		}

		file.Close()
	} else {
		die("%s not found\n", FilePath)
	}
}

func hash(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func main() {
	// initialize SQL database
	initSQL()

	for _, url := range RSS {
		feed, err := rss.Fetch(url)
		checkErr(err)

		err = feed.Update()
		checkErr(err)

		for _, element := range feed.Items {
			writeToFile(FilePath+hash(element.Title), element.Content)
			insertSQL(feed.Title, element.Title, element.Link,
				element.Date, element.Read)
		}
	}
}

/* vim: set noet sw=4 sts=4: */
