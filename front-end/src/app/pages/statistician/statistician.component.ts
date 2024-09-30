import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';  // Import FormsModule

interface Document {
  ID: string;              // MongoDB'deki ID
  OriginalName: string;    // Orijinal dosya adı
  Path: string;            // Dosya yolu
  Owner: string;           // Sahip bilgisi
}

@Component({
  selector: 'app-statistician',
  standalone: true,
  templateUrl: './statistician.component.html',
  styleUrls: ['./statistician.component.css'],
  imports: [CommonModule, FormsModule]  // Include FormsModule in imports
})
export class StatisticianComponent implements OnInit {
  documents: Document[] = [];  // Documents array
  selectedDocument: Document | null = null; // Seçilen belge
  isDocumentLocked: boolean = false; // Belge seçildikten sonra kilitlenip kilitlenmediğini takip etmek için

  selectedOption: string | null = null; // ML veya Report seçimi
  mlQuestions = [
    'Question 1 for ML',
    'Question 2 for ML',
    'Question 3 for ML',
    'Question 4 for ML',
    'Question 5 for ML',
    'Question 6 for ML',
    'Question 7 for ML'
  ];

  reportQuestions = [
    'Question 1 for Report',
    'Question 2 for Report'
  ];

  mlAnswers: string[] = Array(this.mlQuestions.length).fill('');
  reportAnswers: string[] = Array(this.reportQuestions.length).fill('');

  currentQuestionIndex: number = 0;

  constructor(private http: HttpClient) {}

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

  selectDocument(document: Document): void {
    this.selectedDocument = document;
  }

  // Dökümanı seçtikten sonra kilitler ve sorular açılabilir hale gelir
  lockSelectedDocument(): void {
    if (this.selectedDocument) {
      this.isDocumentLocked = true;
      this.currentQuestionIndex = 0;  // İlk soruyu göster
    }
  }

  onSelectOption(option: string): void {
    this.selectedOption = option;
    this.currentQuestionIndex = 0;  // İlk soruyu göster
  }

  showNextQuestion(): void {
    this.currentQuestionIndex += 1;
  }

  canProceed(): boolean {
    if (this.selectedOption === 'ml') {
      return this.mlAnswers.every(answer => answer.trim() !== '');
    } else if (this.selectedOption === 'report') {
      return this.reportAnswers.every(answer => answer.trim() !== '');
    }
    return false;
  }

  proceed(): void {
    if (this.selectedOption === 'ml') {
      console.log('ML Prediction with answers:', this.mlAnswers);
    } else if (this.selectedOption === 'report') {
      console.log('Creating report with answers:', this.reportAnswers);
    }
  }
  sendAnswerToBackend(questionId: number, answer: string): void {
    if (!this.selectedDocument) {
      console.error('No document selected.');
      return;
    }
  
    const body = { 
      question_id: questionId, 
      answer: answer,
      document_path: this.selectedDocument.Path,   // Seçili dosyanın path'i
      document_id: this.selectedDocument.ID        // Seçili dosyanın ID'si
    };
  
    console.log('Sending answer with document details:', body);
  
    this.http.post('http://localhost:8080/process-question', body)
      .subscribe(response => {
        console.log('Answer processed:', response);
        this.showNextQuestion();  // Bir sonraki soruyu göster
      }, error => {
        console.error('Error processing question:', error);
      });
  }
  
}
