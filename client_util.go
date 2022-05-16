package main

import (
	"github.com/3bl3gamer/tgclient/mtproto"
)

func EncodeBool(b bool) mtproto.TL {
	if b {
		return mtproto.TL_boolTrue{}
	} else {
		return mtproto.TL_boolFalse{}
	}
}
