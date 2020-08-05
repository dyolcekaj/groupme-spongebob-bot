package groupme

import (
	"encoding/json"
	"testing"
)

const searchResultJSON = `{
	"response": {
		"count": 14476,
		"messages": [{
			"attachments": [{
				"loci": [
					[18, 11]
				],
				"type": "mentions",
				"user_ids": ["10033239"]
			}],
			"avatar_url": "https://i.groupme.com/760x1344.png.96e6be7b8a6747278605aff796d6b6e9",
			"created_at": 1578193122,
			"favorited_by": [],
			"group_id": "27041965",
			"id": "157819312246955708",
			"name": "Jake Cloyd",
			"sender_id": "13328643",
			"sender_type": "user",
			"source_guid": "8d0a1def4554eb04ed6f18703aeef367",
			"system": false,
			"text": "test message with @Calvin Mak ",
			"user_id": "13328643",
			"platform": "gm"
		}]
	}
}`

func TestGroupMessageSearchUnmarshalling(t *testing.T) {
	var result groupMessageSearchResult

	if err := json.Unmarshal([]byte(searchResultJSON), &result); err != nil {
		t.Errorf("exp no err, got %v", err)
	}
}
