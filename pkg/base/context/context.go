package context

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

// newContext 创建一个新的上下文，将指定的键值对存储为持久化值
// 参数：
//   - ctx: 原始上下文
//   - key: 要存储的键
//   - value: 要存储的值
//
// 返回值：
//   - context.Context: 包含新键值对的上下文
func newContext(ctx context.Context, key string, value string) context.Context {
	return metainfo.WithPersistentValue(ctx, key, value) // 使用metainfo将键值对存储为持久化值
}

// fromContext 从上下文中获取指定键的值
// 参数：
//   - ctx: 包含键值对的上下文
//   - key: 要获取的键
//
// 返回值：
//   - string: 键对应的值
//   - bool: 是否成功获取到值
func fromContext(ctx context.Context, key string) (string, bool) {
	return metainfo.GetPersistentValue(ctx, key) // 使用metainfo从上下文中获取键对应的值
}
