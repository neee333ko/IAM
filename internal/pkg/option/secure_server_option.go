package option

import (
	"path/filepath"

	"github.com/neee333ko/IAM/internal/pkg/server"
	"github.com/neee333ko/component-base/pkg/cli"
	"github.com/neee333ko/component-base/pkg/validation/field"
	"github.com/spf13/pflag"
)

type SecureOption struct {
	Address   string  `json:"address" mapstructure:"address"`
	TLSConfig TLSInfo `json:"tls" mapstructure:"tls"`
}

type TLSInfo struct {
	CertKey   CertKey `json:"certkey" mapstructure:"certkey"`
	Directory string  `json:"dir" mapstructure:"dir"`
	PairName  string  `json:"pair" mapstructure:"pair"`
}

type CertKey struct {
	CertFile string `json:"cert" mapstructure:"cert"`
	CertKey  string `json:"key" mapstructure:"key"`
}

func (o *SecureOption) Flags(nfs *cli.NamedFlagSets) {
	fs := pflag.NewFlagSet("SecureOption", pflag.ExitOnError)

	fs.StringVar(&o.Address, "secure.address", "127.0.0.1:443", "address of https service")
	fs.StringVar(&o.TLSConfig.CertKey.CertFile, "secure.tls.certkey.cert", "", "cert of TLS")
	fs.StringVar(&o.TLSConfig.CertKey.CertKey, "secure.tls.certkey.key", "", "key of TLS")
	fs.StringVar(&o.TLSConfig.Directory, "secure.tls.dir", "/etc/iam/cert", "directory of cert")
	fs.StringVar(&o.TLSConfig.PairName, "secure.tls.pair", "", "name of cert and key")

	nfs.AddFlagSet("SecureOption", fs)
}

func (o *SecureOption) Validate() field.ErrorList {
	return nil
}

func (o *SecureOption) Complete() error {
	if o.Address == "" {
		o.Address = "127.0.0.1:443"
	}

	if o.TLSConfig.CertKey.CertFile == "" || o.TLSConfig.CertKey.CertKey == "" {
		o.TLSConfig.CertKey.CertFile = filepath.Join(o.TLSConfig.Directory, o.TLSConfig.PairName+".pem")
		o.TLSConfig.CertKey.CertKey = filepath.Join(o.TLSConfig.Directory, o.TLSConfig.PairName+"-key.pem")
	}

	return nil
}

func (o *SecureOption) ApplyTo(cfg *server.Config) error {
	cfg.SecureService.Address = o.Address
	cfg.SecureService.CertKey.CertFile = o.TLSConfig.CertKey.CertFile
	cfg.SecureService.CertKey.KeyFile = o.TLSConfig.CertKey.CertKey

	return nil
}
