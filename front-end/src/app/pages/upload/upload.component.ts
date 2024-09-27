import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
// Dosya uzantısını kontrol eden bir yardımcı fonksiyon
function isValidFileExtension(fileName: string): boolean {
  const allowedExtensions = ['csv', 'xls', 'xlsx'];  // İzin verilen uzantılar
  const fileExtension = fileName.split('.').pop()?.toLowerCase(); // Dosya uzantısını al
  return allowedExtensions.includes(fileExtension || '');
}

@Component({
  selector: 'app-upload',
  standalone: true,
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent {
  selectedFile: File | null = null;

  constructor(private http: HttpClient, private router: Router) {}

  // Dosya seçildiğinde çağrılan fonksiyon
  onFileSelected(event: any): void {
    this.selectedFile = event.target.files[0];
    const file = this.selectedFile;
    if (file) {
      const fileNameElement = document.getElementById('file-name');
      fileNameElement!.textContent = file.name;

      // Dosya uzantısını kontrol ediyoruz
      if (!isValidFileExtension(file.name)) {
        alert('Invalid file type. Please upload a CSV or Excel file.');
        this.selectedFile = null; // Geçersiz dosya seçildiğinde sıfırla
        return;
      }
    }
  }

  // Form submit edildiğinde çağrılan fonksiyon
  onSubmit(): void {
    console.log("Submitting the form without page refresh");
  
    if (this.selectedFile) {
      const formData = new FormData();
      formData.append('file', this.selectedFile);
  
      // Local storage'dan JWT token'ını alıyoruz
      const token = localStorage.getItem('token');
      if (!token) {
        alert('No token found. Please log in again.');
        return;
      }
  
      // Authorization başlığına token'ı ekliyoruz
      const headers = { 'Authorization': `Bearer ${token}` };
  
      this.http.post('http://localhost:8080/upload', formData, { headers }).subscribe(
        response => {
          console.log('File uploaded successfully', response);
          alert('File uploaded successfully!');
        },
        error => {
          console.error('Error uploading file', error);
          alert('Error uploading file. Please try again.');
        }
      );
    } else {
      alert('Please select a file to upload.');
    }
  }

  goToMyDocuments(): void {
    this.router.navigate(['/mydocuments']);
  }
  
  // Logout fonksiyonu
  logout(): void {
    localStorage.removeItem('token'); // Local storage'daki token'ı sil
    this.router.navigate(['/login']); // Kullanıcıyı login sayfasına yönlendir
  }
}
