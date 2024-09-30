import { Component ,OnInit} from '@angular/core';
import { CommonModule } from '@angular/common'; // Diğer temel modüller
import { RouterModule } from '@angular/router'; // RouterModule'ü dahil et

@Component({
  selector: 'app-navbar',
  standalone: true,
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css'],
  imports: [CommonModule, RouterModule] // RouterModule'ü buraya ekle
})
export class NavbarComponent implements OnInit {
  isLoggedIn = false;

  ngOnInit(): void {
    // Token'ın olup olmadığını kontrol et
    const token = localStorage.getItem('token');
    this.isLoggedIn = !!token; // Token varsa true, yoksa false
  }
}