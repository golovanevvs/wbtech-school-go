MAKEFLAGS += --no-print-directory

ENV_FILES := .env ../app/.env ../deploy/.env

$(foreach file,$(ENV_FILES), \
  $(if $(wildcard $(file)), $(eval include $(file))) \
)
export

BOT_TOKEN ?= $(TELEGRAM_TOKEN)
WEBHOOK_URL = $(APP_TRANSPORT_HTTP_PUBLIC_HOST)/telegram/webhook
DEPLOY_WEBHOOK_URL = https://$(PUBLIC_HOST_DOCKER)/$(APP_CONTAINER)/telegram/webhook

.PHONY: install-deploy install check delete test

install-deploy:
	@echo "Setting up Telegram webhook for deploy..."
	@echo "BOT_TOKEN: $(BOT_TOKEN)"
	@echo "WEBHOOK_URL: $(DEPLOY_WEBHOOK_URL)"
	@curl -s -X POST \
		-H "Content-Type: application/json" \
		-d "{\"url\": \"$(DEPLOY_WEBHOOK_URL)\", \"max_connections\": 40, \"allowed_updates\": [\"message\", \"callback_query\"]}" \
		"https://api.telegram.org/bot$(BOT_TOKEN)/setWebhook"
	@echo "Telegram webhook for deploy installed successfully!"

install:
	@echo "Setting up Telegram webhook..."
	@echo BOT_TOKEN: $(BOT_TOKEN)
	@echo WEBHOOK_URL: $(WEBHOOK_URL)
	@curl -s -X POST \
		-H "Content-Type: application/json" \
		-d '{
			"url": "$(WEBHOOK_URL)",
			"max_connections": 40,
			"allowed_updates": ["message", "callback_query"]
		}' \
		"https://api.telegram.org/bot$(BOT_TOKEN)/setWebhook"
	@echo "Telegram webhook installed successfully!"

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
		-d "{\"update_id\": 10000, \"message\": {\"message_id\": 123, \"from\": {\"id\": 123456789, \"first_name\": \"Test\", \"last_name\": \"User\", \"username\": \"testuser\"}, \"chat\": {\"id\": 123456789, \"type\": \"private\"}, \"date\": 1500000000, \"text\": \"/start\"}}" \
		"$(WEBHOOK_URL)"

# ========================
# TELEGRAM
# ========================

# TELEGRAM = providers/telegram

# .PHONY: telegram-install-deploy telegram-install telegram-check telegram-delete telegram-test

# telegram-install-deploy:
# 	$(MAKE_CMD) -C $(TELEGRAM) install-deploy
# telegram-install:
# 	$(MAKE_CMD) -C $(TELEGRAM) install
# telegram-check:
# 	$(MAKE_CMD) -C $(TELEGRAM) check
# telegram-delete:
# 	$(MAKE_CMD) -C $(TELEGRAM) delete
# telegram-test:
# 	$(MAKE_CMD) -C $(TELEGRAM) test