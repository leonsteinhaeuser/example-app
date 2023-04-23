IMG_PREFIX?=ghcr.io/leonsteinhaeuser
IMG_NUMBER_SERVICE=$(IMG_PREFIX)/number-service
IMG_VIEW_SERVICE=$(IMG_PREFIX)/view-service
IMG_ARTICLE_SERVICE=$(IMG_PREFIX)/article-service
VERSION?=latest

run:
	docker-compose up --build

build-docker-and-push: build-docker push-docker

push-docker:
	docker push $(IMG_NUMBER_SERVICE):${VERSION}
	docker push $(IMG_VIEW_SERVICE):${VERSION}
	docker push $(IMG_ARTICLE_SERVICE):${VERSION}

build-docker: build-number-service build-view-service build-article-service

build-number-service:
	docker build -t $(IMG_NUMBER_SERVICE):${VERSION} -f number-service/Dockerfile .

build-view-service:
	docker build -t $(IMG_VIEW_SERVICE):${VERSION} -f view-service/Dockerfile .

build-article-service:
	docker build -t $(IMG_ARTICLE_SERVICE):${VERSION} -f view-service/Dockerfile .