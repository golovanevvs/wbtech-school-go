SRC_FILE := ../rpRedisSaveNotice.go

PACKAGE := rpRedisSaveNotice

OUT_DIR := .

OUT_FILE := $(OUT_DIR)/mock_rd_savenotice.go

INTERFACE := RedisClient

.PHONY: mock-rd-savenotice-gen mock-rd-savenotice-clean

mock-rd-savenotice-gen:
	@echo "Generating mock for interface $(INTERFACE)..."
	@mockgen -source=$(SRC_FILE) -destination=$(OUT_FILE) -package=$(PACKAGE) $(INTERFACE)
	@echo "Mock generated at $(OUT_FILE)"

mock-rd-savenotice-clean:
	@rm -f $(OUT_FILE)
	@echo "Mocks cleaned"