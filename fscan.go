// fscan.go - display new rss feeds
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dhn/feedme/lib/config"
	"github.com/dhn/feedme/thirdparty/getopt"
	"github.com/dhn/feedme/thirdparty/sqlite3"
	"github.com/dhn/feedme/thirdparty/strftime"
)

// A Args set getopt arguments
type Args struct {
	limit int
	all   bool
}

var cursor *sqlite3.Conn

func checkFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func initSQL() {
	if !checkFileExist(config.SQLDatabase) {
		fmt.Printf("%s not found, please run first ffetch\n", config.SQLDatabase)
	} else {
		cursor, _ = sqlite3.Open(config.SQLDatabase)
	}
}

func setSQLQuery(all bool) string {
	var query string

	if all {
		query = "SELECT * from feed ORDER BY date " +
			"DESC LIMIT $limit;"
	} else {
		query = "SELECT * from feed WHERE read = 0 " +
			"ORDER BY date DESC LIMIT $limit;"
	}

	return query
}

func showSQL(limit int, all bool) {
	var query = setSQLQuery(all)

	sql := sqlite3.NamedArgs{"$limit": limit}
	row := make(sqlite3.RowMap)

	for s, err := cursor.Query(query, sql); err == nil; err = s.Next() {
		var symb string
		var rowid int64
		s.Scan(&rowid, row)

		date := strftime.Format("%Y-%m-%d %H:%M:%S", time.Unix(row["date"].(int64), 0))
		title := row["title"]
		author := row["site"]
		read := row["read"].(int64)

		// Use * for read
		if read == 0 {
			symb = " "
		} else {
			symb = "*"
		}

		fmt.Printf(" %3.d%s %s  %-21s  %s\n", rowid, symb, date, author, title)
	}
}

func usage() {
	fmt.Println("usage: fscan [options]")
	fmt.Println("   -l,  Show n new RSS feeds (default 10).")
	fmt.Println("   -a,  Show all RSS feeds (default only unreaded).")
	fmt.Println("   -h,  Show this help and exit.")
}

func main() {
	args := Args{10, false}
	var c int

	getopt.OptErr = 0
	for {
		if c = getopt.Getopt("l:ah"); c == getopt.EOF {
			break
		}
		switch c {
		case 'a':
			args.all = true
		case 'l':
			args.limit, _ = strconv.Atoi(getopt.OptArg)
		case 'h':
			usage()
			os.Exit(0)
		}
	}

	initSQL()
	showSQL(args.limit, args.all)
}

/* vim: set noet sw=4 sts=4: */
