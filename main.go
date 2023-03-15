package main

import (
	"fmt"
	"github.com/FirwoodLin/BGM_tyro/model"
	"github.com/FirwoodLin/BGM_tyro/router"
	"github.com/FirwoodLin/BGM_tyro/setting"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	err := initSettings()
	if err != nil {
		fmt.Println(err)
		// TODO:错误处理
		//return err
	}
	model.InitDB()
	r := gin.Default()
	router.NewRouter(r)

}

func initSettings() error {
	vp := viper.New()
	vp.SetConfigFile("config.yaml") /// ./config/
	if err := vp.ReadInConfig(); err != nil {
		return err
		//fmt.Println(err)
	}
	if err := vp.UnmarshalKey("mysql", &setting.DatabaseSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	if err := vp.UnmarshalKey("JWT", &setting.JWTSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	fmt.Println(setting.DatabaseSettings)
	return nil
}
