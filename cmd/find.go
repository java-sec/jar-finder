package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	if_expression "github.com/golang-infrastructure/go-if-expression"
	"github.com/java-sec/jar-finder/pkg/jar_util"
	"github.com/java-sec/jar-finder/pkg/models"
	"github.com/java-sec/jar-finder/pkg/pom_util"
	"github.com/olekukonko/tablewriter"
	"github.com/scagogogo/sonatype-central-sdk/pkg/api"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (

	// 要查询的jar所在的目录
	dir string

	// 要查询的jar文件的路径
	jar string

	// 查询结果输出为pom.xml文件的路径
	pom string
)

func init() {
	// 指定要扫描的目录
	custom.PersistentFlags().StringVarP(&dir, "dir", "d", "./", "Specify the directory where the jar packages you want to look in the maven central repository are located")
	custom.PersistentFlags().StringVarP(&jar, "jar", "j", "", "Specify the path to the jar packages you want to find in Maven's central repository")
	custom.PersistentFlags().StringVarP(&jar, "pom", "p", "", "Output the pom.xml to path, example: pom.xml")
	rootCmd.AddCommand(custom)
}

var custom = &cobra.Command{
	Use:   "find",
	Short: "Find jar in maven central repository",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Find(cmd.Context())
	},
}

func Find(ctx context.Context) error {

	// 参数检查
	if dir == "" && jar == "" {
		tips := "the dir and jar parameters specify at least one"
		return errors.New(tips)
	}

	// 按目录查找，如果指定了目录的话
	if dir != "" {
		err := findByDirectory(ctx, dir)
		if err != nil {
			return err
		}
	}

	if jar != "" {
		return findByJar(ctx, jar)
	}

	return nil
}

// 查找指定的jar包
func findByJar(ctx context.Context, jarPath string) error {

	stat, err := os.Stat(jarPath)
	if err != nil {
		if os.IsNotExist(err) {
			color.Red("File %s not exists!", jarPath)
		}
		return err
	}

	jarSha1, err := jar_util.JarSha1(jarPath)
	if err != nil {
		color.Red("File %s sha1 error", jarPath)
		return err
	}
	bySha1, err := api.SearchBySha1(ctx, jarSha1, 1000)
	if err != nil {
		color.Red("File %s query error", jarPath)
		return err
	}
	files := make([]*models.File, 0)
	for _, x := range bySha1 {
		file := &models.File{
			Path:    jarPath,
			Info:    stat,
			Version: x,
		}
		files = append(files, file)
	}
	renderTable(files)
	return nil
}

// 查找指定目录下的jar包
func findByDirectory(ctx context.Context, dir string) error {

	files := make([]*models.File, 0)
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			//color.Red("File %s read info error: %s", info.Name(), err.Error())
			return err
		} else if info.IsDir() {
			// 忽略目录
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".jar") {
			color.Red("File %s not jar, ignored it.", info.Name())
			return nil
		}

		// 计算jar包的sha1
		jarSha1, err := jar_util.JarSha1(path)
		if err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("File %s sha1 %s, begin search in maven central...", info.Name(), jarSha1))

		versions, err := api.SearchBySha1(ctx, jarSha1, 1)
		if err != nil {
			return err
		}
		file := &models.File{
			Path: path,
			Info: info,
		}
		if len(versions) > 0 {
			file.Version = versions[0]
			v := file.Version
			color.Green("File %s find version in maven central, groupId %s, artifactId %s, version %s", info.Name(), v.GroupId, v.ArtifactId, v.Version)
		} else {
			color.Red("File %s not found in maven central", info.Name())
		}
		files = append(files, file)
		return nil
	})
	if err != nil {
		return err
	}

	if len(files) != 0 {
		renderTable(files)

		if pom != "" {
			fmt.Println(fmt.Sprintf("Result save to %s", pom))
			return OutputPomXMl(pom, files)
		}
	}

	return nil
}

func renderTable(files []*models.File) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetHeader([]string{"File", "In Maven Repo"})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER})
	for _, f := range files {
		table.Append([]string{
			f.Info.Name(),
			if_expression.Return(f.ExistsMaven(), color.GreenString("Y"), color.RedString("N")),
		})
	}
	table.Render()
}

// OutputPomXMl 把查询结果输出为pom.xml文件
func OutputPomXMl(pomPath string, files []*models.File) error {

	err := os.MkdirAll(filepath.Dir(pomPath), os.ModePerm)
	if err != nil {
		return err
	}

	pomXMl := pom_util.BuildPomXml(files)
	file, err := os.OpenFile(pom, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(pomXMl))
	return err
}
