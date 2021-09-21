package internal

type EncryptedChoice struct {
	EncryptedMessage string `json:"encrypted_message,omitempty"`
	Nonce            string `json:"nonce,omitempty"`
	PublicKey        string `json:"public_key,omitempty"`
}

type VotePayload struct {
	VotingId        string           `json:"voting_id,omitempty"`
	DistrictId      int              `json:"district_id,omitempty"`
	EncryptedChoice *EncryptedChoice `json:"encrypted_choice"`
}

type VoteTransaction struct {
	Hash    string
	Payload VotePayload
}
