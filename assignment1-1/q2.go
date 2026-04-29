package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
        )

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// HINT: use for loop over `nums`
	sum := 0
	for n := range nums {
		sum += n
                }
	out <- sum
        }

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	nums, err := readInts(file)
	checkError(err)

	workerChans := make([]chan int, num)
	for i := 0; i < num; i++ {
		workerChans[i] = make(chan int, len(nums))
                }

	out := make(chan int, num)

	for i := 0; i < num; i++ {
		go sumWorker(workerChans[i], out)
                }

	for i, n := range nums {
		workerChans[i%num] <- n
                }

	for i := 0; i < num; i++ {
		close(workerChans[i])
                        }

	total := 0
	for i := 0; i < num; i++ {
		total += <-out
                }
	return total
        }

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
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
