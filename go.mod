module diskcheck

go 1.19

require (
	github.com/diskfs/go-diskfs v1.3.0
	github.com/openebs/node-disk-manager v1.9.0
)

require (
	github.com/google/uuid v1.1.2 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/xattr v0.4.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	gopkg.in/djherbis/times.v1 v1.2.0 // indirect

)

// replace github.com/openebs/node-disk-manager v1.9.0 => ./third_party/github.com/openebs/node-disk-manager
