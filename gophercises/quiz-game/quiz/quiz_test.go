package quiz

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReadProblemsFromCSV(t *testing.T) {
	str := "1+1,2\n2+1,3\n9+9,18\n"
	problems, err := readProblemsFromCSV(strings.NewReader(str))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, problems[0].Question, "1+1")
	assert.Equal(t, problems[0].Answer, "2")
	assert.Equal(t, problems[1].Question, "2+1")
	assert.Equal(t, problems[1].Answer, "3")
	assert.Equal(t, problems[2].Question, "9+9")
	assert.Equal(t, problems[2].Answer, "18")
}

func TestEachQuestion(t *testing.T) {
	done := make(chan bool)
	answered := make(chan bool)
	input := make(chan string)
	problems := make([]Problem, 0)
	problems = append(problems, Problem{Question: "1+1", Answer: "2"})
	problems = append(problems, Problem{Question: "2+1", Answer: "3"})
	problems = append(problems, Problem{Question: "9+9", Answer: "18"})
	scores := make([]int, 0)
	go func() {
		for i, problem := range problems {
			score := eachQuestion(i, problem, done, 5, answered, input)
			scores = append(scores, score)
			if i == 2 {
				assert.Equal(t, score, 0) //wrong answer
			} else {
				assert.Equal(t, score, 1)
			}
		}
		done <- true
	}()
	input <- "2"
	input <- "3"
	input <- "15"
}
