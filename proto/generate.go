package geometryproto

//go:generate protoc --go_out=. --go_opt=paths=import --go_opt=module=github.com/hansmi/dossier/proto geometry.proto report.proto sketch.proto
