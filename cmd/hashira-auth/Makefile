# retrieve list of app versions to delete
VERSIONS = $(shell gcloud app versions list --service default --format="value(version.id)" --filter="traffic_split=0.0")

all:

devserver:
	DATASTORE_EMULATOR_HOST=localhost:8081 GO111MODULE=on python2.7 `which dev_appserver.py` $(CURDIR)/app.yaml

deploy:
	GO111MODULE=on gcloud app deploy -q --promote --stop-previous-version

datastore-emulator:
	GO111MODULE=on gcloud beta emulators datastore start --data-dir=$(CURDIR)

list:
	gcloud app versions list

delete:
	gcloud app versions delete $(VERSIONS)

versions:
	@echo $(VERSIONS)
