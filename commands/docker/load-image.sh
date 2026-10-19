docker build -t test-api:latest ./test-api

kind load docker-image test-api:latest --name cka-lab

