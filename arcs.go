package aegirdungeons

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/types"
)

// ErrMetadataNotFound ...
var ErrMetadataNotFound = errors.New("metadata: no metadata found in transaction history")

// Arc69 ...
type Arc69 struct {
	Standard    string     `json:"standard"`
	Description string     `json:"description"`
	ExternalURL string     `json:"external_url"`
	MimeType    string     `json:"mime_type"`
	Properties  Properties `json:"properties"`
}

// Properties ...
type Properties struct {
	Level  string `json:"Level"`
	Rarity string `json:"Rarity"`
}

func FetchArc69Metadata(ctx context.Context, ai *indexer.Client, asaID uint64, senderAddr string) (uint64, Arc69, error) {
	nextToken := ""
	var tempNote Arc69
	var note Arc69
	var highestRound uint64

	for {
		txnResp, err := ai.LookupAssetTransactions(asaID).TxType(string(types.AssetConfigTx)).NextToken(nextToken).Do(ctx)
		if err != nil {
			return highestRound, note, err
		}

		for _, txn := range txnResp.Transactions {
			if err := json.Unmarshal(txn.Note, &tempNote); err != nil {
				continue
			}
			if tempNote.Standard != "arc69" {
				continue
			}
			if highestRound < txn.ConfirmedRound {
				highestRound = txn.ConfirmedRound
				note = tempNote
			}
		}

		nextToken = txnResp.NextToken
		if nextToken == "" {
			break
		}
	}

	if highestRound == 0 {
		return highestRound, note, ErrMetadataNotFound
	}

	return highestRound, note, nil
}
