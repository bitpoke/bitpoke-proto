include ../build/makelib/common.mk

ifeq ($(CI),true)
PUBLISH_REPO := https://github.com/presslabs/dashboard-proto.git
else
PUBLISH_REPO := git@github.com:presslabs/dashboard-proto.git
endif

PUBLISH_BRANCH ?= $(BRANCH_NAME)

include ../build/makelib/git-publish.mk
