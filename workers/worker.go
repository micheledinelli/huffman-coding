package workers

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

func StartWorker(lines <-chan string) <-chan map[string]int {
	finished := make(chan map[string]int)
	go func() {
		defer close(finished)
		for line := range lines {
			tokens := strings.Split(line, "")
			tokens = append(tokens, "\n")
			dic := make(map[string]int)
			for _, token := range tokens {
				dic[token]++
			}

			finished <- dic
		}
	}()
	return finished
}

func Merge(cs ...<-chan map[string]int) <-chan map[string]int {
	var wg sync.WaitGroup
	out := make(chan map[string]int)

	output := func(c <-chan map[string]int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func ProcessFile(file *os.File) <-chan map[string]int {
	lines := make(chan string)

	wc1 := StartWorker(lines)
	wc2 := StartWorker(lines)
	wc3 := StartWorker(lines)
	wc4 := StartWorker(lines)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	go func() {
		defer close(lines)
		for scanner.Scan() {
			lines <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return Merge(wc1, wc2, wc3, wc4)
}
