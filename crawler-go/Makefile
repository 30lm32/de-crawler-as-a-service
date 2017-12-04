DEFAULT_CONTAINER_NAME := test-crawler
DEFAULT_IMAGE := test-crawler
DOCKER_FILE := Dockerfile

kill_container:
	docker kill ${DEFAULT_CONTAINER_NAME}

rm_container:
	docker rm ${DEFAULT_CONTAINER_NAME}

build_image:
	docker build --no-cache -t ${DEFAULT_IMAGE} -f ${DOCKER_FILE} .

run_container: rm_container
	docker run -it --name ${DEFAULT_CONTAINER_NAME} ${DEFAULT_IMAGE}