# Basic Commands

This example demonstrates how to execute using the commands the client has by default.

## Commands
- kubernetes
- prometheus
- helm

## How to run
1. Registered commands must exist on the server as a Template in order to be used.

    Insert [default templates](../../pkg/client/default_templates.sql) into your server database.

2. Run client.
    ```bash
    go run main.go \
        --server=localhost:18099 \
        --clusteruuid=<CLUSTER_UUID> \
        --token=<CLUSTER_TOKEN>
    ```

3. Request the server to create a service using the default templates.
    * Kubernetes list pods
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"k8s pods list","template_uuid":"00000000000000000000000000000002","inputs":{}}'
        ```
        
    * Prometheus Instant queries
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"p8s instant queries","template_uuid":"10000000000000000000000000000001","inputs":{"url":"http://kps-kube-prometheus-stack-prometheus.monitor.svc.cluster.local:9090","query":"{job=\"apiserver\"}"}}'
        ```
    
    * Helm install tomcat chart from bitnami repo
        ```bash
        curl -X POST http://localhost:8099/api/v1/service \
            -H 'Content-Type: application/json' \
            -d '{"cluster_uuid":"<CLUSTER_UUID>","name":"helm install tomcat","template_uuid":"20000000000000000000000000000003","inputs":{"name":"my-tomcat","chart_name":"tomcat","repo_url":"https://charts.bitnami.com/bitnami","namespace":"tomcat-ns"}}'
        ```