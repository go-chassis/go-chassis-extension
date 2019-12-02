module github.com/go-chassis/go-chassis-plugins

require (
	github.com/DataDog/zstd v1.3.5 // indirect
	github.com/Shopify/sarama v1.20.1 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chassis/go-chassis v1.7.7-0.20191202025124-badc15f68f1c
	github.com/go-chassis/huawei-apm v0.0.0-20190315045100-5b80092faa2d
	github.com/go-mesh/openlogging v1.0.1
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c // indirect
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/hashicorp/golang-lru v0.5.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/stretchr/testify v1.3.0
	golang.org/x/time v0.0.0-20180412165947-fbb02b2291d2 // indirect
	google.golang.org/genproto v0.0.0-20181221175505-bd9b4fb69e2f // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20180925152912-a191abe0b71e
	k8s.io/apimachinery v0.0.0-20181108045954-261df694e725
	k8s.io/client-go v9.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20181106182614-a9a16210091c // indirect
)

replace github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92

go 1.13
