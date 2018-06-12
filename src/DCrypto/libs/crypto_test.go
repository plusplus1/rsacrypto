package libs

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestEncrypt(t *testing.T) {

	pk, _ := LoadPemPubKey(strPubKey)
	sk, _ := LoadPemPKCSPriKey(strPriKey)

	srcText := "hello the world"
	for i := 0; i < 13; i++ {
		srcText += srcText
	}

	t.Logf("length = %.2f K", float64(len(srcText))/1024.)
	st := time.Now()

	if enText, e := Encrypt(srcText, pk); e != nil {
		t.Error("encrypt fail, ", e)
	} else {
		t.Log("encrypt ok cost  = ", time.Since(st).Seconds())

		st = time.Now()
		if deText, e := Decrypt(enText, sk, true); e != nil {
			t.Error("decrypt fail, ", e)
		} else {
			t.Log("decrypt ok cost = ", time.Since(st).Seconds())
			if deText == srcText {
				t.Log("decrypt right ")
			} else {
				t.Error("decrypt wrong ")
			}
		}
	}

}

func TestEncrypt2(t *testing.T) {

	pk, _ := LoadPemPubKey(strPubKey)
	sk, _ := LoadPemPKCSPriKey(strPriKey)

	srcText := "hello the wo"
	for i := 0; i < 16; i++ {
		srcText += srcText
	}

	t.Logf("length = %.2f K", float64(len(srcText))/1024.)
	st := time.Now()

	if enText, e := Encrypt(srcText, pk); e != nil {
		t.Error("encrypt fail, ", e)
	} else {
		t.Log("encrypt ok cost  = ", time.Since(st).Seconds())

		st = time.Now()
		if deText, e := Decrypt(enText, sk, true); e != nil {
			t.Error("decrypt fail, ", e)
		} else {
			t.Log("decrypt ok cost = ", time.Since(st).Seconds())
			if deText == srcText {
				t.Log("decrypt right ")
			} else {
				t.Error("decrypt wrong ")
			}
		}
	}

}

func TestEncrypt3(t *testing.T) {

	pk, _ := LoadPemPubKey(strPubKey)
	sk, _ := LoadPemPKCSPriKey(strPriKey)

	srcText := "hello the world"
	for i := 0; i < 17; i++ {
		srcText += srcText
	}

	t.Logf("length = %.2f M", float64(len(srcText))/1024./1024)

	st := time.Now()

	if enText, e := Encrypt(srcText, pk); e != nil {
		t.Error("encrypt fail, ", e)
	} else {
		t.Log("encrypt ok cost  = ", time.Since(st).Seconds())

		st = time.Now()
		if deText, e := Decrypt(enText, sk, true); e != nil {
			t.Error("decrypt fail, ", e)
		} else {
			t.Log("decrypt ok cost = ", time.Since(st).Seconds())
			if deText == srcText {
				t.Log("decrypt right ")
			} else {
				t.Error("decrypt wrong ")
			}
		}
	}

}
func TestSign(t *testing.T) {

	fmt.Println(runtime.NumCPU())

	srcText := "hello "
	for i := 0; i < 25; i++ {
		srcText += srcText
	}

	t.Logf("length  = %.2f M", float64(len(srcText))/1024./1024)

	pk, _ := LoadPemPubKey(strPubKey)
	sk, _ := LoadPemPKCSPriKey(strPriKey)
	st := time.Now()
	if sig, e := Sign(srcText, sk); e != nil {
		t.Error("sign fail", e)
	} else {
		t.Log("sign ok , cost = ", time.Since(st).Seconds())
		t.Logf(sig)
		st = time.Now()

		if b := Verify(srcText, sig, pk); b {
			t.Log("verify ok , cost = ", time.Since(st).Seconds())
		} else {
			t.Error("verify fail")
		}
	}
}
