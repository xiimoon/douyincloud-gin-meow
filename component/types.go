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
)

type AdTracking interface {
	// SetClickIDIfAbsent 设置 openId 对应的 clickId，如果 openId 已经存在则不改变
	SetClickIDIfAbsent(ctx context.Context, openId, clickId string) error

	// GetClickIDByOpenID 读取 openId 对应的 clickId
	GetClickIDByOpenID(ctx context.Context, openId string) (string, error)

	// IncImpression 增加 openId 对应的 impression 值
	IncImpression(ctx context.Context, openId string) error

	// GetImpressionByOpenID 读取 openId 对应的 impression 值
	GetImpressionByOpenID(ctx context.Context, openId string) (int, error)
}

// const Mongo = "mongodb"
const Redis = "redis"

var (
	//mongoHelloWorld *mongoComponent
	redisHelloWorld *redisComponent
)

// GetComponent 通过传入的component的名称返回实现了HelloWorldComponent接口的component
func GetComponent() (AdTracking, error) {
	return redisHelloWorld, nil
}

func InitComponents() {
	//mongoHelloWorld = NewMongoComponent()
	redisHelloWorld = NewRedisComponent()

	// ctx := context.TODO()
	// // err := mongoHelloWorld.SetName(ctx, "name", "mongodb")
	// // if err != nil {
	// // 	panic(err)
	// // }
	// err := redisHelloWorld.SetName(ctx, "name", "redis")
	// if err != nil {
	// 	panic(err)
	// }
}
