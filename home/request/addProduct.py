import requests
import json

# 請求 URL
url = "http://localhost:81/member/addProduct"  # 請替換為實際的 API endpoint

# 請求的數據
data = {
    "name": "product1",
    "category": "american coin",
    "price": 100,
    "minBidPrice": 123,
    "startDate": "2024-06-30T12:00:00Z",
    "endDate": "2024-07-10T12:00:00Z",
    "productDescription": "best!!"
}

# Authorization token
token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiMTIyNTAwY2UtMGM0YS00NjhmLWJlNmQtNzgwNzgzZGMyYjRlIiwiZXhwIjoxNzIyMzQ3MjczLCJzdWIiOiJuaXJvd3UifQ.-GE1BgbmLkZHCD9e_IEVNyUVItVEMoEM_lzvkTa_RdhX1hYwxIcIhVAhnnjMvsiQAmw0cTYTufB61iIYQmQ3CA"  # 請替換為實際的 token

# 請求頭
headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

# 發送 POST 請求
response = requests.post(url, headers=headers, data=json.dumps(data))

# 輸出響應狀態碼和內容
print(f"Status Code: {response.status_code}")
print(f"Response Content: {response.content.decode('utf-8')}")