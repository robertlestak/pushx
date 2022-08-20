package flags

var (
	HTTPMethod                = FlagSet.String("http-method", "POST", "HTTP method")
	HTTPURL                   = FlagSet.String("http-url", "", "HTTP url")
	HTTPContentType           = FlagSet.String("http-content-type", "", "HTTP content type")
	HTTPSuccessfulStatusCodes = FlagSet.String("http-successful-status-codes", "", "HTTP successful status codes")
	HTTPHeaders               = FlagSet.String("http-headers", "", "HTTP headers")
	HTTPEnableTLS             = FlagSet.Bool("http-enable-tls", false, "HTTP enable tls")
	HTTPTLSInsecure           = FlagSet.Bool("http-tls-insecure", false, "HTTP tls insecure")
	HTTPTLSCertFile           = FlagSet.String("http-tls-cert-file", "", "HTTP tls cert file")
	HTTPTLSKeyFile            = FlagSet.String("http-tls-key-file", "", "HTTP tls key file")
	HTTPTLSCAFile             = FlagSet.String("http-tls-ca-file", "", "HTTP tls ca file")
)
