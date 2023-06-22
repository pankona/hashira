test:
	go test ./...

deploy-call:
	gcloud functions deploy call --entry-point Call --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
