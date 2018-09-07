package libs

import (
	"fmt"
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

func TestDecodePem2(t *testing.T) {

	sk := `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsILIfkNfpEpd01JQ9/7WZXrG+z4sECuY5JSUk4DnV1ec63SM
CGM2Ocr4rIw8WYPcYlXNzxHsZqXQL6wcGL0F2JwI3HEmWQy1kgeYl41k2qBVBTxp
0SQn7kFiPGBO3t9OY2iDjMZ+sLkEFWF6FcVMgeeo/ydJt3akv/YWbE7J4VKOPpF+
goHEs55OZoSFXwnC6OXMTCYI/oPNCp/ilQPmdOhZj/7JkUamZVi6FXpw0Ha1igel
eHf+cYyQTeaA9+F1zicEw99zUK+gvsEdFZf4Jjs9HzF68FKGLJwRMsfrwBsMiXFb
hm0+Iydzfin/PGIMhGRDEkYZoYhyd2GnvY1KkQIDAQABAoIBAFGatUzCo0YWZ5ha
dzqgY8iY4cOoM8DqFw6erq4fK1VTGSY2RMNf8uL3Ns63KvOLpekEmFkChVPTlxNj
M74wLbK+OM6JtVv4Jz2zK1XttdHpzETWcGmTmK5rJhqZv/0EKn6FBESiRFI5x2yH
G2VTSRLJ4nTMWiVqwbJJ+IXsSyNTPggNwpCESv6xmappjbpsr95VQAH3pMNXSUw+
hN9XCSe046fYkQqVYX9CBeSIvJV0X62bOsFANILEqbw97Lcb0CQH6ziyaZg1/Lku
Ahhb0U3w9q1Aa0ds+CU6BNuKF9AgqjQPgJaA2UJkykb+/Lq/he+ba8s7FDaV7flD
/wUqC/0CgYEA1FbW6FG4iGMnhcFO1auAtAwZsVIu4dULrpEh7XkO+GWzwL6LXFfr
aI+2iG/+ry/0Tw3rzpYDOZdURwXQzsg7P0D7fHex+wlKUv7A7461b4J9uGhMyxup
eqGJrQxEz42AAfSXKRxBvBZ2FvHQMj8vMzVT5Nzu7Sbzls1qzzGQyfsCgYEA1M4A
xtrtYpScDuqmtYXfcH3wHwE8EPmHScwVcNkKXHUz8XxKUG2cVTmIF+OHSdVtNXRn
Z8H2DQ8RUrYr13AjEE23dMqyqMwERkeRtAaATpzx5/1kWT7STB5FZdo5taprKnIH
DD//JvDNRZGE+Dm3snVrgVwSh4ELOb5GtF/ww+MCgYEA0tHVVZuth1xK5TXkO6in
hBtvduqNuZJ4or1d3hUPk/gF3BxS6UAxbgMhy5zGVNFb0xGTSe0PDqL8/fb8NhH/
IPI3voBoqUG0FWrxy+b4pNn+UJTdidrRDfxMYQ+JUji/GzHo0txHN6NlY6p6dyjc
iA8uWFAyuCqwNs2EempPwvkCgYApC+codBfvuNx1IBuxzbWwfAox4MoWOBs5R80m
7CZMeSrgvGmVI99QrP0sJm309t7OvbooUEhGco6T1qGUN4P75BGrzGEzn85/Q9S9
1bcv8qPSbJTLhguINRqLr9EL0mhDUqU7Xqi3eSEO9yefvXpKSM3IsNq+GkzxAVuf
9RV6owKBgH576Dx0nsrxnw+V789A8jxwALnmRGIYXNkIS/sObNhkwCM/+hKZqeaS
jwMxFfm6+wB6v3u94JwZCyVNnXNO6xJrN3zeQyi+qSJlnjA7gB6Nhpx/WRTL1NTy
hJr9xoOkO4i7AR+1eqaYoWc84whYmbFhvKpwCAU16Q34swYyRZGE
-----END RSA PRIVATE KEY-----`
	psk, _ := LoadPemPKCSPriKey(sk)
	fmt.Println(psk.N.BitLen())

	pk := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsILIfkNfpEpd01JQ9/7W
ZXrG+z4sECuY5JSUk4DnV1ec63SMCGM2Ocr4rIw8WYPcYlXNzxHsZqXQL6wcGL0F
2JwI3HEmWQy1kgeYl41k2qBVBTxp0SQn7kFiPGBO3t9OY2iDjMZ+sLkEFWF6FcVM
geeo/ydJt3akv/YWbE7J4VKOPpF+goHEs55OZoSFXwnC6OXMTCYI/oPNCp/ilQPm
dOhZj/7JkUamZVi6FXpw0Ha1igeleHf+cYyQTeaA9+F1zicEw99zUK+gvsEdFZf4
Jjs9HzF68FKGLJwRMsfrwBsMiXFbhm0+Iydzfin/PGIMhGRDEkYZoYhyd2GnvY1K
kQIDAQAB
-----END PUBLIC KEY-----`

	ppk, _ := LoadPemPubKey(pk)
	fmt.Println(ppk.N.BitLen())

}

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
		fmt.Println(p.N.BitLen())
	}
}

func TestLoadPemPubKey(t *testing.T) {
	if p, e := LoadPemPubKey(strings.Trim(strPubKey, " \t\r\n")); e != nil {
		t.Errorf("LoadPemPubKey fail %v", e)
	} else {
		t.Log("LoadPemPubKey ok", p)
		fmt.Println(p.N.BitLen())
	}
}
