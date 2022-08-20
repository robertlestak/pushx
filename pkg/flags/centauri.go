package flags

var (
	CentauriPeerURL     = FlagSet.String("centauri-peer-url", "", "Centauri peer URL")
	CentauriChannel     = FlagSet.String("centauri-channel", "default", "Centauri channel")
	CentauriKey         = FlagSet.String("centauri-public-key", "", "Centauri public key")
	CentauriKeyBase64   = FlagSet.String("centauri-public-key-base64", "", "Centauri public key base64")
	CentauriMessageType = FlagSet.String("centauri-message-type", "bytes", "Centauri message type. One of: bytes, file")
	CentauriFilename    = FlagSet.String("centauri-filename", "", "Centauri filename")
)
