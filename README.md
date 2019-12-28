# OpenFaaS NATS Connector Example/Tester

This repo provides a very simple way to verify that an installation of  the [`nats-connector`](https://github.com/openfaas-incubator/nats-connector) is working as expected.


## Running Locally

1. Deploy a test function an annotate it with the nats topic you will publish messages to:
   ```sh
   faas-cli deploy --name pycho --image theaxer/pycho:latest --fprocess='python index.py'  --annotation topic="faas-req"
   ```
2. Install the `nats-connector`
   ```sh
   kubectl apply -f https://raw.githubusercontent.com/LucasRoesler/nats-connector-example/master/yaml/connector-dep.yaml
   ```
3. For simplicity, forward the nats port to your localhost
   ```sh
   kubectl port-forward -n openfaas svc/nats 4222 &
   ```
4. Run the tester
   ```sh
   go run main.go
   ```

## Running in cluster
1. Deploy a test function an annotate it with the nats topic you will publish messages to:
   ```sh
   faas-cli deploy --name pycho --image theaxer/pycho:latest --fprocess='python index.py'  --annotation topic="faas-req"
   ```
2. Install the `nats-connector`
   ```sh
   kubectl apply -f https://raw.githubusercontent.com/LucasRoesler/nats-connector-example/master/yaml/connector-dep.yaml
   ```
3. Run the tester as a Job
   ```sh
   kubectl apply -f https://raw.githubusercontent.com/LucasRoesler/nats-connector-example/master/yaml/tester-job.yaml
   ```
