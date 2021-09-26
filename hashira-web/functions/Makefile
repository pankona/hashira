test:
	go test -count=1 ./...

deploy:
	gcloud functions deploy test-access-token --entry-point TestAccessToken --runtime go116 --memory 128MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
	gcloud functions deploy upload --entry-point Upload --runtime go116 --memory 128MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
	gcloud functions deploy download --entry-point Download --runtime go116 --memory 128MB --trigger-http --allow-unauthenticated --region asia-northeast1 --max-instances 10
