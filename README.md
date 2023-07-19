# claiflow

## What is claiflow

<!-- ***claiflow***는 다중 Kubernetes 클러스터에서 작업을 관리하기 위한 오픈소스 프로젝트입니다. -->
***claiflow*** is an open-source project for managing tasks across multiple Kubernetes clusters.
<!-- ***claiflow***는 Kubernetes 및 Kubernetes에서 동작중인 서비스에 표준화된 템플릿을 제공하여, 작업을 실행 하게 하여,
예측 가능한 운영과 사용자-에러의 가능성을 방지할 수 있습니다.  -->
***claiflow*** provides standardized templates for Kubernetes and services running on Kubernetes to execute tasks. This helps prevent unpredictable operations and user errors.
<!-- Kubernetes를 자동화 하기 위해서, REST API를 이용하여 ***claiflow*** 리소스를 조작하세요. -->
Manipulate ***claiflow*** resources using the REST API for Kubernetes automation.

## How to install claiflow?

### Server

1. Requirements
    <!-- 10.6 이상의 MariaDB가 필요합니다. -->
    You need MariaDB version 10.6 or above.

1. Installation
    <!-- 자세한 정보는 deploy 패키지의 README.md 파일을 확인 하세요. -->
    Please refer to the [README.md](deploy/README.md) file in the deploy package for detailed information.

### Client

1. How to make own claiflow-client?

    <!-- 사용자가 윈하는 구성의 ***claiflow-client***를 빌드 하세요. -->
    Build the claiflow-client with the configuration that you prefer.
    <!-- 미리 구성된 템플릿-명령 핸들러를 등록하여, 손쉽게 자신만의 ***claiflow-client***을 빌드할 수 있습니다. -->
    You can easily build your own ***claiflow-client*** by registering predefined *template-command* handlers.

    <!-- Examples를 통해 ***claiflow-client***의 작성법과 실행법을 확인하세요. -->
    Check the [examples](examples) to see how to write and run claiflow-client.

## How to use?

1. Kubernetes port-forward.
    <!-- curl command for create cluster resource -->
    ```sh
    kubectl -n claiflow port-forward service/claiflow-service 8099:8099
    ```

1. Declare *cluster* resource.
    <!-- curl command for create cluster resource -->
    ```sh
    curl --request POST \
        --url http://localhost:8099/api/v1/cluster \
        --header 'content-type: application/json' \
        --data '{"uuid": "00000000000000000000000000000000","name": "test cluster"}'
    ```

    output:
    > {"uuid":"00000000000000000000000000000000","name":"test cluster","created":"2023-07-06T01:23:54.060601557Z","updated":"2023-07-06T01:23:54.060601557Z"}

    - cluster's *uuid* is optional.

        If it is marked as optional, it means that if no information is available, it will be generated as a random UUID.

1. Declare *cluster-client-token* resource.
    <!-- curl command for create cluster-client-token resource -->
    ```sh
    curl --request POST \
        --url http://localhost:8099/api/v1/cluster_token \
        --header 'content-type: application/json' \
        --data '{"cluster_uuid": "00000000000000000000000000000000","uuid": "00000000000000000000000000000001","name": "test cluster client token","summary": "test cluster client token summary","token": "CLIENT_AUTHENTICATION_TOKEN"}'
    ```

    output:
    > {"uuid":"00000000000000000000000000000001","name":"test cluster client token","summary":"test cluster client token summary","cluster_uuid":"00000000000000000000000000000000","token":"CLIENT_AUTHENTICATION_TOKEN","issued_at_time":"2023-07-06T01:29:23.53447061Z","expiration_time":"2024-07-05T00:00:00Z","created":"2023-07-06T01:29:23.53447061Z","updated":"2023-07-06T01:29:23.53447061Z"}

    - cluster's *uuid*

        The cluster's *uuid* refers to the value generated and outputted when creating the cluster resource.

    - cluster-client-token's *uuid* is optional.

    - cluster-client-token's *token* is optional.

        The *token* of the cluster-client-token is the *assertion* information used by the ***claiflow-client*** for authentication.

1. Deploy the ***claiflow-client*** to Kubernetes cluster.

    <!-- server url, cluster uuid, cluster token 값을 확인하세요. -->
    1. Check the server url, cluster uuid, and cluster token values.
        ```yaml
        server grpc url: localhost:18099
        cluster uuid:    00000000000000000000000000000000
        cluster token:   CLIENT_AUTHENTICATION_TOKEN
        ```
    <!-- 위에서 확인한 값들을 ***claiflow-client*** 실행 시 적용해 주세요. 자세한 사항은 [examples/README.md](examples/README.md) 파일을 통해 확인하세요. -->
    2. Apply the values you found above when running claiflow-client. For more information, see the [examples/README.md](examples/README.md).
        ```sh
        # move examples directory
        cd examples
        # build client binary
        make go-build example=multistep_features
        # build docker image
        make docker-build example=multistep_features image=SOME_REPOSITORY_URL
        # push docker image
        make docker-push example=multistep_features image=SOME_REPOSITORY_URL
        # change your configuration in configmap: server url, cluster uuid, cluster token
        vi environment.yaml
        # create namesapce, configmap
        kubectl apply -f environment.yaml 
        # create serviceaccount, clusterrole, clusterrolebinding
        kubectl apply -f service_account.yaml
        # change image_name to your docker image
        vi application.yaml 
        # create deployment
        kubectl apply -f application.yaml
        ```

