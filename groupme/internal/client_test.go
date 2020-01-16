package internal

import (
	"encoding/json"
	"github.com/dyolcekaj/groupme-spongebob-bot/assertions"
	"testing"
)

const listBotJSON = `{
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
	}{}

	assertions.Ok(t, json.Unmarshal([]byte(listBotJSON), &result))
}
