
build:
	cd $(CURDIR) && go run .

watch:
	cd $(CURDIR) && go run . -watch

clean:
	cd $(CURDIR) && go run . -clean

deploy-firebase-rules:
	firebase deploy --only firestore:rules

install-tools:
	yarn global add firebase-tools

