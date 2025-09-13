package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/temoto/robotstxt"
)

const (
	maxDepth  = 2
	siteDir   = "site"
	workers   = 8
	userAgent = "VSiteDownloader"
)

type job struct {
	url   string
	depth int
}

type downloader struct {
	visited   map[string]bool
	visitedMu sync.Mutex
	robots    *robotstxt.RobotsData
}

func main() {
	startURL := "https://gazeta.ru"
	base, err := url.Parse(startURL)
	if err != nil {
		fmt.Printf("Failed to parse start URL: %v\n", err)
		return
	}

	c := &downloader{
		visited: make(map[string]bool),
		robots:  loadRobots(base),
	}

	jobs := make(chan job, 100)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			defer wg.Done()
			for j := range jobs {
				c.downloadPage(j.url, base.String(), j.depth, jobs)
			}
		})
	}

	jobs <- job{url: startURL, depth: 0}

	go func() {
		wg.Wait()
		close(jobs)
	}()

	wg.Wait()
	fmt.Printf("Site download complete\n")
}

func loadRobots(base *url.URL) *robotstxt.RobotsData {
	robotsURL := base.Scheme + "://" + base.Host + "/robots.txt"
	resp, err := http.Get(robotsURL)
	if err != nil {
		fmt.Printf("robots.txt not found: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	data, err := robotstxt.FromResponse(resp)
	if err != nil {
		fmt.Printf("failed to parse robots.txt: %v\n", err)
		return nil
	}

	return data
}

func (c *downloader) downloadPage(currentURL, baseURL string, depth int, jobs chan<- job) {
	if depth > maxDepth {
		return
	}

	if c.robots != nil && !c.robots.TestAgent(currentURL, userAgent) {
		fmt.Printf("Blocked by robots.txt: %s\n", currentURL)
		return
	}

	c.visitedMu.Lock()
	if c.visited[currentURL] {
		c.visitedMu.Unlock()
		return
	}
	c.visited[currentURL] = true
	c.visitedMu.Unlock()

	resp, err := http.Get(currentURL)
	if err != nil {
		fmt.Printf("Request failed %s: %v\n", currentURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP status not OK: %s: %s\n", currentURL, resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read body %s: %v\n", currentURL, err)
		return
	}

	localPath := getLocalPath(currentURL)
	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("Failed to parse HTML %s: %v\n", currentURL, err)
		return
	}

	rewriteLinks(doc, currentURL, baseURL)

	htmlFile, _ := os.Create(localPath)
	html, _ := doc.Html()
	htmlFile.WriteString(html)
	htmlFile.Close()

	fmt.Printf("Downloaded page: %s -> %s\n", currentURL, localPath)

	selectors := map[string]string{
		"img":                    "src",
		"script":                 "src",
		"link[rel='stylesheet']": "href",
		"link[rel='icon']":       "href",
	}

	for sel, attr := range selectors {
		doc.Find(sel).Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr(attr)
			if !exists {
				return
			}
			resURL := resolveURL(src, currentURL)
			go c.downloadResource(resURL)
		})
	}

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		nextURL := resolveURL(href, currentURL)
		if strings.HasPrefix(nextURL, baseURL) {
			jobs <- job{url: nextURL, depth: depth + 1}
		}
	})
}

func (c *downloader) downloadResource(resourceURL string) {
	c.visitedMu.Lock()
	if c.visited[resourceURL] {
		c.visitedMu.Unlock()
		return
	}
	c.visited[resourceURL] = true
	c.visitedMu.Unlock()

	resp, err := http.Get(resourceURL)
	if err != nil {
		fmt.Printf("Failed to download resource: %s\n", resourceURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	localPath := getLocalPath(resourceURL)
	os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
	file, _ := os.Create(localPath)
	io.Copy(file, resp.Body)
	file.Close()

	fmt.Printf("Downloaded resource: %s -> %s\n", resourceURL, localPath)
}

func resolveURL(href, base string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}

	baseU, err := url.Parse(base)
	if err != nil {
		return href
	}

	return baseU.ResolveReference(u).String()
}

func getLocalPath(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return filepath.Join(siteDir, "unknown")
	}

	path := u.Path
	if path == "" || path == "/" {
		path = "/index.html"
	}

	if strings.HasSuffix(path, "/") {
		path += "index.html"
	}

	return filepath.Join(siteDir, path)
}

func rewriteLinks(doc *goquery.Document, pageURL, baseURL string) {
	attrs := []struct {
		selector string
		attr     string
	}{
		{"img", "src"},
		{"script", "src"},
		{"link[rel='stylesheet']", "href"},
		{"link[rel='icon']", "href"},
		{"a", "href"},
	}

	for _, a := range attrs {
		doc.Find(a.selector).Each(func(i int, s *goquery.Selection) {
			val, exists := s.Attr(a.attr)
			if !exists {
				return
			}

			abs := resolveURL(val, pageURL)
			if strings.HasPrefix(abs, baseURL) {
				local := getLocalPath(abs)
				rel, _ := filepath.Rel(siteDir, local)
				s.SetAttr(a.attr, rel)
			}
		})
	}
}
