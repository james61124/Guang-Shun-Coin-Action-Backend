import requests

url = "http://localhost:8000/delete_user/"

data = {
    'username': 'james'
}

response = requests.post(url, data=data)

print("Response Content:", response.content)
print("Status Code:", response.status_code)