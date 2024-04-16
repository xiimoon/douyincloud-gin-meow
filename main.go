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
package main

import (
	"douyincloud-gin-meow/component"
	"douyincloud-gin-meow/service"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func scheduleTask() {
	c := cron.New(cron.WithSeconds()) // 创建一个新的cron实例，注意我们不一定需要到秒的精度，但这里为了通用性包含了

	// 每天早上7点执行的任务
	_, err := c.AddFunc("0 0 7 * * *", func() {
		fmt.Println("Executing daily morning task at", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Printf("Error scheduling daily morning task: %s\n", err)
		return
	}

	// 每小时的05分执行的任务
	_, err = c.AddFunc("0 5 * * * *", func() {
		fmt.Println("Executing hourly task at", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err != nil {
		fmt.Printf("Error scheduling hourly task: %s\n", err)
		return
	}

	c.Start() // 启动cron调度器
}

func main() {
	component.InitComponents()

	//scheduleTask() // 设置定时任务

	r := gin.Default()

	r.GET("/api/impression", service.Impression)
	r.POST("/api/active", service.Active)

	r.Run(":8000")
}
