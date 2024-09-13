package impl

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var CachedAppDir string = ""
var CachedBuildDir string = ""

func listFiles(root string) []string {
	files := make([]string, 0)

	if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		// make sure we resolve the path first for symbolic links before
		// we proceed.
		// this is where we ensure that only the actual files are being
		// included in the final artefact.
		path, info = resolveFile(path, info)
		if err == nil && !info.IsDir() {
			files = append(files, path)
		}

		return err
	}); err != nil {
		log.Fatalln("Failed to read root directory:", err.Error())
	}

	return files
}

func stringListContains(list []string, key string) bool {
	for _, x := range list {
		if x == key {
			return true
		}
	}

	return false
}

func hasMatchInList(list []string, key string) bool {
	for _, x := range list {
		if x == key || strings.HasPrefix(key, x) {
			return true
		}
	}

	return false
}

func getFileAbsPath(path string) (string, error) {
	if absPath, err := filepath.Abs(path); err == nil {
		return absPath, nil
	} else {
		return "", err
	}
}

func resolveFile(file string, info fs.FileInfo) (string, fs.FileInfo) {
	if info.Mode()&os.ModeSymlink != 0 {
		// is symlink
		if link, err := os.Readlink(file); err == nil {
			if !filepath.IsAbs(link) {
				link = path.Join(path.Dir(file), link)
			}

			if stat, err := os.Stat(link); err == nil {
				return link, stat
			}

			return link, info
		}
	}

	return file, info
}

func trimRoot(path string, root string) string {
	return strings.TrimLeft(strings.ReplaceAll(path, root, ""), "/\\")
}

func getPkgExeName(exeOs string, arch string) string {
	var filepath string
	if exeOs == "windows" {
		filepath = path.Join(GetAppDir(), "pkg", fmt.Sprintf("wrapper-windows-%s.exe", arch))
	} else {
		filepath = path.Join(GetAppDir(), "pkg", fmt.Sprintf("wrapper-%s-%s", exeOs, arch))
	}

	if !FileExists(filepath) {
		log.Fatalf("Unsupported packaging combination %s/%s. File %s cannot be located!", exeOs, arch, filepath)
	}

	return filepath
}

func getEmbedExeName(cmd CommandLine, config Config) string {
	var filepath string
	if config.TargetOs == "windows" {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s.exe", AppEmbedExeName))
	} else {
		filepath = path.Join(getBuildDir(cmd), AppEmbedExeName)
	}

	return filepath
}

func getTargetExeName(cmd CommandLine, config Config) string {
	var filepath string
	if config.TargetOs == "windows" {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s.exe", config.TargetName))
	} else {
		filepath = path.Join(getBuildDir(cmd), config.TargetName)
	}

	return filepath
}

func getTargetBaseName(cmd CommandLine, config Config) string {
	var filepath string
	if config.TargetOs == "windows" {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s-base.exe", config.TargetName))
	} else {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s-base", config.TargetName))
	}

	return filepath
}

func getTargetWrapperName(cmd CommandLine, config Config) string {
	var filepath string
	if config.TargetOs == "windows" {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s-wrap.exe", config.TargetName))
	} else {
		filepath = path.Join(getBuildDir(cmd), fmt.Sprintf("%s-wrap", config.TargetName))
	}

	return filepath
}

func getResourcesDirectory() string {
	return path.Join(GetAppDir(), "Resources")
}

func getPkgExeFromConfig(config Config) string {
	return getPkgExeName(config.TargetOs, config.TargetArch)
}

func getLaunchScriptTempName() string {
	return fmt.Sprintf("%s.json", EmbededLaunchScript)
}

func getLaunchScriptForDarwinApp(config Config) string {
	return fmt.Sprintf("%s.launch", config.TargetName)
}

func getSetupScriptTempName() string {
	return fmt.Sprintf("%s.json", EmbededSetupScript)
}

func getBuildLaunchScriptName(cmd CommandLine) string {
	return path.Join(getBuildDir(cmd), getLaunchScriptTempName())
}

func getBuildSetupScriptName(cmd CommandLine) string {
	return path.Join(getBuildDir(cmd), getSetupScriptTempName())
}

func getBuildDir(cmd CommandLine) string {
	if filepath.IsAbs(cmd.BuildDirectory) {
		return cmd.BuildDirectory
	}

	if CachedBuildDir == "" {
		if dir, err := os.Getwd(); err == nil {
			CachedBuildDir = path.Join(dir, cmd.BuildDirectory)
		}
		CachedBuildDir = cmd.BuildDirectory
	}

	return CachedBuildDir
}

func getTargetBuildArchive(config Config, cmd CommandLine) string {
	if config.TargetOs == "darwin" {
		return path.Join(getBuildDir(cmd), fmt.Sprintf("%s.app", config.TargetName))
	} else {
		return path.Join(getBuildDir(cmd), AppArchiveName)
	}
}

func copyFile(srcpath, dstpath string) error {
	if FileExists(dstpath) {
		return nil
	}

	r, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer r.Close() // ignore error: file was opened read-only.

	w, err := os.Create(dstpath)
	if err != nil {
		return err
	}

	defer func() {
		// Report the error, if any, from Close, but do so
		// only if there isn't already an outgoing error.
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

func getIconFile(config Config) string {
	if config.Icon != "" {
		ext := ""
		switch config.TargetOs {
		case "windows":
			ext = "ico"
			break
		case "darwin":
			ext = "icns"
			break
		case "linux":
			ext = "svg"
			break
		}

		return fmt.Sprintf("%s.%s", config.Icon, ext)
	}

	return ""
}
