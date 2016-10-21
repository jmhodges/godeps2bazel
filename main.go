// Copyright 2016 Jeffrey M Hodges
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/go/vcs"
)

type godeps struct {
	Deps []*dep
}

type dep struct {
	ImportPath string
	Rev        string
}

func main() {
	file := strings.TrimSpace(os.Args[1])
	if file == "" {
		log.Fatal("usage: godeps2bazel GODEPS_JSON_FILE")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("unable to read %#v: %s", file, err)
	}
	g := &godeps{}
	err = json.Unmarshal(b, g)
	if err != nil {
		log.Fatalf("unable to unmarshal Godeps file %#v: %s", file, err)
	}
	if err != nil {
		log.Fatalf("unable to create VCS detection object: %s", err)
	}

	// We don't need to recurse through the dependencies, because we can
	// reasonably assume that the user has kept their Godeps up to date.
	remotes := make(map[string]bool)
	repos := make(sortableRepos, 0)

	for _, d := range g.Deps {
		repo, err := vcs.RepoRootForImportPath(d.ImportPath, false)
		if err != nil {
			log.Fatalf("unable to detect actual git repo URL for %#v: %s", d.ImportPath, err)
		}

		if remotes[repo.Root] {
			continue
		}
		remotes[repo.Root] = true
		repos = append(repos, bzrepo{repoName(repo.Root), repo.Repo, d.Rev})
	}

	sort.Sort(repos)

	out := []string{}
	for _, r := range repos {
		msg := fmt.Sprintf(`new_go_repository(
    name = "%s",
    importpath = "%s",
    commit = "%s",
)`, r.name, r.remote, r.commit)
		out = append(out, msg)
	}
	fmt.Printf("%s\n", strings.Join(out, "\n\n"))
}

// mergeImportPathInto adds the path to allImportPaths and returns true if f the
// path or its prefix is not already in allImportPaths
func mergeImportPathInto(path string, allImportPaths map[string]bool) bool {
	for k, _ := range allImportPaths {
		if k == path {
			// We already have the new path in the map
			return false
		} else if strings.HasPrefix(k, path) {
			// We already have the new path's prefix in place.
			return false
		} else if strings.HasPrefix(path, k) {
			delete(allImportPaths, k)
			allImportPaths[path] = true
			return true
		}
	}
	allImportPaths[path] = true
	return true
}

func repoName(prefix string) string {
	components := strings.Split(prefix, "/")
	labels := strings.Split(components[0], ".")
	var reversed []string
	for i := range labels {
		l := labels[len(labels)-i-1]
		reversed = append(reversed, l)
	}
	repo := strings.Join(append(reversed, components[1:]...), "_")
	repo = strings.NewReplacer("-", "_", ".", "_").Replace(repo)
	return repo
}

type bzrepo struct {
	name   string
	remote string
	commit string
}

type sortableRepos []bzrepo

func (s sortableRepos) Len() int {
	return len(s)
}

func (s sortableRepos) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortableRepos) Less(i, j int) bool { return s[i].name < s[j].name }
