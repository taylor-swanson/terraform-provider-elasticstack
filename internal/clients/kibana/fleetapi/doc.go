package fleetapi

//go:generate go run generate.go -o fleet-filtered.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -package=fleetapi -generate=types -o ./fleetapi.gen.go fleet-filtered.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -package=fleetapi -generate=client -o ./client.gen.go fleet-filtered.json
