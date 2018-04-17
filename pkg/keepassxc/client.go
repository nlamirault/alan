// Copyright (C) 2018 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keepassxc

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/tobischo/gokeepasslib"

	pkgalan "github.com/nlamirault/alan/pkg/alan"
	pkgcmd "github.com/nlamirault/alan/pkg/cmd"
)

// Client define a client to manage KeepassXC database
type Client struct {
	filename string
	db       *gokeepasslib.Database
}

// NewClient create a new KeepassXC database client
func NewClient(filename string) (*Client, error) {
	return &Client{
		filename: filename,
	}, nil
}

func (client *Client) Open() error {
	glog.V(2).Infof("Open database from file: %s", client.filename)
	if _, err := os.Stat(client.filename); os.IsNotExist(err) {
		return fmt.Errorf("Database file not exists")
	}
	file, err := os.Open(client.filename)
	if err != nil {
		return err
	}
	client.db = gokeepasslib.NewDatabase()

	password, err := pkgcmd.ReadPassword("Please input your password: ")
	if err != nil {
		return err
	}
	client.db.Credentials = gokeepasslib.NewPasswordCredentials(password)
	if err := gokeepasslib.NewDecoder(file).Decode(client.db); err != nil {
		return err
	}

	glog.V(2).Infof("Database Metadata: %#v", client.db.Content.Meta)

	glog.V(2).Info("Unlock database entries")
	client.db.UnlockProtectedEntries()

	return nil
}

func (client *Client) Save() error {
	glog.V(2).Infof("Output file for database: %s", client.filename)

	file, err := os.Create(client.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	keepassEncoder := gokeepasslib.NewEncoder(file)
	if err := keepassEncoder.Encode(client.db); err != nil {
		return err
	}

	glog.V(1).Infof("Databased save into: %s", client.filename)
	return nil
}

func (client *Client) Close() error {
	glog.V(2).Infof("Close KeepassXC database: %s", client.filename)
	return client.db.LockProtectedEntries()
}

func (client *Client) Load() (map[string][]pkgalan.Secret, error) {
	secrets := map[string][]pkgalan.Secret{}
	root := client.db.Content.Root
	if len(root.Groups) == 0 {
		return secrets, nil
	}
	for _, group := range root.Groups {
		groupSecrets, err := client.manageGroupEntries(group)
		if err != nil {
			return nil, err
		}
		secrets[group.Name] = groupSecrets
		for _, subgroup := range group.Groups {
			glog.V(2).Infof("Manage group: %s", subgroup.Name)
			subgroupSecrets, err := client.manageGroupEntries(subgroup)
			if err != nil {
				return nil, err
			}
			secrets[subgroup.Name] = subgroupSecrets
		}
	}
	return secrets, nil
}

func (client *Client) Create(secrets map[string][]*pkgalan.Secret) error {
	glog.V(2).Infof("Add secrets to database")

	rootGroup := gokeepasslib.NewGroup()
	rootGroup.Name = "Root"
	rootGroup.EnableSearching = "true"
	rootGroup.EnableAutoType = "true"
	for groupName, groupSecrets := range secrets {
		glog.V(2).Infof("Add group secrets: %s %s", groupName, groupSecrets)
		subGroup := gokeepasslib.NewGroup()
		subGroup.Name = strings.Replace(groupName, "/", "", -1)
		subGroup.EnableSearching = "true"
		subGroup.EnableAutoType = "true"
		for _, secret := range groupSecrets {
			entry := gokeepasslib.NewEntry()
			entry.Values = append(entry.Values, mkValue(pkgalan.Title, secret.Title))
			entry.Values = append(entry.Values, mkValue(pkgalan.Username, secret.Username))
			entry.Values = append(entry.Values, mkValue(pkgalan.URL, secret.URL))
			entry.Values = append(entry.Values, mkProtectedValue(pkgalan.Password, secret.Password))
			subGroup.Entries = append(subGroup.Entries, entry)
		}
		rootGroup.Groups = append(rootGroup.Groups, subGroup)
	}

	glog.V(2).Info("Create a new database")
	password, err := pkgcmd.ReadPassword("Please input your password: ")
	if err != nil {
		return err
	}

	meta := gokeepasslib.NewMetaData()
	meta.Generator = pkgalan.Generator
	meta.HistoryMaxItems = 10
	meta.MaintenanceHistoryDays = "365"
	meta.HistoryMaxSize = 6291456
	client.db = &gokeepasslib.Database{
		Signature:   &gokeepasslib.DefaultSig,
		Headers:     gokeepasslib.NewFileHeaders(),
		Credentials: gokeepasslib.NewPasswordCredentials(password),
		Content: &gokeepasslib.DBContent{
			Meta: meta,
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{rootGroup},
			},
		},
	}
	return nil
}

func (client *Client) manageGroupEntries(group gokeepasslib.Group) ([]pkgalan.Secret, error) {
	secrets := []pkgalan.Secret{}
	for _, entry := range group.Entries {
		if len(entry.GetTitle()) == 0 {
			glog.Infof("Skipping entry: %s", entry.GetContent(pkgalan.URL))
		} else {
			glog.V(1).Infof("Add entry: %s %s", entry.GetTitle(), entry.GetContent(pkgalan.URL))
			secret := pkgalan.Secret{
				Title:    entry.GetTitle(),
				Username: entry.GetContent(pkgalan.Username),
				Password: entry.GetContent(pkgalan.Password),
				URL:      entry.GetContent(pkgalan.URL),
			}
			secrets = append(secrets, secret)
		}
	}
	return secrets, nil
}

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value, Protected: true}}
}
