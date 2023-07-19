package api

//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/client
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/cluster
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/cluster_client_session
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/cluster_client_token
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/global_variables
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/service
//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag@v1.8.7 fmt --dir v1/webhook
