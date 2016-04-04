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
package oauth2

import (
	"github.com/RangelReale/osin"
	"github.com/globalways/hong/g"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var (
	server *osin.Server
)

func InitServer() {
	cfg := g.Config()
	apicfg := cfg.API

	db, err := gorm.Open("mysql", apicfg.AuthStorageAddr)
	if err != nil {
		log.Fatalf("open database error: %v", err)
	}
	if err := db.DB().Ping(); err != nil {
		log.Fatalf("ping to database error: %v", err)
	}

	db.DB().SetMaxIdleConns(apicfg.AuthStorageMaxIdle)
	db.DB().SetMaxOpenConns(apicfg.AuthStorageMaxConn)

	// Disable table name's pluralization
	db.SingularTable(true)
	if cfg.Debug {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	storage := NewDefaultStorage(db)
	if err := storage.SyncDB(); err != nil {
		log.Fatalf("sync database error: %v", err)
	}

	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	sconfig.AllowClientSecretInParams = true

	server = osin.NewServer(sconfig, storage)
}
