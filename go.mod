module fabric-client

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.26.4 // indirect
	github.com/cloudflare/cfssl v0.0.0-20180223231731-4e2dcbde5004
	github.com/fsouza/go-dockerclient v1.6.5 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hyperledger/fabric v1.4.3

	github.com/hyperledger/fabric-amcl v0.0.0-20200424173818-327c9e2cf77a // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0-alpha5.0.20190411180201-5a9a0e749e4f
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190411180201-5a9a0e749e4f

	github.com/jonluo94/cool v0.0.0-20200518072032-4e2c0df52183
	github.com/pkg/errors v0.8.1
	github.com/sykesm/zap-logfmt v0.0.3 // indirect
	go.uber.org/zap v1.15.0 // indirect
)

// replace github.com/hyperledger/fabric => github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190411180201-5a9a0e749e4f
