# Build the Docker image
docker build -t gcr.io/$PROJECT_ID/vertex-api .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/vertex-api

# Deploy to Cloud Run
gcloud run deploy vertex-api \
  --image gcr.io/$PROJECT_ID/vertex-api \
  --platform managed \
  --region $REGION \
  --service-account vertex-ai-service-account@$PROJECT_ID.iam.gserviceaccount.com \
  --set-env-vars PROJECT_ID=$PROJECT_ID,REGION=$REGION,ENDPOINT_ID=$ENDPOINT_ID 