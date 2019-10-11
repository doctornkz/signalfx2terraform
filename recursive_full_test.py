import time
import requests
import subprocess

'''
Tool for testing exporter through all your company dashboards.
It allows to check all cases and catch null pointers, types exceptions, etc


'''

TOKEN='<YOUR SFX TOKEN>'
API='https://api.eu0.signalfx.com/v2/'
LIMIT='1000'      
COOLDOWNTIMER = 2 # Pause to prevent SFX API DDOS

headers = {'X-SF-TOKEN': TOKEN}

r = requests.get(API + 'dashboard/?limit='+LIMIT, headers=headers )
dashboards = r.json()

count = 0
full = len(dashboards["results"]) 
for dashboard in dashboards["results"]:
    count +=1 
    time.sleep(COOLDOWNTIMER)
    print("Dashboard processing: % s, % d from %d" %(dashboard['id'], count, full))
    subprocess.check_call(["go", "run", "main.go", "--token", TOKEN, "--dashboard",dashboard["id"]])
    