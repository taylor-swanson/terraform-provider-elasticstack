package fleetapi

//go:generate go run generate.go -v v8.9.0 -o fleet-filtered.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13 -package=fleetapi -generate=types -o ./fleetapi_gen.go fleet-filtered.json
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13 -package=fleetapi -generate=client -o ./client_gen.go fleet-filtered.json
