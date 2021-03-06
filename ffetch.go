// ffetch.go - fetch and store content to SQL
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/dhn/feedme/lib/config"
	"github.com/dhn/feedme/thirdparty/html2text"
	"github.com/dhn/feedme/thirdparty/rss"
	"github.com/dhn/feedme/thirdparty/sqlite3"
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
	cursor, _ = sqlite3.Open(config.SQLDatabase)

	query := "CREATE TABLE IF NOT EXISTS feed(" +
		"id        INTEGER PRIMARY KEY," +
		"site      TEXT," +
		"title     TEXT UNIQUE," +
		"link      TEXT," +
		"date      INTEGER," +
		"read      INTEGER);"

	cursor.Exec(query)
}

func insertSQL(site, title, link string, date time.Time, read bool) {
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

func convertHTMLToText(html string) string {
	text, err := html2text.FromString(html)
	checkErr(err)

	return text
}

func writeToFile(filename, content string) {
	convertedText := convertHTMLToText(content)
	file, err := os.Create(filename)
	checkErr(err)

	n, err := io.WriteString(file, convertedText)
	if err != nil {
		die("%s", n, err)
	}

	file.Close()
}

func hash(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func main() {
	if checkIfExist(config.FilePath) {
		// initialize SQL database
		initSQL()

		for _, url := range config.RSS {
			feed, err := rss.Fetch(url)
			checkErr(err)

			err = feed.Update()
			checkErr(err)

			for _, element := range feed.Items {
				writeToFile(config.FilePath+hash(element.Title), element.Content)
				insertSQL(feed.Title, element.Title, element.Link,
					element.Date, element.Read)
			}
		}
	} else {
		die("%s not found\n", config.FilePath)
	}
}

/* vim: set noet sw=4 sts=4: */
