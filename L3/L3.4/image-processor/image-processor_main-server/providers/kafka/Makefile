MAKEFLAGS += --no-print-directory

ENV_FILE ?= .env
ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

.PHONY: create-topic

VORTEX = ../../../../../../../vortex

create-topic:
	$(MAKE_CMD) -C $(VORTEX) kafka-create-topic TOPIC=image-processor

# ========================
# KAFKA
# ========================

# KAFKA = providers/kafka

# .PHONY: kafka-create-topic

# kafka-create-topic:
# 	$(MAKE_CMD) -C $(KAFKA) create-topic