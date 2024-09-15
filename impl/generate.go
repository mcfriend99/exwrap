package impl

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func makeAttachements(config Config, files []string) map[string]string {
	attachments := make(map[string]string, 0)

MainLoop:
	for _, file := range files {
		if !hasMatchInList(config.ExcludeFiles, file) && !hasMatchInList(config.ExcludeDirectories, file) {
			mainKey := trimRoot(file, config.Root)

			if val, ok := config.PathOverrides[file]; ok {
				if val != "" {
					attachments[val] = file
				}
				continue
			}

			if val, ok := config.ExtraFiles[file]; ok {
				if val != "" {
					attachments[val] = file
				}
				continue
			}

			for base, override := range config.PathOverrides {
				if strings.HasPrefix(file, base) {
					attachments[trimRoot(file, base)] = strings.ReplaceAll(file, base, override)
					continue MainLoop
				}
			}

			for dir, targetDir := range config.ExtraDirectories {
				if strings.HasPrefix(file, dir) {
					newKey := strings.ReplaceAll(file, dir, targetDir)
					attachments[newKey] = file
					continue MainLoop
				}
			}

			if mainKey != "" {
				attachments[mainKey] = file
			}
		}
	}

	return attachments
}

func generateAttachments(config Config) map[string]string {
	attachments := make(map[string]string, 0)

	for k, v := range makeAttachements(config, listFiles(config.Root)) {
		attachments[v] = k
	}
	for dir := range config.ExtraDirectories {
		for k, v := range makeAttachements(config, listFiles(dir)) {
			attachments[v] = k
		}
	}
	for file := range config.ExtraFiles {
		for k, v := range makeAttachements(config, []string{file}) {
			attachments[v] = k
		}
	}

	return attachments
}

func Generate(config Config, cmd CommandLine) string {
	// ensure we're trying to build a supported os/arch combination.
	failFormat := "Unsupported Os/Arch combination: %s/%s"
	if combo, ok := BuildCombinations[OSArch{config.TargetOs, config.TargetArch}]; ok {

		// For now, we're only supporting first-class build targets.
		// TODO: Support non first-class targets
		if !combo.FirstClass {
			log.Fatalf(failFormat, config.TargetOs, config.TargetArch)
		}
	} else {
		log.Fatalf(failFormat, config.TargetOs, config.TargetArch)
	}

	_ = os.RemoveAll(getBuildDir(cmd))

	if config.TargetOs == "darwin" && config.Darwin.CreateApp {
		return GenerateDarwin(config, cmd)
	} else {
		return GenerateDefault(config, cmd)
	}
}

func GenerateDefault(config Config, cmd CommandLine) string {
	if err := os.MkdirAll(getBuildDir(cmd), os.ModePerm); err == nil {
		attachments := generateAttachments(config)

		// create build archive target
		targetArchive := getTargetBuildArchive(config, cmd)

		zipfile, err := os.Create(targetArchive)
		if err != nil {
			log.Fatalln("Failed to create application archive.")
		}
		defer zipfile.Close()

		// init zip file
		archive := zip.NewWriter(zipfile)

		// write files into it.
		for src, dest := range attachments {
			fmt.Printf("File discovered: %s => %s\n", src, dest)

			if file, err := os.Open(src); err == nil {
				if zf, err := archive.Create(dest); err == nil {
					_, err = io.Copy(zf, file)
				}

				file.Close()
			}
		}
		archive.Close()

		clear(attachments)

		srcWrapper := getPkgExeFromConfig(config)
		targetBase := getTargetBaseName(cmd, config)
		err = copyFile(srcWrapper, targetBase)
		if err != nil {
			log.Fatalln("Failed to copy application wrapper:", err.Error())
		}
		attachments[EmbededArchiveName] = targetArchive

		targetExe := getTargetExeName(cmd, config)
		if FileExists(targetExe) {
			os.Remove(targetExe)
		}

		// Create the setup script
		setupScript := SetupScript{
			InstallDirectory:    config.InstallPath,
			Executables:         config.Executables,
			ExeName:             config.TargetName,
			PreInstallCommands:  config.PreInstallCommands,
			PostInstallCommands: config.PostInstallCommands,
		}
		setupName := getBuildSetupScriptName(cmd)
		if data, err := json.Marshal(setupScript); err == nil {
			if err = os.WriteFile(setupName, data, fs.ModePerm); err != nil {
				log.Fatalln("Failed to create setup script.")
			}
		} else {
			log.Fatalln("Failed to create setup script.")
		}
		attachments[EmbededSetupScript] = setupName

		// Create the launch script
		launchScript := LaunchScript{
			EntryPoint: config.EntryPoint,
		}
		launchName := getBuildLaunchScriptName(cmd)
		if data, err := json.Marshal(launchScript); err == nil {
			if err = os.WriteFile(launchName, data, fs.ModePerm); err != nil {
				log.Fatalln("Failed to create launch script.")
			}
		} else {
			log.Fatalln("Failed to create launch script.")
		}
		attachments[EmbededLaunchScript] = launchName

		err = Embed(targetBase, targetExe, attachments)
		if err != nil {
			log.Fatalln(err.Error())
		}

		// delete redundant files...
		os.Remove(targetArchive)
		os.Remove(targetBase)
		os.Remove(setupName)
		os.Remove(launchName)

		return targetExe
	} else {
		log.Fatalln("Failed to create build directory!")
	}

	return ""
}

