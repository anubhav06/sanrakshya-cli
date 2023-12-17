package options

import (
	"github.com/anchore/clio"
)


type Credentials struct {
	Account_ID string `yaml:"account_id" json:"account_id" mapstructure:"account_id"`
	SecretKey  string `yaml:"secret_key" json:"secret_key" mapstructure:"secret_key"`
}

var _ interface {
	clio.FlagAdder
} = (*Credentials)(nil)


var _ clio.FlagAdder = (*Credentials)(nil)

func (o Credentials) AddFlags(flags clio.FlagSet) {
	flags.StringVarP(&o.Account_ID, "account-id", "a", "the account-id of web-app")
	flags.StringVarP(&o.SecretKey, "secret-key", "k", "the secret-key of web-app")
}

func DefaultCredentials() Credentials {
	return Credentials{
		Account_ID: "",
		SecretKey:  "",
	}
}
