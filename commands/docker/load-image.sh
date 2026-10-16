docker build -t go-api:latest ./test-api

kind load docker-image go-api:latest --name cka-lab

