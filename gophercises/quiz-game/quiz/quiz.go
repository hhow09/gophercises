package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var timer *time.Timer

type Quiz struct {
	Problems []Problem
	Score    int
	Limit    int
}

type Problem struct {
	Question string
	Answer   string
}

func readProblemsFromCSV(f io.Reader) ([]Problem, error) {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	problems := make([]Problem, 0)
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		problem := Problem{Question: record[0], Answer: record[1]}
		problems = append(problems, problem)
	}
	return problems, nil
}

func shuffle(problems []Problem) {
	//shuffle the questions
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

func getInput(input chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\n') // read text until delim `\n` (Enter) occured
		if err != nil {
			log.Fatal(err)
		}
		input <- result
	}
}

func blockWithEnter() {
	fmt.Println("Press Enter to start the quiz")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n') // read text until delim `\n` (Enter) occured
	if err != nil {
		log.Fatal(err)
	}
}

func GetQuiz(filePath string, limit int) *Quiz {
	quiz := Quiz{Score: 0, Limit: limit}
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	problems, err := readProblemsFromCSV(f)
	if err != nil {
		log.Fatalln(err)
	}
	shuffle(problems)
	quiz.Problems = problems
	return &quiz
}

func (q *Quiz) Play() {
	done := make(chan bool)     // channel for done signal
	answered := make(chan bool) // channel for answered signal
	input := make(chan string)  // channel for io input content
	fmt.Println("Start Play Game")
	fmt.Printf("You have %d second(s) for Each question. \n", q.Limit)
	blockWithEnter()
	go getInput(input)
	go func() {
		for i, problem := range q.Problems {
			score := eachQuestion(i, problem, done, q.Limit, answered, input)
			q.Score += score
		}
		done <- true
	}()
	// block with done channel.
	// either timer timeout or all question answered will end the game.
	<-done
}

func eachQuestion(i int, problem Problem, done chan bool, limit int, answered chan bool, input chan string) (score int) {
	go callTimer(done, limit, answered) //start timer
	fmt.Printf("Question %d: %s \n", i, problem.Question)
	userAnswer := <-input //get the input from io
	answered <- true      // notify timer should stop
	score = 0
	if strings.TrimSpace(userAnswer) == problem.Answer {
		score = 1
	}
	return
}

func callTimer(done chan bool, limit int, answered chan bool) {
	timer = time.NewTimer(time.Duration(limit) * time.Second)
	select {
	case <-timer.C:
		// timer timeout, end game
		fmt.Println("time out")
		done <- true
	case <-answered:
		timer.Stop()
		// already answered, cancel timer
	}
}
