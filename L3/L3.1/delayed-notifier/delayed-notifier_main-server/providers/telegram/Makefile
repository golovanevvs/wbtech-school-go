MAKEFLAGS += --no-print-directory

PROJECT_DIR := $(abspath ../..)

ENV_FILES := .env ../app/.env

$(foreach file,$(ENV_FILES), \
  $(if $(wildcard $(file)), $(eval include $(file))) \
)
export

BOT_TOKEN ?= $(TELEGRAM_TOKEN)
WEBHOOK_URL = https://incompletely-elemental-moonfish.cloudpub.ru/delayed-notifier_main-server/telegram/webhook
PORT = 80

.PHONY: install check delete test help

install:
	@echo "Setting up webhook..."
	@curl -s -X GET \
		"https://api.telegram.org/bot$(BOT_TOKEN)/setWebhook?url=$(WEBHOOK_URL)"
	@echo "Webhook installed successfully!"

check:
	@echo "Checking webhook info..."
	@curl -s -X GET \
		"https://api.telegram.org/bot$(BOT_TOKEN)/getWebhookInfo"

delete:
	@echo "Deleting webhook..."
	@curl -s -X GET \
		"https://api.telegram.org/bot$(BOT_TOKEN)/deleteWebhook"
	@echo "Webhook deleted successfully!"

test:
	@echo "Testing webhook with sample data..."
	@curl -s -X POST \
		-H "Content-Type: application/json" \
		-d '{
			"update_id": 10000,
			"message": {
				"message_id": 123,
				"from": {
					"id": 123456789,
					"first_name": "Test",
					"last_name": "User",
					"username": "testuser"
				},
				"chat": {
					"id": 123456789,
					"type": "private"
				},
				"date": 1500000000,
				"text": "/start"
			}
		}' \
		"$(WEBHOOK_URL)"

install-detailed:
	@echo "Setting up detailed webhook..."
	@curl -s -X POST \
		-H "Content-Type: application/json" \
		-d '{
			"url": "$(WEBHOOK_URL)",
			"max_connections": 40,
			"allowed_updates": ["message", "callback_query"]
		}' \
		"https://api.telegram.org/bot$(BOT_TOKEN)/setWebhook"
	@echo "Detailed webhook installed successfully!"

help:
	@echo "Telegram Bot Webhook Management"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  install          - Set up webhook with basic settings"
	@echo "  install-detailed - Set up webhook with detailed settings"
	@echo "  check            - Check webhook information"
	@echo "  delete           - Remove webhook"
	@echo "  test             - Test webhook with sample data"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Configuration:"
	@echo "  BOT_TOKEN: $(BOT_TOKEN)"
	@echo "  WEBHOOK_URL: $(WEBHOOK_URL)"
	@echo ""
	@echo "Example:"
	@echo "  $$ make BOT_TOKEN=123456:ABC-DEF install"
	@echo "  $$ make check"

# ========================
# TELEGRAM
# ========================

# TELEGRAM = providers/telegram

# .PHONY: telegram-install telegram-check telegram-delete telegram-test telegram-help

# telegram-install:
# 	$(MAKE_CMD) -C $(TELEGRAM) install
# telegram-check:
# 	$(MAKE_CMD) -C $(TELEGRAM) check
# telegram-delete:
# 	$(MAKE_CMD) -C $(TELEGRAM) delete
# telegram-test:
# 	$(MAKE_CMD) -C $(TELEGRAM) test
# telegram-install-detailed:
# 	$(MAKE_CMD) -C $(TELEGRAM) install-detailed
# telegram-help:
# 	$(MAKE_CMD) -C $(TELEGRAM) help