package fs

type Config struct {
	CertificatePath string `yaml:"certificate_path"`
	PrivateKeyPath  string `yaml:"private_key_path"`
}
