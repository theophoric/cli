package plugin_installer

import (
	"github.com/theophoric/cf-cli/cf/actors/plugin_repo"
	"github.com/theophoric/cf-cli/cf/models"
	"github.com/theophoric/cf-cli/cf/terminal"
	"github.com/theophoric/cf-cli/downloader"
	"github.com/theophoric/cf-cli/utils"
)

type PluginInstaller interface {
	Install(inputSourceFilepath string) string
}

type PluginInstallerContext struct {
	Checksummer    utils.Sha1Checksum
	FileDownloader downloader.Downloader
	GetPluginRepos pluginReposFetcher
	PluginRepo     plugin_repo.PluginRepo
	RepoName       string
	Ui             terminal.UI
}

type pluginReposFetcher func() []models.PluginRepo

func NewPluginInstaller(context *PluginInstallerContext) (installer PluginInstaller) {
	pluginDownloader := &PluginDownloader{Ui: context.Ui, FileDownloader: context.FileDownloader}
	if context.RepoName == "" {
		installer = &PluginInstallerWithoutRepo{
			Ui:               context.Ui,
			PluginDownloader: pluginDownloader,
			RepoName:         context.RepoName,
		}
	} else {
		installer = &PluginInstallerWithRepo{
			Ui:               context.Ui,
			PluginDownloader: pluginDownloader,
			RepoName:         context.RepoName,
			Checksummer:      context.Checksummer,
			PluginRepo:       context.PluginRepo,
			GetPluginRepos:   context.GetPluginRepos,
		}
	}
	return installer
}
