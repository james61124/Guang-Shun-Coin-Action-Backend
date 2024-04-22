import requests

url = "http://localhost:8000/login/"

data = {
    'username': 'james',
    'password': 'secret'
}

jwt_token = "james"
headers = {
    'Authorization': f'Bearer {jwt_token}',
    'User-Agent': 'MyApp/1.0',
    'Content-Type': 'application/x-www-form-urlencoded',
}

response = requests.post(url, headers=headers, data=data)

print("Response Content:", response.content)
print("Status Code:", response.status_code)