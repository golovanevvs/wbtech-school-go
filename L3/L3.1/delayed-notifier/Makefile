# ========================
# SETUP
# ========================

ifeq ($(OS),Windows_NT)
	MAKE_CMD = mingw32-make
else
	MAKE_CMD = make
endif

# ========================
# ALL
# ========================

.PHONY: up up-no-cache down

up: main-server-up web-client-up

up-no-cache: main-server-up-no-cache web-client-up-no-cache

down: main-server-down web-client-down

# ========================
# MAIN-SERVER
# ========================

MAIN-SERVER = delayed-notifier_main-server

.PHONY: main-server-up main-server-up-no-cache main-server-down

main-server-up:
	$(MAKE_CMD) -C $(MAIN-SERVER) all-up
main-server-up-no-cache:
	$(MAKE_CMD) -C $(MAIN-SERVER) all-up-no-cache
main-server-down:
	$(MAKE_CMD) -C $(MAIN-SERVER) all-down

# ========================
# WEB-CLIENT
# ========================

WEB-CLIENT = delayed-notifier_web-client

.PHONY: web-client-up web-client-up-no-cache web-client-down

web-client-up:
	$(MAKE_CMD) -C $(WEB-CLIENT) up
web-client-up-no-cache:
	$(MAKE_CMD) -C $(WEB-CLIENT) up-no-cache
web-client-down:
	$(MAKE_CMD) -C $(WEB-CLIENT) down
