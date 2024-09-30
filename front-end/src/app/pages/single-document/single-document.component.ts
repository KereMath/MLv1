import { Component, OnInit } from '@angular/core';
import { ActivatedRoute,Router } from '@angular/router';  
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';  
import * as Papa from 'papaparse';  
import * as XLSX from 'xlsx';       

@Component({
  selector: 'app-single-document',
  templateUrl: './single-document.component.html',
  styleUrls: ['./single-document.component.css'],
  standalone: true,
  imports: [CommonModule]  
})
export class SingleDocumentComponent implements OnInit {
  documentId: string = '';  
  documentData: any = {};   
  tableData: any[][] = [];  
  displayedData: any[][] = []; // To store paginated data
  currentPage: number = 0;
  rowsPerPage: number = 20; // Number of rows per page (excluding the header)

  constructor(private route: ActivatedRoute, private http: HttpClient) {}

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.documentId = params.get('docID') || '';
      this.fetchDocumentData();
    });
  }

  fetchDocumentData(): void {
    const token = localStorage.getItem('token');
    if (!token) {
      alert('You need to be logged in to view the document.');
      return;
    }

    const headers = { 'Authorization': `Bearer ${token}` };
    this.http.get(`http://localhost:8080/documents/${this.documentId}`, { headers })
      .subscribe((response: any) => {
        this.documentData = response;
        console.log("Fetched Document Data: ", this.documentData);

        if (this.documentData.path.endsWith('.csv')) {
          this.parseCSVBackend();
        } else if (this.documentData.path.endsWith('.xls') || this.documentData.path.endsWith('.xlsx')) {
          this.parseExcelBackend();
        }
      }, error => {
        console.error('Error fetching document data:', error);
      });
  }

  // Backend'den dosya içeriğini alarak CSV olarak parse etme
  parseCSVBackend(): void {
    this.http.get(`http://localhost:8080/document-content/${this.documentId}`, { responseType: 'text' }).subscribe(csvData => {
      Papa.parse(csvData, {
        complete: (result: Papa.ParseResult<string[]>) => {
          this.tableData = result.data.filter(row => row.some(cell => cell && cell.trim() !== ''));  // Sadece boş olmayan satırları ekle
          this.updateDisplayedData(); // Update paginated data
        }
      });
    });
  }

  // Backend'den Excel dosyasını alıp işleme
  parseExcelBackend(): void {
    this.http.get(`http://localhost:8080/document-content/${this.documentId}`, { responseType: 'arraybuffer' }).subscribe(arrayBuffer => {
      const data = new Uint8Array(arrayBuffer);
      const workbook = XLSX.read(data, { type: 'array' });
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];
      const jsonData = XLSX.utils.sheet_to_json<any[]>(worksheet, { header: 1 });  // Satırların dizisi olarak JSON verisi
      this.tableData = jsonData.filter((row: any[]) => row.some(cell => cell && cell.toString().trim() !== ''));  // Sadece boş olmayan satırları ekle
      this.updateDisplayedData(); // Update paginated data
    });
  }

  // Paginated data updates
  updateDisplayedData(): void {
    const start = this.currentPage * this.rowsPerPage + 1; // Skip header row for pagination
    const end = start + this.rowsPerPage;
    this.displayedData = this.tableData.slice(start, end);
  }

  // Navigation functions
  nextPage(): void {
    if ((this.currentPage + 1) * this.rowsPerPage+1 < this.tableData.length) {
      this.currentPage++;
      this.updateDisplayedData();
    }
  }

  previousPage(): void {
    if (this.currentPage > 0) {
      this.currentPage--;
      this.updateDisplayedData();
    }
  }

  getTotalPages(): number {
    return Math.ceil((this.tableData.length - 2) / this.rowsPerPage); // Exclude the header row from pagination
  }
  goToStatistician(): void {
    window.location.href = '/statistician';  // Redirect using href
  }
  
}
