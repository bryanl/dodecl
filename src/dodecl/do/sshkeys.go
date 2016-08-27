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
	"strconv"

	"github.com/digitalocean/godo"
)

// SSHKey wraps godo Key.
type SSHKey struct {
	*godo.Key
}

// SSHKeys is a slice of SSHKey
type SSHKeys []SSHKey

// KeysService is the godo KeysService interface.
type KeysService interface {
	List() (SSHKeys, error)
	Get(id string) (*SSHKey, error)
	Create(kcr *godo.KeyCreateRequest) (*SSHKey, error)
	Update(id string, kur *godo.KeyUpdateRequest) (*SSHKey, error)
	Delete(id string) error
}

type keysService struct {
	client *godo.Client
}

var _ KeysService = &keysService{}

// NewKeysService builds an instance of KeysService.
func NewKeysService(client *godo.Client) KeysService {
	return &keysService{
		client: client,
	}
}

func (ks *keysService) List() (SSHKeys, error) {
	f := func(opt *godo.ListOptions, out chan interface{}) (*godo.Response, error) {
		list, resp, err := ks.client.Keys.List(opt)
		if err != nil {
			return nil, err
		}

		for _, d := range list {
			out <- d
		}

		return resp, nil
	}

	resp, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	items := resp.([]interface{})
	list := make(SSHKeys, len(items))
	for i := range items {
		d := items[i].(godo.Key)
		list[i] = SSHKey{Key: &d}
	}

	return list, nil
}

func (ks *keysService) Get(id string) (*SSHKey, error) {
	var err error
	var k *godo.Key

	if i, aerr := strconv.Atoi(id); aerr == nil {
		k, _, err = ks.client.Keys.GetByID(i)
	} else {
		if len(id) > 0 {
			k, _, err = ks.client.Keys.GetByFingerprint(id)
		} else {
			err = fmt.Errorf("missing key id or fingerprint")
		}
	}

	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Create(kcr *godo.KeyCreateRequest) (*SSHKey, error) {
	k, _, err := ks.client.Keys.Create(kcr)
	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Update(id string, kur *godo.KeyUpdateRequest) (*SSHKey, error) {
	var k *godo.Key
	var err error
	if i, aerr := strconv.Atoi(id); aerr == nil {
		k, _, err = ks.client.Keys.UpdateByID(i, kur)
	} else {
		k, _, err = ks.client.Keys.UpdateByFingerprint(id, kur)
	}

	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Delete(id string) error {
	var err error

	if i, aerr := strconv.Atoi(id); aerr == nil {
		_, err = ks.client.Keys.DeleteByID(i)
	} else {
		_, err = ks.client.Keys.DeleteByFingerprint(id)
	}

	return err
}