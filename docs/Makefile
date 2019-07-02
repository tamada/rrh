PWD := ${CURDIR}
DOCKER_IMAGE_NAME := "wwwrrh"
CONTAINER_REMOVE_FLAG := "--rm"

start:
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v "$(PWD)"/public:/output -e HUGO_THEME=cayman-hugo-theme -e HUGO_WATCH="true" jojomi/hugo:latest

stop:
	docker stop wwwrrh

build:
	echo $(PWD)
	docker run ${CONTAINER_REMOVE_FLAG} -d --name ${DOCKER_IMAGE_NAME} -p 1313:1313 -v "${PWD}":/src -v ${PWD}/public:/output -e HUGO_THEME=cayman-hugo-theme jojomi/hugo:latest
