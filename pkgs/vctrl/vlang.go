package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver/v3"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/downloader"
	"github.com/moqsien/gvc/pkgs/utils"
)

type Vlang struct {
	Conf *config.GVConfig
	*downloader.Downloader
	env *utils.EnvsHandler
}

func NewVlang() (vl *Vlang) {
	vl = &Vlang{
		Conf:       config.New(),
		Downloader: &downloader.Downloader{},
		env:        utils.NewEnvsHandler(),
	}
	return
}

func (that *Vlang) download(force bool) string {
	that.Url = that.Conf.Vlang.VlangGiteeUrls[runtime.GOOS]
	if that.Url != "" {
		fpath := filepath.Join(config.VlangFilesDir, "vlang.zip")
		if ok, _ := utils.PathIsExist(fpath); !ok || force {
			if size := that.GetFile(fpath, os.O_CREATE|os.O_WRONLY, 0644); size > 0 {
				return fpath
			} else {
				os.RemoveAll(fpath)
			}
		}
	}
	return ""
}

func (that *Vlang) Install(force bool) {
	zipFilePath := that.download(force)
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok && !force {
		fmt.Println("Vlang is already installed.")
		return
	} else {
		os.RemoveAll(config.VlangRootDir)
	}
	if err := archiver.Unarchive(zipFilePath, config.VlangFilesDir); err != nil {
		os.RemoveAll(config.VlangRootDir)
		os.RemoveAll(zipFilePath)
		fmt.Println("[Unarchive failed] ", err)
		return
	}
	if ok, _ := utils.PathIsExist(config.VlangRootDir); ok {
		that.CheckAndInitEnv()
	}
}

func (that *Vlang) CheckAndInitEnv() {
	if runtime.GOOS != utils.Windows {
		vlangEnv := fmt.Sprintf(utils.VlangEnv, config.VlangRootDir)
		that.env.UpdateSub(utils.SUB_VLANG, vlangEnv)
	} else {
		envList := map[string]string{
			"PATH": config.VlangRootDir,
		}
		that.env.SetEnvForWin(envList)
	}
}