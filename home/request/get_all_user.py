import requests

url = "http://localhost:8000/get_all_user/"

response = requests.get(url)

print("Response Content:", response.content)
print("Status Code:", response.status_code)