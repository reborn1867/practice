#!/usr/bin/env python3

import requests
from requests.auth import HTTPBasicAuth

def get_active_sprint(board_id):
    r = requests.get("https://sapjira.wdf.sap.corp/rest/agile/1.0/board/{}/sprint?state=active".format(board_id), verify=False, auth=HTTPBasicAuth("di-monitor", "Test123!"))
    return r.json()["values"]

def get_report(sprint_id):
    r = requests.get("https://sapjira.wdf.sap.corp/rest/agile/1.0/sprint/{}/issue?fields=summary".format(sprint_id), verify=False, auth=HTTPBasicAuth("di-monitor", "Test123!"))
    
    for issue in r.json()["issues"]:
        print("* [{key}](https://sapjira.wdf.sap.corp/browse/{key}) {summary}".format(key=issue["key"],summary=issue["fields"]["summary"]))

def main():
    active_sprint = get_active_sprint("15293")
    print(active_sprint)
    get_report(active_sprint[0]["id"])
    

if __name__ == "__main__":
    main()