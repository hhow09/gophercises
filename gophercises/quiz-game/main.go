package main

import (
	"flag"
	"fmt"

	"github.com/hhow09/gophercises/gophercises/quiz-game/quiz"
)

func main() {
	csvPtr := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limitPtr := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	quiz := quiz.GetQuiz(*csvPtr, *limitPtr)
	quiz.Play()

	fmt.Printf("Correct answers: %d out of total %d Questions \n", quiz.Score, len(quiz.Problems))
}
