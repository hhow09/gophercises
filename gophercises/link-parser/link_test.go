package link

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var htmlFiles = []string{"ex1.html", "ex2.html", "ex3.html"}

func getAns(htmlFiles []string) [][]Link {
	ans := make([][]Link, len(htmlFiles))
	for i := 0; i < len(htmlFiles); i++ {
		ans[i] = make([]Link, 0)
	}
	ans[0] = append(ans[0], Link{Href: "/other-page", Text: "A link to another page"})

	ans[1] = append(ans[1], Link{Href: "https://www.twitter.com/joncalhoun", Text: "Check me out on twitter"})
	ans[1] = append(ans[1], Link{Href: "https://github.com/gophercises", Text: "Gophercises is on Github !"})

	ans[2] = append(ans[2], Link{Href: "#", Text: "Login"})
	ans[2] = append(ans[2], Link{Href: "/lost", Text: "Lost? Need help?"})
	ans[2] = append(ans[2], Link{Href: "https://twitter.com/marcusolsson", Text: "@marcusolsson"})

	return ans
}

func TestParser(t *testing.T) {
	ans := getAns(htmlFiles)
	for i, filePath := range htmlFiles {
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		links, err := Parse(f)

		assert.EqualValues(t, links, ans[i])
	}

}
