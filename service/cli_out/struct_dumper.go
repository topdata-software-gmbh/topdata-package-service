package cli_out

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/topdata-software-gmbh/topdata-package-service/model"
	"os"
)

func DumpGitBranchInfoList(gitBranchInfoList model.GitBranchInfoList) {

	gitBranchInfoList.SortByPackageVersionAsc()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Git Branch", "Package Version", "Shopware Version", "Commit Id"})

	for _, b := range gitBranchInfoList.GitBranchInfos {
		t.AppendRow([]interface{}{b.Name, b.PackageVersion, b.ShopwareVersion, b.CommitId})
	}

	t.Render()
}

func DumpPkgsTable(pkgInfos []model.PkgInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Package Name", "Release Branch Names", "Other Branch Names" /*, "URL"*/})

	for _, p := range pkgInfos {
		t.AppendRow([]interface{}{p.Name, p.ReleaseBranchNames, p.OtherBranchNames /*, p.URL*/})
	}

	t.Render()
}
