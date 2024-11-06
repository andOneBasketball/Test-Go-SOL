package main

/*
浏览器手动获取钱包签名测试
const message = "hello";
const encodedMessage = new TextEncoder().encode(message);
window.solana.signMessage(encodedMessage, 'utf8');

curl -l -H "Content-type: application/json" -X POST -d '{"message":"hello","public_key":"EVhXY9R5Ztyg5Jw9zAEGQ2C1tvJC5HHGi2HFHQf9HqCK", "signature": "tHEc8LTVAdmCK9hionDLG5BanhA4PzqTNvE2uXwy5BeYL4EEbDmJKLHUKe2DH1hGFxS7QuAX9jr674hZSQaSbxj"}' http://localhost:8080/login
*/

import (
	"net/http"

	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
)

// UserLoginRequest 表示用户登录请求的结构
type UserLoginRequest struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
	Message   string `json:"message"` // 可以是任何需要用户签名的消息
}

// main 函数启动 Gin 服务器
func main() {
	router := gin.Default()

	// 登录接口
	router.POST("/login", loginHandler)

	// 启动服务器
	router.Run(":8080")
}

// loginHandler 处理用户登录请求
func loginHandler(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 将公钥解析为 PublicKey
	pubKey, err := solana.PublicKeyFromBase58(req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid public key"})
		return
	}

	// 验证签名
	/*
		decodedBytes, err := base58.Decode(req.Message)
		if err != nil {
			log.Fatalf("Error decoding base58: %v", err)
		}
		fmt.Println(decodedBytes)
	*/

	msg := []byte(req.Message) // 将消息转为字节数组
	sig, err := solana.SignatureFromBase58(req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// 验证签名是否有效
	if !sig.Verify(pubKey, msg) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	// 登录成功，返回成功消息
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
