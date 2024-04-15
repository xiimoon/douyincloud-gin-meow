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
package component

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	//Redis地址
	redisAddr = ""
	//Redis用户名
	redisUserName = ""
	//Redis密码
	redisPassword = ""
)

type redisComponent struct {
	client *redis.Client
}

// SetClickIDIfAbsent 设置 openId 对应的 clickId，如果 openId 已经存在则不改变
func (r *redisComponent) SetClickIDIfAbsent(ctx context.Context, openId, clickId string) error {
	_, err := r.client.HSetNX(ctx, "open:"+openId, "clickId", clickId).Result()
	return err
}

// GetClickIDByOpenID 读取 openId 对应的 clickId
func (r *redisComponent) GetClickIDByOpenID(ctx context.Context, openId string) (string, error) {
	return r.client.HGet(ctx, "open:"+openId, "clickId").Result()
}

// IncImpression 增加 openId 对应的 impression 值
func (r *redisComponent) IncImpression(ctx context.Context, openId string) error {
	_, err := r.client.HIncrBy(ctx, "open:"+openId, "impression", 1).Result()
	return err
}

// GetImpressionByOpenID 读取 openId 对应的 impression 值
func (r *redisComponent) GetImpressionByOpenID(ctx context.Context, openId string) (int, error) {
	result, err := r.client.HGet(ctx, "open:"+openId, "impression").Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(result)
}

// NewRedisComponent 初始化一个实现了HelloWorldComponent接口的RedisComponent
func NewRedisComponent() *redisComponent {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: redisUserName,
		Password: redisPassword,
		DB:       0, // use default DB
	})
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		fmt.Printf("redisClient init error. err %s", err)
		panic(fmt.Sprintf("redis init failed. err %s\n", err))
	}
	return &redisComponent{
		client: rdb,
	}
}

// init 项目启动时会从环境变量中获取
func init() {
	redisAddr = os.Getenv("REDIS_ADDRESS")
	redisUserName = os.Getenv("REDIS_USERNAME")
	redisPassword = os.Getenv("REDIS_PASSWORD")
}
