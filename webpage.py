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

@app.route("/docs")
def docs():
    return render_template("docs.html")

# ===============================================
# START APPLICATION
# ===============================================

if __name__ == "__main__":
    
    app.run(host="0.0.0.0", port=5000)  # allow external connections