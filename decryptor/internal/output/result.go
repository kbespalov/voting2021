package output

type VotingResult struct {
	CandidateId uint32 `json:"candidate_id,omitempty"`
	VotesCount  int    `json:"votes,omitempty"`
}
