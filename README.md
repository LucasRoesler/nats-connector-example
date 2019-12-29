# OpenFaaS NATS Connector Example/Tester

This repo provides a very simple way to verify that an installation of  the [`nats-connector`](https://github.com/openfaas-incubator/nats-connector) is working as expected.

To verify the behavior, it will deploy two functions

1. `connector-test` is a functiont that sends a pre-configured message to a NATS subject `faas-req` and then listens for an echo response on another subject `faas-resp`. It then tests the response matches the request.
2. `republish` accepts an arbitrary payload and then publishes it on the `faas-resp` subject. This is used as the target of the `nats-connector`

The test message flow looks like this

```
connector-test --> nats --> nats-connector --> republish --> nats --> connector-test
```

## Running

1. Deploy the republish function:
   ```sh
   faas-cli deploy
   ```
2. Install the `nats-connector`
   ```sh
   kubectl apply -f https://raw.githubusercontent.com/LucasRoesler/nats-connector-example/master/yaml/connector-dep.yaml
   ```
3. Invoke the test
   ```sh
   faas-cli invoke connector-test <<< "test message"
   ```
