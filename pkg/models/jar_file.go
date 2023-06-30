package models

import (
	"github.com/scagogogo/sonatype-central-sdk/pkg/response"
	"io/fs"
	"strings"
)

// File 表示一个文件
type File struct {

	// 文件的路径
	Path string

	// 文件系统信息
	Info fs.FileInfo

	// 在Maven仓库查找到的信息
	Version *response.Version
}

// IsJarFile 是否是Jar文件
func (x *File) IsJarFile() bool {
	return strings.HasSuffix(strings.ToLower(x.Path), ".jar")
}

// ExistsMaven 是否在Maven中央仓库存在
func (x *File) ExistsMaven() bool {
	return x.Version != nil
}

// IsDirectory 此文件是否是目录
func (x *File) IsDirectory() bool {
	return x.Info.IsDir()
}
