module github.com/chunhui2001/go-starter

go 1.23.0

toolchain go1.24.2

exclude github.com/gomodule/redigo v2.0.0+incompatible

replace (
	github.com/99designs/gqlgen => github.com/99designs/gqlgen v0.17.29
	github.com/JGLTechnologies/gin-rate-limit => github.com/JGLTechnologies/gin-rate-limit v1.5.2
)

require (
	github.com/99designs/gqlgen v0.17.73
	github.com/JGLTechnologies/gin-rate-limit v1.5.6
	github.com/alitto/pond v1.9.2
	github.com/aws/aws-sdk-go v1.55.7
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/dgraph-io/ristretto v0.2.0
	github.com/dustin/go-humanize v1.0.1
	github.com/elastic/go-elasticsearch v0.0.0
	github.com/elastic/go-elasticsearch/v8 v8.18.0
	github.com/fatih/structs v1.1.0
	github.com/fogleman/fauxgl v0.0.0-20250110135958-abf826acbbbd
	github.com/fogleman/ribbon v0.0.0-20191101191537-568057efb726
	github.com/foolin/goview v0.3.0
	github.com/gfremex/logrus-kafka-hook v0.0.0-20180109031623-f62e125fcbfe
	github.com/gin-contrib/cache v1.4.0
	github.com/gin-contrib/gzip v1.2.3
	github.com/gin-contrib/sessions v1.0.4
	github.com/gin-contrib/static v1.1.5
	github.com/gin-gonic/gin v1.10.1
	github.com/go-errors/errors v1.5.1
	github.com/go-latex/latex v0.0.0-20250304174226-2790903426af
	github.com/go-playground/validator/v10 v10.26.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.9.2
	github.com/go-zookeeper/zk v1.0.4
	github.com/gobuffalo/events v1.4.3
	github.com/gobwas/ws v1.4.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/hhsnopek/etag v0.0.0-20171206181245-aea95f647346
	github.com/jinzhu/copier v0.4.0
	github.com/joho/godotenv v1.5.1
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/opensearch-project/opensearch-go v1.1.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron v1.2.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/shirou/gopsutil/v3 v3.24.5
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.20.1
	github.com/thoas/go-funk v0.9.3
	github.com/ubiq/go-ubiq v3.0.1+incompatible
	github.com/vektah/gqlparser/v2 v2.5.27
	github.com/xuri/excelize/v2 v2.9.1
	go.mongodb.org/mongo-driver v1.17.3
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6
	golang.org/x/oauth2 v0.30.0
	google.golang.org/api v0.234.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
	moul.io/http2curl v1.0.0
)

require (
	cloud.google.com/go/auth v0.16.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	git.sr.ht/~sbinet/gg v0.6.0 // indirect
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20250403215159-8d39553ac7cf // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/campoy/embedmd v1.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.7.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fogleman/ease v0.0.0-20170301025033-8da417bf1776 // indirect
	github.com/fogleman/simplify v0.0.0-20170216171241-d32f302d5046 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.14.2 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250317134145-8bc96cf8fc35 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/memcachier/mc/v3 v3.0.3 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/robfig/go-cache v0.0.0-20130306151617-9fc39e0dbf62 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/smartystreets/goconvey v1.7.2 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tiendc/go-deepcopy v1.6.0 // indirect
	github.com/tj/assert v0.0.3 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.14 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.36.0 // indirect
	go.opentelemetry.io/otel/metric v1.36.0 // indirect
	go.opentelemetry.io/otel/trace v1.36.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.17.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250519155744-55703ea1f237 // indirect
	google.golang.org/grpc v1.72.2 // indirect
)
