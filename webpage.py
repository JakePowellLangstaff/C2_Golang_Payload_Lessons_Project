from flask import Flask, render_template, send_from_directory, abort
import os

app = Flask(__name__)

# Folder this file lives in — all other paths build from here.
# os.path.realpath resolves symlinks and OneDrive quirks on Windows.
BASE = os.path.realpath(os.path.dirname(os.path.abspath(__file__)))

# Word docs (.docx) — downloadables/lessonDocs/
DOCS = os.path.join(BASE, "downloadables", "lessonDocs")

# Source files (.go, .exe, .py) — downloadables/lessonPayloads/
PAYLOADS = os.path.join(BASE, "downloadables", "lessonPayloads")


# Print paths at startup so you can confirm they look right in the terminal
print("\n[paths]")
print(f"  base:     {BASE}")
print(f"  docs:     {DOCS}  {'OK' if os.path.isdir(DOCS) else '*** FOLDER NOT FOUND ***'}")
print(f"  payloads: {PAYLOADS}  {'OK' if os.path.isdir(PAYLOADS) else '*** FOLDER NOT FOUND ***'}")
print()


# Serve a Word doc — called by /download/<filename> in every HTML page
@app.route("/download/<path:filename>")
def download_file(filename):
    full = os.path.join(DOCS, filename)
    if not os.path.isfile(full):
        print(f"[404] doc not found: {full}")
        abort(404)
    return send_from_directory(DOCS, filename, as_attachment=True)


# Serve a source file or binary — called by /payload/<filename> in every HTML page
@app.route("/payload/<path:filename>")
def download_payload(filename):
    full = os.path.join(PAYLOADS, filename)
    if not os.path.isfile(full):
        print(f"[404] payload not found: {full}")
        abort(404)
    return send_from_directory(PAYLOADS, filename, as_attachment=True)


# ── Pages ────────────────────────────────────────────────────────
@app.route("/")
def home():
    return render_template("home.html")

@app.route("/lessons")
def lessons():
    return render_template("lessons.html")

@app.route("/foundation")
def foundation():
    return render_template("foundation.html")

@app.route("/recon")
def recon():
    return render_template("recon.html")

@app.route("/c2")
def c2():
    return render_template("c2.html")

@app.route("/obfuscation")
def obfuscation():
    return render_template("obfuscation.html")

@app.route("/defender")
def defender():
    return render_template("defender.html")

@app.route("/zoo")
def zoo():
    return render_template("zoo.html")


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8000, debug=True)
