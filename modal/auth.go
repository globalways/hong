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

import "time"

type Client struct {
	Id          string `gorm:"primary_key; not null"`
	Secret      string `gorm:"not null"`
	Extra       string `gorm:"not null"`
	RedirectUri string `gorm:"not null"`
}

type Authorize struct {
	Client      string    `gorm:"not null"`
	Code        string    `gorm:"primary_key; not null"`
	ExpiresIn   int32     `gorm:"not null"`
	Scope       string    `gorm:"not null"`
	RedirectUri string    `gorm:"not null"`
	State       string    `gorm:"not null"`
	Extra       string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
}

type Access struct {
	Client       string    `gorm:"not null"`
	Authorize    string    `gorm:"not null"`
	Previous     string    `gorm:"not null"`
	AccessToken  string    `gorm:"primary_key; not null"`
	RefreshToken string    `gorm:"not null"`
	ExpiresIn    int32     `gorm:"not null"`
	Scope        string    `gorm:"not null"`
	RedirectUri  string    `gorm:"not null"`
	Extra        string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
}

type Refresh struct {
	Token  string `gorm:"primary_key; not null"`
	Access string `gorm:"not null"`
}
