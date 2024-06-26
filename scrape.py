import os
import requests
import time

def fetch_plugin_list(page_number):
    url = f"https://api.wordpress.org/plugins/info/1.2/?action=query_plugins&request[page]={page_number}"
    response = requests.get(url)
    response.raise_for_status()
    return response.json()

def download_plugin(slug, version, download_url):
    response = requests.get(download_url)
    response.raise_for_status()
    file_name = f"{slug}-{version}.zip"
    with open(file_name, 'wb') as f:
        f.write(response.content)

def main():
    page_number = 1
    while True:
        try:
            plugin_list = fetch_plugin_list(page_number)
            if not plugin_list['plugins']:
                break  # No more plugins to fetch
            for plugin in plugin_list['plugins']:
                slug = plugin['slug']
                version = plugin['version']
                download_url = plugin['download_link']
                print(f"Downloading {slug} version {version} from {download_url}")
                download_plugin(slug, version, download_url)
                time.sleep(1)  # Sleep to avoid overwhelming the server
            page_number += 1
        except requests.exceptions.RequestException as e:
            print(f"An error occurred: {e}")
            time.sleep(5)
        except Exception as e:
            print(f"An unexpected error occurred: {e}")
            break

if __name__ == "__main__":
    main()
