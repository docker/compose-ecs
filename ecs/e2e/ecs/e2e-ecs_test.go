/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/icmd"
	"gotest.tools/v3/poll"
)

func TestSecrets(t *testing.T) {
	startTime := strconv.Itoa(int(time.Now().UnixNano()))
	secretName := "secret" + strings.ToLower(t.Name()) + startTime
	t.Run("create secret", func(t *testing.T) {
		secretFile := filepath.Join(t.TempDir(), "secret.txt")
		err := os.WriteFile(secretFile, []byte("pass1"), 0644)
		assert.Check(t, err == nil)
		res := icmd.RunCommand("compose-ecs", "secret", "create", secretName, secretFile)
		assert.Check(t, strings.Contains(res.Stdout(), secretName), res.Stdout())
	})

	t.Run("list secrets", func(t *testing.T) {
		res := icmd.RunCommand("compose-ecs", "secret", "list")
		assert.Check(t, strings.Contains(res.Stdout(), secretName), res.Stdout())
	})

	t.Run("inspect secret", func(t *testing.T) {
		res := icmd.RunCommand("compose-ecs", "secret", "inspect", secretName)
		assert.Check(t, strings.Contains(res.Stdout(), `"Name": "`+secretName+`"`), res.Stdout())
	})

	t.Run("rm secret", func(t *testing.T) {
		icmd.RunCommand("compose-ecs", "secret", "rm", secretName)
		res := icmd.RunCommand("secret", "list")
		assert.Check(t, !strings.Contains(res.Stdout(), secretName), res.Stdout())
	})
}

func TestCompose(t *testing.T) {
	t.Run("compose-ecs up", func(t *testing.T) {
		res := icmd.RunCommand("compose-ecs", "up")
		res.Assert(t, icmd.Success)
	})

	var webURL, wordsURL, secretsURL string
	t.Run("compose-ecs ps", func(t *testing.T) {
		res := icmd.RunCommand("compose-ecs", "ps")
		lines := strings.Split(strings.TrimSpace(res.Stdout()), "\n")

		assert.Equal(t, 5, len(lines))

		var dbDisplayed, wordsDisplayed, webDisplayed, secretsDisplayed bool
		for _, line := range lines {
			fields := strings.Fields(line)
			containerID := fields[0]
			serviceName := fields[1]
			switch serviceName {
			case "db":
				dbDisplayed = true
				assert.DeepEqual(t, fields, []string{containerID, serviceName, "Running"})
			case "words":
				wordsDisplayed = true
				assert.Check(t, strings.Contains(fields[3], ":8080->8080/tcp"),
					"Got -> %q. All fields -> %#v", fields[3], fields)
				wordsURL = "http://" + strings.Replace(fields[3], "->8080/tcp", "", 1) + "/noun"
			case "web":
				webDisplayed = true
				assert.Check(t, strings.Contains(fields[3], ":80->80/tcp"),
					"Got -> %q. All fields -> %#v", fields[3], fields)
				webURL = "http://" + strings.Replace(fields[3], "->80/tcp", "", 1)
			case "websecrets":
				secretsDisplayed = true
				assert.Check(t, strings.Contains(fields[3], ":90->90/tcp"),
					"Got -> %q. All fields -> %#v", fields[3], fields)
				secretsURL = "http://" + strings.Replace(fields[3], "->90/tcp", "", 1)
			}
		}

		assert.Check(t, dbDisplayed)
		assert.Check(t, wordsDisplayed)
		assert.Check(t, webDisplayed)
		assert.Check(t, secretsDisplayed)
	})

	t.Run("Words GET validating cross service connection", func(t *testing.T) {
		out := HTTPGetWithRetry(t, wordsURL, http.StatusOK, 5*time.Second, 300*time.Second)
		assert.Assert(t, strings.Contains(out, `"word":`))
	})

	t.Run("web app GET", func(t *testing.T) {
		out := HTTPGetWithRetry(t, webURL, http.StatusOK, 3*time.Second, 120*time.Second)
		assert.Assert(t, strings.Contains(out, "Docker Compose demo"))

		out = HTTPGetWithRetry(t, webURL+"/words/noun", http.StatusOK, 2*time.Second, 60*time.Second)
		assert.Assert(t, strings.Contains(out, `"word":`))
	})

	t.Run("access secret", func(t *testing.T) {
		out := HTTPGetWithRetry(t, secretsURL+"/mysecret1", http.StatusOK, 3*time.Second, 120*time.Second)
		out = strings.ReplaceAll(out, "\r", "")
		assert.Equal(t, out, "myPassword1\n")
	})

	t.Run("compose-ecs down", func(t *testing.T) {
		res := icmd.RunCommand("compose-ecs", "down")

		checkUp := func(t poll.LogT) poll.Result {
			out := res.Combined()
			if !strings.Contains(out, "DeleteComplete") {
				return poll.Continue("current status \n%s\n", out)
			}
			return poll.Success()
		}
		poll.WaitOn(t, checkUp, poll.WithDelay(2*time.Second), poll.WithTimeout(60*time.Second))
	})
}

// HTTPGetWithRetry performs an HTTP GET on an `endpoint`, using retryDelay also as a request timeout.
// In the case of an error or the response status is not the expected one, it retries the same request,
// returning the response body as a string (empty if we could not reach it)
func HTTPGetWithRetry(
	t testing.TB,
	endpoint string,
	expectedStatus int,
	retryDelay time.Duration,
	timeout time.Duration,
) string {
	t.Helper()
	var (
		r   *http.Response
		err error
	)
	client := &http.Client{
		Timeout: retryDelay,
	}
	fmt.Printf("\t[%s] GET %s\n", t.Name(), endpoint)
	checkUp := func(t poll.LogT) poll.Result {
		r, err = client.Get(endpoint)
		if err != nil {
			return poll.Continue("reaching %q: Error %s", endpoint, err.Error())
		}
		if r.StatusCode == expectedStatus {
			return poll.Success()
		}
		return poll.Continue("reaching %q: %d != %d", endpoint, r.StatusCode, expectedStatus)
	}
	poll.WaitOn(t, checkUp, poll.WithDelay(retryDelay), poll.WithTimeout(timeout))
	if r != nil {
		b, err := io.ReadAll(r.Body)
		assert.NilError(t, err)
		return string(b)
	}
	return ""
}
