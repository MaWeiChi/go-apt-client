//
//  This file is part of go-apt-client library
//
//  Copyright (C) 2017  Arduino AG (http://www.arduino.cc/)
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package apt

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAPTConfigFolder(t *testing.T) {
	repos, err := ParseAPTConfigFolder("testdata/apt")
	require.NoError(t, err, "running List command")

	expectedData, err := ioutil.ReadFile("testdata/TestParseAPTConfigFolder.json")
	require.NoError(t, err, "Reading test data")
	expected := []*Repository{}
	err = json.Unmarshal(expectedData, &expected)
	require.NoError(t, err, "Decoding expected data")

	for i, repo := range repos {
		require.EqualValues(t, expected[i], repo, "Comparing element %d", i)
	}
}

func TestAddAndRemoveRepository(t *testing.T) {
	// test cleanup
	defer os.Remove("testdata/apt2/sources.list.d/managed.list")

	repo1 := &Repository{
		Enabled:      true,
		SourceRepo:   false,
		URI:          "http://ppa.launchpad.net/webupd8team/java/ubuntu",
		Distribution: "zesty",
		Components:   "main",
		Comment:      "",
	}
	repo2 := &Repository{
		Enabled:      false,
		SourceRepo:   true,
		URI:          "http://ppa.launchpad.net/webupd8team/java/ubuntu",
		Distribution: "zesty",
		Components:   "main",
		Comment:      "",
	}
	err := AddRepository(repo1, "testdata/apt2")
	require.NoError(t, err, "Adding repository")
	err = AddRepository(repo2, "testdata/apt2")
	require.NoError(t, err, "Adding repository")

	repos, err := ParseAPTConfigFolder("testdata/apt2")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo2)

	err = AddRepository(repo2, "testdata/apt2")
	require.Error(t, err, "Adding repository again")

	err = RemoveRepository(repo2, "testdata/apt2")
	require.NoError(t, err, "Removing repository")

	repos, err = ParseAPTConfigFolder("testdata/apt2")
	require.NoError(t, err, "running List command")
	require.True(t, repos.Contains(repo1), "Configuration contains: %#v", repo1)
	require.False(t, repos.Contains(repo2), "Configuration contains: %#v", repo2)

	err = RemoveRepository(repo2, "testdata/apt2")
	require.Error(t, err, "Removing repository again")
}
