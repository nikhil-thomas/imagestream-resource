# imagestream-resource
Openshift Imagestream resource for Tekton Pipelines

## Develppment

### As local process

1. Set OPENSHIFT_IMAGE_REGISTRY environment variable
    ```bash
    export OPENSHIFT_IMAGE_REGISTRY=$(oc get route default-route -n openshift-image-registry -o jsonpath='{.spec.host}')
    ```

2. Run ImagestreResource Controller
    ```bash
    make local-dev
    ```

### With ko

1. Set OPENSHIFT_IMAGE_REGISTRY environment variable in controller deployment manifest
    ```yaml
              env:
                - name: OPENSHIFT_IMAGE_REGISTRY
                  value: "<openshift-image-registry-url>"
               ...
    ```
2. Run ImagestreResource Controllermake local-dev
    ```bash
    make ko-app
    ```
