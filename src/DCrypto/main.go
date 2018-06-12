package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

func main() {

	strK := `-----BEGIN RSA PRIVATE KEY-----
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

	conn, _ := net.Dial("tcp", "127.0.0.1:4006")
	defer conn.Close()

	bufLen := make([]byte, 4)
	bufK := []byte(strK)
	binary.BigEndian.PutUint32(bufLen, uint32(len(bufK)))

	conn.Write(bufLen)
	conn.Write(bufK)

	reader := bufio.NewReader(conn)

	reader.Read(bufLen)

	fmt.Println("result, len = ", binary.BigEndian.Uint32(bufLen))

}
