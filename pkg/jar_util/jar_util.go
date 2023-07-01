package jar_util

import (
	"crypto/sha1"
	"fmt"
	"os"
)

// JarSha1 计算给定路径的jar包的sha1
func JarSha1(jarPath string) (string, error) {

	// 读取jar包
	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return "", err
	}

	// 计算sha1
	h := sha1.New()
	h.Write([]byte(jarBytes))
	bs := h.Sum(nil)
	jarSha1 := fmt.Sprintf("%x", bs)
	return jarSha1, nil
}
