package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"voting2021/decryptor/internal"
	"voting2021/decryptor/internal/input"
	"voting2021/decryptor/internal/output"
	"voting2021/decryptor/internal/processor"
)

func main() {

	// cmd {token} {path_to_csv}
	privateKey := os.Args[1]
	pathToFile := os.Args[2]

	// Sample: По одномандатному избирательному округу

	//https://observer.mos.ru/all/servers/1/txs
	// Private Key: 54e3cf70f712b2ff727bde3849772fa811a9d5de796aa7d788d205aa86af04ad
	// Published in TX: b7f30f549dbef79e96b744d44c638e72a429ce949839c54990f65d4306a9d13d

	println("File: ", pathToFile)
	println("Key: ", privateKey)
	println("Output: out_aggregated.json")
	println("Output: out_full.json")
	println("Output: out_full.csv")

	processor := processor.NewProcessor(privateKey)
	input.ProcessVotesCsvFile(pathToFile, processor)

	votes := processor.GetVotes()
	result := make([]output.VotingResult, 0)

	for candidateId, transactions := range votes {
		fmt.Printf("Candidate ID: %d, Votes: %d\n", candidateId, len(transactions))
		result = append(result, output.VotingResult{CandidateId: candidateId, VotesCount: len(transactions)})
	}

	writeJsonAggregatedOutput(&result)
	writeJsonFullOutput(&votes)
	writeCsvFullOutput(&votes)
}

func writeJsonAggregatedOutput(result *[]output.VotingResult) {
	fileName := "out_tiny.json"

	file, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)
}

func writeJsonFullOutput(result *map[uint32][]*internal.VoteTransaction) {
	fileName := "out_full.json"

	file, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile(fileName, file, 0644)
}

func writeCsvFullOutput(result *map[uint32][]*internal.VoteTransaction) {
	fileName := "out_full.csv"

	file, _ := os.Create(fileName)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for candidateId, transactions := range *result {
		for _, transaction := range transactions {
			err := writer.Write([]string{
				strconv.Itoa(int(candidateId)),
				transaction.Hash,
				strconv.Itoa(transaction.Payload.DistrictId)})
			if err != nil {
				log.Fatalf("Failed to write into file %s", fileName)
				return
			}
		}
	}
}
