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

var cachedInstallExtractDir string = ""
var cachedInstallDir string = ""
var cachedLaunchCommand []string = []string{}

func GetInstallExtractDir() string {
	if cachedInstallExtractDir == "" {
		if dir, err := os.Getwd(); err == nil {
			if runtime.GOOS == "windows" {
				cachedInstallExtractDir = path.Join(dir, WinTmpExtractDir)
			} else {
				cachedInstallExtractDir = path.Join(dir, TmpExtractDir)
			}
		} else {
			cachedInstallExtractDir = WinTmpExtractDir
		}
	}

	return cachedInstallExtractDir
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
	if cachedInstallDir == "" {
		if filepath.IsAbs(installPath) {
			cachedInstallDir = installPath
		} else if runtime.GOOS == "windows" {
			if home, err := os.UserHomeDir(); err == nil {
				cachedInstallDir = path.Join(filepath.VolumeName(home), "Program Files", installPath)
			}
		} else {
			if home, err := os.UserHomeDir(); err == nil {
				cachedInstallDir = filepath.Join(home, installPath)
			}
		}

		// if it is still empty, default to as is.
		if cachedInstallDir == "" {
			cachedInstallDir = installPath
		}
	}

	return cachedInstallDir
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info != nil && !info.IsDir()
}

func GetAppDir() string {
	if cachedAppDir == "" {
		if ex, err := os.Executable(); err == nil {
			cachedAppDir = filepath.Dir(ex)
		} else {
			log.Fatalln("Failed to get application directory")
		}
	}

	return cachedAppDir
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
	if len(cachedLaunchCommand) == 0 {
		launchFile := GetLaunchScript(installPath)
		if data, err := os.ReadFile(launchFile); err == nil {
			_ = json.Unmarshal(data, &cachedLaunchCommand)
		}
	}

	return cachedLaunchCommand
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
