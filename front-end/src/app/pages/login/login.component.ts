import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms'; // Import FormsModule for ngModel

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [FormsModule],  // Include FormsModule here
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {
  username: string = '';
  password: string = '';

  constructor(private http: HttpClient, private router: Router) { }

  // Token kontrolü için ngOnInit lifecycle hook'u ekleyelim
  ngOnInit(): void {
    const token = localStorage.getItem('token');
    if (token) {
      // Eğer token varsa, kullanıcıyı /upload sayfasına yönlendir
      this.router.navigate(['/upload']);
    }
  }

  // Handle the form submission
  onLogin(): void {
    const loginData = { username: this.username, password: this.password };
    this.http.post('http://localhost:8080/login', loginData).subscribe(
      (response: any) => {
        console.log('Login successful:', response);

        // Token varsa, token'i kaydet ve yönlendirme yap
        if (response && response.token) {
          localStorage.setItem('token', response.token); // Token'i localStorage'a kaydet
          alert('Login successful!');
          this.router.navigate(['/upload']); // Yönlendirme yap
        } else {
          console.error('No token found in the response');
          alert('Login failed: No token received.');
        }
      },
      error => {
        console.error('Login failed:', error);
        alert('Login failed. Please check your username and password.');
      }
    );
  }
}
