test:
	go test ./...

deploy-all:
	make deploy-ping
	make deploy-upload
	make deploy-download
	make deploy-test-access-token

deploy-ping:
	gcloud functions deploy ping --entry-point Ping --runtime go116 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

deploy-upload:
	gcloud functions deploy upload --entry-point Upload --runtime go116 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

deploy-download:
	gcloud functions deploy download --entry-point Download --runtime go116 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

deploy-add:
	gcloud functions deploy add --entry-point Add --runtime go116 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

deploy-test-access-token:
	gcloud functions deploy test-access-token --entry-point TestAccessToken --runtime go116 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
