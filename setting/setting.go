package setting

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewSetting() (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigFile("/config/config.yaml")

	err := vp.ReadInConfig()
	if err != nil {
		fmt.Println("err in read conf")
		return nil, err
	}
	return vp, nil
}
