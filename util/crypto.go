package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

type CryptWriter struct {
	md5f  bool
	sha1f bool
	sha2f bool

	md5  string
	sha1 string
	sha2 string

	md5d  hash.Hash
	sha1d hash.Hash
	sha2d hash.Hash
}

func (wc *CryptWriter) Init(p []byte) {
	if wc.md5f {
		wc.md5d = md5.New()
	}
	if wc.sha1f {
		wc.sha1d = sha1.New()
	}
	if wc.sha2f {
		wc.sha2d = sha256.New()
	}
}

func (wc *CryptWriter) Write(p []byte) (int, error) {
	if wc.md5f {
		wc.md5d.Write(p)
	}
	if wc.sha1f {
		wc.sha1d.Write(p)
	}
	if wc.sha2f {
		wc.sha2d.Write(p)
	}
}

func (wc *CryptWriter) Sum(p []byte) (int, error) {
	if wc.md5f {
		wc.md5 = string(wc.md5d.Sum(nil)[:])
	}
	if wc.sha1f {
		wc.sha1 = string(wc.sha1d.Sum(nil)[:])
	}
	if wc.sha2f {
		wc.sha2 = string(wc.sha2d.Sum(nil)[:])
	}
}
