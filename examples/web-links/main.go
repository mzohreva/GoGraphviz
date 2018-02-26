package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/mzohreva/GoGraphviz/graphviz"
	"golang.org/x/net/html"
)

func main() {
	startPage := "https://github.com/topics/go"
	fetchLimit := 200
	pattern := regexp.MustCompile("^https://github\\.com/[a-zA-Z0-9]+/[a-zA-Z0-9]+$")
	shouldFollow := func(link string) bool {
		return !strings.HasPrefix(link, "https://github.com/topics/") &&
			!strings.HasPrefix(link, "https://github.com/trending/") &&
			link != "https://github.com/trending" &&
			link != "https://github.com/site/privacy" &&
			pattern.MatchString(link)
	}
	w := crawl(startPage, fetchLimit, shouldFollow)
	title := fmt.Sprintf("start page: %#v, fetch limit: %v", startPage, fetchLimit)
	visualizeWebGraph(w, "repo-web", title)
}

func visualizeWebGraph(w *webGraph, filename, title string) error {
	G := graphviz.Graph{}
	G.MakeDirected()
	G.DefaultNodeAttribute(graphviz.Shape, graphviz.ShapePoint)
	G.DefaultEdgeAttribute("arrowhead", "vee")
	G.DefaultEdgeAttribute("arrowsize", "0.2")

	nodeID := make(map[int]int)
	for from, toList := range w.links {
		if _, ok := nodeID[from]; !ok {
			nodeID[from] = G.AddNode(w.page[from])
		}
		for _, to := range toList {
			if !w.hasVisited(w.page[to]) {
				continue
			}
			if _, ok := nodeID[to]; !ok {
				nodeID[to] = G.AddNode(w.page[to])
			}
			if from != to {
				G.AddEdge(nodeID[from], nodeID[to], "")
			}
		}
	}
	G.SetTitle("\n\n" + title)
	f, err := os.Create(filename + ".dot")
	if err != nil {
		return err
	}
	G.GenerateDOT(f)
	f.Close()
	log.Println("Saved graph description in", f.Name())
	err = G.GenerateImage("neato", filename+".png", "png")
	log.Println("Saved graph image in", filename+".png")
	return err
}

func crawl(startPage string, fetchLimit int, shouldFollow func(string) bool) *webGraph {
	w := newWebGraph()
	Q := queue{}
	Q.push(startPage)
	for len(w.visited) < fetchLimit && !Q.empty() {
		page := Q.pop()
		if w.hasVisited(page) {
			continue
		}
		w.markVisited(page)
		linkCount, err := fetch(page)
		if err != nil {
			log.Println(err)
			continue
		}
		list := make([]string, 0, len(linkCount))
		for k := range linkCount {
			if shouldFollow(k) {
				list = append(list, k)
			}
		}
		log.Printf("[%4d] Found %4d URLs on %v\n", len(w.visited), len(list), truncateString(page, 60))
		sort.Slice(list, func(i, j int) bool {
			if linkCount[list[i]] != linkCount[list[j]] {
				return linkCount[list[i]] > linkCount[list[j]]
			}
			return list[i] < list[j]
		})
		for i := range list {
			w.createLink(page, list[i])
			if !w.hasVisited(list[i]) {
				Q.push(list[i])
			}
		}
	}
	return w
}

func fetch(page string) (map[string]int, error) {
	baseURL, err := url.Parse(page)
	if err != nil {
		return nil, err
	}
	linkedPages := make(map[string]int)
	var dfs func(*html.Node)
	dfs = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i := range n.Attr {
				if n.Attr[i].Key == "href" {
					u := n.Attr[i].Val
					parsed, err := baseURL.Parse(u)
					if err != nil {
						continue
					}
					parsed.Fragment = ""
					if parsed.Scheme == "http" || parsed.Scheme == "https" {
						linkedPages[parsed.String()]++
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			dfs(c)
		}
	}
	resp, err := http.Get(page)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	dfs(doc)
	return linkedPages, nil
}

type webGraph struct {
	number  map[string]int   // map a page URL to a unique id
	page    map[int]string   // map ids to page URLs
	visited map[int]struct{} // set of visited page ids
	links   map[int][]int    // links between pages
}

func newWebGraph() *webGraph {
	number := make(map[string]int)
	page := make(map[int]string)
	visited := make(map[int]struct{})
	links := make(map[int][]int)
	return &webGraph{number, page, visited, links}
}

func (w *webGraph) id(page string) int {
	num, inTheMap := w.number[page]
	if !inTheMap {
		num = len(w.number)
		w.number[page] = num
		w.page[num] = page
	}
	return num
}

func (w *webGraph) hasVisited(page string) bool {
	num, inTheMap := w.number[page]
	if !inTheMap {
		return false
	}
	_, inTheMap = w.visited[num]
	return inTheMap
}

func (w *webGraph) markVisited(page string) {
	w.visited[w.id(page)] = struct{}{}
}

func (w *webGraph) createLink(fromPage, toPage string) {
	fromNum, toNum := w.id(fromPage), w.id(toPage)
	w.links[fromNum] = append(w.links[fromNum], toNum)
}

type queueNode struct {
	value string
	next  *queueNode
}

type queue struct {
	head, tail *queueNode
	size       int
}

func (Q *queue) push(value string) {
	node := &queueNode{value, nil}
	Q.size++
	if Q.head == nil {
		Q.head, Q.tail = node, node
		return
	}
	Q.tail.next = node
	Q.tail = node
}

func (Q *queue) pop() string {
	value := Q.head.value
	Q.head = Q.head.next
	if Q.head == nil {
		Q.tail = nil
	}
	Q.size--
	return value
}

func (Q *queue) empty() bool {
	return Q.head == nil
}

func truncateString(str string, maxLength int) string {
	if len(str) <= maxLength {
		return str
	}
	return str[0:maxLength-3] + "..."
}
