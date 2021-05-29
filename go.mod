module github.com/spongeprojects/kubebigbrother

go 1.15

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/gin-gonic/gin v1.6.3
	github.com/muesli/termenv v0.8.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/spongeprojects/magicconch v0.0.6
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	gorm.io/driver/mysql v1.0.5
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.3
	k8s.io/apimachinery v0.21.0
	k8s.io/cli-runtime v0.21.0
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
)

replace k8s.io/klog/v2 => github.com/spongeprojects/klog/v2 v2.9.1-0.20210527182053-d23af0d5cbe7
