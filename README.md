# EVO-test-task

## Build and start application using docker-compose

### Build containers and start application

```
docker-compose up
```

On first run application will wait for PostgreSQL database initialization.

### Stop and remove containers

```
docker-compose down
```

<br />

## Build and start application using Makefile

### Pull postgres:latest Docker image

```
make image
```

### Make and run postgre container

```
make postgres
```

### Compile and run application

```
make server
```

<br />

## API documentation

[http://localhost:8080/swagger/](http://localhost:8080/swagger/)
