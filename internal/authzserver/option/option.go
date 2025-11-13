package option

import "github.com/neee333ko/IAM/internal/pkg/option"

type Option struct {
	InsecureOp *option.InsecureOption `json:"insecure" mapstructure:"insecure"`
	SecureOp   *option.SecureOption   `json:"secure" mapstructure:"secure"`
	RunOp      *option.RunOption      `json:"run" mapstructure:"run"`
	Feature    *option.Feature        `json:"featrue" mapstructure:"feature"`
}
