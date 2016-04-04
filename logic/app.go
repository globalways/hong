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
	"github.com/aiwuTech/devKit/random"
	"github.com/globalways/hong/g"
	"github.com/globalways/hong/modal"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
)

type AppAdapter interface {
	SyncDB() error
	NewApp(name, desc string) (*modal.App, error)
	GetApp(key string) (*modal.App, error)
}

type AppDefault struct {
	db *gorm.DB
}

func NewAppDefault(db *gorm.DB) *AppDefault {
	return &AppDefault{
		db: db,
	}
}

func (this *AppDefault) NewApp(name, desc string) (*modal.App, error) {
	if name == "" {
		return nil, errors.New(g.INVALID_PARAM_ERROR)
	}

	app := &modal.App{
		Key:    random.RandomAlphanumeric(16),
		Secret: random.RandomAlphanumeric(32),
		Name:   name,
		Desc:   desc,
	}

	if err := this.db.Create(app).Error; err != nil {
		return nil, errors.New(err)
	}

	return app, nil
}

func (this *AppDefault) GetApp(key string) (*modal.App, error) {
	if key == "" {
		return nil, errors.New(g.INVALID_PARAM_ERROR)
	}

	app := &modal.App{}
	if err := this.db.Where(&modal.App{Key: key}).First(app).Error; err != nil {
		return nil, errors.New(err)
	}

	return app, nil
}

func (this *AppDefault) SyncDB() error {
	tx := this.db.Begin()
	{
		app := &modal.App{}
		if tx.HasTable(app) {
			if err := tx.AutoMigrate(app).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(app).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	return nil
}
