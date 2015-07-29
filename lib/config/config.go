// config.go - config file for RSS feeds etc.
package config

// Default SQL Database
// and default file path.
const (
	SQLDatabase = "feedme.db"
	FilePath    = "/tmp/feedme/"
)

/* RSS feed's */
var RSS = map[string]string{
	"miau":    "http://www.lets-hack.it/feed/",
	"marmaro": "http://marmaro.de/lue/feed.rss",
	"kuchen":  "https://kuchen.io/feed",
	"g0tmi1k": "https://blog.g0tmi1k.com/atom.xml",
}

/* vim: set noet sw=4 sts=4: */
