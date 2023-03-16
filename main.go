package main

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/FirwoodLin/BGM_tyro/router"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置信息
	err := setting.InitSettings()
	if err != nil {
		fmt.Println(err)
	}
	// 初始化数据库
	model.InitDB()
	// 开始监听
	r := gin.Default()
	router.NewRouter(r)
}
