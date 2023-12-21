# Message streaming application

### Technologies

- Golang
- Redis streams
- Docker & Docker Compose

This application demonstrates a simple producer / consumer stream.

Run with docker.

```bash
docker compose up --build
```

By default this is running with one producer and one consumer.

To increase the number of consumers just change replicas number on consumer service at `docker-compose.yml` file.

```yml
deploy:
  replicas: 1
```

### Graph representation

```mermaid
flowchart LR;
    A[Producer]-->B[Redis];
    C[Consumer]-->B[Redis];
    D[Consumer]-->B[Redis];
    E[Consumer]-->B[Redis];
```

### Producing messages

You can publish messages by doing POST requests on `localhost:8000/produce`

Example:

```bash
curl -X POST localhost:8000/produce -d '{"message": "Hello from terminal"}'
```
