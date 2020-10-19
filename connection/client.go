package connection

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type TencentCloudClient struct {
	Credential *common.Credential
	Region     string
	Protocol   string
	Domain     string

	cvmConn *cvm.Client
}

// NewClientProfile returns a new ClientProfile
func (me *TencentCloudClient) NewClientProfile(timeout int) *profile.ClientProfile {
	cpf := profile.NewClientProfile()

	// all request use method POST
	cpf.HttpProfile.ReqMethod = "POST"
	// request timeout
	cpf.HttpProfile.ReqTimeout = timeout
	// request protocol
	cpf.HttpProfile.Scheme = me.Protocol
	// request domain
	cpf.HttpProfile.RootDomain = me.Domain
	// default language
	cpf.Language = "en-US"

	return cpf
}

// UseCvmClient returns cvm client for service
func (me *TencentCloudClient) UseCvmClient() *cvm.Client {
	if me.cvmConn != nil {
		return me.cvmConn
	}

	cpf := me.NewClientProfile(300)
	me.cvmConn, _ = cvm.NewClient(me.Credential, me.Region, cpf)
	return me.cvmConn
}
