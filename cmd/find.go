package cmd

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/fatih/color"
	if_expression "github.com/golang-infrastructure/go-if-expression"
	"github.com/java-sec/jar-finder/pkg/crawler"
	"github.com/java-sec/jar-finder/pkg/models"
	"github.com/olekukonko/tablewriter"
	"github.com/scagogogo/sonatype-central-sdk/pkg/api"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	dir string
	jar string
)

func init() {
	// 指定要扫描的目录
	custom.PersistentFlags().StringVarP(&dir, "dir", "d", "./", "Specify the directory where the jar packages you want to look in the maven central repository are located")
	custom.PersistentFlags().StringVarP(&jar, "jar", "", "", "Specify the path to the jar packages you want to find in Maven's central repository")
	rootCmd.AddCommand(custom)
}

var custom = &cobra.Command{
	Use:   "find",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Find(cmd.Context())
	},
}

func Find(ctx context.Context) error {

	// 参数检查
	if dir == "" && jar == "" {
		return errors.New("the dir and jar parameters specify at least one")
	}

	// 按目录查找，如果指定了目录的话
	if dir != "" {
		err := findByDirectory(ctx, dir)
		if err != nil {
			return err
		}
	}

	if jar != "" {

	}

}

// 查找指定的jar包
func findByJar(ctx context.Context, jar string) error {

}

// 查找指定目录下的jar包
func findByDirectory(ctx context.Context, dir string) error {

	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			color.Red("File %s read info error: %s", info.Name(), err.Error())
			return err
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".jar") {
			color.White("File %s not jar, ignored it.", info.Name())
			return nil
		}

		// 计算jar包的sha1
		jarBytes, err := os.ReadFile(path)
		if err != nil {
			color.Red("File %s, read error: %s", info.Name(), err.Error())
			return err
		}
		h := sha1.New()
		h.Write([]byte(jarBytes))
		bs := h.Sum(nil)
		jarSha1 := fmt.Sprintf("%x", bs)
		versions, err := api.SearchBySha1(ctx, jarSha1, 1)
		if err != nil {
			color.Red("")
			return err
		}
		if len(versions) > 0 {
			file.Version = versions[0]
		}
	})

	err := crawler.FindByDirectory(ctx, dir, func(ctx context.Context, file *models.File) error {

	})
	return
}

func renderTable(files []*models.File) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetHeader([]string{"File", "In Maven Repo"})
	for _, f := range files {
		table.Append([]string{
			f.Info.Name(),
			if_expression.Return(f.ExistsMaven(), "Y", "N"),
		})
	}
	table.Render()
}
