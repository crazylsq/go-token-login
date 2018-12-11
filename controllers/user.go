package controllers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"opserver/common"
	"time"
	"net/http"
	"log"
)
//用户信息
type User struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Username_cn  string `json:"username_cn"`
	Is_superuser int    `json:"is_superuser"`
	Email        string `json:"email"`
	Is_active    int    `json:"is_active"`
	Date_joined  int64  `json:"date_joined"`
}

/*
@Title CreateUser
@Description create users
@Param header Authorization string true
@Param body body username string true
@Param body password string true
@Param body Username_cn string false
@Param body Is_superuser int false
@Param body Email string true
@Param body Is_active int false
@Param body Date_joined int64 false
@Success 200 {object} 用户添加成功
@Failure 201 {object} 用户已存在
@router /user/add [POST]
*/
func Add(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	passwd := common.Encrypt(user.Password)
	db, _ := common.Db_Conn()
	date_joined := time.Now().Unix()
	stmt, _ := db.Prepare("insert into op_auth_user (username, password, email, date_joined) values (?, ?, ?, ?)")
	if _, err := stmt.Exec(user.Username, passwd, "allen.liu@meetsocial.cn", date_joined); err != nil {
		common.ResponseHandle(http.StatusCreated, "用户已存在", c)
	} else {
		common.ResponseHandle(200, "用户添加成功", c)
	}
}

/*
@Title delete user
@Description delete user
@Param header Authorization string true
@Param username string true
@Success 200 {object} 已删除
@Failure 255 {object} 用户不存在
@router /user/delete [DELETE]
*/
func Delete(c *gin.Context) {
	username := c.Query("username")
	db, _ := common.Db_Conn()
	stmt, _ := db.Prepare("delete from op_auth_user where username = ?")
	if _, err := stmt.Exec(username); err != nil {
		log.Print(err)
		common.ResponseHandle(255, "用户不存在", c)
	} else {
		common.ResponseHandle(200, "已删除", c)
	}

}

/*
@Title update user
@Description update user
@Param header Authorization string true
@Param body username string true
@Param body email string true
@Param body password string true
@Success 200 {object} 已更新
@Failure 255 {object} 用户不存在
@router /user/delete [PUT]
*/
func Update(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	passwd := common.Encrypt(user.Password)
	db, _ := common.Db_Conn()
	stmt, _ := db.Prepare("update op_auth_user set password = ?, email = ? where username = ?")
	if _, err := stmt.Exec(passwd, user.Email, user.Username); err != nil {
		log.Print(err)
		common.ResponseHandle(255, "用户不存在", c)
	} else {
		common.ResponseHandle(200, "已更新", c)
	}
}

/*
@Title getAll
@Description view user info
@Param header Authorization string true
@Param username string true
@Success 200 {object} User
@Failure 255 用户不存在
@router /user/delete [PUT]
*/
func View(c *gin.Context) {
	username := c.Query("username")
	db, _ := common.Db_Conn()
	var email, username_cn string
	var is_superuser, is_active, date_joined int
	row := db.QueryRow("select email, username_cn, is_superuser, is_active, date_joined"+
		" from op_auth_user where username = ?", username)
	if err := row.Scan(&email, &username_cn, &is_active, &is_superuser, &date_joined); err != nil {
		log.Printf("数据库查询err: %s", err)
		common.ResponseHandle(255, "用户不存在", c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"username":     username,
			"email":        email,
			"username_cn":  username_cn,
			"is_active":    is_active,
			"is_superuser": is_superuser,
			"date_joined":  date_joined,
		})
	}

}
