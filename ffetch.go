/* See LICENSE file for license details. */
package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/SlyMarbo/rss"
	"strings"
	"time"
)

/* RSS feed's */
var RSS map[string]string = map[string]string{
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

func init_sql() {
	cursor, _ = sqlite3.Open("feedme.db")

	sql := "CREATE TABLE IF NOT EXISTS feed(" +
		"id        INTEGER PRIMARY KEY," +
		"site      VARCHAR(50)," +
		"title     VARCHAR(50) UNIQUE," +
		"link      VARCHAR(100)," +
		"date      INTEGER);"

	cursor.Exec(sql)
}

func insert_sql(site string, title string, link string, date time.Time) {
	sql := sqlite3.NamedArgs{"$site": site, "$title": title, "$link": link, "$date": date}
	cursor.Exec("INSERT INTO feed (site, title, link, date) VALUES($site, $title, $link, $date)", sql)
}

func main() {
	// initialize SQL database
	init_sql()

	for _, url := range RSS {
		feed, err := rss.Fetch(url)
		checkErr(err)

		err = feed.Update()
		checkErr(err)

		for _, element := range feed.Items {
			title := strings.Replace(feed.Title, "'", "\\'", -1)
			insert_sql(title, element.Title, element.Link, element.Date)
		}
	}
}

/* vim: set noet sw=4 sts=4: */