func GenerateDarwin(config Config, cmd CommandLine) string {
	if err := os.MkdirAll(getBuildDir(cmd), os.ModePerm); err == nil {
		attachments := generateAttachments(config)

		// create build archive target
		targetArchive := getTargetBuildArchive(config, cmd)

		macosDir := path.Join(targetArchive, "Contents", "MacOS")
		resourcesDir := path.Join(targetArchive, "Contents", "Resources")
		frameworksDir := path.Join(targetArchive, "Contents", "Frameworks")

		// create required dirs
		os.MkdirAll(macosDir, os.ModePerm)
		os.MkdirAll(resourcesDir, os.ModePerm)
		os.MkdirAll(frameworksDir, os.ModePerm)

		// write files into it.
		for src, dest := range attachments {
			fmt.Printf("File discovered: %s => %s\n", src, dest)
			tmpDst := dest

			dest = path.Join(resourcesDir, dest)
			os.MkdirAll(filepath.Dir(dest), os.ModePerm)

			if file, err := os.Open(src); err == nil {
				mode := os.ModePerm
				if stat, err := file.Stat(); err == nil {
					mode = stat.Mode()
				}

				if zf, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode); err == nil {
					_, err = io.Copy(zf, file)

					// process executables list
					if stringListContains(config.Executables, tmpDst) {
						os.Chmod(dest, mode|0111)
					}
				}

				file.Close()
			}
		}

		clear(attachments)

		// Create the launch script
		launchScript := path.Join(macosDir, getLaunchScriptForDarwinApp(config))

		if data, err := json.Marshal(config.EntryPoint); err == nil {
			if err = os.WriteFile(launchScript, data, fs.ModePerm); err != nil {
				log.Fatalln("Failed to create launch script:", err.Error())
			}
		} else {
			log.Fatalln("Failed to create launch script:", err.Error())
		}

		// indicate this is a darwin app
		_ = os.WriteFile(path.Join(macosDir, DarwinAppLockfile), []byte{}, os.ModePerm)

		srcWrapper := getPkgExeFromConfig(config)
		targetExe := path.Join(macosDir, config.TargetName)

		if err = copyFile(srcWrapper, targetExe); err != nil {
			log.Fatalln("Failed to create launch file:", err.Error())
		} else {
			if stat, err := os.Stat(targetExe); err == nil {
				os.Chmod(targetExe, stat.Mode()|0111)
			}
		}

		// Create the info.plist file
		if data, err := os.ReadFile(config.Darwin.PlistFile); err == nil {
			data = []byte(strings.ReplaceAll(string(data), "${EXE}", config.TargetName))

			if err = os.WriteFile(path.Join(targetArchive, "Contents", "Info.plist"), data, os.ModePerm); err != nil {
				log.Fatalln("Plist creation failed:", err.Error())
			}
		} else {
			log.Fatalln("Plist read failed:", err.Error())
		}

		// add the icon file if set
		icon := getIconFile(config)
		if icon != "" {
			if err = copyFile(icon, path.Join(resourcesDir, "icon.icns")); err != nil {
				// do nothing (because apps will still run without their icons)...
			}
		}

		return targetArchive
	} else {
		log.Fatalln("Failed to create build directory:", err.Error())
	}

	return ""
}
