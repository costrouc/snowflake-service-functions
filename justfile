DEFAULT_ROLE := "NATIVE_APP_PUBLISHER"

venv:
	uv venv .venv --python 3.11
	uv pip install -r requirements.txt

docker-login:
	.venv/bin/snow spcs image-registry login --connection posit-software-pbc --role {{DEFAULT_ROLE}}

bootstrap:
	.venv/bin/snow sql -f scripts/01-bootstrap.sql

install:
	docker build -t duloftf-posit-software-pbc.registry.snowflakecomputing.com/cno_service_functions/data/image_repository/service-functions:latest .
	docker push duloftf-posit-software-pbc.registry.snowflakecomputing.com/cno_service_functions/data/image_repository/service-functions:latest
	.venv/bin/snow sql -f scripts/02-install.sql 

uninstall:
	.venv/bin/snow sql -f scripts/03-uninstall.sql 

debug:
	.venv/bin/snow sql -f scripts/04-debug.sql

teardown:
	.venv/bin/snow sql -f scripts/05-teardown.sql
