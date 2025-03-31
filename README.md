
# 🚀 rmrf - A Multi-threaded Alternative to `rm -rf`

`rmrf` is a **fast, multi-threaded alternative** to `rm -rf`, written in **Go**.  
It efficiently deletes directories and their contents using **goroutines**, making it significantly faster for large file structures.

## 📌 Features
- ✅ Multi-threaded directory deletion using **goroutines**
- ✅ Automatically **scales with available CPU cores**
- ✅ Works on **Linux, macOS, and Windows**
- ✅ Uses **efficient system calls** (`os.Remove`, `filepath.Walk`)
- ✅ Built-in **documentation server** and **static documentation generator**
- ✅ Includes **automatic linting** to enforce Go best practices

---

## 🚀 Quick Start

### **1️⃣ Install Go (If Not Already Installed)**
Ensure you have **Go 1.20+** installed.  
To check:
```sh
go version
```
If not installed, download from [golang.org/dl](https://golang.org/dl).

---

### **2️⃣ Clone the Repository**
```sh
git clone https://github.com/yourusername/rmrf.git
cd rmrf
```

---

### **3️⃣ One-Step Setup and Build**
Run this command to install **everything needed**, build the project, and generate documentation:
```sh
make
```
This will:
✅ Install all required dependencies (`go mod tidy`)  
✅ Install documentation & linting tools (`godoc`, `golangci-lint`)  
✅ Build the project  
✅ Run the linter  
✅ Generate the static documentation

---

### **4️⃣ Running `rmrf`**
Once built, you can use `rmrf` to delete directories:
```sh
./rmrf my_directory
```
Or install it globally:
```sh
make install
rmrf my_directory
```

---

## 🛠 Developer Guide

### **🔹 Running the Linter**
To ensure code follows Go best practices:
```sh
make lint
```

---

### **📖 Generating Documentation**
#### **1️⃣ Start Local Documentation Server**
```sh
make doc-server
```
- Open **`http://localhost:6060/pkg/`** in your browser.

#### **2️⃣ Generate Static Documentation**
```sh
make doc
```
- This creates **`docs/index.html`**, which you can open in any browser.

---

### **🧹 Cleaning Up**
To remove build artifacts:
```sh
make clean
```

---

## ⚡ **Project Structure**
```
rmrf/
│── main.go        # The main Go program
│── go.mod         # Go module dependencies
│── Makefile       # Build, install, and documentation automation
│── README.md      # This documentation
│── docs/          # Auto-generated static documentation
```

---

## 🚀 **Contributing**
1. Fork the repo & clone your fork.
2. Create a new feature branch.
3. Run `make lint` to ensure best practices.
4. Push your changes & create a PR!

---

## 📜 **License**
MIT License. Feel free to modify and use.
# rmrf
