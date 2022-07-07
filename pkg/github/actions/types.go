package actions

type RefType string

const (
	RefTypeTag    RefType = "tag"
	RefTypeBranch RefType = "branch"
)

func (r RefType) String() string {
	return string(r)
}

type OS string

const (
	OSLinux   OS = "Linux"
	OSWindows OS = "Windows"
	OSDarwin  OS = "macOS"
)

func (o OS) String() string {
	return string(o)
}

type Arch string

const (
	ArchX86   Arch = "X86"
	ArchX64   Arch = "X64"
	ArchARM   Arch = "ARM"
	ArchARM64 Arch = "ARM64"
)

func (a Arch) String() string {
	return string(a)
}
