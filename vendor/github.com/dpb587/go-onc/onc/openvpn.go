package onc

// https://chromium.googlesource.com/chromium/src/+/master/components/onc/docs/onc_spec.md#openvpn-type
type OpenVPN struct {
	Auth                   string   `json:"Auth,omitempty"`
	AuthRetry              string   `json:"AuthRetry,omitempty"`
	AuthNoCache            bool     `json:"AuthNoCache,omitempty"`
	Cipher                 string   `json:"Cipher,omitempty"`
	ClientCertPKCS11Id     string   `json:"ClientCertPKCS11Id,omitempty"`
	ClientCertPattern      string   `json:"ClientCertPattern,omitempty"`
	ClientCertRef          string   `json:"ClientCertRef,omitempty"`
	ClientCertType         string   `json:"ClientCertType,omitempty"`
	CompLZO                string   `json:"CompLZO,omitempty"`
	CompNoAdapt            bool     `json:"CompNoAdapt,omitempty"`
	ExtraHosts             []string `json:"ExtraHosts,omitempty"`
	IgnoreDefaultRoute     bool     `json:"IgnoreDefaultRoute,omitempty"`
	KeyDirection           string   `json:"KeyDirection,omitempty"`
	NsCertType             string   `json:"NsCertType,omitempty"`
	OTP                    string   `json:"OTP,omitempty"`
	Password               string   `json:"Password,omitempty"`
	Port                   int      `json:"Port,omitempty"`
	Proto                  string   `json:"Proto,omitempty"`
	PushPeerInfo           bool     `json:"PushPeerInfo,omitempty"`
	RemoteCertEKU          string   `json:"RemoteCertEKU,omitempty"`
	RemoteCertKU           []string `json:"RemoteCertKU,omitempty"`
	RemoteCertTLS          string   `json:"RemoteCertTLS,omitempty"`
	RenegSec               int      `json:"RenegSec,omitempty"`
	SaveCredentials        bool     `json:"SaveCredentials,omitempty"`
	ServerCAPEMs           []string `json:"ServerCAPEMs,omitempty"`
	ServerCARefs           []string `json:"ServerCARefs,omitempty"`
	ServerCARef            string   `json:"ServerCARef,omitempty"`
	ServerCertRef          string   `json:"ServerCertRef,omitempty"`
	ServerPollTimeout      int      `json:"ServerPollTimeout,omitempty"`
	Shaper                 int      `json:"Shaper,omitempty"`
	StaticChallenge        string   `json:"StaticChallenge,omitempty"`
	TLSAuthContents        string   `json:"TLSAuthContents,omitempty"`
	TLSRemote              string   `json:"TLSRemote,omitempty"`
	TLSVersionMin          string   `json:"TLSVersionMin,omitempty"`
	UserAuthenticationType string   `json:"UserAuthenticationType,omitempty"`
	Username               string   `json:"Username,omitempty"`
	Verb                   string   `json:"Verb,omitempty"`
	VerifyHash             string   `json:"VerifyHash,omitempty"`
	VerifyX509             string   `json:"VerifyX509,omitempty"`
}
