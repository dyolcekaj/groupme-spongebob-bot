package groupme

import (
	"encoding/json"
	"github.com/dyolcekaj/groupme-spongebob-bot/assertions"
	"testing"
)

const searchResultJson = `{
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

	assertions.Ok(t, json.Unmarshal([]byte(searchResultJson), &result))
}

const listBotJson = `{
	"meta": {
		"code": 200
	},
	"response": [{
		"name": "spongebob",
		"bot_id": "939b3463663dd9a70a86b8b876",
		"group_id": "27041965",
		"group_name": "Let's Bitch About Nebraska",
		"avatar_url": "https://i.groupme.com/529x379.jpeg.6b4c6928f02f40b283bf56e4606b32b5",
		"callback_url": "https://hg9qpvogfj.execute-api.us-east-1.amazonaws.com/default/groupme-spongebob",
		"dm_notification": false
	}, {
		"name": "testbot",
		"bot_id": "52c2c25545d72f52fc4b9f9f9c",
		"group_id": "56570638",
		"group_name": "My Test Chat",
		"avatar_url": "https://i.groupme.com/529x379.jpeg.6b4c6928f02f40b283bf56e4606b32b5",
		"callback_url": "https://hg9qpvogfj.execute-api.us-east-1.amazonaws.com/default/groupme-spongebob",
		"dm_notification": false
	}]
}`

func TestListBotUnmarshalling(t *testing.T) {
	var result = struct {
		Response []Bot
	} {}

	assertions.Ok(t, json.Unmarshal([]byte(listBotJson), &result))
}