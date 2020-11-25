module github.com/nevzatcirak/go-examples/config-server

go 1.15

require (
	github.com/nevzatcirak/go-examples/config-server/config v0.0.0
	github.com/Piszmog/cfservices v1.4.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/subosito/gotenv v1.2.0
)

replace github.com/nevzatcirak/go-examples/config-server/config => ./config
