# feedme
A simple RSS reader written in Go

[![Build Status](https://travis-ci.org/dhn/feedme.svg?branch=master)](https://travis-ci.org/dhn/feedme)

## Dependencies

- github.com/SlyMarbo/rs
- github.com/axgle/mahonia
- github.com/mattn/go-getopt
- github.com/jehiah/go-strftime
- code.google.com/p/go-sqlite/go1/sqlite3

## Programms

> Write programs that do one thing and do it well. [Doug McIlroy]

### ffetch

Fetch content from RSS feed and store site, title, link and date (publish date) to a SQLite3 Database.

### fscan

Show the last (default 10) new RSS feeds.

- ```-l```:  Show n new RSS feeds (default 10).

```
$ fscan -l 5
  21* 2015-06-24 23:22:20  miau's blog            Tell github to use another ssh key
  22  2015-06-23 13:37:25  miau's blog            Rename file with git-mv
  31  2015-06-09 00:00:00  http://marmaro.de/lue  2015-06-09
  32  2015-06-08 00:00:00  http://marmaro.de/lue  2015-06-08
  33  2015-06-01 00:00:00  http://marmaro.de/lue  2015-06-01
```
