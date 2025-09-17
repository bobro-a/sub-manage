# Command Examples (curl)

## Create subscription
```bash
curl -X POST http://localhost:9090 \
  -H "Content-Type: application/json" \
  -d '{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",  
  "start_date": "2025-09",
  "end_date": "2025-10"
  }'
```

## Read subscription
```bash
curl -X GET "http://localhost:9090/?id=1"
```
## Read all subscriptions
```bash
curl -X GET "http://localhost:9090/"
```

## Delete subscription by id
```bash
curl -X DELETE "http://localhost:9090/?id=5"
```

## Change subscription by id
```bash
curl -X PUT http://localhost:9090 \
-H "Content-Type: application/json" \
-d '{
"id": 1,
"service_name": "Netflix",
"price": 1200,
"user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
"start_date": "2025-09",
"end_date": "2026-09"
}'
```

## Sum price with conditions
```bash
curl -X GET "http://localhost:9090/sum?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_date=2024-10&end_date=2024-02"
```