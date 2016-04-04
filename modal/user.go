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
package modal

import (
	"github.com/jinzhu/gorm"
	"github.com/pquerna/ffjson/ffjson"
)

type UserType int

const (
	UserType_TEST UserType = iota + 1
	UserType_NORMAL
)

type UserGender int

const (
	UserGender_NONE UserGender = iota
	UserGender_MAN
	UserGender_WOMAN
)

type User struct {
	Hong     string `gorm:"unique;not null"`
	Nick     string
	Tel      string `gorm:"unique;not null"`
	Password string
	Type     UserType `gorm:"index;not null"`
	Avatar   string
	Age      int
	Gender   UserGender `gorm:"index;not null"`
	gorm.Model
}

func (this *User) MarshalJSON() ([]byte, error) {
	type InnterUser struct {
		Hong       string     `json:"hong"`
		Nick       string     `json:"nick,omitempty"`
		Tel        string     `json:"tel"`
		UserType   UserType   `json:"userType"`
		Avatar     string     `json:"avatar,omitempty"`
		Age        int        `json:"age,omitempty"`
		UserGender UserGender `json:"userGender"`
	}

	return ffjson.Marshal(&InnterUser{
		Hong:       this.Hong,
		Nick:       this.Nick,
		Tel:        this.Tel,
		UserType:   this.Type,
		Avatar:     this.Avatar,
		Age:        this.Age,
		UserGender: this.Gender,
	})
}
