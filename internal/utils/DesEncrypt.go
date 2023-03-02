package utils

import "github.com/golang-module/dongle"

func DesECBEncrypt(message, key string) string {
	decode := dongle.Decode.FromString(key).ByBase64().ToString()
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.ECB)      // CBC、ECB、CFB、OFB、CTR
	cipher.SetPadding(dongle.PKCS7) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(decode)           // key 长度必须是 8 字节
	encode := dongle.Encrypt.FromString(message).ByDes(cipher).ToBase64String()
	return encode
}
