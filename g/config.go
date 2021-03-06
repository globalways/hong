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
package g

import (
	"github.com/pquerna/ffjson/ffjson"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type APIConfig struct {
	Addr            string `json:"addr"`
	UserDSN         string `json:"userDSN"`
	AppDSN          string `json:"appDSN"`
	AuthStorageAddr string `json:"authStorageAddr"`
	MaxIdle         int    `json:"maxIdle"`
	MaxOpen         int    `json:"maxOpen"`
}

type GlobalConfig struct {
	Debug    bool       `json:"debug"`
	Logstash string     `json:"logstash"`
	API      *APIConfig `json:"api"`
}

var (
	config *GlobalConfig
	lock   = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) *GlobalConfig {
	if cfg == "" {
		log.Fatal("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalf("config file: %s is not existent. maybe you need `mv cfg.example.json cfg.json`", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalf("read config file: %s fail: %v", cfg, err)
	}

	var c GlobalConfig
	err = ffjson.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalf("parse config file: %s fail: %v", cfg, err)
	}

	lock.Lock()
	config = &c
	lock.Unlock()

	log.Printf("g.ParseConfig ok, file: %s", cfg)
	return config
}
