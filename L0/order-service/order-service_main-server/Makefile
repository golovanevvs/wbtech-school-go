MAKEFLAGS += --no-print-directory

kafka-up:
	@$(MAKE) -f Makefile.kafka up
kafka-down:
	@$(MAKE) -f Makefile.kafka down
kafka-clean:
	@$(MAKE) -f Makefile.kafka clean

redis-up:
	@$(MAKE) -f Makefile.redis up
redis-down:
	@$(MAKE) -f Makefile.redis down
redis-clean:
	@$(MAKE) -f Makefile.redis clean

postgres-up:
	@$(MAKE) -f Makefile.postgres up
postgres-down:
	@$(MAKE) -f Makefile.postgres down
postgres-clean:
	@$(MAKE) -f Makefile.postgres clean

migrate-up:
	@$(MAKE) -f Makefile.migrate up
migrate-down:
	@$(MAKE) -f Makefile.migrate down

main-server-producer-build-up:
	@$(MAKE) -f Makefile.main server-producer-build-up
main-server-build-up:
	@$(MAKE) -f Makefile.main server-build-up
main-server-up:
	@$(MAKE) -f Makefile.main server-up
main-producer-build-up:
	@$(MAKE) -f Makefile.main producer-build-up
main-producer-up:
	@$(MAKE) -f Makefile.main producer-up

dep-up: kafka-up redis-up postgres-up migrate-up

main-build-up: main-server-producer-build-up

main-up: main-server-up main-producer-up

all-up: dep-up main-build-up

all-down: kafka-down redis-down postgres-down

all-clean: kafka-clean redis-clean postgres-clean