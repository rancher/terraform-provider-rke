#!/usr/bin/env bash

set -x

TESTACC_DOCKER_LIST=${TESTACC_DOCKER_LIST:-"rke-dind-tf-testacc1 rke-dind-tf-testacc2"}

DOCKER_BIN=${DOCKER_BIN:-$(which docker)}


if [ "${TESTACC_DOCKER_LIST}" != "" ]; then
	echo Cleaning up testacc docker list ${TESTACC_DOCKER_LIST}
	for DOCKER_TEST in ${TESTACC_DOCKER_LIST}
	do
		DOCKER_ID=$(${DOCKER_BIN} ps -q -f name=${DOCKER_TEST})
		if [ "${DOCKER_ID}" != "" ]; then
			${DOCKER_BIN} rm -fv ${DOCKER_TEST}
		fi
	done
fi

