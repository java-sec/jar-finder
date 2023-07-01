package pom_util

import (
	"fmt"
	"github.com/java-sec/jar-finder/pkg/models"
	"strings"
)

// BuildPomXml 为依赖构建pom.xml文件
func BuildPomXml(files []*models.File) string {

	dependencies := strings.Builder{}
	// 对依赖按照本地还是云端分组，这样的话产出的pom.xml文件会比较容易维护一些
	cloudDependencies := make([]*models.File, 0)
	localDependencies := make([]*models.File, 0)
	for _, file := range files {
		if file.Version == nil {
			localDependencies = append(localDependencies, file)
		} else {
			cloudDependencies = append(cloudDependencies, file)
		}
	}
	for _, file := range cloudDependencies {
		dependencies.WriteString(BuildDependencyXml(file))
	}

	dependencies.WriteString(`

<-- Below are local dependencies -->

`)

	for _, file := range localDependencies {
		dependencies.WriteString(BuildDependencyXml(file))
	}

	template := `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>


    <dependencies>
%s
    </dependencies>


</project>`
	return fmt.Sprintf(template, dependencies.String())
}

// BuildDependencyXml 根据依赖文件构建pom.xml文件的dependency块
func BuildDependencyXml(file *models.File) string {
	if file.Version != nil {
		// 此文件在中央仓库存在，则直接云端依赖即可
		v := file.Version
		template := `
        <!-- %s -->
        <dependency>
            <groupId>%s</groupId>
            <artifactId>%s</artifactId>
            <version>%s</version>
        </dependency>
`
		return fmt.Sprintf(template, BuildMavenCentralUrl(v.GroupId, v.ArtifactId, v.Version), v.GroupId, v.ArtifactId, v.Version)
	} else {
		// 此文件在中央仓库不存在，则使用本地依赖
		template := `
        <dependency>
            <groupId>%s</groupId>
            <artifactId>%s</artifactId>
            <version>%s</version>
            <scope>system</scope>
            <systemPath>%s</systemPath>
        </dependency>
`
		return fmt.Sprintf(template, file.Path, file.Path, file.Path, file.Path)
	}
}

// BuildMavenCentralUrl 构建Maven中央仓库的URL
func BuildMavenCentralUrl(groupId, artifactId, version string) string {
	return fmt.Sprintf("https://mvnrepository.com/artifact/%s/%s/%s", groupId, artifactId, version)
}
