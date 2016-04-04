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

// https://api.hongid.com/v1/oauth2/authorize

package api

import (
	"github.com/globalways/dvip/api/oauth2"
	"github.com/globalways/dvip/api/user"
	"github.com/globalways/dvip/g"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"time"
	"github.com/globalways/dvip/api/app"
)

func Start() {
	cfg := g.Config()
	apicfg := cfg.API

	// init oauth2
	oauth2.InitServer()
	// init user
	user.InitUser()
	// init app
	app.InitApp()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	if cfg.Debug {
		martini.Env = martini.Dev
	} else {
		martini.Env = martini.Prod
	}

	m.Group("/v1", func(r martini.Router) {
		// oauth2 router
		r.Group("/oauth2", func(r martini.Router) {
			r.Post("/authorize", oauth2.Authorize)
			r.Post("/token", oauth2.Token)
			r.Post("/info", oauth2.Info)
		})

		// user router
		r.Group("/users", func(r martini.Router) {
			r.Post("", binding.Bind(user.UserNewRequestParam{}), user.NewUser)
			r.Get("/h/:username", user.GetUser)
			r.Get("/t/:username", user.GetUser)
			r.Post("/login", )
		})

		// app router
		r.Group("/apps", func(r martini.Router) {
			r.Post("", binding.Bind(app.AppNewRequestParam{}), app.NewApp)
			r.Get("/:appKey", app.GetApp)
		})
	})

	m.RunOnAddr(apicfg.Addr)
}
