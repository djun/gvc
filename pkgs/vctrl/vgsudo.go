package vctrl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	tui "github.com/moqsien/goutils/pkgs/gtui"
	"github.com/moqsien/goutils/pkgs/request"
	config "github.com/moqsien/gvc/pkgs/confs"
	"github.com/moqsien/gvc/pkgs/utils"
)

type VGSudo struct {
	Conf    *config.GVConfig
	fetcher *request.Fetcher
	env     *utils.EnvsHandler
	checker *SumChecker
}

func NewGSudo() (gs *VGSudo) {
	gs = &VGSudo{
		Conf:    config.New(),
		fetcher: request.NewFetcher(),
		env:     utils.NewEnvsHandler(),
	}
	gs.checker = NewSumChecker(gs.Conf)
	gs.env.SetWinWorkDir(config.GVCWorkDir)
	return
}

func (that *VGSudo) Install(force bool) {
	if runtime.GOOS != utils.Windows {
		return
	}
	that.fetcher.Url = that.Conf.GSudo.GitlabUrl
	if that.fetcher.Url != "" {
		that.fetcher.Timeout = 20 * time.Minute
		that.fetcher.SetThreadNum(2)
		fPath := filepath.Join(config.GsudoFilePath, "gsudo.zip")
		dstDir := filepath.Join(config.GsudoFilePath, "gsudo")
		if err := that.fetcher.DownloadAndDecompress(fPath, dstDir, force); err == nil {
			that.CheckAndInitEnv(dstDir)
			tui.PrintSuccess(fPath)
		} else {
			os.RemoveAll(fPath)
			os.RemoveAll(dstDir)
			tui.PrintError(err)
		}
	}
}

func (that *VGSudo) CheckAndInitEnv(dstDir string) {
	binPath := that.GetBinPath(dstDir)
	if binPath == "" {
		return
	}
	if runtime.GOOS != utils.Windows {
		protoEnv := fmt.Sprintf(utils.ProtoEnv, binPath)
		that.env.UpdateSub(utils.SUB_PROTOC, protoEnv)
	} else {
		envList := map[string]string{
			"PATH": binPath,
		}
		that.env.SetEnvForWin(envList)
	}
}

func (that *VGSudo) GetBinPath(dstDir string) string {
	var binPath string
	if dirList, err := os.ReadDir(dstDir); err == nil {
		for _, d := range dirList {
			if d.IsDir() && d.Name() == "x64" && runtime.GOARCH == "amd64" {
				binPath = filepath.Join(dstDir, d.Name())
				break
			}
			if d.IsDir() && d.Name() == "arm64" && runtime.GOARCH == "arm64" {
				binPath = filepath.Join(dstDir, d.Name())
				break
			}
		}
	}
	return binPath
}