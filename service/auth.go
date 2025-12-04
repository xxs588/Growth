package service

// Authable 接口定义了：任何想要被验证的东西，都必须有 CheckPassword 方法
// 这就是接口定义在“使用者”这边：VerifyLogin 函数需要用到 CheckPassword，所以它定义了这个接口
type Authable interface {
	CheckPassword(password string) bool
}

// VerifyLogin 是一个通用的登录验证函数
// 它不关心传进来的是 User 还是 Admin，只要能 CheckPassword 就行
func VerifyLogin(user Authable, password string) bool {
	return user.CheckPassword(password)
}
