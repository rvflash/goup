module github.com/rvflash/goup

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3
	github.com/gemnasium/logrus-graylog-hook v2.0.7+incompatible
	github.com/gin-gonic/gin v1.3.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/mock v1.4.0
	github.com/golang/protobuf v1.3.2
	github.com/google/wire v0.3.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.4
	github.com/stretchr/testify v1.3.0
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2
	google.golang.org/appengine v1.6.0 // indirect
	google.golang.org/grpc v1.23.1
)

exclude github.com/golang/mock v1.5.0

replace google.golang.org/grpc => ../tree
