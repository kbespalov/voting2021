package input

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"voting2021/decryptor/internal"
	"voting2021/decryptor/internal/processor"
)

func ProcessVotesCsvFile(filePath string, handler processor.TransactionProcessor) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	for {
		csvRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		parsedRecord, err := parseCsvRecord(csvRecord)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = handler.ExtractChoices(parsedRecord)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

// CSV with two columns [transaction hash, json payload]
func parseCsvRecord(record []string) (*internal.VoteTransaction, error) {
	hash := record[0]
	payloadRaw := record[1]
	payload := internal.VotePayload{}
	err := json.Unmarshal([]byte(payloadRaw), &payload)
	if err != nil {
		return nil, err
	}
	return &internal.VoteTransaction{Hash: hash, Payload: payload}, nil
}
