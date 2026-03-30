from flask import Flask, render_template

app = Flask(__name__)

@app.route("/")
def home():
     return render_template("home.html")
            

@app.route("/zoo")
def zoo():
    return render_template("zoo.html")

@app.route("/c2")
def c2():
    return render_template("c2.html")


@app.route("/lessons")
def lessons():
    return render_template("lessons.html")

@app.route("/recon")
def recon():
    return render_template("recon.html")

@app.route("/foundation")
def foundation():
    return render_template("foundation.html")

@app.route("/obfuscation")
def obfuscation():
    return render_template("obfuscation.html")

# ===============================================
# START APPLICATION
# ===============================================

if __name__ == "__main__":
    
    app.run(host="0.0.0.0", port=8000)  # allow external connections