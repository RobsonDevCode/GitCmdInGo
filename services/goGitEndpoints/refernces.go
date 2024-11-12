package goGitEndpoints

import (
	. "GoGitCLI.RobsonDevCode.Com/models/gogitApiModels"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

// type of Reference types git can take in Go doesn't support Enums so we make do, could also be done with constants
var (
	TypeBranch = ReferenceType{Name: "branch", ID: 0}
	TypeTag    = ReferenceType{Name: "tag", ID: 1}
	TypeCommit = ReferenceType{Name: "commit", ID: 2}
)

func FetchRefs(url string) ([]GoGitReference, error) {
	// Make HTTP GET request to /info/refs
	url = url + "/info/refs"

	log.Infof("Fetching references from %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Error fetching references from %s : %s", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	var refsResponse []GoGitReference
	if err := json.NewDecoder(resp.Body).Decode(&refsResponse); err != nil {
		return nil, err
	}
	log.Infof("Api Responded with %s", refsResponse)

	if refsResponse == nil || len(refsResponse) == 0 {
		errMessage := fmt.Sprintf("Repository %s cannot be found", url)
		log.Error(errMessage)
		return nil, errors.New(errMessage)
	}

	for i := range refsResponse {
		if len(refsResponse[i].Head.Name) > 0 {
			refsResponse[i].Type = TypeBranch
		} else if len(refsResponse[i].Tags.Name) > 0 {
			refsResponse[i].Type = TypeTag
		} else {
			refsResponse[i].Type = TypeCommit
		}
	}

	return refsResponse, nil
}
