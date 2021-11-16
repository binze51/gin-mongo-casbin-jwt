module gitee.com/binze/binzekeji

go 1.16

require (
	github.com/apache/pulsar-client-go v0.6.0
	github.com/casbin/casbin/v2 v2.39.0
	github.com/casbin/mongodb-adapter/v3 v3.2.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.4
	github.com/h2non/filetype v1.1.1
	github.com/json-iterator/go v1.1.10
	github.com/lestrrat-go/jwx v1.2.9
	github.com/qiniu/qmgo v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/tjfoc/gmsm v1.4.1
	go.mongodb.org/mongo-driver v1.7.4
	go.uber.org/ratelimit v0.2.0
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v2 v2.4.0

)

replace github.com/qiniu/qmgo => ../qmgo
