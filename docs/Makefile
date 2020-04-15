PWD := ${CURDIR}
DOCKER_IMAGE_NAME := "wwwrrh"
CONTAINER_REMOVE_FLAG := "--rm"
BASE_URL := "https://tamada.github.io/rrh"
JOJOMI_VERSION := 0.69.0

build:
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v ${PWD}/public:/output -e HUGO_THEME=cayman-hugo-theme -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

start:
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v "$(PWD)"/public:/output -e HUGO_THEME=cayman-hugo-theme -e HUGO_WATCH="true" -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

stop:
	docker stop wwwrrh

