/*
Copyright 2016 The Doctl Authors All rights reserved.
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

package do

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/digitalocean/godo"
)

const maxFetchPages = 5

var perPage = 200

var fetchFn = fetchPage

type paginatedList struct {
	list []interface{}
	mu   sync.Mutex
}

func (pl *paginatedList) append(items ...interface{}) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.list = append(pl.list, items...)
}

// Generator is a function that generates the list to be paginated.
type Generator func(*godo.ListOptions, chan interface{}) (*godo.Response, error)

// PaginateResp paginates a Response.
func PaginateResp(gen Generator) (interface{}, error) {
	opt := &godo.ListOptions{Page: 1, PerPage: perPage}

	fetchChan := make(chan int, maxFetchPages)
	out := make(chan interface{}, 100)

	var wg sync.WaitGroup
	for i := 0; i < maxFetchPages-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range fetchChan {
				if err := fetchFn(gen, page, out); err != nil {
					// TODO we should do something here
					fmt.Printf("something went wrong with the fetches: %s", err)
					return
				}
			}
		}()
	}

	// fetch first page to get page count (x)

	resp, err := gen(opt, out)
	if err != nil {
		return nil, err
	}

	// find last page
	lp, err := lastPage(resp)
	if err != nil {
		return nil, err
	}

	// start with second page
	opt.Page++
	for ; opt.Page <= lp; opt.Page++ {
		fetchChan <- opt.Page
	}
	close(fetchChan)

	wg.Wait()

	<-fetchChan
	close(out)

	list := []interface{}{}
	for item := range out {
		list = append(list, item)
	}

	return list, nil
}

func fetchPage(gen Generator, page int, out chan interface{}) error {
	opt := &godo.ListOptions{Page: page, PerPage: 200}
	_, err := gen(opt, out)
	return err
}

func lastPage(resp *godo.Response) (int, error) {
	if resp.Links == nil || resp.Links.Pages == nil {
		// no other pages
		return 1, nil
	}

	uStr := resp.Links.Pages.Last
	if uStr == "" {
		return 1, nil
	}

	u, err := url.Parse(uStr)
	if err != nil {
		return 0, fmt.Errorf("could not parse last page: %v", err)
	}

	pageStr := u.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, fmt.Errorf("could not find page param: %v", err)
	}

	return page, err
}
