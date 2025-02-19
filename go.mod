module packetbeat6_dpdk

go 1.23.6

require (
        github.com/elastic/beats v7.6.2+incompatible
        github.com/elastic/gosigar v0.14.3
        github.com/golang/snappy v0.0.4
        github.com/insomniacslk/dhcp v0.0.0-20250109001534-8abf58130905
        github.com/magefile/mage v1.15.0
        github.com/miekg/dns v1.1.63
        github.com/njcx/gopacket_dpdk v0.0.0-20250217093055-cac11ccca30f
        github.com/njcx/libbeat_v6 v0.0.0-20250218075542-7b5fa43d54f5
        github.com/njcx/packetbeat6_dpdk v0.0.0-20250218073813-88519116329a
        github.com/pkg/errors v0.9.1
        github.com/samuel/go-thrift v0.0.0-20210915161234-7b67f98e972f
        github.com/spf13/cobra v1.9.1
        github.com/spf13/pflag v1.0.6
        github.com/stretchr/testify v1.10.0
        go.uber.org/zap v1.27.0
        golang.org/x/net v0.35.0
        golang.org/x/sys v0.30.0
        gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
        github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
        github.com/IBM/sarama v1.45.0 // indirect
        github.com/Microsoft/go-winio v0.6.2 // indirect
        github.com/containerd/containerd v1.7.25 // indirect
        github.com/containerd/errdefs v0.3.0 // indirect
        github.com/davecgh/go-spew v1.1.1 // indirect
        github.com/distribution/reference v0.6.0 // indirect
        github.com/dlclark/regexp2 v1.11.4 // indirect
        github.com/docker/distribution v2.7.1+incompatible // indirect
        github.com/docker/docker v20.10.24+incompatible // indirect
        github.com/docker/go-connections v0.5.0 // indirect
        github.com/docker/go-units v0.5.0 // indirect
        github.com/dop251/goja v0.0.0-20250125213203-5ef83b82af17 // indirect
        github.com/dop251/goja_nodejs v0.0.0-20250217171036-ba90ff8d8790 // indirect
        github.com/dustin/go-humanize v1.0.1 // indirect
        github.com/eapache/go-resiliency v1.7.0 // indirect
        github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
        github.com/eapache/queue v1.1.0 // indirect
        github.com/elastic/go-lumber v0.1.1 // indirect
        github.com/elastic/go-seccomp-bpf v1.5.0 // indirect
        github.com/elastic/go-structform v0.0.12 // indirect
        github.com/elastic/go-sysinfo v1.15.1 // indirect
        github.com/elastic/go-txfile v0.0.8 // indirect
        github.com/elastic/go-ucfg v0.8.8 // indirect
        github.com/elastic/go-windows v1.0.0 // indirect
        github.com/ericchiang/k8s v1.2.0 // indirect
        github.com/fatih/color v1.18.0 // indirect
        github.com/garyburd/redigo v1.6.4 // indirect
        github.com/ghodss/yaml v1.0.0 // indirect
        github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
        github.com/gofrs/flock v0.7.1 // indirect
        github.com/gofrs/uuid v4.4.0+incompatible // indirect
        github.com/gogo/protobuf v1.3.2 // indirect
        github.com/golang/protobuf v1.5.4 // indirect
        github.com/gorilla/mux v1.8.1 // indirect
        github.com/hashicorp/errwrap v1.1.0 // indirect
        github.com/hashicorp/go-multierror v1.1.1 // indirect
        github.com/hashicorp/go-uuid v1.0.3 // indirect
        github.com/inconshreveable/mousetrap v1.1.0 // indirect
        github.com/jcmturner/aescts/v2 v2.0.0 // indirect
        github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
        github.com/jcmturner/gofork v1.7.6 // indirect
        github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
        github.com/jcmturner/rpc/v2 v2.0.3 // indirect
        github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
        github.com/jstemmer/go-junit-report v1.0.0 // indirect
        github.com/klauspost/compress v1.17.11 // indirect
        github.com/kr/pretty v0.3.1 // indirect
        github.com/mattn/go-colorable v0.1.14 // indirect
        github.com/mattn/go-isatty v0.0.20 // indirect
        github.com/mitchellh/hashstructure v1.1.0 // indirect
        github.com/morikuni/aec v1.0.0 // indirect
        github.com/opencontainers/go-digest v1.0.0 // indirect
        github.com/opencontainers/image-spec v1.1.0 // indirect
        github.com/pierrec/lz4/v4 v4.1.22 // indirect
        github.com/pmezard/go-difflib v1.0.0 // indirect
        github.com/prometheus/procfs v0.15.1 // indirect
        github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a  // indirect
        github.com/rogpeppe/go-internal v1.13.1 // indirect
        github.com/samuel/go-parser v0.0.0-20170131185712-99744db8e45c // indirect
        github.com/sirupsen/logrus v1.9.3 // indirect
        github.com/urso/go-bin v0.0.0-20180220135811-781c575c9f0e // indirect
        github.com/urso/magetools v0.0.0-20190919040553-290c89e0c230 // indirect
        go.uber.org/multierr v1.10.0 // indirect
        golang.org/x/crypto v0.33.0 // indirect
        golang.org/x/mod v0.18.0 // indirect
        golang.org/x/sync v0.11.0 // indirect
        golang.org/x/term v0.29.0 // indirect
        golang.org/x/text v0.22.0 // indirect
        golang.org/x/time v0.10.0 // indirect
        golang.org/x/tools v0.22.0 // indirect
        golang.org/x/tools/go/vcs v0.1.0-deprecated // indirect
        google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
        google.golang.org/grpc v1.69.4 // indirect
        google.golang.org/protobuf v1.36.3 // indirect
        gopkg.in/yaml.v2 v2.4.0 // indirect
        gopkg.in/yaml.v3 v3.0.1 // indirect
        gotest.tools v2.2.0+incompatible // indirect
        howett.net/plist v0.0.0-20181124034731-591f970eefbb // indirect
)

replace (
        github.com/dop251/goja => github.com/andrewkroh/goja v0.0.0-20190128172624-dd2ac4456e20
        github.com/dop251/goja_nodejs => github.com/dop251/goja_nodejs v0.0.0-20171011081505-adff31b136e6
        github.com/insomniacslk/dhcp => github.com/elastic/dhcp v0.0.0-20200227161230-57ec251c7eb3
        github.com/samuel/go-thrift => github.com/samuel/go-thrift v0.0.0-20140522043831-2187045faa54
        github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20210823122811-11c3ef800752
        github.com/rcrowley/go-metrics => github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a
)