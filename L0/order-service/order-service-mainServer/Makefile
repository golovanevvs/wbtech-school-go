kafka-up:
	$(MAKE) -f Makefile.kafka up

kafka-clean:
	$(MAKE) -f Makefile.kafka clean

redis-up:
	$(MAKE) -f Makefile.redis up

redis-clean:
	$(MAKE) -f Makefile.redis clean

postgres-up:
	$(MAKE) -f Makefile.postgres up

postgres-clean:
	$(MAKE) -f Makefile.postgres clean

migrate-up:
	$(MAKE) -f Makefile.migrate up

main-build-run-server-producer:
	$(MAKE) -f Makefile.main build-run-server-producer

run-all: kafka-up redis-up postgres-up migrate-up main-build-run-server-producer

run-main: main-build-run-server-producer

clean-all: kafka-clean redis-clean postgres-clean