1. Deploy the pre-defined *templates* used by ***claiflow-client***
    <!-- sql command for create pre-defined templates -->
    ```sql
    REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_simple', 'example_simple', '', '[{"$id":"step1","$command":"helloworld","inputs":"$inputs"}]', '{}', 'userdefined', NOW());
    REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_iter', 'example_iter', '', '[{"$id":"step1","$range":"$inputs.x_list","$steps":[{"$id":"step2","$command":"math_pow","inputs":{"x":"$step1.val","y":2}}]}]', '{}', 'userdefined', NOW());
    REPLACE INTO `template` (`uuid`, `name`, `summary`, `flow`, `inputs`, `origin`, `created`) VALUES ('example_pass_val', 'example_pass_val', '', '[{"$id":"step1","$command":"swap_command","inputs":{"param1":"$inputs.input1","param2":"$inputs.input2"}},{"$id":"step2","$command":"swap_command","inputs":{"param1":"$step1.outputs.value1","param2":"$step1.outputs.value2"}}]', '{}', 'userdefined', NOW());
    ```

1. Create the *tasks* using the *template* you want to execute.
    <!-- curl command for create service resource -->
    ```sh
    curl --request POST \
        --url http://localhost:8099/api/v1/service \
        --header 'Content-Type: application/json' \
        --data '{"cluster_uuids":["00000000000000000000000000000000"],"uuid":"c8590e761af64c6891a47b6570c0f93e","name":"example_simple","template_uuid":"example_simple","inputs":{"name":"world"}}'
    ```

    output:
    > [{"cluster_uuid":"00000000000000000000000000000000","uuid":"c8590e761af64c6891a47b6570c0f93e","name":"example_simple","summary":"","template_uuid":"example_simple","flow":"[{\"$id\":\"step1\",\"$command\":\"helloworld\",\"inputs\":\"$inputs\"}]","inputs":{"name":"world"},"step_max":1,"priority":0,"created":"2023-07-06T02:08:17.580567016Z","statuses":[{"step_seq":0,"status":0,"created":"2023-07-06T02:08:17.580567016Z"}]}]

    - task's *uuid* is optional.

    - cluster's *uuid* is string or string array.

        ```json
        { "cluster_uuid": "00000000000000000000000000000000"}
        { "cluster_uuid": ["00000000000000000000000000000000","00000000000000000000000000000001"]}
        ```

1. Check the *task* result.
    <!-- curl command for get task resource -->
    ```sh
    curl --request GET \
        --url http://localhost:8099/api/v1/cluster/00000000000000000000000000000000/service/c8590e761af64c6891a47b6570c0f93e \
        --header   "accept: application/json"
    ```

    output:
    > {"cluster_uuid":"00000000000000000000000000000000","uuid":"c8590e761af64c6891a47b6570c0f93e","name":"example_simple","summary":"","template_uuid":"example_simple","flow":"[{\"$id\":\"step1\",\"$command\":\"helloworld\",\"inputs\":\"$inputs\"}]","inputs":{"name":"world"},"step_max":1,"priority":0,"created":"2023-07-06T02:08:17.580567Z","statuses":[{"step_seq":0,"status":0,"created":"2023-07-06T02:08:17.580567Z"}]}

## Etc

1. More information about the APIs can be found on the Swagger page.

    [http://localhost:8099/swagger/index.html](http://localhost:8099/swagger/index.html)

1. You can access the metric information for the server in Prometheus metric format.

    [http://localhost:8099/metrics](http://localhost:8099/metrics)

1. Kubernetes version compatibility

    <!-- claiflow-client는 storage.k8s.io/v1 CSIStorageCapacity 리소스 정보를 사용하기 때문에 Kubernetes 버전 v1.24.0 이상이 필요합니다. -->
    Kubernetes version v1.24.0 or later is required because ***claiflow-client*** uses storage.k8s.io/v1 CSIStorageCapacity resource information.

    <!-- claiflow 클라이언트는 [kubernetes](https://github.com/kubernetes/client-go), [helm](https://github.com/helm/helm), [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) 라이브러리를 사용하여 Kubernetes API와 통신합니다. -->
    ***claiflow-client*** uses the [kubernetes/client-go](https://github.com/kubernetes/client-go), [helm](https://github.com/helm/helm), and [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) libraries to communicate with the Kubernetes API.

    <!-- 각 라이브러리들의 버전 호환성은 아래에서 찾을 수 있습니다. -->
    Version compatibility for each library can be found below.
    * [kubernetes/client-go](https://github.com/kubernetes/client-go#compatibility-matrix)
    * [helm](https://helm.sh/docs/topics/version_skew/)
    * [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/compatibility.md)

    <!-- 현재 클라이언트는 다음 버전들을 사용하고 있습니다. -->
    Currently, ***claiflow-client*** are using the following versions
    * kubernetes/client-go: v0.26.0
    * helm: v3.11.1
    * prometheus-operator: v0.55.0
