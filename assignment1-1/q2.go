package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

// sumWorker receives numbers from nums channel, sums them, sends total to out
func sumWorker(nums chan int, out chan int) {
	sum := 0
	for n := range nums {
		sum += n
	}
	out <- sum
}

// sum reads integers from fileName and sums using num goroutines
func sum(num int, fileName string) int {
	// Open file
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	// Read all integers from file
	nums, err := readInts(file)
	checkError(err)

	// Create one buffered channel per worker
	workerChans := make([]chan int, num)
	for i := 0; i < num; i++ {
		workerChans[i] = make(chan int, len(nums))
	}

	// Channel to collect partial sums from workers
	out := make(chan int, num)

	// Launch goroutines
	for i := 0; i < num; i++ {
		go sumWorker(workerChans[i], out)
	}

	// Distribute numbers across workers round-robin
	for i, n := range nums {
		workerChans[i%num] <- n
	}

	// Close worker channels so sumWorker for-range loops finish
	for i := 0; i < num; i++ {
		close(workerChans[i])
	}

	// Collect all partial sums and add them up
	total := 0
	for i := 0; i < num; i++ {
		total += <-out
	}
	return total
}

func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
