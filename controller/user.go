package controller

import (
	"mygo/config"
	"mygo/model"
	"mygo/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 发送验证码
func SendCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求体参数错了啊喂",
		})
		return
	}
	verificationCode := utils.GenerateCode()
	emailBody := "Ciallo～ (∠・ω< )⌒★,您的验证码是：" + verificationCode + "，有效期10分钟喵~..."
	err = utils.SendEmail(req.Email, "[Growth]xxs的服务 注册验证码", emailBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送邮件失败",
		})
		return
	}
	verificationRecord := model.VerificationCode{
		Email:     req.Email,
		Code:      verificationCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	result := config.DB.Create(&verificationRecord)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "保存验证码失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "验证码发送成功",
	})
}

// 注册
func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求体参数错了啊喂",
		})
		return
	}

	var existingUser model.User
	result := config.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"msg": "该邮箱已被注册",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "密码加密失败",
		})
		return
	}

	newUser := model.User{
		Email:    req.Email,
		Password: hashedPassword,
		Name:     req.Name,
	}

	result = config.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "注册用户失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":     "注册成功",
		"user_id": newUser.ID,
	})
}

// 登录
func Login(c *gin.Context) {

	var user model.User

	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name"`
	}

	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号不存在"})
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求体参数错了啊喂",
		})
		return
	}
	result := config.DB.Where("email = ? AND password = ?", req.Email, req.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "邮箱或密码错误",
		})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "密码错误"})
		return
	}

	token, err := utils.GenerateAuthToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "未携带token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
