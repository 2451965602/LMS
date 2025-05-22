package context

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
)

func newContext(ctx context.Context, key string, value string) context.Context {
	return metainfo.WithPersistentValue(ctx, key, value)
}

func fromContext(ctx context.Context, key string) (string, bool) {
	return metainfo.GetPersistentValue(ctx, key)
}
