package impl

type SetupScript struct {
	InstallDirectory string   `json:"install_dir"`
	Executables      []string `json:"executables"`
}

type LaunchScript struct {
	EntryPoint []string `json:"entrypoint"`
}
