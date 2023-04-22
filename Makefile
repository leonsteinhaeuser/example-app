IMG_PREFIX?=ghcr.io/leonsteinhaeuser
IMG_NUMBER_SERVICE=$(IMG_PREFIX)/number-service
IMG_VIEW_SERVICE=$(IMG_PREFIX)/view-service
VERSION?=latest

run:
	docker-compose up --build

docker-build-and-push: docker-build docker-push

docker-push:
	docker push $(IMG_NUMBER_SERVICE):${VERSION}
	docker push $(IMG_VIEW_SERVICE):${VERSION}

build-docker: build-number-service build-view-service

build-number-service:
	docker build -t $(IMG_NUMBER_SERVICE):${VERSION} -f number-service/Dockerfile .

build-view-service:
	docker build -t $(IMG_VIEW_SERVICE):${VERSION} -f view-service/Dockerfile .