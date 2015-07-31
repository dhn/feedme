// fshow.go - show RSS Content
package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/dhn/feedme/lib/config"
	"github.com/dhn/feedme/thirdparty/getopt"
	"github.com/dhn/feedme/thirdparty/sqlite3"
)

// A Args set getopt arguments
type Args struct {
	read   int
	editor string
}

var cursor *sqlite3.Conn
var args = Args{0, os.Getenv("EDITOR")}

func die(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkSystemEnv(env string) bool {
	if len(os.Getenv(env)) == 0 {
		return false
	}

	return true
}

func checkFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

func initSQL() {
	if !checkFileExist(config.SQLDatabase) {
		die("%s not found, please run first ffetch\n", config.SQLDatabase)
	} else {
		cursor, _ = sqlite3.Open(config.SQLDatabase)
	}
}

func getTitle(id int) string {
	var row sqlite3.RowMap
	var result *sqlite3.Stmt
	var err error

	// If id is empty (0) show the last row entry.
	if id == 0 {
		query := "SELECT title from feed ORDER BY date DESC LIMIT 1;"
		row = make(sqlite3.RowMap)

		result, err = cursor.Query(query)
	} else {
		query := "SELECT title from feed WHERE id = $id;"
		sql := sqlite3.NamedArgs{"$id": id}
		row = make(sqlite3.RowMap)

		result, err = cursor.Query(query, sql)
		if err != nil {
			die("Id: %d not found.\n", id)
		}
	}

	checkErr(err)
	result.Scan(row)

	return row["title"].(string)
}

func hash(str string) string {
	hash := sha1.New()
	hash.Write([]byte(str))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func openArticle(article string) {
	cmd := exec.Command(args.editor, article)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("usage: fshow [options]")
	fmt.Println("   -r,  Show article n in your $EDITOR (default the newest).")
	fmt.Println("   -h,  Show this help and exit.")
}

func main() {
	if checkSystemEnv("EDITOR") {
		var c int

		getopt.OptErr = 0
		for {
			if c = getopt.Getopt("r:h"); c == getopt.EOF {
				break
			}
			switch c {
			case 'r':
				args.read, _ = strconv.Atoi(getopt.OptArg)
			case 'h':
				usage()
				os.Exit(0)
			}
		}

		if checkFileExist(config.FilePath) {
			initSQL()
			article := getTitle(args.read)
			file := config.FilePath + hash(article)

			if checkFileExist(file) {
				openArticle(file)
			}
		} else {
			die("%s not found\n", config.FilePath)
		}
	} else {
		die("Your $EDITOR value is not set\n")
	}
}

/* vim: set noet sw=4 sts=4: */
