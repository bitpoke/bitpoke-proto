include ../build/makelib/common.mk

ifeq ($(CI),true)
PUBLISH_REPO := https://github.com/bitpoke/bitpoke-proto.git
else
PUBLISH_REPO := git@github.com:bitpoke/bitpoke-proto.git
endif

PUBLISH_BRANCH ?= $(BRANCH_NAME)

include ../build/makelib/git-publish.mk
