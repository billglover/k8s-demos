# Persistent Volume Claims under Node Failure

> "Show what happens to a `PersistentVolumeClaim` in the event of a node failure. The behaviour may differ based on the storage driver and underlying storage provider used."

## Scenario

A `Deployment` is created with a single instance of a `Pod` that mounts persistent storage via a `PersistentVolumeClaim`. The `Pod` contains an app that utilises the persistent storage that is external to the cluster. In this example the storage is used to persist the state of a counter.

The node underneath the running pod suffers an ungraceful shutdown.

## Expected Behaviour

1. The application is able to persist data to the underlying storage
2. If the `Pod` is terminated, it is rescheduled and storage is reattached.
3. If the `Node` is ungracefully terminated, there is an outage while the `Pod` is rescheduled
4. The `Pod` is rescheduled and storage is reattached.

## Demonstration

Note: This demonstration is currently set to run against any cluster with a `standard` `StorageClass` available (e.g. GKE). In your cluster this may be `default`. If this is the case you'll need to modify the `02_pvc.yaml` definition before running the demonstration.

1. Create the Kubernetes resources `kubectl apply -f ./k8s`.
2. Query the service endpoint to show the counter incrementing on each request.
3. Terminate the `Pod` and wait for a new `Pod` to be scheduled.
4. Query the service endpoint to show the counter incrementing from previous value.
5. Identify the `Node` where the `Pod` has been scheduled.
6. Use your IaaS control plane to terminate the `Node`.
7. Wait for a new `Pod` to be scheduled.
8. Query the service endpoint to show the counter incrementing from previous value.

## Gotchas

* Look out for the `Pod` failing to start after `Node` failure. It may be waiting on the `PersistentVolumeClaim` to be released.

## Building the Container Image

1. Build and tag the container image.

   ```
   export APP_TAG="latest"
   export APP_IMG="harbor.shared.12factor.xyz/bglover/counter"
   docker build -t "$APP_IMG:$APP_TAG" .
   ```

2. Push the image to an accessible repository.

   ```
   docker push "$APP_IMG:$APP_TAG"
   ```

3. Modify the `k8s/03_deployment.yaml` file if you want to use your image instead of the default.