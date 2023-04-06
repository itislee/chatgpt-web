package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"log"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

const (
	appID       = "101570536"
	appSecret   = ""
	redirectURI = "https://gpt.itislee.com/callback"
)
//DB 
const (
	username = ""
	password = ""
	dbname = "qqauth"
	tablename = "auth_openid"
)

func main() {
	log.Printf("server started")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//init database
	sql := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", username, password, dbname)
	//db, err := NewMySQLDatabase("username:password@tcp(localhost:3306)/dbname")
	db, err := NewMySQLDatabase(sql)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router.GET("/callback", func(c *gin.Context) { handleCallback(c, db) })

	router.GET("/check", func(c *gin.Context) { checkAccessToken(c, db) })

	router.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})

	log.Printf("server running")
	router.Run("127.0.0.1:8080")
}

func handleCallback(c *gin.Context, db Database) {
	log.Printf("handle callback")

	code := c.Query("code")
	if code == "" {
		c.String(http.StatusBadRequest, "Code is missing.")
		return
	}

	accessToken, err := getAccessToken(code)
	if err != nil {
		log.Printf("code=%s %v", code, err)
		c.String(http.StatusInternalServerError, "Failed to get access token.")
		return
	}

	openID, err := getOpenID(accessToken)
	if err != nil {
		log.Printf("accesstoken=%s %v", accessToken, err)
		c.String(http.StatusInternalServerError, "Failed to get openID.")
		return
	}

	c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	c.SetCookie("open_id", openID, 3600, "/", "", false, true)

	log.Printf("openid=%s token=%s", openID, accessToken)

	// set to database
	if err := db.UpdateAccessToken(openID, accessToken); err != nil {
		log.Printf("Failed to update access token for openID %s: %v", openID, err)
		c.String(http.StatusInternalServerError, "请联系yaoli添加权限.openid=%s", openID)
		return
	}
	log.Printf("db ok.login ok")

	//c.String(http.StatusOK, "Authorization successful.")
	log.Printf("Host=%s", c.Request.Host)
	c.Redirect(302, "https://"+c.Request.Host)
}

func checkAccessToken(c *gin.Context, db Database) {
	log.Printf("handle check")

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Missing access token.")
		return
	}

	openID, err := c.Cookie("open_id")
	if err != nil {
		c.String(http.StatusUnauthorized, "Missing openID.")
		return
	}

	// 这里可以实现验证 accessToken 和 openID 的逻辑
	if accessToken == "" || openID == "" {
		c.String(http.StatusUnauthorized, "Invalid access token or openID.")
		return
	}

	found, err := db.IsOpenIDExists(openID, accessToken)
	if err != nil {
		log.Println("Failed to query database:", err)
		c.String(http.StatusInternalServerError, "权限检查失败 请联系Yaoli添加权限.")
		return
	}

	if found {
		log.Printf("check ok.openid=%s cgi=%s %s", openID, c.Request.Header["X-Original-URI"], c.Request.URL.Path)
		c.String(http.StatusOK, "OpenID found in database.")
	} else {
		c.String(http.StatusUnauthorized, "OpenID not found in database.")
	}

	//c.String(http.StatusOK, "Access token and openID are valid.")
}

func getAccessToken(code string) (string, error) {
	tokenURL := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s", appID, appSecret, code, url.QueryEscape(redirectURI))

	log.Printf("tokenURL=%s", tokenURL)
	resp, err := http.Get(tokenURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	return values.Get("access_token"), nil
}

func getOpenID(accessToken string) (string, error) {
	openIDURL := fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", accessToken)

	resp, err := http.Get(openIDURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("before response=%s abc", string(body))
	response := strings.TrimPrefix(string(body), "callback(")
	response = strings.TrimSuffix(response, ");\n")

	var openIDResponse struct {
		ClientID string `json:"client_id"`
		OpenID string `json:"openid"`
	}

	log.Printf("after response=%s", response)
	if err := json.Unmarshal([]byte(response), &openIDResponse); err != nil {
		log.Printf("get openid unmarshal failed", err)
		return "", err
	}

	return openIDResponse.OpenID, nil
}
