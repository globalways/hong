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
	"github.com/aiwuTech/devKit/convert"
	"github.com/globalways/common/utils"
	"github.com/globalways/hong/g"
	"github.com/globalways/hong/modal"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
)

type UserAdapter interface {
	NewUser(userType modal.UserType, tel, password string) (*modal.User, error)
	NextHong(userType modal.UserType) string
	GetUser(username string) (*modal.User, error)
	SyncDB() error
}

type UserDefault struct {
	db *gorm.DB
}

func NewUserDefault(db *gorm.DB) *UserDefault {
	return &UserDefault{
		db: db,
	}
}

func (this *UserDefault) NewUser(userType modal.UserType, tel, password string) (*modal.User, error) {
	if tel == "" || password == "" {
		return nil, errors.New(g.INVALID_PARAM_ERROR)
	}

	user := &modal.User{
		Hong:     this.NextHong(userType),
		Nick:     tel,
		Tel:      tel,
		Password: utils.EncodePassword(password),
		Type:     userType,
	}

	if err := this.db.Create(user).Error; err != nil {
		return nil, errors.New(err)
	}

	return user, nil
}

func (this *UserDefault) NextHong(userType modal.UserType) string {
	user := &modal.User{}
	if err := this.db.Where(&modal.User{Type: userType}).Last(user).Error; err != nil {
		nextHong := ""
		switch userType {
		case modal.UserType_TEST:
			nextHong = "1000"
		case modal.UserType_NORMAL:
			nextHong = "100000"
		}

		return nextHong
	}

	return convert.Int642str(convert.Str2Int64(user.Hong) + 1)
}

func (this *UserDefault) GetUser(username string) (*modal.User, error) {
	if username == "" {
		return nil, errors.New(g.INVALID_PARAM_ERROR)
	}

	user := &modal.User{}
	if err := this.db.Where(&modal.User{Hong: username}).Or(&modal.User{Tel: username}).First(user).Error; err != nil {
		return nil, errors.New(err)
	}

	return user, nil
}

func (this *UserDefault) SyncDB() error {
	tx := this.db.Begin()
	{
		user := &modal.User{}
		if tx.HasTable(user) {
			if err := tx.AutoMigrate(user).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.CreateTable(user).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	return nil
}
