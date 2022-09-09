package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}
func isP(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Vol                     int
	Ref, Img, Title, StrVol string
}

var CoinList = []Item{}
var VolumeList = []string{}

//func readItem(item *html.Node) *Item {
//	if a := item; isElem(a, "tr") {
//
//		if cs := getChildren(a); len(cs) == 2 && isElem(cs[0], "time") && isText(cs[1]) {
//			return &Item{
//				Ref:   getAttr(a, "href"),
//				Time:  getAttr(cs[0], "title"),
//				Title: cs[1].Data,
//			}
//		}
//	}
//	return nil
//}

func MakeItem(ref string, img string, title string) *Item {
	return &Item{
		Ref:   ref,
		Img:   img,
		Title: title,
	}
}

func search(node *html.Node) []*Item {

	if isDiv(node, "sc-16r8icm-0 escjiH") {

		var ref string
		ref = "https://coinmarketcap.com" + getAttr(node.FirstChild, "href")
		var img string
		img = getAttr(node.FirstChild.FirstChild.FirstChild, "src")

		var title string = getChildren(node.FirstChild.FirstChild)[1].FirstChild.FirstChild.Data

		CoinList = append(CoinList, *MakeItem(ref, img, title))

	} else if isDiv(node, "sc-16r8icm-0 j3nwcd-0 cRcnjD") {

		VolumeList = append(VolumeList, node.FirstChild.FirstChild.FirstChild.Data)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []Item {
	log.Info("sending request to lenta.ru")
	if response, err := http.Get("https://coinmarketcap.com/"); err != nil {
		log.Error("request to lenta.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from lenta.ru", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from lenta.ru", "error", err)
			} else {
				log.Info("HTML from lenta.ru parsed successfully")

				search(doc)

				for i := 0; i < len(CoinList); i++ {
					var vol string = strings.ReplaceAll(strings.ReplaceAll(VolumeList[i], ",", ""), "$", "")
					k, err := strconv.Atoi(vol)
					if err == nil {
						CoinList[i].Vol = k
						CoinList[i].StrVol = VolumeList[i]
					}

				}
				sort.Slice(CoinList, func(i, j int) (less bool) {
					return CoinList[i].Vol > CoinList[j].Vol
				})
				return CoinList
			}
		}
	}
	return nil
}
