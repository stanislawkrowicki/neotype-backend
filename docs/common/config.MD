## Config directory is used for almost all configuration. 

#### Hierarchy:
    - service_name
        - config_local.yaml
        - config_production.yaml

#### Service decides which config file to use by presence (or absence) of environment variable named _**envKey**_.
When **envKey** equals _**local**_ or is absent then the local config will be used.

When envKey is equal to _**production**_ the production config will be used.

In any other case, the service will panic at startup.

## Use of config by other services
Every service is expected to have at least two variables: _**base**_ and _**port**_.
A special function is used to return the baseURL of service. If _**base**_ is set to _http://localhost_, 
then port is added automatically to the end of baseURL.

