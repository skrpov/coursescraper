package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	DependencyTypeNone = 0
	DependencyTypeAll  = 1
	DependencyTypeAny  = 2
)

type Course struct {
	Desc    string
	DepType int
	Edges   []string
}

type CourseGraph = map[string]Course

func makeKeyFromTitle(title string) string {
	toks := strings.Fields(strings.ToLower(title))
	dept := toks[0]
	dept = dept[:len(dept)-2]
	code := toks[1]
	return dept + "_" + code
}

func getDescriptionFromTitle(title string) string {
	toks := strings.Fields(title)
	return strings.Join(toks[3:], " ")
}

func extractDependencies(g CourseGraph, courseDesc string, key string, body string) {
	t := MakeTokenizer(body)
	dep := Course{
		DepType: DependencyTypeNone,
		Desc:    courseDesc,
	}

	seenPrereq := false

	// TODO: There is for sure a better way to do this iterator.
	for {
		tok, res := t.NextToken()
		if !res {
			break
		}

		switch tok.TType {
		case TokenTypeCourseCode:
			if seenPrereq {
				if dep.DepType == DependencyTypeNone {
					dep.DepType = DependencyTypeAll
				}
				dep.Edges = append(dep.Edges, tok.Val)
			}
		case TokenTypePrereq:
			seenPrereq = true
		case TokenTypeCorreq:
			seenPrereq = false
		case TokenTypeAny:
			if seenPrereq {
				dep.DepType = DependencyTypeAny
			}
		case TokenTypeAll:
			if seenPrereq {
				dep.DepType = DependencyTypeAll
			}
		}
	}

	g[key] = dep
}

// TODO: Have the below functions return errors instead of direct exit.

func scrapeIntoCourseGraph(courseGraph CourseGraph, url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".node__content").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").Text()
		body := s.Find("p").Text()
		key := makeKeyFromTitle(title)
		courseDesc := getDescriptionFromTitle(title)
		extractDependencies(courseGraph, courseDesc, key, body)
	})
}

func saveToDotFile(courseGraph CourseGraph, path string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, "digraph courses {"); err != nil {
		log.Fatal(err)
	}

	for k, v := range courseGraph {
		if _, err := fmt.Fprintf(f, "  %s [label=\"%s\n%s\"];\n", k, k, v.Desc); err != nil {
			log.Fatal(err)
		}
	}

	for k, v := range courseGraph {
		for _, e := range v.Edges {
			// NOTE: Invert the edges when saving so that when viewing the viewer gets the courses
			// that allow them to take course x rather than the courses that course x depends on
			var color string = "black"
			if v.DepType == DependencyTypeAny {
				color = "indigo"
			}
			if _, err := fmt.Fprintf(f, "  %s -> %s [color=%s];\n", e, k, color); err != nil {
				log.Fatal(err)
			}
		}
	}
	if _, err := fmt.Fprintln(f, "}"); err != nil {
		log.Fatal(err)
	}
}

// Feature map
// // 1. Have the graph somehow reflect any vs all (color the edges?)
// // 2. Storing the course name along with the course code in the outputted graph.
// 3. highlight courses which are part of a degree requirements, as well as courses which have already been taken.
// 4. Remove edges for courses which have already been taken or mark them as green or something.

func main() {
	var urls = []string{
		"https://okanagan.calendar.ubc.ca/course-descriptions/subject/cosco",
		// "https://okanagan.calendar.ubc.ca/course-descriptions/subject/matho",
		// "https://okanagan.calendar.ubc.ca/course-descriptions/subject/englo",
		// "https://okanagan.calendar.ubc.ca/course-descriptions/subject/physo",
	}

	courseGraph := make(map[string]Course)
	for _, url := range urls {
		scrapeIntoCourseGraph(courseGraph, url)
	}
	saveToDotFile(courseGraph, "save.dot")
}
