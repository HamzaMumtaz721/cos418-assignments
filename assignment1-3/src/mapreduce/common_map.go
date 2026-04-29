package mapreduce

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
                )

func doMap(jobName string,
	mapTaskNumber int, inFile string,
	nReduce int,
	mapF func(file string, contents string) []KeyValue,
) {
	// You can find the filename for this map task's input to reduce task number
	// r using reduceName(jobName, mapTaskNumber, r). The ihash function (given
	// below doMap) should be used to decide which file a given key belongs into.
	data, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal("doMap: cannot read ", inFile, err)
                }

	kvPairs := mapF(inFile, string(data))

	files := make([]*os.File, nReduce)
	encoders := make([]*json.Encoder, nReduce)
	for i := 0; i < nReduce; i++ {
		files[i], err = os.Create(
			reduceName(jobName, mapTaskNumber, i))
		if err != nil {
			log.Fatal("doMap: create error ", err)
                        }
		encoders[i] = json.NewEncoder(files[i])
                }

	// One format often used for serializing data to a byte stream that the
	// other end can correctly reconstruct is JSON.
	for _, kv := range kvPairs {
		bucket := ihash(kv.Key) % uint32(nReduce)
		err = encoders[bucket].Encode(&kv)
		if err != nil {
			log.Fatal("doMap: encode error ", err)
                                }
                }

	// Remember to close the file after you have written all the values!
	for i := 0; i < nReduce; i++ {
		files[i].Close()
                        }
        }

func ihash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
                        }
