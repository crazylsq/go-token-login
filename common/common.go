package common

import (
	"fmt"
	"database/sql"
	"opserver/config"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"time"
	"strconv"
)


var commonIV = "meetsocial@,2018"
var key_text = "astaxie12798akljzmkme.ahkjkljl;k"    //aes的加密字符串

func Db_Conn() (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		config.GetValue("mysql","user") + ":" + config.GetValue("mysql","password") +
			"@tcp(" + config.GetValue("mysql", "host") + ":" + config.GetValue("mysql","port") + ")/" + config.GetValue("mysql", "db") + "?charset=" + config.GetValue("mysql", "charset"))
	if err != nil {
		fmt.Printf("数据库连接err: %s", err)
		return db, err
	} else {
		return db,err
	}
}

//加密字符串
func AesEncrypt(base64EncodeString string) (err error , encodeContent string) {
	content := base64EncodeString
	text := []byte(content)

	//创建加密算法aes
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		return err, ""
	}

	//加密字符串
	ciphertext := make([]byte, len(text))
	cfb := cipher.NewCFBDecrypter(c, []byte(commonIV))

	cfb.XORKeyStream(ciphertext, text)
	fmt.Printf("%s=>%x\n", text, ciphertext)
	AesEncryptString := fmt.Sprintf("%s", ciphertext)
	return err, AesEncryptString
}

func AesDecrypt(content string) (err error, passwd string) {
	c, err := aes.NewCipher([]byte(key_text))
	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(key_text), err)
		return err, ""
	}

	//解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, []byte(commonIV))
	plaintextCopy := make([]byte, len([]byte(content)))
	//cfbdec.XORKeyStream(plaintextCopy, []byte(content))
	cfbdec.XORKeyStream(plaintextCopy, make([]byte, len(content)))
	fmt.Printf("%x=>%s\n", []byte(content), plaintextCopy)
	EncryptContent := fmt.Sprintf("%s",plaintextCopy)
	return nil, EncryptContent
}

func Base64Encode(content string) string {
	input := []byte(content)
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

func Base64Decode(content string) (error, string) {
	decodeBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		fmt.Printf("base64 decode err: %s", err)
		return err, ""
	}

	return nil, string(decodeBytes)
}

func Encrypt(content string) string {
	return Base64Encode(content)
}

func Decrypt(EncryptString string) (error, string) {
	return Base64Decode(EncryptString)
}

func ResponseHandle(code int, message string, c *gin.Context)  {
	c.JSON(code, gin.H{
		"code": code,
		"msg": message,
	})
	c.AbortWithStatus(code)
}

func CreateToken(username string) string {
	mySigningKey := []byte(config.GetValue("token", "secret"))

	type MyCustomClaims struct {
		User string `json:"user"`
		jwt.StandardClaims
	}

	// Create the Claims
	expire_time, _ := strconv.Atoi(config.GetValue("token", "expire_time"))
	claims := MyCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Unix()) + int64(expire_time),
			IssuedAt: int64(time.Now().Unix()),
			Issuer:    username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(mySigningKey)
	//fmt.Printf("%v %v", tokenString, err)
	return tokenString
}
