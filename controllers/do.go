package controllers

import (
	"opserver/common"
	"opserver/config"
	"log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//登陆
type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

/*
@Title Login
@Description Login
@Param body body username string true
@Param body password string true
@Success 200 {object} token
@Failure 201 {object} error info
@router /do/login [POST]
*/
func Login(c *gin.Context) {
	var logininfo LoginInfo
	c.BindJSON(&logininfo)
	if logininfo.Username == "" || logininfo.Password == "" {
		common.ResponseHandle(255, "用户名/密码不能为空", c)
	} else {
		db, _ := common.Db_Conn()
		var passwd string
		row := db.QueryRow("select password from op_auth_user where username = ?", logininfo.Username)
		if err := row.Scan(&passwd); err != nil {
			log.Printf("数据库查询err: %s", err)
			common.ResponseHandle(255, "用户不存在", c)
		} else {
			_, decryptString := common.Decrypt(passwd)
			if logininfo.Password == decryptString {
				token := common.CreateToken(logininfo.Username)
				common.ResponseHandle(200, token, c)
			} else {
				common.ResponseHandle(255, "用户名/密码不正确", c)
			}
		}
	}
}

/*
@Title verify token
@Description verify token
@Param header Authorization string true
@Success 200 {object} ok
@Failure 401 {object} token error
@router /do/verifytoken [GET]
*/
func LoginTokenVerify(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		common.ResponseHandle(401, "token不存在", c)
	} else {
		t, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetValue("token", "secret")), nil
		})
		if t.Valid {
			common.ResponseHandle(200, "ok", c)
		} else {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "token error",
			})
		}
	}
}
