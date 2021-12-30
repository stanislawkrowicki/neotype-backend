# To run services on your machine and Docker images, you will need to create an .env file in *docker* directory.
### The file should look like that. You can change values in square brackets to whatever you wish.
### Values that are not in brackets can be changed only if you change them in config directory as well.

```
MYSQL_DATABASE={{database}}
MYSQL_USER={{username}}
MYSQL_PASSWORD={{password}}
MYSQL_ROOT_PASSWORD={{root_password}}
MYSQL_PORT=3306

RABBITMQ_USER={{username}}
RABBITMQ_PASSWORD={{password}}
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT={{15672}}

REDIS_PORT=6379

JWT_SECRET_KEY={{some_secret_key}}
```