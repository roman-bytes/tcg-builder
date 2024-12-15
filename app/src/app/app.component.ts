import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from './api.service';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  standalone: true,
  imports: [CommonModule],
  providers: [CookieService],
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  randomCard: any;
  storedCards: any[] = [];
  error: any;

  constructor(private apiService: ApiService, private cookieService: CookieService) { }

  ngOnInit(): void {
    this.fetchRandomCard();
    this.getCookie();
  }

  getCookie(): void {
    const cookieCards = this.cookieService.get('storedCards');
    if (cookieCards) {
      try {
        this.storedCards = JSON.parse(cookieCards);
      } catch (error) {
        console.error('Error parsing stored cards from cookie', error);
      }
    } else {
      console.log('No stored cards found in cookie')
    }
  }

  // ** Get a random card */
  fetchRandomCard(): void {
    this.apiService.getRandomCard()
      .then(card => {
        this.randomCard = card;
      })
      .catch(error => {
        console.error('Error fetching random card', error);
      });
  }

  // ** Store a card */
  storeCard(card: any): void {
    this.apiService.storeCard({
      id: card.id,
      name: card.name,
      image: card.images.large
    })
      .then(response => {
        this.getCookie();
        this.fetchRandomCard();
        if (response === 'Card already stored') {
          this.error = 'Card already stored';
        }
      })
      .catch(error => {
        console.error('Error storing card', error);
      })
  }

  // Next button click method
  nextRandomCard(): void {
    this.fetchRandomCard();
  }

  // Save the current card
  saveCurrentCard(): void {
    if (this.randomCard && this.randomCard.length > 0) {
      this.storeCard(this.randomCard[0]);
    } else {
      console.error('No card to store');
    }
  }
}
