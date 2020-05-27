package fabsdk

import (
	"fmt"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	// "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	// "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	// "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	// "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	// "fmt"
	// "time"
	// cfg "github.com/jonluo94/cool/config"
	// "github.com/jonluo94/cool/log"
	"encoding/hex"
	"strings"
)

func (f *FabricClient) GetKeyFile(id msp.SigningIdentity) (string, string) {
	priFile := hex.EncodeToString(id.PrivateKey().SKI()) + "_sk"
	pubFile := id.Identifier().ID + "@" + id.Identifier().MSPID + "-cert.pem"
	return priFile, pubFile
}

func (f *FabricClient) RegisterUser(userName string) (priFile string, pubFile string, err error) {
	//secret is userName+userOrg
	secret := userName + f.DefaultOrg
	mspClient, err := mspclient.New(f.sdk.Context(), mspclient.WithOrg(f.DefaultOrg))
	if err != nil {
		return
	}
	//判断是否存在
	id, err := mspClient.GetSigningIdentity(userName)
	if err == nil {
		priFile, pubFile = f.GetKeyFile(id)
		return
	}
	//注册用户
	request := &mspclient.RegistrationRequest{Name: userName, Type: "client", Secret: secret}
	_, err = mspClient.Register(request)
	if err != nil && !strings.Contains(err.Error(), "is already registered") {
		return
	}
	//登记保存证书到stores
	err = mspClient.Enroll(userName, mspclient.WithSecret(secret))
	if err != nil {
		return
	}

	id, _ = mspClient.GetSigningIdentity(userName)
	priFile, pubFile = f.GetKeyFile(id)
	fmt.Println(priFile, pubFile)
	return
}
