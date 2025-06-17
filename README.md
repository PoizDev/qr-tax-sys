<div align="center">

# <strong>QR-TAX-SYS</strong>

*Transforming Financial Workflows with Seamless QR Solutions*

![Last Commit](https://img.shields.io/github/last-commit/PoizDev/qr-tax-sys?style=flat-square)
![Language](https://img.shields.io/github/languages/top/PoizDev/qr-tax-sys?style=flat-square)
![Languages Count](https://img.shields.io/github/languages/count/PoizDev/qr-tax-sys?style=flat-square)
![Repo Size](https://img.shields.io/github/repo-size/PoizDev/qr-tax-sys?style=flat-square)

### Built with the tools and technologies:

<img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
<img src="https://img.shields.io/badge/Flutter-02569B?style=for-the-badge&logo=flutter&logoColor=white" />
<img src="https://img.shields.io/badge/Gin-00B2A9?style=for-the-badge&logoColor=white" />
<img src="https://img.shields.io/badge/Dart-0175C2?style=for-the-badge&logo=dart&logoColor=white" />
<img src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white" />
<img src="https://img.shields.io/badge/YAML-CB171E?style=for-the-badge&logo=yaml&logoColor=white" />
<img src="https://img.shields.io/badge/Markdown-000000?style=for-the-badge&logo=markdown&logoColor=white" />

</div>

---

## Table of Contents

* [Overview](#overview)
* [Getting Started](#getting-started)
* [Features](#features)
* [Project Structure](#project-structure)
* [API Endpoints](#api-endpoints)
* [License](#license)

---

## 📖 Overview

QR-TAX-SYS is a full-stack system designed to streamline financial workflows using QR code technology. The backend is built with Go and Gin, while the frontend is a cross-platform Flutter application.

---

## ⚙️ Getting Started

### Backend (Go)

```bash
cd api
cp .env.example .env # create and configure your environment file
go mod tidy
go run main.go
```

Example `.env`:

```env
DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=3306
DB_NAME=qrfatura
```

### Frontend (Flutter)

```bash
flutter pub get
flutter run
```

> 💡 Use an emulator or connect a physical device

---

## ✨ Features

* ✅ User authentication (JWT)
* ✅ Invoice creation and listing
* ✅ Product management
* ✅ QR Code generation for invoices
* ✅ Cross-platform support (Web, Android, iOS, Desktop)

---

## 📁 Project Structure

```
qr-tax-sys/
├── lib/                 # Flutter UI and logic
│   ├── main.dart
│   └── screen/
│       ├── login.dart
│       ├── signup.dart
│       ├── homescreen.dart
│       └── scan.dart
├── api/                 # Go backend
│   ├── main.go
│   ├── controllers/
│   ├── models/
│   ├── db/
│   └── initializers/
├── pubspec.yaml
├── go.mod, go.sum
└── platform folders     # android/, ios/, web/, etc.
```

---

## 📡 API Endpoints

| Method | Endpoint             | Description      |
| ------ | -------------------- | ---------------- |
| POST   | `/signup`            | Register user    |
| POST   | `/login`             | User login       |
| GET    | `/users`             | Get all users    |
| POST   | `/products`          | Add new product  |
| GET    | `/products`          | List products    |
| POST   | `/fatura`            | Create invoice   |
| GET    | `/fatura`            | List invoices    |
| GET    | `/qrcode/:fatura_id` | Generate QR code |

---

## 📄 License

This project is licensed under the **MIT License**.

---

## 🤝 Contribution

Feel free to fork the project and submit pull requests.
