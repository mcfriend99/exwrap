// copied from golang internal/platform module.

package impl

type OsArchInfo struct {
	CgoSupported bool
	FirstClass   bool
	Broken       bool
}

type OSArch struct {
	GOOS, GOARCH string
}

var BuildCombinations = map[OSArch]OsArchInfo{
	{"aix", "ppc64"}:       {CgoSupported: true},
	{"android", "386"}:     {CgoSupported: true},
	{"android", "amd64"}:   {CgoSupported: true},
	{"android", "arm"}:     {CgoSupported: true},
	{"android", "arm64"}:   {CgoSupported: true},
	{"darwin", "amd64"}:    {CgoSupported: true, FirstClass: true},
	{"darwin", "arm64"}:    {CgoSupported: true, FirstClass: true},
	{"dragonfly", "amd64"}: {CgoSupported: true},
	{"freebsd", "386"}:     {CgoSupported: true},
	{"freebsd", "amd64"}:   {CgoSupported: true},
	{"freebsd", "arm"}:     {CgoSupported: true},
	{"freebsd", "arm64"}:   {CgoSupported: true},
	{"freebsd", "riscv64"}: {CgoSupported: true},
	{"illumos", "amd64"}:   {CgoSupported: true},
	{"ios", "amd64"}:       {CgoSupported: true},
	{"ios", "arm64"}:       {CgoSupported: true},
	{"js", "wasm"}:         {},
	{"linux", "386"}:       {CgoSupported: true, FirstClass: true},
	{"linux", "amd64"}:     {CgoSupported: true, FirstClass: true},
	{"linux", "arm"}:       {CgoSupported: true, FirstClass: true},
	{"linux", "arm64"}:     {CgoSupported: true, FirstClass: true},
	{"linux", "loong64"}:   {CgoSupported: true},
	{"linux", "mips"}:      {CgoSupported: true},
	{"linux", "mips64"}:    {CgoSupported: true},
	{"linux", "mips64le"}:  {CgoSupported: true},
	{"linux", "mipsle"}:    {CgoSupported: true},
	{"linux", "ppc64"}:     {},
	{"linux", "ppc64le"}:   {CgoSupported: true},
	{"linux", "riscv64"}:   {CgoSupported: true},
	{"linux", "s390x"}:     {CgoSupported: true},
	{"linux", "sparc64"}:   {CgoSupported: true, Broken: true},
	{"netbsd", "386"}:      {CgoSupported: true},
	{"netbsd", "amd64"}:    {CgoSupported: true},
	{"netbsd", "arm"}:      {CgoSupported: true},
	{"netbsd", "arm64"}:    {CgoSupported: true},
	{"openbsd", "386"}:     {CgoSupported: true},
	{"openbsd", "amd64"}:   {CgoSupported: true},
	{"openbsd", "arm"}:     {CgoSupported: true},
	{"openbsd", "arm64"}:   {CgoSupported: true},
	{"openbsd", "mips64"}:  {CgoSupported: true, Broken: true},
	{"openbsd", "ppc64"}:   {},
	{"openbsd", "riscv64"}: {Broken: true},
	{"plan9", "386"}:       {},
	{"plan9", "amd64"}:     {},
	{"plan9", "arm"}:       {},
	{"solaris", "amd64"}:   {CgoSupported: true},
	{"wasip1", "wasm"}:     {},
	{"windows", "386"}:     {CgoSupported: true, FirstClass: true},
	{"windows", "amd64"}:   {CgoSupported: true, FirstClass: true},
	{"windows", "arm"}:     {},
	{"windows", "arm64"}:   {CgoSupported: true},
}
