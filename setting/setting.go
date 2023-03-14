package setting

//package main

import (
	"fmt"
	"github.com/spf13/viper"
)

//var vp *viper.Viper

func NewSetting() (*viper.Viper, error) {
	vp := viper.New()
	//vp.AddConfigPath("./config")
	//vp.SetConfigFile("config")
	//vp.SetConfigType("json")
	vp.SetConfigFile("/config/config.yaml")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("err in read")
		return nil, err
	}
	return vp, nil
}
