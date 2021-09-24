test:
	go test -count=1 ./...

deploy:
	gcloud functions deploy TestAccessToken --runtime go116 --trigger-http --allow-unauthenticated --region asia-northeast1
