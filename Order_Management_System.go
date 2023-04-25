package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
)

func main() {

	db := initDB()

	defer db.Close()

	r := gin.Default()

	// 用户登录√
	r.POST("/login", func(c *gin.Context) {
		loginUser(c, db)
	})

	// 用户注册√(无用户去重)
	r.POST("/register", func(c *gin.Context) {
		registerUser(c, db)
	})

	// 修改密码√
	r.PUT("/update_password", func(c *gin.Context) {
		updatePassword(c, db)
	})

	// 删除用户（逻辑删除）√
	r.PUT("/delete_user", func(c *gin.Context) {
		deleteUser(c, db)
	})

	// 加入购物车/外卖下单(合并，方便)√
	r.POST("/add_order", jwtAuthMiddleware(), func(c *gin.Context) {
		addOrder(c, db)
	})

	// 订单送达√
	r.PUT("/order_delivered", func(c *gin.Context) {
		orderDelivered(c, db)
	})

	// 取消订单√
	r.PUT("/cancel_order", func(c *gin.Context) {
		cancelOrder(c, db)
	})

	// 添加评论
	r.POST("/add_review", func(c *gin.Context) {
		addReview(c, db)
	})
	r.Run(":8080")
}

// 初始化
func initDB() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/go_db")
	if err != nil {
		panic(err)
	}
	return db
}

// 登录加密
func generateToken(username string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{}
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

// // 验证登录加密
// func parseToken(tokenString string) (jwt.MapClaims, error) {
// 	// Parse token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte("your-secret-key"), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Verify token and get claims
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	} else {
// 		return nil, fmt.Errorf("invalid token")
// 	}
// }

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		fmt.Println(tokenString)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}

// 登录
func loginUser(c *gin.Context, db *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var storedPassword string
	var userid int
	err := db.QueryRow("SELECT userid, password FROM userinfo WHERE username = ? AND isdelete = 0", username).Scan(&userid, &storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	if password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	// 生成token
	token, err := generateToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "登录失败"})
		return
	}

	// 将token添加到响应头
	c.Writer.Header().Set("Authorization", "Bearer "+token)
	fmt.Println(token)

	c.JSON(http.StatusOK, gin.H{"userid": userid, "username": username, "message": "登录成功"})
}

// 改密码
func updatePassword(c *gin.Context, db *sql.DB) {
	userid := c.PostForm("userid")
	newPassword := c.PostForm("new_password")

	_, err := db.Exec("UPDATE userinfo SET password = ? WHERE userid = ?", newPassword, userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码已更新"})
}

// 注册
func registerUser(c *gin.Context, db *sql.DB) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, err := db.Exec("INSERT INTO userinfo (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 删除用户
func deleteUser(c *gin.Context, db *sql.DB) {
	userid := c.PostForm("userid")

	_, err := db.Exec("UPDATE userinfo SET isdelete = 1 WHERE userid = ?", userid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// 添加订单/购物车
func addOrder(c *gin.Context, db *sql.DB) {
	// 验证登录
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
		return
	}

	userid := c.PostForm("userid")
	shopid := c.PostForm("shopid")

	// 先查询shops表中是否有该shopid对应的数据
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM shops WHERE shopid = ?", shopid).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询店铺信息失败"})
		return
	}

	if count == 0 { // 如果没有，则插入该数据
		_, err = db.Exec("INSERT INTO shops (shopid, shopname, rating) VALUES (?, ?, ?)", shopid, "未知店铺", 5.0)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "添加店铺信息失败"})
			return
		}
	}

	// 插入订单数据
	_, err = db.Exec("INSERT INTO orders (userid, shopid) VALUES (?, ?)", userid, shopid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已创建", "username": username})
}

// 订单送达
func orderDelivered(c *gin.Context, db *sql.DB) {
	orderid := c.PostForm("orderid")

	_, err := db.Exec("UPDATE orders SET status = '已完成' WHERE orderid = ?", orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "订单送达失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已送达"})
}

// 取消订单
func cancelOrder(c *gin.Context, db *sql.DB) {
	orderid := c.PostForm("orderid")

	_, err := db.Exec("UPDATE orders SET status = '已取消' WHERE orderid = ?", orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "取消订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "订单已取消"})
}

func addReview(c *gin.Context, db *sql.DB) {
	orderid := c.PostForm("orderid")
	userid := c.PostForm("userid")
	content := c.PostForm("content")
	rating := c.PostForm("rating")

	_, err := db.Exec("INSERT INTO reviews (orderid, userid, content, rating) VALUES (?, ?, ?, ?)", orderid, userid, content, rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "添加评价失败"})
		return
	}

	// 更新订单状态为已完成
	_, err = db.Exec("UPDATE orders SET status = '已完成' WHERE orderid = ?", orderid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新订单状态失败"})
		return
	}
	///////////////////////////////////////
	// 查询订单所属的店铺
	var shopid int
	err = db.QueryRow("SELECT shopid FROM orders WHERE orderid = ?", orderid).Scan(&shopid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询订单所属店铺失败"})
		return
	}

	// 计算店铺的评分
	var count int
	var sum int
	err = db.QueryRow("SELECT COUNT(reviewid), SUM(rating) FROM reviews WHERE orderid IN (SELECT orderid FROM orders WHERE shopid = ? AND status = '已完成')", shopid).Scan(&count, &sum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "计算店铺评分失败"})
		return
	}

	var score float32
	if count > 0 {
		score = float32(sum) / float32(count)
	}

	// 更新店铺评分
	_, err = db.Exec("UPDATE shops SET rating = ? WHERE shopid = ?", score, shopid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新店铺评分失败"})
		return
	}
	/////////////////////////////////
	c.JSON(http.StatusOK, gin.H{"message": "评价已添加"})
}
