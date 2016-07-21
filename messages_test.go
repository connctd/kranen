package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var examplePayload = `{
  "callback_url": "https://registry.hub.docker.com/u/connctd/gate/hook/25040jj1i4e2a4jc1ehbabb41hj45h0ef/",
  "push_data": {
    "images": [],
    "pushed_at": 1.469096372e+09,
    "pusher": "connctddev",
    "tag": "latest-development"
  },
  "repository": {
    "comment_count": 0,
    "date_created": 1.465379301e+09,
    "description": "",
    "full_description": null,
    "is_official": false,
    "is_private": true,
    "is_trusted": false,
    "name": "testhook",
    "namespace": "connctd",
    "owner": "connctd",
    "repo_name": "connctd/gate",
    "repo_url": "https://hub.docker.com/r/connctd/gate",
    "star_count": 0,
    "status": "Active"
  }
}`

func TestPayloadStruct(t *testing.T) {
	assert := assert.New(t)
	payload := Payload{}
	err := json.Unmarshal([]byte(examplePayload), &payload)
	assert.Nil(err)

	assert.NotNil(payload.PushData)
	assert.NotNil(payload.Repo)
	assert.Equal("testhook", payload.Repo.Name)
}
