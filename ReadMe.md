# C2 Zoo — Documentation Website

ReadMe
# C2 Zoo — Documentation Website

Educational C2 framework documentation portal.

## Requirements

Python 3.11+

## Setup

### 1. Clone the repo
git clone <REPO>
cd C2_zooSetup_early

### 2. Create virtual environment
python -m venv venv

### 3. Activate virtual environment

Windows:
venv\Scripts\activate

Linux/Mac:
source venv/bin/activate

### 4. Install dependencies
pip install flask gunicorn

That is all that is needed — no other packages required.

Note: Gunicorn does not run on Windows. On Windows use:
python webpage.py

Gunicorn is used automatically on the AWS Linux server deployment.
### 5. Run locally
python webpage.py

Site runs at http://localhost:8000

### 6. Run with Gunicorn (production)
gunicorn --bind 0.0.0.0:8000 --workers 2 webpage:app

## Note on Ports
Default port is 8000. If port 8000 is already in use on your machine:

python webpage.py
# Flask will show the port in terminal output

# Or with Gunicorn on a different port:
gunicorn --bind 0.0.0.0:5001 --workers 2 webpage:app
When deployed to Kubernetes the site is exposed on NodePort 30081.
## Project Structure
├── webpage.py              # Main Flask app
├── templates/              # HTML pages
├── static/                 # CSS and assets
├── downloadables/
│   ├── lessonDocs/         # Word document walkthroughs
│   └── lessonPayloads/     # Go source files
└── uploads/                # C2 beacon uploads (only if implemented from front end c2 section   and configured not in form of downloadable due to ethical considerations)
## Ethical Notice
All materials are for educational use in a controlled lab environment only.
No compiled binaries are distributed publicly.
