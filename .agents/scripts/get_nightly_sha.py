#!/usr/bin/env python3
import sys
import os
import subprocess
import datetime
import urllib.request
import urllib.error
import json
import argparse

DAYS = {
    "monday": 0, "mon": 0,
    "tuesday": 1, "tue": 1, "tues": 1,
    "wednesday": 2, "wed": 2,
    "thursday": 3, "thu": 3, "thur": 3, "thurs": 3,
    "friday": 4, "fri": 4,
    "saturday": 5, "sat": 5,
    "sunday": 6, "sun": 6,
}

def get_la_timezone():
    try:
        from zoneinfo import ZoneInfo
        return ZoneInfo("America/Los_Angeles")
    except ImportError:
        return datetime.timezone(datetime.timedelta(hours=-7))

def get_current_time_la():
    dt = datetime.datetime.now(datetime.timezone.utc)
    return dt.astimezone(get_la_timezone())

def get_cutoff_datetime(target_day_name, now_la):
    target_weekday = DAYS[target_day_name.lower()]
    today = now_la.date()
    days_diff = today.weekday() - target_weekday
    if days_diff < 0:
        days_diff += 7
    candidate_date = today - datetime.timedelta(days=days_diff)
    cutoff = datetime.datetime.combine(candidate_date, datetime.time(23, 59, 59, 999999), tzinfo=get_la_timezone())
    if now_la < cutoff:
        candidate_date -= datetime.timedelta(days=7)
        cutoff = datetime.datetime.combine(candidate_date, datetime.time(23, 59, 59, 999999), tzinfo=get_la_timezone())
    return cutoff


def format_date(raw_date):
    if not raw_date:
        return "UNKNOWN"
    try:
        dt = datetime.datetime.strptime(raw_date, "%Y%m%dT%H%M%S%z")
        try:
            from zoneinfo import ZoneInfo
            dt_la = dt.astimezone(ZoneInfo("America/Los_Angeles"))
        except ImportError:
            dt_la = dt.astimezone(datetime.timezone(datetime.timedelta(hours=-7)))
        return dt_la.strftime("%a %b %d %H:%M:%S %Z %Y")
    except Exception:
        return raw_date

def run_git(args):
    result = subprocess.run(["git"] + args, capture_output=True, text=True, check=True)
    return result.stdout.strip()

import ssl

def query_teamcity(sha_short, token):
    url = f"https://hashicorp.teamcity.com/app/rest/builds?locator=project:(id:TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS),branch:(name:refs/heads/nightly-test),number:{sha_short},count:1000&fields=build(id,status,state,finishDate,webUrl)"
    req = urllib.request.Request(url)
    req.add_header("Authorization", f"Bearer {token}")
    req.add_header("Accept", "application/json")
    
    context = ssl.create_default_context()
    try:
        ctx = ssl._create_unverified_context()
        with urllib.request.urlopen(req, context=ctx) as response:
            body = response.read().decode("utf-8")
            data = json.loads(body)
    except urllib.error.HTTPError as e:
        print(f"STATUS=QUERY_FAILED (HTTP {e.code})")
        return None
    except Exception as e:
        print(f"STATUS=QUERY_FAILED ({str(e)})")
        return None

    builds = data.get("build", [])
    if not builds:
        return {
            "status": "NO_BUILD_FOUND",
            "finish_date": "UNKNOWN",
            "web_url": "https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?branch_TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS=refs%2Fheads%2Fnightly-test"
        }

    success_count = 0
    failure_count = 0
    running_count = 0
    queued_count = 0
    finish_dates = []
    compute_web_url = ""
    compute_status = "UNKNOWN"
    first_web_url = ""

    for b in builds:
        state = b.get("state")
        status = b.get("status")
        f_date = b.get("finishDate")
        web_url = b.get("webUrl", "")
        
        if not first_web_url:
            first_web_url = web_url
        if "PACKAGE_COMPUTE" in web_url:
            if not compute_web_url: # Take the first (newest) one
                compute_web_url = web_url
                compute_status = status if state == "finished" else state.upper()
            
        if f_date:
            finish_dates.append(f_date)
            
        if state == "finished":
            if status == "SUCCESS":
                success_count += 1
            else:
                failure_count += 1
        elif state == "running":
            running_count += 1
        elif state == "queued":
            queued_count += 1

    if running_count > 0 or queued_count > 0:
        overall_status = f"IN_PROGRESS (succeeded: {success_count}, failed: {failure_count}, running: {running_count}, queued: {queued_count}, total: {len(builds)})"
    else:
        overall_status = f"COMPLETED (succeeded: {success_count}, failed: {failure_count}, total: {len(builds)})"

    latest_finish = max(finish_dates) if finish_dates else ""

    return {
        "status": overall_status,
        "finish_date": format_date(latest_finish),
        "web_url": "https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?branch_TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS=refs%2Fheads%2Fnightly-test",
        "compute_status": compute_status,
        "compute_web_url": compute_web_url if compute_web_url else "https://hashicorp.teamcity.com/project/TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS?branch_TerraformProviders_GoogleCloud_GOOGLE_NIGHTLYTESTS=refs%2Fheads%2Fnightly-test"
    }

