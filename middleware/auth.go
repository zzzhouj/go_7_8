package middleware

import (
	"filesys/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 cookie 里获取 sid
		sid, err := c.Cookie("sid")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		// 2. 校验 sid 是否在 session 表中存在
		session, err := dao.Q.Session.WithContext(c.Request.Context()).
			Where(dao.Session.SessionID.Eq(sid)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效sid"})
			c.Abort()
			return
		}

		// 3. 校验 session 对应的用户是否存在
		user, err := dao.Q.User.WithContext(c.Request.Context()).
			Where(dao.User.ID.Eq(session.UserID)).First()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		// 4. 把用户信息注入 gin.Context，后续业务可以直接用
		c.Set("user_id", int(user.ID))
		c.Set("user_name", user.Name)
		c.Next()
	}
}
