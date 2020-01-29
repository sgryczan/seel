# Seel
Sealed Secrets Converter

### Purpose: Converts normal K8S secrets into SealedSecrets for a given environment

### Usage:

Example:
```
kubectl -n dev create secret generic basic-auth \
--from-literal=user=admin \
--from-literal=password=admin \
--dry-run \
-o json | curl -X POST -d @- http://seel.sre-dev.solidfire.net/convert
```

Response:
```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: myAwesomeSecret
  namespace: destinationNamespace
spec:
  encryptedData:
    admin: AgAfYT9jLehpMf2CI3Z9ikqQQv0eaWwQoHP1Lq4NclodmFHIl8z0gUDp+aKFLEHtwbiQs8m1R3EPgyr/w6rdEXrPcc7Je6UNd2IUSfRToOO6bG0+oakDhlxl+BM61x+NCJY6/cbC/nhQqv/P8YEDXL+HSMpKfshMgiCRd84BbSHiY2w+h9F+4JyVZH7NhK22MQRPPbe4H9o1r5me5wq3Dqo92IzpwjgAJs/LKOZIFnQJTXtw4Hn45WGf+8N2HJjEDsZrHamnm2PFdxZUVXLNEmYLarNmaYhbJa3bsk5QZ8g9t+JyM1FcfKOF6D0s/t9UawX+k+o/C0oTrdklAa+GsPhVPllZAWEus+654kOJcMWMAxNifVRqV9eukIqtxkZAUdX5YEc2b2fK2xCjpJYvO9EYExheL5EHmZ4EbHUak0BKRTcac6zmFRWB36uRSBP/SUNleueP9hv4YrwxXfv64fAWfSLCniHb0CtxfRuWntg+XMQ75aFdCj7EE3SbbLJtpeLOVAWhHoMZL7NU/CIch0v1KT6vOcWrlZzymz+/rHPJu2+TPf/iGTY7gmj09eV/oFds9snHVx8EdQbG4lH3oBWri+d36ArPurl2KXb3927SGNcSMsvBDHv0ybcsJPtjyE5nuZe/8YXxQNgqBaUs+H1UBRZnqtJqO2dxLrBfxUGzkPsPli3bx1lD7okGFA8iGFg6RJDWvA==
    password: AgCGjrksF/M4s6HZPH+CZA74V9+s26/YE2ZHG5FgKYgo9UZ/RUGkxJBDtCJsKgG7rbBHCJ6g+Z3i8Wnfd32t8/M3LT7p+KbG9MiW9v1Ui02EZFrkOtGBZt3+cGYB0+74JvXCMvkKnE7LgtCa3BaguueQcdNHNtwtKtq8T/xBx44hNoM2ARTR//EXrmwNua5kriGvXjgfaOD+cXJad/SOkwrQRK/ghQ5SsAHFgjxtVVe6uFKRsK54/q9MLmcKlNr8DNkluVqvPUeXH7fmVXcDSytKByF975GYfBJF0KmuRHM+opTMrRaYRrZ7BhjKi/mhB0UYR3xzbi7LI+XBiu3h0LZxw2NZhRnDNRh/RmyVyTLZJcB+alWfBpVFWndjia88INyskFZB7jw3vcg52vyLE8iNkK53kBBQcdQR9OKhPgLtvT9OY8ASPCzI+mjQ6yDbNUGguY7xISoOUlyedUjA1MTcoM9CbfH3MfFPIxxlQwntu9fQpcWJ0lDMOtPgcLAsM7Yzq5rr7nEzMa1vdI8d9Hug+NVrPY+wmMO0ss3Ga0APSFh1sFKmzWHa2YGYWAD6TOrFBD5EVXIFd+Zbg2HKSMWE0DLmjQVOdNO9prSvRm2W+qpci4GfELUN+PA1zP6jAHjw4UWY+MIRp8unPSqU4iDoHqPaysYwcjwJ5LlCZ8ifXXlBuXQVgKOLw/VnDDyDsmh3BFPgc584qQ==
  template:
    metadata:
      creationTimestamp: null
      name: myAwesomeSecret
      namespace: destinationNamespace
status: {}
```

This file can then be added to the proper repository to be synced to the cluster by Flux.


## Deployment to K8S
The deployment manifest can be found in the `k8s` folder. 

It should be deployed into the same namespace as the sealed secrets controller(usually `adm`)