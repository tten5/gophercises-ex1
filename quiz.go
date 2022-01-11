package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Quiz struct {
	content string
	ans     int
}

func main() {
	// parse the arguments from command line
	var csvFileName = flag.String("csv", "problems.csv", "a csv filename in the format of 'question,answer'")
	var limitTime = flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// open csv file and take out quiz info
	fmt.Println("csvFileName is", *csvFileName)
	fmt.Printf("Time limit for each quiz is %v seconds \n", *limitTime)
	csvFile, err := os.Open(*csvFileName)
	if err != nil {
		fmt.Println(err)
	}

	// init quizArr
	var quizArr = make([]Quiz, 0)
	var correctAns = 0

	// read from csv
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, line := range csvLines {
			if i, err := strconv.Atoi(line[1]); err == nil {
				q := Quiz{
					content: line[0],
					ans:     i,
				}
				quizArr = append(quizArr, q)
			} else {
				fmt.Println(err)
			}
		}
	}
	csvFile.Close()

	fmt.Println("Press Enter to start the quiz ")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	over := make(chan int) // keep track of timeover
	endProgram := false

	for i, q := range quizArr {
		go showQuiz(i, q.content, q.ans, &correctAns, over)
		go countTime(*limitTime, over, i+1)

		for {
			signal := <-over
			// fmt.Println("Signal is", signal)
			if signal == 0 {
				// only when signal = 0 which mean it is from showQuiz
				// then we can move to next quiz
				break
			} else if signal == (i + 1) {
				// signal from countTime return the quizNumber that has time over
				// if signal == current quizNum then end the program
				endProgram = true
				// fmt.Println("Time Out")
				break
			}
		}

		if endProgram {
			break
		}
	}

	fmt.Printf("\nYou scored %v out of %v.", correctAns, len(quizArr))
}

func showQuiz(quizNum int, quiz string, ans int, correctAns *int, over chan<- int) {
	fmt.Printf("Problem #%v: %s = ", quizNum+1, quiz)
	var a int
	fmt.Scan(&a)
	// invalid answer is considered incorrect
	if a == ans {
		(*correctAns)++
	}

	over <- 0
}

func countTime(limitTime int, over chan int, quizNum int) {
	// Calling Sleep method
	time.Sleep(time.Duration(limitTime) * time.Second)

	// fmt.Println("Time over! for quizNum", quizNum)
	over <- quizNum
}
