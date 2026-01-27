Docker Commands 

# Build and start everything
docker-compose up --build -d

# View build process
docker-compose logs api

# Inspect the built image
docker images | grep go-domain-api

# See the running process inside container
docker exec go-domain-api ps aux