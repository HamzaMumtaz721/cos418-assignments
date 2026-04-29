package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
                )

func doReduce(jobName string,
	reduceTaskNumber int,
	nMap int,
	reduceF func(key string, values []string) string,
) {
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	kvMap := make(map[string][]string)

	for i := 0; i < nMap; i++ {
		fname := reduceName(jobName, i, reduceTaskNumber)
		f, err := os.Open(fname)
		if err != nil {
			log.Fatal("doReduce: open error ", fname, err)
                        }
		// If you chose to use JSON, you can read out multiple decoded values
		// by creating a decoder, and then repeatedly calling .Decode() on it
		// until Decode() returns an error.
		dec := json.NewDecoder(f)
		for {
			var kv KeyValue
			err = dec.Decode(&kv)
			if err != nil {
				break
                                }
			kvMap[kv.Key] = append(kvMap[kv.Key], kv.Value)
                        }
		f.Close()
                }

	keys := make([]string, 0, len(kvMap))
	for k := range kvMap {
		keys = append(keys, k)
                        }
	sort.Strings(keys)

	// We require you to use JSON here because that is what the merger that
	// combines the output from all the reduce tasks expects.
	outFile, err := os.Create(
		mergeName(jobName, reduceTaskNumber))
	if err != nil {
		log.Fatal("doReduce: create error ", err)
                }
	defer outFile.Close()

	enc := json.NewEncoder(outFile)
	for _, k := range keys {
		enc.Encode(KeyValue{k, reduceF(k, kvMap[k])})
                        }
        }
