package tools

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"time"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

func GetCurrentDate(format string) string {
	result := fmt.Sprintf("Current date: %s", time.Now())
	return result
}

func OpenWebAddress(address string) string {

	client := &http.Client{}
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}

	// Set User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}

	links := strings.Builder{}
	content := strings.Builder{}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode && n.Data != "style" && n.Data != "script" {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				content.WriteString(fmt.Sprintf("%s\n", text))
				fmt.Println(text)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Data != "script" && c.Data != "style" {
				f(c)
			}
		}
	}
	f(doc)
	return fmt.Sprintf("Links: %s\nContent: %s", links.String(), content.String())
}
