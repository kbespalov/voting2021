package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"voting2021/decryptor/internal/input"
	"voting2021/decryptor/internal/output"
	"voting2021/decryptor/internal/processor"
)

// cmd {token} {path_to_csv}

func main() {

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

	processor := processor.NewProcessor(privateKey)
	input.ProcessVotesCsvFile(pathToFile, processor)
	votes := processor.GetVotes()

	result := make([]output.VotingResult, 0)
	for candidateId, transactions := range votes {
		fmt.Printf("Candidate ID: %d, Votes: %d\n", candidateId, len(transactions))
		result = append(result, output.VotingResult{CandidateId: candidateId, VotesCount: len(transactions)})
	}

	file_tiny, _ := json.MarshalIndent(result, "", " ")
	_ = ioutil.WriteFile("out_tiny.json", file_tiny, 0644)

	file_full, _ := json.MarshalIndent(votes, "", " ")
	_ = ioutil.WriteFile("out_full.json", file_full, 0644)
}
