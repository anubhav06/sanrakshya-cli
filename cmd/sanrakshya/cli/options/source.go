package options

type sourceCfg struct {
	File    fileSource `json:"file" yaml:"file" mapstructure:"file"`
}

type fileSource struct {
	Digests []string `json:"digests" yaml:"digests" mapstructure:"digests"`
}

func defaultSourceCfg() sourceCfg {
	return sourceCfg{
		File: fileSource{
			Digests: []string{"sha256"},
		},
	}
}
