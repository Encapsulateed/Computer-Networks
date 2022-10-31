package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mmcdole/gofeed"
)

const (
	User  = "iu9networkslabs"
	Pass  = "Je2dTYr6"
	Host  = "students.yss.su"
	Port  = "3306"
	DB    = "iu9networkslabs"
	Table = "IU9Mitroshkin"

	URL = "http://static.feed.rbc.ru/rbc/logical/footer/news.rss"
	//http://static.feed.rbc.ru/rbc/logical/footer/news.rss
)

func main() {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", User, Pass, Host, Port, DB))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(URL)
	if err != nil {
		panic(err)
	}

	//fmt.Println(feed.Items[0].Title)
	//fmt.Println(feed.Items[0].Description)
	//fmt.Println(feed.Items[0].Link)
	//fmt.Println(feed.Items[0].Authors[0].Name)

	for _, item := range feed.Items {

		var autor = ""

		if (len(item.Authors) != 0) {
			autor = item.Authors[0].Name
		} else if (len(item.Authors) == 0) {
			autor = "none"

		}
		if _, err := db.Exec(
			fmt.Sprintf("INSERT IGNORE INTO `%s` (title, description, link, autor) VALUES ( ?, ?, ?, ?)", Table),
			item.Title,
			item.Description,
			item.Link,
			autor,
		); err != nil {
			panic(err)
		}
	}
}
