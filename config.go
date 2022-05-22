package kcp

import (
	xkcp "github.com/xtaci/kcp-go/v5"
)

const (
	Normal = "normal"
	Fast   = "fast"
	Fast2  = "fast2"
	Fast3  = "fast3"
	Manual = "manual"
)

// Config KCP
// mode: normal, fast, fast2, fast3，manual
// key: the password needed to decrypt for AES, must given if crypt is not none
// crypt: aes, aes-128, aes-192, salsa20, blowfish, twofish, cast5, 3des, tea, xtea, xor, sm4, none (default: "aes")
// for more config options please refer to kcptun:https://github.com/xtaci/kcptun
type Config struct {
	Key          string `json:"key"`
	Crypt        string `json:"crypt"`
	Mode         string `json:"mode"`
	MTU          int    `json:"mtu"`
	SndWnd       int    `json:"sndwnd"`
	RcvWnd       int    `json:"rcvwnd"`
	DataShard    int    `json:"datashard"`
	ParityShard  int    `json:"parityshard"`
	DSCP         int    `json:"dscp"`
	AckNodelay   bool   `json:"acknodelay"`
	NoDelay      int    `json:"nodelay"`
	Interval     int    `json:"interval"`
	Resend       int    `json:"resend"`
	NoCongestion int    `json:"nc"`
	SockBuf      int    `json:"sockbuf"`
	TCP          bool   `json:"tcp"`
}

// Default config
func (c *Config) fixture() {
	if c.Mode == "" {
		c.Mode = Fast
	}
	if c.Crypt == "" {
		c.Crypt = "aes"
	}
	if c.Key == "" {
		c.Key = "test"
	}
	switch c.Mode {
	case Normal:
		c.NoDelay, c.Interval, c.Resend, c.NoCongestion = 0, 40, 2, 1
	case Fast:
		c.NoDelay, c.Interval, c.Resend, c.NoCongestion = 0, 30, 2, 1
	case Fast2:
		c.NoDelay, c.Interval, c.Resend, c.NoCongestion = 1, 20, 2, 1
	case Fast3:
		c.NoDelay, c.Interval, c.Resend, c.NoCongestion = 1, 10, 2, 1
	case Manual:
		// nothing to do
	default:
		// default is fast mode
		c.NoDelay, c.Interval, c.Resend, c.NoCongestion = 0, 30, 2, 1
	}
	if c.MTU == 0 {
		c.MTU = 1350
	}
	if c.SndWnd == 0 {
		c.SndWnd = 2048
	}
	if c.RcvWnd == 0 {
		c.RcvWnd = 2048
	}
	if c.DataShard == 0 {
		c.DataShard = 10
	}
	if c.ParityShard == 0 {
		c.ParityShard = 3
	}
	if c.SockBuf == 0 {
		c.SockBuf = 4194304
	}
}

func getCryptBlock(crypt, key string) (block xkcp.BlockCrypt, err error) {
	pass := []byte(key)
	switch crypt {
	case "sm4":
		block, err = xkcp.NewSM4BlockCrypt(pass[:16])
	case "tea":
		block, _ = xkcp.NewTEABlockCrypt(pass[:16])
	case "xor":
		block, _ = xkcp.NewSimpleXORBlockCrypt(pass)
	case "none":
		block, _ = xkcp.NewNoneBlockCrypt(pass)
	case "aes-128":
		block, _ = xkcp.NewAESBlockCrypt(pass[:16])
	case "aes-192":
		block, _ = xkcp.NewAESBlockCrypt(pass[:24])
	case "blowfish":
		block, _ = xkcp.NewBlowfishBlockCrypt(pass)
	case "twofish":
		block, _ = xkcp.NewTwofishBlockCrypt(pass)
	case "cast5":
		block, _ = xkcp.NewCast5BlockCrypt(pass[:16])
	case "3des":
		block, _ = xkcp.NewTripleDESBlockCrypt(pass[:24])
	case "xtea":
		block, _ = xkcp.NewXTEABlockCrypt(pass[:16])
	case "salsa20":
		block, _ = xkcp.NewSalsa20BlockCrypt(pass)
	default:
		block, _ = xkcp.NewAESBlockCrypt(pass)
	}
	return block, err
}
