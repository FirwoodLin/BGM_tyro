package setting

import (
	"github.com/spf13/viper"
)

func InitSettings() error {
	vp := viper.New()
	vp.SetConfigFile("config.yaml") /// ./config/
	if err := vp.ReadInConfig(); err != nil {
		return err
		//fmt.Println(err)
	}
	if err := vp.UnmarshalKey("mysql", &DatabaseSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	if err := vp.UnmarshalKey("JWT", &JWTSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	if err := vp.UnmarshalKey("mail", &MailSettings); err != nil {
		// if err := vp.UnmarshalKey("mysql.username", &s); err != nil {
		//fmt.Println(err)
		return err
	}
	//fmt.Printf("in initsettings below ")
	//fmt.Println(DatabaseSettings)
	return nil
}
