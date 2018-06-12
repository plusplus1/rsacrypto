package libs

import (
	"strings"
	"testing"
)

var (
	strPubKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCB9E/o59hwdJn8QS4lQLEdg4SX
rScpvEYqSGLWeMcgif6qn/9WWbyLUMuxJojBYPeyO7H/wHkfBvAsTGTOT4JBpO1I
F5U+msOSoq392NQXHeIsoaxGBJPOpfRzHyOHR9I/y6I616aEs2xK8QNo07fDJRXw
BErEz4GXNIe77+3rCwIDAQAB
-----END PUBLIC KEY-----`

	strPriKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCB9E/o59hwdJn8QS4lQLEdg4SXrScpvEYqSGLWeMcgif6qn/9W
WbyLUMuxJojBYPeyO7H/wHkfBvAsTGTOT4JBpO1IF5U+msOSoq392NQXHeIsoaxG
BJPOpfRzHyOHR9I/y6I616aEs2xK8QNo07fDJRXwBErEz4GXNIe77+3rCwIDAQAB
AoGAEmK2S4VBoeddQcsW9D5K0Joi0DK3L7hrVOAY268KvRfI5+rq/RSFEFDGALIq
9vD5GkZH0J8yg6k/oYTNXkLxcpMbm4alhwEBVAcDv3Xk2U/HSdNVPcGlkYi7ZiGZ
VCyfV0B8tMEmyQ5O2VF0p2Bhznp1zIeQrZQV5PXiWwx/WwECQQCQu+zrdqyRaq1F
G+O585CbQzMm6SlTOnjl/7rf9XR0D5qrk1JMhz00DtlRvpZzHpMDpi2p6NM5vrvH
3f5CN0YRAkEA5duxQNpiJWPNDyocEC4ThlwA6sYsoHhkJ2mYvmv7P+T+zqpDYmHs
kz8xjvrN0hkb6agci368qcPuwfEPNwjTWwJAAsYXPEwB8qeAuppWOvIYC2G2UUCW
simkt4O3KSOjH7ZM2IzyPtU4rw65y39DkuE7IA7HQUJdCfZF0wbGIK6+gQJAevSf
B7MKBzgwq+j5pAoRtbCnaO7jVl+wK4kIBOycNNyZFRHtA8agF1AZgYNV8AowbTfZ
NSFxaFp/8Eyzt9vHuwJAULmctS4mcY+vkUN50tJ8ZAD/BTXUthtkf3JdurMz3Fxr
wbklj/kzHjP3Hw6eDhZmHsne1fA9ek3/HATFT2zADA==
-----END RSA PRIVATE KEY-----`
)

func TestDecodePem(t *testing.T) {

	if b, e := DecodePem(strPubKey); e != nil {
		t.Errorf("decode public pem fail, %v", e)
	} else {
		t.Log("decode public pem ok", b.Type)

		if p, e := LoadPubKeyBytes(b.Bytes); e != nil {
			t.Errorf("load bytes public key fail %v", e)
		} else {
			t.Log("load bytes public key ok", p)
		}
	}

	if b, e := DecodePem(strPriKey); e != nil {
		t.Errorf("decode private pem fail, %v", e)
	} else {
		t.Log("decode private pem ok", b.Type)

		if p, e := LoadPKCSPriKeyBytes(b.Bytes); e != nil {
			t.Errorf("load bytes private key fail %v", e)
		} else {
			t.Log("load bytes private key ok", p)
		}

	}

}

func TestLoadPemPKCSPriKey(t *testing.T) {

	if p, e := LoadPemPKCSPriKey(strPriKey); e != nil {
		t.Errorf("LoadPemPKCSPriKey fail %v", e)
	} else {
		t.Log("LoadPemPKCSPriKey ok", p)
	}
}

func TestLoadPemPubKey(t *testing.T) {
	if p, e := LoadPemPubKey(strings.Trim(strPubKey, " \t\r\n")); e != nil {
		t.Errorf("LoadPemPubKey fail %v", e)
	} else {
		t.Log("LoadPemPubKey ok", p)
	}
}
