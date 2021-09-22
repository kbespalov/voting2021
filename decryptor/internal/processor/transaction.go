package processor

import (
	"log"
	"voting2021/decryptor/internal"
	"voting2021/decryptor/internal/crypto"
)

type TransactionProcessor interface {
	ExtractChoices(transaction *internal.VoteTransaction) ([]uint32, error)
	GetVotes() map[uint32][]*internal.VoteTransaction
}

type processor struct {
	privateKey     string
	processedCount int
	votes          map[uint32][]*internal.VoteTransaction
}

func (receiver *processor) GetVotes() map[uint32][]*internal.VoteTransaction {
	return receiver.votes
}

func (receiver *processor) ExtractChoices(transaction *internal.VoteTransaction) ([]uint32, error) {
	choices, err := crypto.DecryptVoteMessage(transaction.Payload.EncryptedChoice, receiver.privateKey)
	if err != nil {
		log.Printf("Failed to decrypt transaction with id %s", transaction.Hash)
		return nil, nil
	}
	receiver.processedCount = receiver.processedCount + 1
	if receiver.processedCount%10000 == 0 {
		log.Printf("Processed %d transactions", receiver.processedCount)
	}

	if len(choices) > 1 {
		log.Printf("Found multiple choices in transaction %s", transaction.Hash)
	}

	for _, candidateId := range choices {

		if candidateVotes, ok := receiver.votes[candidateId]; ok {
			receiver.votes[candidateId] = append(candidateVotes, transaction)
		} else {
			newCandidateVotes := make([]*internal.VoteTransaction, 0)
			newCandidateVotes = append(newCandidateVotes, transaction)
			receiver.votes[candidateId] = newCandidateVotes
		}

	}
	return choices, err
}

func NewProcessor(privateKey string) *processor {
	return &processor{
		privateKey:     privateKey,
		votes:          map[uint32][]*internal.VoteTransaction{},
		processedCount: 0,
	}
}
