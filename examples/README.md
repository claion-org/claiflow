## Example directory

Each example demonstrates a few features of the client along with a test.

- [**Basic Commands**](./basic_commands): Simple examples of usage for Commands and Templates provided by the client by default.

- [**Multistep Features**](./multistep_features): Simple examples of registering custom commands and creating and using templates with them.

## How to run
1. Local
    
    Move each directory

2. Docker 
    ```bash
    make docker-build example=basic_commands
    docker run --rm -it \
        --env CLAIFLOW_SERVER_URL=localhost:18099 \
        --env CLAIFLOW_CLUSTER_UUID=00001 \
        --env CLAIFLOW_TOKEN=user-token \
        claiflow-client:basic_commands
    ```

3. Kubernetes
    
    docker build / push
    ```bash
    make docker-build example=basic_commands
    make docker-push example=basic_commands
    ```
    deploy in kubernetes
    ```bash
    kubectl apply -f environment.yaml # change your configuration: server url, cluster uuid, cluster token
    kubectl apply -f service_account.yaml
    kubectl apply -f application.yaml # change image_name to your docker image
    ```