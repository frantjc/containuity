package runtime

const (
	crtsDir     = "/etc/ssl"
	hostsFile   = "/etc/hosts"
	resolveConf = "/etc/resolv.conf"
)

var NetworkMounts = []*Mount{
	{
		Source:      crtsDir,
		Destination: crtsDir,
		Type:        MountTypeBind,
		Options:     []string{MountOptReadOnly},
	},
	{
		Source:      hostsFile,
		Destination: hostsFile,
		Type:        MountTypeBind,
		Options:     []string{MountOptReadOnly},
	},
	{
		Source:      resolveConf,
		Destination: resolveConf,
		Type:        MountTypeBind,
		Options:     []string{MountOptReadOnly},
	},
}
