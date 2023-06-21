package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csv := flag.String("csv", "problems.csv", ".csv file with quiz Q&A")
	secTime := flag.Int("time", 30, "quiz time limit") // time in seconds
	r := flag.Bool("random", false, "randomize questions")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	questions := csvLoader(*csv)
	correct := 0
	total := len(questions) - 1

	if *r {
		questions = random(questions)
	}

	fmt.Println("Total Questions:", total)
	fmt.Println("Duration [s]:", *secTime)

	done := make(chan bool, 1)

	go func() {
		for i := 0; i < total; i++ {
			fmt.Printf("Question #%d %s = ", i+1, questions[i][0])

			answer, _ := reader.ReadString('\n')
			// CRLF (0x0D0A, \r) to LF (0x0A, \n)
			answer = strings.Replace(answer, "\n", "", -1)
			answer = strings.ToLower(answer)
			answer = strings.TrimSpace(answer)

			// compare answer
			if strings.Compare(questions[i][1], answer) == 0 {
				correct++
			}

		}
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("Good Job!")

	case <-time.After(time.Duration(*secTime) * time.Second):
		fmt.Println("\nTime's over!")
	}

	fmt.Println("Score:", correct, "/", total)
}

// randomize questions (q)
func random(q [][]string) [][]string {

	sec := rand.NewSource(time.Now().UnixNano())
	r := rand.New(sec)

	for i := range q {
		np := r.Intn(len(q) - 1)
		q[i], q[np] = q[np], q[i]
	}

	return q
}

// load CSV file content
func csvLoader(file string) [][]string {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal("error: ", err)
	}

	r := csv.NewReader(bytes.NewReader(content)) // content being slice/array

	data, err := r.ReadAll()

	if err != nil {
		log.Fatal("error: ", err)
	}

	return data[1:len(data)] // iteration
}
