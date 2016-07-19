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
  "callback_url": "https://registry.hub.docker.com/u/svendowideit/testhook/hook/2141b5bi5i5b02bec211i4eeih0242eg11000a/",
  "push_data": {
    "images": [
        "27d47432a69bca5f2700e4dff7de0388ed65f9d3fb1ec645e2bc24c223dc1cc3",
        "51a9c7c1f8bb2fa19bcd09789a34e63f35abb80044bc10196e304f6634cc582c",
        "..."
    ],
    "pushed_at": 1.417566161e+09,
    "pusher": "trustedbuilder",
    "tag": "latest"
  },
  "repository": {
    "comment_count": "0",
    "date_created": 1.417494799e+09,
    "description": "",
    "dockerfile": "irrelevant",
    "is_official": false,
    "is_private": true,
    "is_trusted": true,
    "name": "testhook",
    "namespace": "connctd",
    "owner": "connctd",
    "repo_name": "connctd/test",
    "repo_url": "https://registry.hub.docker.com/u/svendowideit/testhook/",
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
	prepare()
	assert := assert.New(t)

	execCommand = fakeExecCommand

	request := &http.Request{}
	bodyBuf := bytes.Buffer{}
	bodyBuf.WriteString(successPayload)
	request.Body = ioutil.NopCloser(&bodyBuf)
	request.Method = "GET"
	request.RequestURI = "/docker/foobaz"

	apiKeyParam := httprouter.Param{
		Key:   "apikey",
		Value: "foobaz",
	}
	w := httptest.NewRecorder()
	hook(w, request, httprouter.Params{apiKeyParam})
	w.Flush()
	assert.Equal(http.StatusOK, w.Code)
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, "I'm just a helper")
	os.Exit(0)
}
