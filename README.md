C2 Zoo — Documentation Website
Author: Jake Powell Langstaff Module: Cyber Security level 8 honors final year 2026 Student Id: C00287662 Course No: CW_KCCYB

C2 Zoo — Documentation Website Educational C2 framework documentation portal.

Requirements Python 3.11+

Setup

Clone the repo git clone cd C2_zooSetup_early

Create virtual environment python -m venv venv

Activate virtual environment Windows: venv\Scripts\activate

Linux/Mac: source venv/bin/activate

Install dependencies pip install flask gunicorn
Note : gunicorn optional if just marking ignore That is all that is needed — no other packages required.

For webpage lessonsRun locally after virtual environment made and dependiency flask is installed python webpage.py
Site runs at http://localhost:8000
Run with Gunicorn (production) gunicorn --bind 0.0.0.0:8000 --workers 2 webpage:app
Note on Ports Default port is 8000. If port 8000 is already in use on your machine

Flask will show the port in terminal output Or with Gunicorn on a different port: gunicorn --bind 0.0.0.0:5001 --workers 2 webpage:app When deployed to Kubernetes the site is exposed on NodePort 30081.
C2 server Simply run: python c2_upload.py
Once listening navigate to downloads, then lesson payloads. Run the exceuctables payload 3 or up(payload 3 is just v3 compiled etc) will be already compiled so only need to recompile if making modifications like changing ip address to different machine. for c2 server.

Note: If you want to test c2 server on custom ip best too use net_scan_v3.go after editing run the command: go build -ldflags="-s -w" -o payload3.exe .\net_scan_v3.go.

this is best because v4 and up are obfuscated making it a pain.

Note: v1 and v2 just save files to local users temp folder c2 is only implemented in v3 upward.

location in code:
Project Structure C2_zooSetup_early ├── webpage.py # Main Flask app ├── templates/ # HTML pages ├── static/ # CSS and assets ├── downloadables/ │ ├── lessonDocs/ # Word document walkthroughs │ └── lessonPayloads/ # Go source files └── uploads/ # C2 beacon uploads (only if implemented from front end c2 section and configured not in form of downloadable due to ethical considerations)
Ethical Notice The Compiled Binaries the .exe will not be uploaded to github or anywhere public. this is only inteded for SETU academic staff for marking purposes and will be uploaded to setus portal assessment section for final ear project to prove functionality of c2 server and other functions.
