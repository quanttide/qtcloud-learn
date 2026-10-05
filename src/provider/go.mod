module github.com/quanttide/qtcloud-learn-provider

go 1.26.0

require github.com/aliyun/aliyun-oss-go-sdk v3.0.2+incompatible

require github.com/quanttide/quanttide-learn-toolkit/packages/go v0.0.0

require (
	golang.org/x/time v0.15.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/quanttide/quanttide-learn-toolkit/packages/go => ../../../../packages/quanttide-learn-toolkit/packages/go
