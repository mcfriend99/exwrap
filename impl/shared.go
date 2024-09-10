package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var CachedInstallExtractDir string = ""
var CachedInstallDir string = ""
var CachedLaunchCommand []string = []string{}

func GetInstallExtractDir() string {
	if CachedInstallExtractDir == "" {
		if dir, err := os.Getwd(); err == nil {
			if runtime.GOOS == "windows" {
				CachedInstallExtractDir = path.Join(dir, WinTmpExtractDir)
			} else {
				CachedInstallExtractDir = path.Join(dir, TmpExtractDir)
			}
		} else {
			CachedInstallExtractDir = WinTmpExtractDir
		}
	}

	return CachedInstallExtractDir
}

func GetSelfInstallExtractFile() string {
	if exe, err := os.Executable(); err == nil {
		return path.Join(GetInstallExtractDir(), path.Base(exe))
	}

	return ""
}

func GetInstallExtractFile() string {
	return path.Join(GetInstallExtractDir(), ExtractDstFile)
}

func GetInstallDir(installPath string) string {
	if CachedInstallDir == "" {
		if filepath.IsAbs(installPath) {
			CachedInstallDir = installPath
		} else if runtime.GOOS == "windows" {
			if home, err := os.UserHomeDir(); err == nil {
				CachedInstallDir = path.Join(filepath.VolumeName(home), "Program Files", installPath)
			}
		} else {
			if home, err := os.UserHomeDir(); err == nil {
				CachedInstallDir = filepath.Join(home, installPath)
			}
		}

		// if it is still empty, default to as is.
		if CachedInstallDir == "" {
			CachedInstallDir = installPath
		}
	}

	return CachedInstallDir
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && !info.IsDir()
}

func GetAppDir() string {
	if CachedAppDir == "" {
		if ex, err := os.Executable(); err == nil {
			CachedAppDir = filepath.Dir(ex)
		} else {
			log.Fatalln("Failed to get application directory")
		}
	}

	return CachedAppDir
}

func GetAppName() string {
	if ex, err := os.Executable(); err == nil {
		return filepath.Base(ex)
	}

	log.Fatalln("Failed to get application directory")
	return ""
}

func GetLaunchScript(installPath string) string {
	return path.Join(GetInstallDir(installPath), fmt.Sprintf("%s.launch", GetAppName()))
}

func GetLaunchCommand(installPath string) []string {
	if len(CachedLaunchCommand) == 0 {
		launchFile := GetLaunchScript(installPath)
		if data, err := os.ReadFile(launchFile); err == nil {
			_ = json.Unmarshal(data, &CachedLaunchCommand)
		}
	}

	return CachedLaunchCommand
}

func GetAbsoluteCommandProgram(cmd string, isDarwin bool) string {
	if filepath.IsAbs(cmd) {
		return cmd
	}

	if isDarwin {
		return path.Join(GetAppDir(), "../Resources", cmd)
	}

	return path.Join(GetAppDir(), cmd)
}
