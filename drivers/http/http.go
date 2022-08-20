package http

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/robertlestak/pushx/pkg/flags"
	"github.com/robertlestak/pushx/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type HTTPRequest struct {
	Method                string
	URL                   string
	ContentType           string
	SuccessfulStatusCodes []int
	Headers               map[string]string
}

type HTTP struct {
	Client      *http.Client
	EnableTLS   *bool
	TLSCA       *string
	TLSCert     *string
	TLSKey      *string
	TLSInsecure *bool
	Request     *HTTPRequest
	Key         *string
}

func (d *HTTP) LoadEnv(prefix string) error {
	l := log.WithFields(log.Fields{
		"pkg": "http",
		"fn":  "LoadEnv",
	})
	l.Debug("Loading environment")
	if d.Request == nil {
		d.Request = &HTTPRequest{}
	}
	if os.Getenv(prefix+"HTTP_REQUEST_METHOD") != "" {
		d.Request.Method = os.Getenv(prefix + "HTTP_REQUEST_METHOD")
	}
	if os.Getenv(prefix+"HTTP_REQUEST_URL") != "" {
		d.Request.URL = os.Getenv(prefix + "HTTP_REQUEST_URL")
	}
	if os.Getenv(prefix+"HTTP_REQUEST_CONTENT_TYPE") != "" {
		d.Request.ContentType = os.Getenv(prefix + "HTTP_REQUEST_CONTENT_TYPE")
	}
	if os.Getenv(prefix+"HTTP_REQUEST_SUCCESSFUL_STATUS_CODES") != "" {
		d.Request.SuccessfulStatusCodes = parseIntSlice(os.Getenv(prefix + "HTTP_REQUEST_SUCCESSFUL_STATUS_CODES"))
	}
	if os.Getenv(prefix+"HTTP_REQUEST_HEADERS") != "" {
		d.Request.Headers = parseHeaderMap(os.Getenv(prefix + "HTTP_REQUEST_HEADERS"))
	}
	if os.Getenv(prefix+"HTTP_ENABLE_TLS") != "" {
		v := os.Getenv(prefix+"HTTP_ENABLE_TLS") == "true"
		d.EnableTLS = &v
	}
	if os.Getenv(prefix+"HTTP_TLS_CERT_FILE") != "" {
		v := os.Getenv(prefix + "HTTP_TLS_CERT_FILE")
		d.TLSCert = &v
	}
	if os.Getenv(prefix+"HTTP_TLS_KEY_FILE") != "" {
		v := os.Getenv(prefix + "HTTP_TLS_KEY_FILE")
		d.TLSKey = &v
	}
	if os.Getenv(prefix+"HTTP_TLS_CA_FILE") != "" {
		v := os.Getenv(prefix + "HTTP_TLS_CA_FILE")
		d.TLSCA = &v
	}
	return nil
}

func parseIntSlice(s string) []int {
	var r []int
	for _, v := range strings.Split(s, ",") {
		i, e := strconv.Atoi(v)
		if e != nil {
			continue
		}
		r = append(r, i)
	}
	return r
}

func parseHeaderMap(s string) map[string]string {
	r := make(map[string]string)
	for _, v := range strings.Split(s, ",") {
		kv := strings.Split(v, ":")
		if len(kv) != 2 {
			continue
		}
		r[kv[0]] = kv[1]
	}
	return r
}

func (d *HTTP) LoadFlags() error {
	l := log.WithFields(log.Fields{
		"pkg": "http",
		"fn":  "LoadFlags",
	})
	l.Debug("Loading flags")
	rr := &HTTPRequest{
		Method:                *flags.HTTPMethod,
		URL:                   *flags.HTTPURL,
		ContentType:           *flags.HTTPContentType,
		SuccessfulStatusCodes: parseIntSlice(*flags.HTTPSuccessfulStatusCodes),
		Headers:               parseHeaderMap(*flags.HTTPHeaders),
	}
	if rr.Method == "" {
		rr.Method = "POST"
	}
	d.Request = rr
	d.EnableTLS = flags.HTTPEnableTLS
	d.TLSCert = flags.HTTPTLSCertFile
	d.TLSKey = flags.HTTPTLSKeyFile
	d.TLSCA = flags.HTTPTLSCAFile
	return nil
}

func (d *HTTP) Init() error {
	l := log.WithFields(log.Fields{
		"pkg": "http",
		"fn":  "Init",
	})
	l.Debug("Initializing http driver")
	d.Client = &http.Client{}
	tc, err := utils.TlsConfig(d.EnableTLS, d.TLSInsecure, d.TLSCA, d.TLSCert, d.TLSKey)
	if err != nil {
		return err
	}
	d.Client.Transport = &http.Transport{
		TLSClientConfig: tc,
	}
	return nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (d *HTTP) Push(r io.Reader) error {
	l := log.WithFields(log.Fields{
		"pkg": "http",
		"fn":  "Push",
	})
	l.Debug("sending http request")
	if d.Request == nil {
		return errors.New("request is nil")
	}
	if d.Request.Method == "" {
		d.Request.Method = "POST"
	}
	if d.Request.URL == "" {
		return errors.New("URL is nil")
	}
	req, err := http.NewRequest(d.Request.Method, d.Request.URL, r)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	for k, v := range d.Request.Headers {
		req.Header.Add(k, v)
	}
	if d.Request.ContentType != "" {
		req.Header.Add("Content-Type", d.Request.ContentType)
	}
	resp, err := d.Client.Do(req)
	if err != nil {
		l.Errorf("%+v", err)
		return err
	}
	if len(d.Request.SuccessfulStatusCodes) > 0 {
		if !contains(d.Request.SuccessfulStatusCodes, resp.StatusCode) {
			l.Errorf("Status code %d not in successful status codes", resp.StatusCode)
			return errors.New("status code not in successful status codes")
		}
	}
	l.Debug("http request sent")
	return nil
}

func (d *HTTP) Cleanup() error {
	return nil
}
