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
    console.log('Getting cookie');
    const cookieCards = this.cookieService.get('storedCards');
    console.log('Cookie cards:', cookieCards);
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

      })
      .catch(error => {
        if (error.message === 'Card already stored') {
          this.error = 'Card already stored';
          return;
        }
        if (error.message === 'Could not store card') {
          this.error = 'Limit of 6 cards reached';
          return;
        }
      })
  }

  // Next button click method
  nextRandomCard(): void {
    this.fetchRandomCard();
  }

  // Save the current card
  saveCurrentCard(): void {
    if (this.randomCard && this.randomCard.length > 0) {
      if (this.storedCards.length == 6) {
        this.error = 'Limit of 6 cards reached';
        return;
      }
      
      this.storeCard(this.randomCard[0]);
    } else {
      console.error('No card to store');
    }
  }
}