def main():
    parser = argparse.ArgumentParser(description="Get nightly branch commit SHA for a specific day of the week.")
    parser.add_argument("remote", nargs="?", default="origin", help="Git remote to fetch from")
    parser.add_argument("--day-of-week", "--day", dest="day_of_week", default="thursday", help="Day of the week to take the cut from (default: thursday)")
    args = parser.parse_args()
    
    remote = args.remote
    day_name = args.day_of_week.lower()
    if day_name not in DAYS:
        print(f"Error: Invalid day of week '{args.day_of_week}'. Supported: {', '.join(sorted(list(set(DAYS.keys()))))}")
        sys.exit(1)
    
    token = os.environ.get("TEAMCITY_TOKEN")
    if not token:
        print("AGENT_INSTRUCTION: TEAMCITY_TOKEN is not set in the environment. Please guide the user to generate and set one:\n"
              "1. Log into TeamCity at https://hashicorp.teamcity.com/\n"
              "2. Go to Profile -> Access Tokens (https://hashicorp.teamcity.com/profile.html?item=accessTokens)\n"
              "3. Click 'Create Access Token', enter a name (e.g. 'release-agent'), and copy the token.\n"
              "4. Export it in shell: export TEAMCITY_TOKEN=\"<token>\" (ensure no spaces around '=').")
        sys.exit(1)

    print(f"Fetching refs/heads/nightly-test from {remote}...")
    try:
        subprocess.run(["git", "fetch", remote, "refs/heads/nightly-test", "--quiet"], check=True)
    except subprocess.CalledProcessError:
        print(f"AGENT_INSTRUCTION: Could not fetch refs/heads/nightly-test from {remote}. Check connection or remote name.")
        sys.exit(1)

    now_la = get_current_time_la()
    cutoff_dt = get_cutoff_datetime(day_name, now_la)
    cutoff_str = cutoff_dt.isoformat()

    print(f"Calculating cutoff for most recent {day_name.capitalize()} night cut...")
    print(f"Cutoff datetime (LA): {cutoff_dt.strftime('%Y-%m-%d %H:%M:%S %Z')}")

    try:
        nightly_sha = run_git(["log", "-n", "1", "--format=%H", f"--before={cutoff_str}", "FETCH_HEAD"])
    except Exception as e:
        print(f"AGENT_INSTRUCTION: Error running git log to find candidate SHA: {str(e)}")
        sys.exit(1)

    if not nightly_sha:
        print(f"AGENT_INSTRUCTION: Could not find any commit on refs/heads/nightly-test on or before {cutoff_str}.")
        sys.exit(1)

    short_sha = nightly_sha[:7]
    print(f"Found Nightly Branch Commit SHA: {nightly_sha} (Short: {short_sha})")

    print(f"Querying TeamCity for build status of short SHA: {short_sha}...")
    result = query_teamcity(short_sha, token)
    if not result:
        sys.exit(1)

    try:
        commit_msg = run_git(["log", "-n", "1", "--format=%s", nightly_sha])
    except Exception:
        commit_msg = "UNKNOWN"

    print("--- TEAMCITY RESULTS ---")
    print(f"Candidate SHA: {nightly_sha}")
    print(f"Commit Message: {commit_msg}")
    print(f"Build Status: {result['status']}")
    print(f"Finish Date: {result['finish_date']}")
    print(f"Web URL: {result['web_url']}")
    print(f"Compute Status: {result['compute_status']}")
    print(f"Compute Web URL: {result['compute_web_url']}")
    print("------------------------")

    if result['status'] == "NO_BUILD_FOUND":
        print(f"AGENT_INSTRUCTION: No TeamCity build found for SHA {short_sha}. Advise the user to manually verify correlation.")

if __name__ == "__main__":
    main()
