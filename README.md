[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-2e0aaae1b6195c2367325f4f02e2d04e9abb55f0b24a779b69b11b9e10269abc.svg)](https://classroom.github.com/online_ide?assignment_repo_id=19954278&assignment_repo_type=AssignmentRepo)
# Graded Challenge 2 - P3

Graded Challenge ini dibuat guna mengevaluasi pembelajaran pada Hacktiv8 Program Fulltime Golang khususnya pada pembelajaran grpc

## Assignment Objective
Graded Challenge 2 ini dibuat guna mengevaluasi pemahaman gRPC sebagai berikut:

- Mampu memahami konsep gRPC
- Mampu membuat service dengan gRPC
- Mampu melakukan pengujian dengan unit testing
- Mampu melakukan deployment ke Google Cloud dengan menggunakan docker image


## Assignment Directions

1. Buat proto file untuk definisi layanan gRPC dengan pesan (message) yang sesuai.
2. Implementasikan server gRPC untuk service Payments yang dibuat pada project sebelumnya menggunakan REST API.
3. Tambahkan service gRPC dan REST API untuk service Payment untuk melihat data payment secara keseluruhan, melihat data payment berdasarkan payment id, dan menghapus data payment berdasarkan payment id.
4. Dokumentasikan keseluruhan service dengan Swagger
5. Lakukan pengujian menggunakan unit testing minimal 1 function
6. Setelah aplikasi running maka lakukan deployment dengan Google Cloud menggunakan docker image


#### Sebagai tambahan dari requirement yang sudah diberikan sebelumnya, Student juga diharapkan untuk memahami dan menerapkan konsep-konsep berikut:
- Cloud Deployment using GCP
Student diharapkan untuk mengimplementasikan Cloud Deployment menggunakan Google Cloud Platform (GCP).
Pastikan aplikasi Anda dapat diakses secara publik setelah deployment.
Sediakan dokumentasi sederhana mengenai langkah-langkah deployment yang Anda lakukan.
- Job Scheduling
Implementasikan konsep job scheduling untuk beberapa proses yang memerlukannya, seperti proses pembaharuan data atau pembersihan data yang tidak diperlukan.
- Unit Test
Buat unit test untuk memastikan bahwa setiap fungsi atau method dalam aplikasi Anda bekerja dengan semestinya.
- Docker
Kontainerisasi aplikasi Anda menggunakan Docker.
Pastikan Anda menyediakan Dockerfile dan dokumentasi singkat tentang bagaimana menjalankan aplikasi Anda menggunakan Docker.

## Database Schema:
Database sesuaikan dengan schema yang ada pada Graded Challenge 1


## Expected Result
- File Proto.
- Deploy menggunakan url Google Cloud hingga mendapat public url,  contoh : http://url-google-cloud.com/
- Dokumentasi API Swagger.
- Menambahkan autentikasi token JWT ke layanan gRPC

