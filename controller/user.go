package controller

import (
	"fmt"
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
	// 1. 先生成验证码
	verificationCode := utils.GenerateCode()

	// 2. 先存入数据库 (这样如果数据库挂了，就不会发邮件了)
	verificationRecord := model.VerificationCode{
		Email:     req.Email,
		Code:      verificationCode,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	result := config.DB.Create(&verificationRecord)
	if result.Error != nil {
		// 打印数据库错误
		println("保存验证码失败:", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "系统错误: 保存验证码失败",
		})
		return
	}

	// 3. 再发送邮件
	// 使用 HTML 模板，保留原有文案风格并增加可爱元素
	emailBody := fmt.Sprintf(`
		<div style="background-color: #fff0f5; padding: 20px; font-family: 'Microsoft YaHei', sans-serif;">
			<div style="max-width: 600px; margin: 0 auto; background: #fff; padding: 30px; border-radius: 15px; box-shadow: 0 4px 15px rgba(255,182,193,0.3); border: 2px solid #ffb6c1;">
				<h4 style="color: #ff69b4; text-align: center;">✨ Ciallo～ (∠・ω< )⌒★ </h4>
				<p style="font-size: 16px; color: #666; line-height: 1.6;">
					亲爱的喵喵：<br>
					这里是小学生的 Growth 服务喵！收到您的注册请求啦~ (QwQ)
				</p>
				<div style="background-color: #fff5f7; padding: 20px; text-align: center; border-radius: 10px; margin: 25px 0; border: 1px dashed #ffb6c1;">
					<p style="color: #ff69b4; margin: 0 0 10px 0; font-size: 14px;">您的喵喵验证码是：</p>
					<span style="font-size: 32px; font-weight: bold; color: #ff1493; letter-spacing: 6px; text-shadow: 1px 1px 2px #ffd1dc;">%s</span>
				</div>
				<p style="font-size: 14px; color: #888;">
					⏰ 有效期只有 10 分钟哦！请尽快使用喵~<br>
					(如非本人操作，请无视这封邮件，继续睡觉觉吧~ 💤)
				</p>
				<hr style="border: none; border-top: 1px dashed #ffb6c1; margin: 20px 0;">
				<p style="font-size: 12px; color: #aaa; text-align: center;">
					Growth  敬上 🐾<br>
					<span style="font-size: 10px;">(邮件由系统自动发送，回复也不会有猫猫理你哦~)</span>
				</p>
			</div>
		</div>
	`, verificationCode)

	err = utils.SendEmail(req.Email, "[Growth] 您的注册验证码来啦 ( >ω<)♡", emailBody)
	if err != nil {
		println("发送邮件失败:", err.Error())
		// 如果邮件发送失败，可以选择把刚才存的验证码删掉，或者不管它
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "发送邮件失败",
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
		Code     string `json:"code" binding:"required"` // 新增验证码字段
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求体参数错了啊喂",
		})
		return
	}

	// 1. 校验验证码
	var verificationCode model.VerificationCode
	// 查询数据库中是否存在该邮箱和验证码的记录
	if err := config.DB.Where("email = ? AND code = ?", req.Email, req.Code).First(&verificationCode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误或已被使用",
		})
		return
	}

	// 检查验证码是否过期
	if time.Now().After(verificationCode.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码已过期，请重新获取",
		})
		return
	}

	// 2. 检查邮箱是否已注册
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

	// 3. 创建用户
	result = config.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "注册用户失败",
		})
		return
	}

	// 4. 注册成功后自动登录（生成 Token）
	token, err := utils.GenerateAuthToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":     "注册成功，但自动登录失败，请手动登录",
			"user_id": newUser.ID,
		})
		return
	}

	// 5.验证码使用后删除，防止重复使用
	config.DB.Delete(&verificationCode)

	c.JSON(http.StatusOK, gin.H{
		"msg":     "注册成功",
		"user_id": newUser.ID,
		"token":   token,
	})
}

// 登录
func Login(c *gin.Context) {
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

	// 2. 根据邮箱找用户
	var user model.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号不存在"})
		return
	}

	// 3. 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "密码错误"})
		return
	}

	// 4. 生成 Token
	token, err := utils.GenerateAuthToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成Token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "登录成功",
		"token": token,
	})
}
