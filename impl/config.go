package impl

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

type DarwinConfig struct {
	// Allows user to point to their own plist file for macos.
	PlistFile string `json:"plist,omitempty"`

	// When true, exwrap generates and app bundle instead of an installer.
	// Default: false
	CreateApp bool `json:"create_app,omitempty"`
}

type Config struct {
	// The root of the entire application.
	// Defaults to the current working directory.
	Root string `json:"root,omitempty"`

	// The entry point describes the command to run when
	// the executable starts. For example, ["python", "app.py"]
	EntryPoint []string `json:"entry_point"`

	// The name of the final executable.
	// Defaults to the name of the root folder.
	TargetName string `json:"target_name,omitempty"`

	// A list of commands to be run in order before installation begins
	PostInstallCommands []string `json:"post_install_cmds,omitempty"`

	// A list of commands to be run in order after installation completes.
	PreInstallCommands []string `json:"pre_install_cmds,omitempty"`

	// The OS on which exwrap is being run on (Defaults to your OS)
	SourceOs string `json:"source_os,omitempty"`

	// The processor architecture on which exwrap is being run on.
	// Defaults to your processor architecture
	SourceArch string `json:"source_arch,omitempty"`

	// The OS for which you are generating an executable for.
	// Defaults to your OS.
	TargetOs string `json:"os,omitempty"`

	// The processor architecture for which you are generating an
	// executable for.
	// Defaults to your processor architecture.
	TargetArch string `json:"arch,omitempty"`

	// A key/value pair that tells exwrap to replace any path the
	// matches or starts with the pattern indicated in the value
	// with the one indicated in the key.
	// In the format "override => match"
	PathOverrides map[string]string `json:"path_overrides,omitempty"`

	// Extra directories to add to the final executable.
	// In the format "source => destination"
	ExtraDirectories map[string]string `json:"extra_dirs"`

	// Extra files to add to the final executable.
	// In the format "source => destination"
	ExtraFiles map[string]string `json:"extra_files"`

	// A list of directories to not add to the final executable.
	ExcludeDirectories []string `json:"exclude_dirs"`

	// A list of files to not add to the final executable.
	ExcludeFiles []string `json:"exclude_files"`

	// The path that the final executable should be installed on.
	//
	// It is advisable that the install path should be a relative path.
	// On windows, it is relative to C:\Program Files\,
	// On Unix, it is relative to /home/$USERNAME/.
	// If empty, it defaults to the app name.
	// If an absolute path is given, it's used as is.
	InstallPath string `json:"install_path,omitempty"`

	// List of files that must be granted execute permission
	// when installation is extracted.
	Executables []string `json:"executables,omitempty"`

	// The application icon. This path should omit the extension as
	// exwrap will add the appropriate extension to the file.
	// Best practice is to have the icon in .icns (MacOS),
	// .ico (Windows), and .svg (Linux) format in a directory with
	// the same name
	Icon string `json:"icon,omitempty"`

	// Darwin (MacOS) specific configurations.
	Darwin DarwinConfig `json:"mac_os,omitempty"`
}

func LoadConfig(cmd CommandLine) Config {
	config := Config{}

	if data, err := os.ReadFile(cmd.ConfigFile); err == nil {
		if err = json.Unmarshal(data, &config); err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		log.Fatalln(err.Error())
	}

	if len(config.EntryPoint) == 0 {
		log.Fatalln("Entrypoint required.")
	}

	if config.Root == "" || config.Root == "." {
		if dir, err := os.Getwd(); err == nil {
			config.Root = dir
		} else {
			log.Fatalln("Could not detect root directory!")
		}
	} else {
		if absPath, err := getFileAbsPath(config.Root); err == nil {
			config.Root = absPath
		} else {
			log.Fatalln("Failed to resolve root directory:", err.Error())
		}
	}

	// ensure some major compatibilities
	config.SourceOs = strings.ToLower(config.SourceOs)
	config.SourceArch = strings.ToLower(config.SourceArch)
	config.TargetOs = strings.ToLower(config.TargetOs)
	config.TargetArch = strings.ToLower(config.TargetArch)

	// set defaults
	if config.TargetName == "" {
		config.TargetName = path.Base(config.Root)
	}
	if config.InstallPath == "" {
		config.InstallPath = config.TargetName
	}

	if config.SourceOs == "" {
		config.SourceOs = runtime.GOOS
	} else {
		config.SourceOs = strings.ToLower(config.SourceOs)
	}
	if config.SourceArch == "" {
		config.SourceArch = runtime.GOARCH
	} else {
		config.SourceArch = strings.ToLower(config.SourceArch)
	}

	if config.TargetOs == "" {
		config.TargetOs = config.SourceOs
	} else {
		config.TargetOs = strings.ToLower(config.TargetOs)
	}
	if config.TargetArch == "" {
		config.TargetArch = config.SourceArch
	} else {
		config.TargetArch = strings.ToLower(config.TargetArch)
	}

	if config.PathOverrides == nil {
		config.PathOverrides = make(map[string]string, 0)
	} else {
		for i, x := range config.PathOverrides {
			if abs, err := getFileAbsPath(i); err == nil {
				delete(config.PathOverrides, i)
				config.PathOverrides[abs] = x
			}
		}
	}

	if config.ExtraDirectories == nil {
		config.ExtraDirectories = make(map[string]string, 0)
	} else {
		for i, x := range config.ExtraDirectories {
			if abs, err := getFileAbsPath(i); err == nil {
				config.ExtraDirectories[abs] = x
			} else {
				delete(config.ExtraDirectories, i)
			}
		}
	}

	if config.ExtraFiles == nil {
		config.ExtraFiles = make(map[string]string, 0)
	} else {
		for i, x := range config.ExtraFiles {
			if abs, err := getFileAbsPath(i); err == nil {
				config.ExtraFiles[abs] = x
			} else {
				delete(config.ExtraFiles, i)
			}
		}
	}

	if config.ExcludeDirectories == nil {
		config.ExcludeDirectories = make([]string, 0)
		config.ExcludeDirectories = append(config.ExcludeDirectories, cmd.BuildDirectory)
	}

	if config.Executables == nil {
		config.Executables = make([]string, 0)
	}

	newExcludedDirList := make([]string, 0)
	for _, x := range config.ExcludeDirectories {
		if abs, err := getFileAbsPath(x); err == nil {
			newExcludedDirList = append(newExcludedDirList, abs)
		}
	}
	config.ExcludeDirectories = newExcludedDirList

	if config.ExcludeFiles == nil {
		config.ExcludeFiles = make([]string, 0)
	} else {
		newList := make([]string, 0)

		for _, x := range config.ExcludeFiles {
			if abs, err := getFileAbsPath(x); err == nil {
				newList = append(newList, abs)
			}
		}

		config.ExcludeFiles = newList
	}

	if config.Darwin.PlistFile == "" {
		config.Darwin.PlistFile = path.Join(getResourcesDirectory(), "Info.plist")
	}

	if config.PreInstallCommands == nil {
		config.PreInstallCommands = make([]string, 0)
	}

	if config.PostInstallCommands == nil {
		config.PostInstallCommands = make([]string, 0)
	}

	return config
}

func LoadDefaultConfig() Config {
	return LoadConfig(CommandLine{
		BuildDirectory: DefaultBuildDirectory,
		ConfigFile:     DefaultConfigFile,
	})
}
