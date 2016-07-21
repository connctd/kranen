package main

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

var (
	successPayload = `{
  "callback_url": "https://registry.hub.docker.com/u/connctd/gate/hook/25040jj1i4e2a4jc1ehbabb41hj45h0ef/",
  "push_data": {
    "images": [],
    "pushed_at": 1.469096372e+09,
    "pusher": "connctddev",
    "tag": "latest"
  },
  "repository": {
    "comment_count": 0,
    "date_created": 1.465379301e+09,
    "description": "",
    "full_description": null,
    "is_official": false,
    "is_private": true,
    "is_trusted": false,
    "name": "test",
    "namespace": "connctd",
    "owner": "connctd",
    "repo_name": "connctd/test",
    "repo_url": "https://hub.docker.com/r/connctd/gate",
    "star_count": 0,
    "status": "Active"
  }
}`

	wrongTagPayload = `{
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
    "name": "test",
    "namespace": "connctd",
    "owner": "connctd",
    "repo_name": "connctd/test",
    "repo_url": "https://hub.docker.com/r/connctd/gate",
    "star_count": 0,
    "status": "Active"
  }
}`
	wrongNamePayload = `{
  "callback_url": "https://registry.hub.docker.com/u/connctd/gate/hook/25040jj1i4e2a4jc1ehbabb41hj45h0ef/",
  "push_data": {
    "images": [],
    "pushed_at": 1.469096372e+09,
    "pusher": "connctddev",
    "tag": "latest"
  },
  "repository": {
    "comment_count": 0,
    "date_created": 1.465379301e+09,
    "description": "",
    "full_description": null,
    "is_official": false,
    "is_private": true,
    "is_trusted": false,
    "name": "bar",
    "namespace": "connctd",
    "owner": "connctd",
    "repo_name": "connctd/bar",
    "repo_url": "https://hub.docker.com/r/connctd/gate",
    "star_count": 0,
    "status": "Active"
  }
}`
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func prepare() {
	configs = []RepoConfig{
		RepoConfig{
			Name:   "connctd/test",
			ApiKey: "foobaz",
			Tag:    "latest",
			Script: "/deploy.sh",
		},
	}
}

func TestSuccessfullCall(t *testing.T) {
	t.Parallel()
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	testHook(successPayload, "foobaz", assert, http.StatusOK)
	testHook(wrongTagPayload, "foobaz", assert, http.StatusBadRequest)
	testHook(successPayload, "wrongapikey", assert, http.StatusNotFound)
	testHook(wrongNamePayload, "foobaz", assert, http.StatusBadRequest)
	testHook("", "foobaz", assert, http.StatusBadRequest)
}

func TestWrontTag(t *testing.T) {
	t.Parallel()
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	testHook(wrongTagPayload, "foobaz", assert, http.StatusBadRequest)
}

func TestWrongApiKey(t *testing.T) {
	t.Parallel()
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	testHook(successPayload, "wrongapikey", assert, http.StatusNotFound)
}

func TestWrongName(t *testing.T) {
	t.Parallel()
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	testHook(wrongNamePayload, "foobaz", assert, http.StatusBadRequest)
}

func TestUnparseablePayload(t *testing.T) {
	t.Parallel()
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	testHook("", "foobaz", assert, http.StatusBadRequest)
}

func testHook(payload, apikey string, assert *assert.Assertions, expectedStatusCode int) {
	request := &http.Request{}
	bodyBuf := bytes.Buffer{}
	bodyBuf.WriteString(payload)
	request.Body = ioutil.NopCloser(&bodyBuf)
	request.Method = "GET"
	request.RequestURI = fmt.Sprintf("/docker/%s", apikey)

	apiKeyParam := httprouter.Param{
		Key:   "apikey",
		Value: apikey,
	}
	w := httptest.NewRecorder()
	hook(w, request, httprouter.Params{apiKeyParam})
	w.Flush()
	assert.Equal(expectedStatusCode, w.Code)
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "I'm just a helper")
	os.Exit(0)
}
