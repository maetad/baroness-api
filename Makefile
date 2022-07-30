# Directories
ROOT_DIR:=$(shell pwd)
MIGRATIONS_DIR:=${ROOT_DIR}/migrations

.PHONY: migrate

migrate:
	@set -e && \
	for name in $(filter-out $@,$(MAKECMDGOALS)); do \
 		touch ${MIGRATIONS_DIR}/$$(date +"%Y%m%d%H%M%S")_$$name.up.sql; \
 		touch ${MIGRATIONS_DIR}/$$(date +"%Y%m%d%H%M%S")_$$name.down.sql; \
	done

# omit error "No rules to make target" when using `make start` without matching targets
%:
	@:

# m:
# 	@set -e && \
# 	for name in $(MAKECMDGOALS); do \
# 		touch ${MIGRATIONS_DIR}/${TIMESTAMP}_${name}.up.sql; \
# 		touch ${MIGRATIONS_DIR}/${TIMESTAMP}_${name}.down.sql; \
# 	done
