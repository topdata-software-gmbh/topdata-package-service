package factory

import (
	"github.com/fatih/color"
	"github.com/topdata-software-gmbh/topdata-package-service/config"
	git_cli_wrapper2 "github.com/topdata-software-gmbh/topdata-package-service/git_cli_wrapper"
	"github.com/topdata-software-gmbh/topdata-package-service/model"
	"github.com/topdata-software-gmbh/topdata-package-service/serializers"
	"github.com/topdata-software-gmbh/topdata-package-service/util"
)

// NewPkgInfo creates a new PkgInfo object (aka constructor)
func NewPkgInfo(pkgConfig model.PkgConfig) model.PkgInfo {

	color.Blue("//////////////// NewPkgInfo: %s", pkgConfig.Name)

	git_cli_wrapper2.RefreshRepo(pkgConfig)
	branchNames := git_cli_wrapper2.GetRemoteBranchNames(pkgConfig)
	regex := `^(main|main-.*|release-.*)$` // TODO: the regex should be part of the service config or even pkgConfig

	return model.PkgInfo{
		PkgConfig: pkgConfig,
		// Name:               pkgConfig.Name,
		// URL:                pkgConfig.URL,
		ReleaseBranchNames: util.FilterStringSlicePositive(branchNames, regex),
		OtherBranchNames:   util.FilterStringSliceNegative(branchNames, regex),
	}

}

func NewPkgInfoList(pkgConfigList model.PkgConfigList) *model.PkgInfoList {
	pkgInfos := make([]model.PkgInfo, len(pkgConfigList.PkgConfigs))

	for i, pkgConfig := range pkgConfigList.PkgConfigs {
		pkgInfos[i] = NewPkgInfo(pkgConfig)
	}

	return &model.PkgInfoList{PkgInfos: pkgInfos}
}

func NewPkgInfoListCached(pkgConfigList model.PkgConfigList) *model.PkgInfoList {
	if util.FileExists(config.PathCacheFile) {
		return serializers.LoadPkgInfoList(config.PathCacheFile)
	} else {
		pkgInfoList := NewPkgInfoList(pkgConfigList)
		serializers.SavePkgInfoList(pkgInfoList, config.PathCacheFile)
		return pkgInfoList
	}

}
