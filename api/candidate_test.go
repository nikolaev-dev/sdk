// +build integration

package api

import (
	"testing"
)

func TestApi_Candidate(t *testing.T) {

	responseCandidates, err := testApi.CandidatesAtHeight(0, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(responseCandidates) == 0 {
		t.Fatal("no candidates")
	}

	response, err := testApi.Candidate(responseCandidates[0].PubKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", response)
}
