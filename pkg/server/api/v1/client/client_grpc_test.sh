grpcurl -plaintext localhost:18099 describe 

grpcurl -plaintext localhost:18099 describe pkg.server.api.client.ClientService

grpcurl -plaintext localhost:18099 describe .pkg.server.api.client.AuthRequest_v1

grpcurl -plaintext -d \
    '{
        "clusterUuid": "3031f7e4fd88437493c52db5afe49acb",
        "assertion": "8577250ae3c742c1aaff3fddbb25d52e",
        "clientVersion": "v1",
        "clientLibVersion": "v1"
    }' \
    localhost:18099 pkg.server.api.client.ClientService.Auth_v1 

# eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbHVzdGVyLXV1aWQiOiIzMDMxZjdlNGZkODg0Mzc0OTNjNTJkYjVhZmU0OWFjYiIsImNsdXN0ZXItY2xpZW50LXRva2VuLXV1aWQiOiJmMjhlMmI0NTU2NTA0M2U1OTQ4OTY4ZWRkZjVmN2ZiYSIsInV1aWQiOiIzNDE0NzhmNjU3ZmU0N2U5ODZjZWNkOWRjMWNjMGI3ZCIsImlhdCI6MTY4NzMzODc5MCwiZXhwIjoxNzAzMTE2ODAwLCJjbGllbnQtdmVyc2lvbiI6InYxIiwiY2xpZW50X2xpYl92ZXJzaW9uIjoidjEifQ.pyMcuEyDOEWeVzI_NEeNthHaB8CI8qqXYXYcYoPAdFQ

grpcurl -plaintext localhost:18099 -rpc-header ""