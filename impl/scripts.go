package impl

type SetupScript struct {
	InstallDirectory    string   `json:"install_dir"`
	ExeName             string   `json:"exe_name"`
	Executables         []string `json:"executables"`
	PreInstallCommands  []string `json:"pre_install_cmds"`
	PostInstallCommands []string `json:"post_install_cmds"`
}

type LaunchScript struct {
	EntryPoint []string `json:"entrypoint"`
}
