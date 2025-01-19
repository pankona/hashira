test:
	go test ./...

deploy-all:
	make deploy-call

deploy-call:
	gcloud functions deploy call --gen2 --entry-point Call --runtime go122 --memory 256Mi --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
