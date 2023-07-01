package crawler

import (
	"context"
	"github.com/java-sec/jar-finder/pkg/jar_util"
	"github.com/java-sec/jar-finder/pkg/models"

	"github.com/scagogogo/sonatype-central-sdk/pkg/api"
	"io/fs"
	"path/filepath"
	"strings"
)

// FindByDirectory 把目录下的所有jar包寻根溯源
func FindByDirectory(ctx context.Context, directory string, handler models.FileHandler) error {
	return filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := &models.File{
			Path: path,
			Info: info,
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".jar") {
			// 计算jar包的sha1
			jarSha1, err := jar_util.JarSha1(info.Name())
			if err != nil {
				return err
			}
			versions, err := api.SearchBySha1(ctx, jarSha1, 1)
			if err != nil {
				return err
			}
			if len(versions) > 0 {
				file.Version = versions[0]
			}
		}

		return handler(ctx, file)
	})
}

//
//func FindPomXmlDependencyByDirectory(ctx context.Context, directory string) {
//	FindByDirectory(ctx, directory, func(ctx context.Context, path string, info fs.FileInfo, version *response.Version) error {
//		if version != nil {
//			// 说明是在Maven中央仓库存在，则构造Maven的依赖xml语句
//		} else {
//			// 说明在中央仓库不存在，则使用本地依赖
//			t := `                        <!-- https://mvnrepository.com/artifact/org.projectlombok/lombok -->
//                        <dependency>
//                            <groupId>{{groupId}}</groupId>
//                            <artifactId>{{artifactId}}</artifactId>
//                            <version>{{version}}</version>
//        //                     <scope>system</scope>
//        //                    <systemPath>${jarFilePath}</systemPath>
//                        </dependency>`
//
//		}
//		return nil
//	})
//}

//// FindMavenRepositoryNotExists 找到给定目录下的jar包中没有托管在Maven中央仓库的
//func FindMavenRepositoryNotExists(ctx context.Context, directory string) ([]*models.File, error) {
//	files := make([]*models.File, 0)
//	FindByDirectory(ctx, directory, func(ctx context.Context, file *models.File) error {
//		//if file.IsDirectory() {
//		//	return nil
//		//} else if !file.IsJarFile() {
//		//	color.Black("File %s is not jar file, so ignored it.", file.Info.Name())
//		//} else if file.ExistsMaven() {
//		//	color.Green("File %s exists in maven repo, groupId = %s, artifactId = %s, version = %s",
//		//		file.Info.Name(), file.Version.GroupId, file.Version.ArtifactId, file.Version.Version)
//		//} else {
//		//	color.Red("File %s not exists maven repo, it may be custom jar.", file.Info.Name())
//		//}
//		files = append(files, file)
//		return nil
//	})
//	return files, nil
//}
