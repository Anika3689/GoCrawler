package main
import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"net/url"
)

type PageData struct {
	URL string
	Heading string
	FirstParagraph string
	OutgoingLinks []string
	ImageURLs []string
}

func extractPageData(html, pageURL string) PageData {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}
	}

	olinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		olinks = []string{}
	}
	ilinks, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		ilinks = []string{}
	}
		
	return PageData {
		URL : pageURL,
		Heading: getHeadingFromHTML(html),	
		FirstParagraph: getFirstParagraphFromHTML(html),	
		OutgoingLinks: olinks,
		ImageURLs: ilinks,
	}
}

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	nodes := doc.Find("h1").Nodes

	if len(nodes) > 0 {
		return extractText(nodes[0])
	} 
	nodes = doc.Find("h2").Nodes
	if len(nodes) > 0 {
                return extractText(nodes[0])
	}
	return ""
}

func extractText(head *html.Node) string {
	if head == nil {
		return ""
	}
	if head.Type == html.TextNode {
		return head.Data
	}
	
	var sb strings.Builder
	for cur := head.FirstChild; cur != nil; cur = cur.NextSibling {
		sb.WriteString(extractText(cur))	
	}
	return sb.String()
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
                 return ""
	}
	if nodes := doc.Find("main"); len(nodes.Nodes) != 0 {
		p := nodes.Find("p").Nodes
		if len(p) > 0 {
			return extractText(p[0])
		}
	}
	p := doc.Find("p").Nodes
	if len(p) > 0 {
		return extractText(p[0])
	}
	return ""
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {	
	return getFromHTML(htmlBody, baseURL, "a[href]", "href")
}


func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	return getFromHTML(htmlBody, baseURL, "img[src]", "src")
}

func getFromHTML(htmlBody string, baseURL *url.URL, tag string, attrName string) ([]string, error) {	
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))	
	urls := make([]string, 0)
	if err != nil {
		return urls, err
	}
	
	doc.Find(tag).Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr(attrName)
		if exists {
			u, err := url.Parse(href)
			if err == nil {
				urls = append(urls, baseURL.ResolveReference(u).String())
			}
		}
	})
	return urls, nil
}

