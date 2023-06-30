package models

import (
	"context"
)

// FileHandler 单个文件应该如何处理
type FileHandler func(ctx context.Context, file *File) error
