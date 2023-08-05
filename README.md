# ogper

```sh
gcloud functions deploy function-1 \
      --gen2 \
      --runtime=go120 \
      --region=asia-northeast1 \
      --source=. \
      --entry-point=httpFuntion \
      --trigger-http \
      --allow-unauthenticated
```
