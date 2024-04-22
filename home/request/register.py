import requests

url = "http://localhost:81/register"

data = {
    'username': 'james',
    'password': 'secret'
}

headers = {'Content-Type': 'application/json'}
response = requests.post(url, json=data, headers=headers)

print("Response Content:", response.content)
print("Status Code:", response.status_code)