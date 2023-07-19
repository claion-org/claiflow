# How to deploy the claiflow server

## Modify Kubernetes manifests

1. Modify deployment.yaml

    Update it with the URL of the Docker image that has been deployed.

    ```sh
    DOCKER_IMAGE_URL="{REGISTRY_URL}/{NAMESPACE}/{IMAGE}" yq -i '.spec.template.spec.containers[0].image = strenv(DOCKER_IMAGE_URL)' deployment.yaml
    ```

## Deploy to kubernetes

1. Deploy to Kubernetes cluster

    ```sh
    make all
    ```

2. Undeploy from Kubernetes cluster

    ```sh
    make clean
    ```

## Use the existing database

### Database configuration

1. create database

    ```sql
    create database `DATABASENAME`;
    ```

1. create user

    ```sql
    CREATE USER 'USERNAME'@'%' IDENTIFIED BY 'PASSWORD';
    ```

1. GRANT PRIVILEGES

    ```sql
    GRANT ALL PRIVILEGES ON *.* TO 'USERNAME'@'%' WITH GRANT OPTION;
    FLUSH PRIVILEGES;
    ```

## Default config values

1. config/application.yaml

    | path | desc | default |
    | --- | --- | --- |
    | .http.port | - | 8099 |
    | .http.urlPrefix | - |  |
    | .http.tls.enable | - | false |
    | .http.tls.certFile | - | 'server.crt' |
    | .http.tls.keyFile | - | 'server.key' |
    | .http.cors.allowOrigins | - | '' |
    | .http.cors.allowMethods | - | '' |
    | .grpc.port | - | 18099 |
    | .grpc.tls.enable | - | false |
    | .grpc.tls.certFile | - | 'server.crt' |
    | .grpc.tls.keyFile | - | 'server.key' |
    | .grpc.maxRecvMsgSize | - | 1073741824 |
    | .migrate.source | - | 'migrations/mysql' |
    | .logger.verbose | - | 0 |
    | .logger.disableCaller | - | false |
    | .logger.disableStacktrace | - | false |

1. config/enigma.yaml

    | path | desc | default |
    | --- | --- | --- |
    | .enigma.blockMethod | - | 'none' |
    | .enigma.blockSize | - | 0 |
    | .enigma.blockKey | - | '' |
    | .enigma.cipherMode | - | 'none' |
    | .enigma.cipherSalt | - | null |
    | .enigma.padding | - | 'none' |
    | .enigma.strconv | - | 'plain' |

1. config/database.yaml

    | path | desc | default |
    | --- | --- | --- |
    | .database.type | - | 'mysql' |
    | .database.protocol | - | 'tcp' |
    | .database.host | - | 'localhost' |
    | .database.port | - | 3306 |
    | .database.dbname | - | 'flow' |
    | .database.username | - | '' |
    | .database.password | - | '' |
    | .database.maxOpenConns | - | 15 |
    | .database.maxIdleConns | - | 5 |
    | .database.connMaxLifetime | - | 1 |
    | .database.password | - | '' |