import requests
import json
import os

# Update the API URL to use the Kubernetes service name
api = "http://localhost:5050/api/v1/resource"

# Define the JSON data 0 to be sent
payload0 = {
    "type": "set vled",
    "data": "0"
}

# Define the JSON data 1 to be sent
payload1 = {
    "type": "set vled",
    "data": "1"
}

def set_vled_0():
    try:
        response = requests.post(api, json=payload0)
        response.raise_for_status()  # Raises an HTTPError for bad responses
        print(response.status_code)
        print(response.text)
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")

def set_vled_1():
    try:
        response = requests.post(api, json=payload1)
        response.raise_for_status()  # Raises an HTTPError for bad responses
        print(response.status_code)
        print(response.text)
    except requests.exceptions.RequestException as e:
        print(f"An error occurred: {e}")