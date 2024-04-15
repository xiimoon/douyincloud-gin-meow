/*
Copyright (year) Bytedance Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package service

import (
	"douyincloud-gin-meow/component"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ActiveData struct {
	ClickID string `json:"clickId" binding:"required"`
}

func Impression(ctx *gin.Context) {
	openId := ctx.GetHeader("X-TT-OPENID")
	if openId == "" {
		Failure(ctx, fmt.Errorf("X-TT-OPENID null"))
	}

	adtracker, err := component.GetComponent()
	if err != nil {
		Failure(ctx, fmt.Errorf("redis component not found"))
		return
	}
	err = adtracker.IncImpression(ctx, openId)
	if err != nil {
		Failure(ctx, err)
		return
	}
	Success(ctx, "")
}

func Active(ctx *gin.Context) {
	var req ActiveData
	err := ctx.Bind(&req)
	if err != nil {
		Failure(ctx, err)
		return
	}

	openId := ctx.GetHeader("X-TT-OPENID")
	if openId == "" {
		Failure(ctx, fmt.Errorf("X-TT-OPENID null"))
	}

	adtracker, err := component.GetComponent()
	if err != nil {
		Failure(ctx, fmt.Errorf("redis component not found"))
		return
	}

	_, err = adtracker.GetClickIDByOpenID(ctx, openId)
	if err != nil {
		if err.Error() == "redis: nil" {
			// new player
			err = adtracker.SetClickIDIfAbsent(ctx, openId, req.ClickID)
			if err != nil {
				Failure(ctx, err)
				return
			}

			// send conversion
		} else {
			Failure(ctx, err)
			return
		}
	}

	Success(ctx, "active")
}

func Failure(ctx *gin.Context, err error) {
	resp := &Resp{
		ErrNo:  -1,
		ErrMsg: err.Error(),
	}
	ctx.JSON(200, resp)
}

func Success(ctx *gin.Context, data string) {
	resp := &Resp{
		ErrNo:  0,
		ErrMsg: "success",
		Data:   data,
	}
	ctx.JSON(200, resp)
}

type Resp struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}
