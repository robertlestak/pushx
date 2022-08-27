package flags

var (
	EtcdHosts       = FlagSet.String("etcd-hosts", "", "Etcd hosts")
	EtcdUsername    = FlagSet.String("etcd-username", "", "Etcd username")
	EtcdPassword    = FlagSet.String("etcd-password", "", "Etcd password")
	EtcdKey         = FlagSet.String("etcd-key", "", "Etcd key")
	EtcdLimit       = FlagSet.Int("etcd-limit", 0, "Etcd limit. 0 for no limit")
	EtcdTLSEnable   = FlagSet.Bool("etcd-tls-enable", false, "Etcd TLS enable")
	EtcdTLSInsecure = FlagSet.Bool("etcd-tls-insecure", false, "Etcd TLS insecure")
	EtcdTLSCert     = FlagSet.String("etcd-tls-cert", "", "Etcd TLS cert")
	EtcdTLSKey      = FlagSet.String("etcd-tls-key", "", "Etcd TLS key")
	EtcdTLSCA       = FlagSet.String("etcd-tls-ca", "", "Etcd TLS ca")
)
