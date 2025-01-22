package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"golang.org/x/net/html"
	"strings"
)
 
func fetchDetails(url string, ch chan string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Error parsing HTML from %s", url)
		return
	}

	var title, metaDesc string
	var h1Tags []string
	
	var extractData func(*html.Node)
	extractData = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			} else if n.Data == "meta" {
				for _, attr := range n.Attr {
					if attr.Key == "name" && attr.Val == "description" {
						for _, attr := range n.Attr { 
							if attr.Key == "content" {
								metaDesc = attr.Val
							}
						}
					}
				}
			} else if n.Data == "h1" && n.FirstChild != nil {
				h1Tags = append(h1Tags, n.FirstChild.Data)
			}
		}


		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractData(c)
		}
	}

	extractData(doc)


	if title == "" {
		title = "No title found"
	}
	if metaDesc == "" {
		metaDesc = "No meta description found"
	}
	if len(h1Tags) == 0 {
		h1Tags = append(h1Tags, "No H1 headings found")
	}

	result := fmt.Sprintf("\nURL: %s\nTitle: %s\nMeta Description: %s\nH1 Headings: %s\n", url, title, metaDesc, strings.Join(h1Tags, ", "))
	ch <- result
}

func main() {
	fmt.Print("Enter a URL: ")
	reader := bufio.NewReader(os.Stdin)
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url) 

	ch := make(chan string)

	go fetchDetails(url, ch)

	fmt.Println(<-ch)
}


