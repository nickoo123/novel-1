package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"novel/app/services"
	"os"
	"os/exec"
	"time"
)

const (
	// 请求失败重试次数
	RETRY = 5
)

type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int8   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Gdrive struct {
	Type         string `json:"type"`
	Scope        string `json:"scope"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Token        struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Expiry      string `json:"expiry"`
	}
}

type OrigAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expiry      int64  `json:"expires_in"`
}

type AccessToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	Expiry      time.Time `json:"expiry"`
}

func Start() {
	c := cron.New()

	client_id := beego.AppConfig.String("GoogleDrive::client_id")
	client_secret := beego.AppConfig.String("GoogleDrive::client_secret")
	// get refreshtoken
	addCronFunc(c, "* * */5 * * *", func() {
		autocode := services.ConfigService.String("AuthCode")
		if autocode != "" {
			var str bytes.Buffer
			str.WriteString("https://accounts.google.com/o/oauth2/token")
			body := make(map[string]interface{})
			body["code"] = autocode
			body["client_id"] = client_id
			body["client_secret"] = client_secret
			body["redirect_uri"] = "https://www.biqugesk.cc/home/code.html"
			body["grant_type"] = "authorization_code"
			if response, err := resty.New().R().SetBody(body).Post(str.String()); err != nil {
				logrus.Error(err)
			} else {
				refresh := RefreshToken{}
				json.Unmarshal(response.Body(), &refresh)
				fmt.Println(string(response.Body()))
				if refresh.RefreshToken != "" {
					services.ConfigService.Set("RefreshToken", refresh.RefreshToken)
					services.ConfigService.Set("AuthCode", "")
				}
			}
		}
	})

	// refreshToken
	addCronFunc(c, "@every 30m", func() {
		refreshToken := services.ConfigService.String("RefreshToken")
		if refreshToken != "" {
			var str bytes.Buffer
			str.WriteString("https://accounts.google.com/o/oauth2/token")
			body := make(map[string]interface{})
			body["client_id"] = client_id
			body["client_secret"] = client_secret
			body["refresh_token"] = refreshToken
			body["grant_type"] = "refresh_token"
			if response, err := resty.New().R().SetBody(body).Post(str.String()); err != nil {
				logrus.Error(err)
			} else {
				origaccessToken := OrigAccessToken{}
				json.Unmarshal(response.Body(), &origaccessToken)
				f, err := os.Create("./rclone/rclone.conf")
				defer f.Close()
				if err != nil {
					fmt.Println(err.Error())
				} else {
					var txtstr bytes.Buffer
					txtstr.WriteString("[gdrive]\n")
					txtstr.WriteString("type = drive\n")
					txtstr.WriteString("scope = drive\n")
					txtstr.WriteString("client_id = ")
					txtstr.WriteString(client_id)
					txtstr.WriteString("\n")
					txtstr.WriteString("client_secret = ")
					txtstr.WriteString(client_secret)
					txtstr.WriteString("\n")
					txtstr.WriteString("token = ")
					currentTime := time.Now()
					m, _ := time.ParseDuration(fmt.Sprintf("%ds", origaccessToken.Expiry))
					endTimeTmp := currentTime.Add(m).Unix()
					endTime := time.Unix(endTimeTmp, 0)
					accessToken := AccessToken{}
					accessToken.AccessToken = origaccessToken.AccessToken
					accessToken.TokenType = origaccessToken.TokenType
					accessToken.Expiry = endTime
					token, _ := json.Marshal(accessToken)
					txtstr.WriteString(string(token))
					_, err = f.Write([]byte(txtstr.String()))
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	})

	c.Start()
}

// 这里为了简化，我省去了stderr和其他信息
func Command(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	// 此处是windows版本
	// c := exec.Command("cmd", "/C", cmd)
	output, err := c.CombinedOutput()
	fmt.Println(string(output))
	return err
}

func addCronFunc(c *cron.Cron, sepc string, cmd func()) {
	err := c.AddFunc(sepc, cmd)
	if err != nil {
		logrus.Error(err)
	}
}
