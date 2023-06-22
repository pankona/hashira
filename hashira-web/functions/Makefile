test:
	go test ./...

deploy-all:
	make deploy-call
	make deploy-ping
	make deploy-upload
	make deploy-download
	make deploy-test-access-token

deploy-call:
	gcloud functions deploy call --entry-point Call --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

# deprecated
deploy-ping:
	gcloud functions deploy ping --entry-point Ping --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

# deprecated
deploy-upload:
	gcloud functions deploy upload --entry-point Upload --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

# deprecated
deploy-download:
	gcloud functions deploy download --entry-point Download --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

# deprecated
deploy-add:
	gcloud functions deploy add --entry-point Add --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10

# deprecated
deploy-test-access-token:
	gcloud functions deploy test-access-token --entry-point TestAccessToken --runtime go120 --memory 256MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
