package wxpay

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strings"
)

/*

 * 对密文进行解密.
 *
 * @param text 需要解密的密文
 * @return 解密得到的明文
 * @throws AesException aes解密失败
 */
func Decrypt(text string) (bool, string) {
	bkey, err := base64.StdEncoding.DecodeString(EncodingAESKey + "=")

	//aesKey := Base64.decodeBase64(encodingAesKey + "=");
	block, err := aes.NewCipher(bkey) //选择加密算法
	if err != nil {
		return false, ""
	}
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(rand.Reader, iv)

	ciphertext, err := base64.StdEncoding.DecodeString(text)

	blockModel := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)

	//fmt.Println(len(plantText))

	buf := bytes.NewBuffer(plantText[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	//fmt.Println(string(plantText[20:20+length]))
	//fmt.Println(string(plantText))
	content := string(plantText[20 : 20+length])

	appIDStart := 20 + length
	mAppID := string(plantText[appIDStart : int(appIDStart)+len(AppID)])
	if !strings.EqualFold(mAppID, AppID) {
		return false, ""
	}

	//plantText = PKCS7UnPadding(plantText, block.BlockSize())
	return true, content
}

func DecryptMsg(msgSignature string, timeStamp string, nonce string, postData string) (bool, string) {
	// 密钥，公众账号的app secret
	// 提取密文
	//Object[] encrypt = XMLParse.extract(postData);

	// 验证安全签名
	signature := getSHA1(Token, timeStamp, nonce, postData)
	//fmt.Println(signature)

	// 和URL中的签名比较是否相等
	// System.out.println("第三方收到URL中的签名：" + msg_sign);
	// System.out.println("第三方校验签名：" + signature);
	if !strings.EqualFold(signature, msgSignature) {
		return false, ""
	}
	// 解密
	//String result = decrypt(encrypt[1].toString());
	//return result;
	return Decrypt(postData)
}

/**
 * 用SHA1算法生成安全签名
 * @param token 票据
 * @param timestamp 时间戳
 * @param nonce 随机字符串
 * @param encrypt 密文
 * @return 安全签名
 * @throws AesException
 */
func getSHA1(token, timestamp, nonce, encrypt string) string {

	array := []string{timestamp, nonce, encrypt, token}
	sb := ""
	// 字符串排序
	sort.Strings(array)
	//fmt.Println(array)
	for i := 0; i < len(array); i++ {
		sb = sb + array[i]
	}
	// SHA1签名生成
	h := sha1.New()
	io.WriteString(h, sb)
	return fmt.Sprintf("%x", h.Sum(nil))
}
