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
package app

import (
	"github.com/globalways/common/resp"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
)

type AppNewRequestParam struct {
	Name string `form:"name" json:"name" binding:"required"`
	Desc string `form:"desc" json:"desc"`
}

func NewApp(reqParam AppNewRequestParam, r render.Render, req *http.Request) {
	newApp, err := appAdapter.NewApp(reqParam.Name, reqParam.Desc)
	if err != nil {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp(err.Error()))
		return
	}

	r.JSON(http.StatusCreated, resp.NewAPIResp(newApp))
}

func GetApp(params martini.Params, r render.Render) {
	appKey := params["appKey"]
	app, err := appAdapter.GetApp(appKey)
	if err != nil {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp(err.Error()))
		return
	}

	r.JSON(http.StatusOK, resp.NewAPIResp(app))
}
