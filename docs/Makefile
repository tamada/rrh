PWD := ${CURDIR}
DOCKER_IMAGE_NAME := "wwwrrh"
CONTAINER_REMOVE_FLAG := "--rm"
BASE_URL := "https://tamada.github.io/rrh"
HUGO_THEME := "cayman-hugo-theme"
JOJOMI_VERSION := 0.83.1

build:
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME}_build -v "${PWD}":/src -v ${PWD}/public:/output -e HUGO_THEME=$(HUGO_THEME) -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

start:
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v "$(PWD)"/public:/output -e HUGO_THEME=$(HUGO_THEME) -e HUGO_WATCH="true" -e HUGO_BASEURL=${BASE_URL} jojomi/hugo:${JOJOMI_VERSION}

build_with_hugo:
	hugo
	rm public/favicon*.png public/favicon.ico public/apple-touch-icon.png

stop:
	docker stop wwwrrh
