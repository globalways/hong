// Copyright 2016 zm@huantucorp.com
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
/*
                   _ooOoo_
                  o8888888o
                  88" . "88
                  (| -_- |)
                  O\  =  /O
               ____/`---'\____
             .'  \\|     |//  `.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||-  \
           |   | \\\  -  /// |   |
           | \_|  ''\---/''  |   |
           \  .-\__  `-`  ___/-. /
         ___`. .'  /--.--\  `. . __
      ."" '<  `.___\_<|>_/___.'  >'"".
     | | :  `- \`.;`\ _ /`;.`/ - ` : | |
     \  \ `-.   \_ __\ /__ _/   .-` /  /
======`-.____`-.___\_____/___.-`____.-'======
                   `=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
         佛祖保佑       永无BUG
*/
package logic

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/globalways/hong/modal"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	"time"
)

type OauthStorage interface {
	osin.Storage

	// CreateClient stores the client in the database and returns an error, if something went wrong.
	CreateClient(client osin.Client) error

	// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
	// Returns an error if something went wrong.
	UpdateClient(client osin.Client) error

	// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
	RemoveClient(id string) error
}

type OauthStorageDefault struct {
	db *gorm.DB
}

// New returns a new mysql storage instance.
func NewDefaultOauthStorage(db *gorm.DB) *OauthStorageDefault {
	return &OauthStorageDefault{db: db}
}

// Clone the storage if needed. For example, using mgo, you can clone the session with session.Clone
// to avoid concurrent access problems.
// This is to avoid cloning the connection at each method access.
// Can return itself if not a problem.
func (this *OauthStorageDefault) Clone() osin.Storage {
	return &OauthStorageDefault{
		db: this.db.New(),
	}
}

// Close the resources the Storage potentially holds (using Clone for example)
func (this *OauthStorageDefault) Close() {
}

// CreateClient stores the client in the database and returns an error, if something went wrong.
func (this *OauthStorageDefault) CreateClient(c osin.Client) error {
	extra, err := assertToString(c.GetUserData())
	if err != nil {
		return err
	}

	client := &modal.Client{
		Id:          c.GetId(),
		Secret:      c.GetSecret(),
		RedirectUri: c.GetRedirectUri(),
		Extra:       extra,
	}
	if err := this.db.Create(client).Error; err != nil {
		return err
	}
	return nil
}

// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
func (this *OauthStorageDefault) UpdateClient(c osin.Client) error {
	extra, err := assertToString(c.GetUserData())
	if err != nil {
		return err
	}

	if err := this.db.Where(&modal.Client{Id: c.GetId()}).Save(&modal.Client{
		Id:          c.GetId(),
		Secret:      c.GetSecret(),
		Extra:       extra,
		RedirectUri: c.GetRedirectUri(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
func (this *OauthStorageDefault) RemoveClient(id string) (err error) {
	if err := this.db.Delete(&modal.Client{Id: id}).Error; err != nil {
		return err
	}
	return nil
}

// GetClient loads the client by id (client_id)
func (this *OauthStorageDefault) GetClient(id string) (osin.Client, error) {
	c := &modal.Client{}
	if err := this.db.Where(&modal.Client{Id: id}).First(c).Error; err != nil {
		return nil, errors.New("Client not found")
	}

	return &osin.DefaultClient{
		Id:          c.Id,
		Secret:      c.Secret,
		RedirectUri: c.RedirectUri,
		UserData:    c.Extra,
	}, nil
}

// SaveAuthorize saves authorize data.
func (this *OauthStorageDefault) SaveAuthorize(data *osin.AuthorizeData) error {
	extra, err := assertToString(data.UserData)
	if err != nil {
		return err
	}

	authorize := &modal.Authorize{
		Client:      data.Client.GetId(),
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		Extra:       extra,
	}

	if err := this.db.Create(authorize).Error; err != nil {
		return errors.New(err)
	}
	return nil
}

// LoadAuthorize looks up AuthorizeData by a code.
// Client information MUST be loaded together.
// Optionally can return error if expired.
func (this *OauthStorageDefault) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	data := &modal.Authorize{}
	if err := this.db.Where(&modal.Authorize{Code: code}).First(data).Error; err != nil {
		return nil, errors.New(err)
	}

	c, err := this.GetClient(data.Client)
	if err != nil {
		return nil, err
	}

	authData := &osin.AuthorizeData{
		Client:      c,
		Code:        data.Code,
		ExpiresIn:   data.ExpiresIn,
		Scope:       data.Scope,
		RedirectUri: data.RedirectUri,
		State:       data.State,
		CreatedAt:   data.CreatedAt,
		UserData:    data.Extra,
	}

	if authData.ExpireAt().Before(time.Now()) {
		return nil, errors.Errorf("Token expired at %s.", authData.ExpireAt().String())
	}

	return authData, nil
}

// RemoveAuthorize revokes or deletes the authorization code.
func (this *OauthStorageDefault) RemoveAuthorize(code string) error {
	if err := this.db.Delete(&modal.Authorize{Code: code}).Error; err != nil {
		return errors.New(err)
	}

	return nil
}

// SaveAccess writes AccessData.
// If RefreshToken is not blank, it must save in a way that can be loaded using LoadRefresh.
func (this *OauthStorageDefault) SaveAccess(data *osin.AccessData) error {
	if data.Client == nil {
		return errors.New("data.Client must not be nil")
	}

	prev := ""
	authorizeData := &osin.AuthorizeData{}

	if data.AccessData != nil {
		prev = data.AccessData.AccessToken
	}

	if data.AuthorizeData != nil {
		authorizeData = data.AuthorizeData
	}

	extra, err := assertToString(data.UserData)
	if err != nil {
		return err
	}

	tx := this.db.Begin()
	{
		if data.RefreshToken != "" {
			refresh := &modal.Refresh{
				Token:  data.RefreshToken,
				Access: data.AccessToken,
			}
			if err := tx.Create(refresh).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		access := &modal.Access{
			Client:       data.Client.GetId(),
			Authorize:    authorizeData.Code,
			Previous:     prev,
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
			ExpiresIn:    data.ExpiresIn,
			Scope:        data.Scope,
			RedirectUri:  data.RedirectUri,
			CreatedAt:    data.CreatedAt,
			Extra:        extra,
		}
		if err := tx.Create(access).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

// LoadAccess retrieves access data by token. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (this *OauthStorageDefault) LoadAccess(token string) (*osin.AccessData, error) {
	access := &modal.Access{}
	if err := this.db.Where(&modal.Access{AccessToken: token}).First(access).Error; err != nil {
		return nil, err
	}
	client, err := this.GetClient(access.Client)
	if err != nil {
		return nil, err
	}
	authorizeData, err := this.LoadAuthorize(access.Authorize)
	if err != nil {
		return nil, err
	}
	prevAccess, err := this.LoadAccess(access.Previous)
	if err != nil {
		return nil, err
	}

	result := &osin.AccessData{
		Client:        client,
		AuthorizeData: authorizeData,
		AccessData:    prevAccess,
		AccessToken:   access.AccessToken,
		RefreshToken:  access.RefreshToken,
		ExpiresIn:     access.ExpiresIn,
		Scope:         access.Scope,
		RedirectUri:   access.RedirectUri,
		CreatedAt:     access.CreatedAt,
		UserData:      access.Extra,
	}

	return result, nil
}

// RemoveAccess revokes or deletes an AccessData.
func (this *OauthStorageDefault) RemoveAccess(token string) error {
	if err := this.db.Delete(&modal.Access{AccessToken: token}).Error; err != nil {
		return err
	}

	return nil
}

// LoadRefresh retrieves refresh AccessData. Client information MUST be loaded together.
// AuthorizeData and AccessData DON'T NEED to be loaded if not easily available.
// Optionally can return error if expired.
func (this *OauthStorageDefault) LoadRefresh(token string) (*osin.AccessData, error) {
	refresh := &modal.Refresh{}
	if err := this.db.Where(&modal.Refresh{Token: token}).Error; err != nil {
		return nil, err
	}

	return this.LoadAccess(refresh.Access)
}

// RemoveRefresh revokes or deletes refresh AccessData.
func (this *OauthStorageDefault) RemoveRefresh(token string) error {
	if err := this.db.Delete(&modal.Refresh{Token: token}).Error; err != nil {
		return err
	}

	return nil
}

// Makes easy to create a osin.DefaultClient
func (this *OauthStorageDefault) CreateClientWithInformation(id string, secret string, redirectUri string, userData interface{}) osin.Client {
	return &osin.DefaultClient{
		Id:          id,
		Secret:      secret,
		RedirectUri: redirectUri,
		UserData:    userData,
	}
}

func assertToString(in interface{}) (string, error) {
	var ok bool
	var data string
	if in == nil {
		return "", nil
	} else if data, ok = in.(string); ok {
		return data, nil
	} else if str, ok := in.(fmt.Stringer); ok {
		return str.String(), nil
	}
	return "", errors.Errorf(`Could not assert "%v" to string`, in)
}

func (this *OauthStorageDefault) SyncDB() error {
	tx := this.db.Begin()
	{
		client := &modal.Client{}
		if tx.HasTable(client) {
			if err := tx.AutoMigrate(client).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(client).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		authorize := &modal.Authorize{}
		if tx.HasTable(authorize) {
			if err := tx.AutoMigrate(authorize).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(authorize).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		access := &modal.Access{}
		if tx.HasTable(access) {
			if err := tx.AutoMigrate(access).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(access).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		refresh := &modal.Refresh{}
		if tx.HasTable(refresh) {
			if err := tx.AutoMigrate(refresh).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(refresh).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	return nil
}
