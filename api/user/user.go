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
package user

import (
	"github.com/globalways/common"
	"github.com/globalways/common/resp"
	"github.com/globalways/dvip/modal"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"regexp"
	"github.com/globalways/common/utils"
)

type UserNewRequestParam struct {
	Type     modal.UserType `form:"userType" json:"userType" binding:"required"`
	Tel      string         `form:"tel" json:"tel" binding:"required"`
	Password string         `form:"password" json:"password" binding:"required"`
}

func (this UserNewRequestParam) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	if matched, err := regexp.MatchString(common.Regexp_Mobile, this.Tel); !matched || err != nil {
		errors = append(errors, binding.Error{
			FieldNames:     []string{"tel"},
			Classification: "invalidParamError",
			Message:        "the giving tel is invalid.",
		})
	}

	return errors
}

func NewUser(reqParam UserNewRequestParam, r render.Render, req *http.Request) {
	newUser, err := userAdapter.NewUser(reqParam.Type, reqParam.Tel, reqParam.Password)
	if err != nil {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp(err.Error()))
		return
	}

	r.JSON(http.StatusCreated, resp.NewAPIResp(newUser))
}

func GetUser(params martini.Params, r render.Render) {
	username := params["username"]
	user, err := userAdapter.GetUser(username)
	if err != nil {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp(err.Error()))
		return
	}

	r.JSON(http.StatusOK, resp.NewAPIResp(user))
}

type UserLoginRequestParam struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func UserLogin(reqParam UserLoginRequestParam, r render.Render) {
	user, err := userAdapter.GetUser(reqParam.Username)
	if err != nil {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp(err.Error()))
		return
	}

	if !utils.DecodePassword(user.Password, reqParam.Password) {
		r.JSON(http.StatusBadRequest, resp.NewAPIErrorResp("user password is wrong"))
		return
	}

	r.JSON(http.StatusOK, )
}