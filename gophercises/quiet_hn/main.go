package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hhow09/gophercises/gophercises/quiet_hn/cache"
	"github.com/hhow09/gophercises/gophercises/quiet_hn/hn"
)

func debug(s string) {
	debug := os.Getenv("DEBUG")
	if debug == "1" {
		fmt.Println(s)
	}
}

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Println("server listening to port:", port)
	fmt.Printf("visit: http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	memcache := cache.NewInMemoryCache()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if v := memcache.Get(cache.TOP_STORIES); v != nil {
			debug("cache hit")
			renderTemplate(w, tpl, v.([]item), start)
			return
		}
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		resChan := make(chan chanItem)
		var stories []chanItem
		for idx, id := range ids {
			go func(id int, idx int) {
				hnItem, err := client.GetItem(id)
				if err != nil {
					resChan <- chanItem{item: *new(item), index: idx, err: err}
					return
				}
				item := parseHNItem(hnItem)

				resChan <- chanItem{item: item, index: idx, err: nil}
			}(id, idx)
		}
		for i := 0; i < len(ids); i++ {
			stories = append(stories, <-resChan)
		}

		sort.Slice(stories, func(i, j int) bool {
			return stories[i].index < stories[j].index
		})

		sortedStories := []item{}
		for _, s := range stories {
			if s.err != nil {
				continue
			}
			if !isStoryLink(s.item) {
				continue
			}
			sortedStories = append(sortedStories, s.item)
			if len(sortedStories) == numStories {
				break
			}
		}
		memcache.Set(cache.TOP_STORIES, sortedStories)
		renderTemplate(w, tpl, sortedStories, start)
	})
}

func renderTemplate(w http.ResponseWriter, tpl *template.Template, stories []item, start time.Time) {
	data := templateData{
		Stories: stories,
		Time:    time.Since(start),
	}
	err := tpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to process the template", http.StatusInternalServerError)
		return
	}
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type chanItem struct {
	item  item
	index int
	err   error
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
