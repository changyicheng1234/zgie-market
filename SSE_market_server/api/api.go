// 调用api的文件，不用管，要看也行
package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tms/v20201229"
)

func GetSuggestion(inStr string) string {
	credential := common.NewCredential(
		viper.GetString("tms.secretId"),
		viper.GetString("tms.secretKey"),
	)

	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tms.NewClient(credential, "ap-guangzhou", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tms.NewTextModerationRequest()
	request.Content = common.StringPtr(base64.StdEncoding.EncodeToString([]byte(inStr)))

	// 返回的resp是一个TextModerationResponse的实例，与请求对象对应
	response, err := client.TextModeration(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	if *response.Response.Label == "Ad" {
		return "Review"
	}
	return *response.Response.Suggestion
}

func ApiTest(c *gin.Context) {
	var inData struct {
		InputVal string `json:"inputVal"`
	}
	c.Bind(&inData)
	suggestion := GetSuggestion(inData.InputVal)
	c.JSON(http.StatusOK, gin.H{"outputVal": suggestion})
}
