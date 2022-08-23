package flags

var (
	SMBHost  = FlagSet.String("smb-host", "", "SMB host")
	SMBPort  = FlagSet.Int("smb-port", 445, "SMB port")
	SMBUser  = FlagSet.String("smb-user", "", "SMB user")
	SMBPass  = FlagSet.String("smb-pass", "", "SMB pass")
	SMBKey   = FlagSet.String("smb-key", "", "SMB key")
	SMBShare = FlagSet.String("smb-share", "", "SMB share")
)
