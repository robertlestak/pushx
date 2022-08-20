package flags

var (
	NFSHost   = FlagSet.String("nfs-host", "", "NFS host")
	NFSKey    = FlagSet.String("nfs-key", "", "NFS key")
	NFSFolder = FlagSet.String("nfs-folder", "", "NFS folder")
	NFSTarget = FlagSet.String("nfs-target", "", "NFS target")
)
