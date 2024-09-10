package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/maja42/ember"
	"github.com/maja42/ember/embedding"
	"github.com/mcfriend99/exwrap/impl"
)

func damaged(err error) {
	log.Fatalln("Damaged executable:", err.Error())
}

func failed(err error) {
	log.Fatalln("Failed executable:", err.Error())
}

func readEmbededConfig(r ember.Reader, v any) {
	if buffer, err := io.ReadAll(r); err == nil {
		if err = json.Unmarshal(buffer, v); err != nil {
			damaged(err)
		}
	} else {
		damaged(err)
	}
}

func extractEmbededFile(r ember.Reader, v any) {
	if buffer, err := io.ReadAll(r); err == nil {
		if err = json.Unmarshal(buffer, v); err != nil {
			damaged(err)
		}
	} else {
		damaged(err)
	}
}

func main() {
	embedding.SkipCompatibilityCheck = true
	attachments, err := ember.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer attachments.Close()

	contents := attachments.List()

	var setup impl.SetupScript
	var launch impl.LaunchScript

	hasArchive := false
	_ = os.RemoveAll(impl.GetInstallExtractDir())

	for _, name := range contents {
		// s := attachments.Size(name)
		// fmt.Printf("\nAttachment %q has %d bytes:\n", name, s)
		r := attachments.Reader(name)

		switch name {
		case impl.EmbededArchiveName:
			hasArchive = true
			os.MkdirAll(impl.GetInstallExtractDir(), 0755)

			if file, err := os.OpenFile(impl.GetInstallExtractFile(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755); err == nil {
				if _, err := io.Copy(file, r); err != nil {
					damaged(err)
				}
			} else {
				damaged(err)
			}

			break
		case impl.EmbededSetupScript:
			readEmbededConfig(r, &setup)
			break
		case impl.EmbededLaunchScript:
			readEmbededConfig(r, &launch)
			break
		}
	}

	if hasArchive {
		install(setup, launch)
	} else {
		launchApp()
	}
}

func install(setup impl.SetupScript, launch impl.LaunchScript) {
	target := impl.GetInstallDir(setup.InstallDirectory)
	_ = os.RemoveAll(target)
	os.MkdirAll(target, 0755)

	if err := impl.Unzip(impl.GetInstallExtractFile(), target); err != nil {
		damaged(err)
	}

	if exe, err := os.Executable(); err == nil {
		exeTarget := path.Join(target, path.Base(exe))
		if impl.FileExists(exeTarget) {
			_ = os.RemoveAll(exeTarget)
		}

		if err = impl.RemoveEmbed(exe, exeTarget); err != nil {
			damaged(err)
		}
	} else {
		failed(err)
	}

	if runtime.GOOS != "windows" {
		for _, file := range setup.Executables {
			file = path.Join(target, file)
			if stat, err := os.Stat(file); err == nil {
				os.Chmod(file, stat.Mode()|0111)
			}
		}
	}

	if data, err := json.Marshal(launch.EntryPoint); err == nil {
		os.WriteFile(impl.GetLaunchScript(target), data, os.ModePerm)
	} else {
		log.Fatalln("Corrupt entrypoint.")
	}

	_ = os.RemoveAll(impl.GetInstallExtractDir())

	log.Println("Installation Completed!")
}

func launchApp() {

	// move to app directory
	appDir := impl.GetAppDir()
	command := impl.GetLaunchCommand(appDir)
	if len(command) == 0 {
		log.Fatalln("Missing entrypoint.")
	}

	type output struct {
		out []byte
		err error
	}

	ch := make(chan output)

	go func() {
		// move into app directory

		runtimeDir := appDir
		if runtime.GOOS == "darwin" {
			runtimeDir = path.Join(appDir, "../Resources")
		}
		os.Chdir(runtimeDir)

		var cmd *exec.Cmd
		program := impl.GetAbsoluteCommandProgram(command[0], runtime.GOOS == "darwin")

		if len(command) > 1 {
			cmd = exec.Command(program, command[1:]...)
		} else {
			cmd = exec.Command(program)
		}

		out, err := cmd.CombinedOutput()
		ch <- output{out, err}
	}()

	select {
	case x := <-ch:
		if x.err != nil {
			log.Fatalln(x.err.Error())
		}
	}
}
