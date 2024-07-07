import requests
import json

url = "http://localhost:81/member/addProduct"
products = [
    {
        "name": "product1",
        "category": "american coin",
        "price": 100,
        "minBidPrice": 123,
        "startDate": "2024-06-30T12:00:00Z",
        "endDate": "2024-07-10T12:00:00Z",
        "productDescription": "best!!"
    },
    {
        "name": "product2",
        "category": "american coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product3",
        "category": "american coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-02T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product4",
        "category": "american coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-03T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product5",
        "category": "american coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-04T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product6",
        "category": "american coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-05T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product7",
        "category": "europe coin",
        "price": 100,
        "minBidPrice": 123,
        "startDate": "2024-06-30T12:00:00Z",
        "endDate": "2024-07-10T12:00:00Z",
        "productDescription": "best!!"
    },
    {
        "name": "product8",
        "category": "europe coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product9",
        "category": "europe coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product10",
        "category": "europe coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product11",
        "category": "europe coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    },
    {
        "name": "product12",
        "category": "europe coin",
        "price": 200,
        "minBidPrice": 223,
        "startDate": "2024-07-01T12:00:00Z",
        "endDate": "2024-07-11T12:00:00Z",
        "productDescription": "great!!"
    }
]

token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiOTRhN2ViNmYtMTQwMC00MzZiLWE0OTEtYzgzNDQyYWY0MWQzIiwiZXhwIjoxNzIyNTMxMzE0LCJzdWIiOiJuaXJvd3UifQ.ufMwPPKmH3HlT2b_E4Lm5iGfRTeCb9-el5mp5qOVamxO2N1HV18kolV0eVQy67pAQ2WsXq0kMa6VbKFaVV3GkQ"

headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

for product in products:
    response = requests.post(url, headers=headers, data=json.dumps(product))
    print(f"Status Code: {response.status_code}")
    print(f"Response Content: {response.content.decode('utf-8')}")