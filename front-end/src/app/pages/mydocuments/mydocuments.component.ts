import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

interface Document {
  ID: string;              // MongoDB'deki ID'yi karşılamak için
  OriginalName: string;    // Orijinal dosya adı
  Path: string;            // Dosya yolu
  Owner: string;           // Sahip bilgisi
}

@Component({
  selector: 'app-mydocuments',
  standalone: true,
  templateUrl: './mydocuments.component.html',
  styleUrls: ['./mydocuments.component.css'],
  imports: [CommonModule]
})
export class MyDocumentsComponent implements OnInit {
  documents: Document[] = []; // Belgeleri obje olarak tutuyoruz

  constructor(private http: HttpClient, private router: Router) { }

  ngOnInit(): void {
    this.getDocuments();
  }

  getDocuments(): void {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('You need to be logged in to see your documents.');
      return;
    }
  
    const headers = { 'Authorization': `Bearer ${token}` };
    this.http.get<{documents: Document[]}>('http://localhost:8080/documents', { headers })
      .subscribe(response => {
        this.documents = response.documents;
        console.log("Documents fetched: ", this.documents);
      }, error => {
        console.error('Error fetching documents:', error);
      });
  }

  // Confirm deletion
  confirmDelete(documentId: string): void {
    const userConfirmed = confirm('Are you sure you want to delete this document?');
    if (userConfirmed) {
      this.deleteDocument(documentId);
    }
  }

  // Dökümanı silme işlemi
  deleteDocument(documentId: string): void {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('You need to be logged in to delete your documents.');
      return;
    }

    const headers = { 'Authorization': `Bearer ${token}` };
    this.http.delete(`http://localhost:8080/delete-file/${documentId}`, { headers })
      .subscribe(() => {
        // Başarıyla silindikten sonra listeden kaldırıyoruz
        this.documents = this.documents.filter(doc => doc.ID !== documentId);
      }, error => {
        console.error('Error deleting document:', error);
      });
  }

  viewInsights(documentId: string): void {
    this.router.navigate([`/singleDoc/${documentId}`]);  // Yönlendirme yapılıyor
  }
}
