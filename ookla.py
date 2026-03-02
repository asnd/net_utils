#!/usr/bin/env python3

# Ookla
# Updated 5/30/18
# Refactored 12/27/25

# This Python script queries a list of available data extract files from Speedtest Intelligence,
# determines what data sets are available, and then downloads the most recent version of each.

import urllib.request as urllib_request
import urllib.error as urllib_error
import json
import os
import base64
import sys
import argparse

EXTRACTS_URL = 'https://intelligence.speedtest.net/extracts'

def get_credentials(args):
    api_key = args.api_key or os.environ.get('OOKLA_API_KEY')
    api_secret = args.api_secret or os.environ.get('OOKLA_API_SECRET')
    
    if not api_key or not api_secret:
        print("Error: API Key and API Secret are required.")
        print("Please provide them via arguments (--api-key, --api-secret) or environment variables (OOKLA_API_KEY, OOKLA_API_SECRET).")
        sys.exit(1)
    
    return api_key, api_secret

def setup_opener(username, password):
    opener = urllib_request.build_opener()
    urllib_request.install_opener(opener)
    opener.addheaders = [('Accept', 'application/json')]

    login_credentials = f'{username}:{password}'
    base64string = base64.b64encode(login_credentials.encode('utf-8')).decode('ascii')
    opener.addheaders = [('Authorization', f'Basic {base64string}')]
    return opener

def fetch_file_list(opener):
    try:
        response = urllib_request.urlopen(EXTRACTS_URL).read()
        return json.loads(response)
    except urllib_error.HTTPError as error:
        if error.code == 401:
            print("Authentication Error\nPlease verify that the API key and secret are correct")
        elif error.code == 404:
            print("The account associated with this API key does not have any files attached to it.\nPlease contact your technical account manager to enable data extracts for this account.")
        elif error.code == 500:
            print("Server Error\nPlease contact your technical account manager")
        else:
            print(f"HTTP Error: {error.code} {error.reason}")
        sys.exit(1)
    except ValueError as err:
        print(f"Error parsing JSON: {err}")
        sys.exit(1)
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
        sys.exit(1)

def sort_files_and_directories(contents, files=None):
    if files is None:
        files = {}
        
    for entry in contents:
        if entry['type'] == 'file' and entry['name'].find('headers') == -1 and '_20' in entry['name']:
            filter_file(entry, files)
        elif entry['type'] == 'dir':
            subdir = EXTRACTS_URL + entry['url']
            try:
                sub_files = json.loads(urllib_request.urlopen(subdir).read())
                sort_files_and_directories(sub_files, files)
            except Exception as e:
                print(f"Error reading subdirectory {subdir}: {e}")

    return files

def filter_file(data_file, files):
    try:
        # identify the dataset by the file name prefix
        dataset = data_file['name'][:data_file['name'].index('_20')]
        if dataset not in files or data_file['mtime'] > files[dataset]['age']:
            files[dataset] = {'name': data_file['name'], 'url': data_file['url'], 'age': data_file['mtime']}
    except ValueError:
        pass # Skip files that don't match the expected pattern

def download(files, storage_dir):
    if not files:
        print("No data extract files found. If this is an error, please contact your technical account manager.")
        return

    if not os.path.exists(storage_dir):
        try:
            os.makedirs(storage_dir)
        except OSError as e:
            print(f"Error creating directory {storage_dir}: {e}")
            sys.exit(1)

    for data_set, file in files.items():
        flocation = os.path.join(storage_dir, file['name'])
        print(f"Downloading: {file['name']}")
        try:
            response = urllib_request.urlopen(file['url'])
            with open(flocation, 'wb') as content:
                content.write(response.read())
        except Exception as e:
            print(f"Failed to download {file['name']}: {e}")

def main():
    parser = argparse.ArgumentParser(description="Ookla Speedtest Intelligence Downloader")
    parser.add_argument("--api-key", help="Ookla API Key")
    parser.add_argument("--api-secret", help="Ookla API Secret")
    parser.add_argument("--output-dir", default=os.getcwd(), help="Directory to store downloaded files")
    
    args = parser.parse_args()
    
    username, password = get_credentials(args)
    opener = setup_opener(username, password)
    
    content = fetch_file_list(opener)
    files = sort_files_and_directories(content)
    download(files, args.output_dir)

if __name__ == "__main__":
    main()