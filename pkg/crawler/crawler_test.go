package crawler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindByDirectory(t *testing.T) {
	//err := FindByDirectory(context.Background(), "C:\\Program Files\\SmartGit\\lib", func(ctx context.Context, path string, info fs.FileInfo, version *response.Version) error {
	//	fmt.Println(path)
	//	fmt.Println(version)
	//	return nil
	//})
	//assert.Nil(t, err)
}

func TestFindMavenRepositoryNotExists(t *testing.T) {
	err := FindMavenRepositoryNotExists(context.Background(), "C:\\Program Files\\SmartGit\\lib")
	assert.Nil(t, err)
}
