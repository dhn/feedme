package main

import (
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

var cursor *sqlite3.Conn

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initSQL() {
	cursor, _ = sqlite3.Open("feedme.db")

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

func main() {
	// initialize SQL database
	initSQL()

	for _, url := range RSS {
		feed, err := rss.Fetch(url)
		checkErr(err)

		err = feed.Update()
		checkErr(err)

		for _, element := range feed.Items {
			insertSQL(feed.Title, element.Title, element.Link,
				element.Date, element.Read)
		}
	}
}

/* vim: set noet sw=4 sts=4: */
