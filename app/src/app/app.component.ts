import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ApiService } from './api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  standalone: true,
  imports: [CommonModule],
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  randomCard: any;
  storedCards: any[] = [];

  constructor(private apiService: ApiService) { }

  ngOnInit(): void {
    this.fetchRandomCard();
    this.fetchStoredCards();

    console.log('randomCard', this.randomCard);
  }

  // ** Get a random card */
  fetchRandomCard(): void {
    this.apiService.getRandomCard()
      .then(card => {
        console.log('Random card fetched successfully', card);
        this.randomCard = card;
      })
      .catch(error => {
        console.error('Error fetching random card', error);
      });
  }

  // ** Get stored cards */
  fetchStoredCards(): void {
    this.apiService.getStoredCards()
      .then(cards => {
        this.storedCards = cards;
      })
      .catch(error => {
        console.error('Error fetching stored cards', error);
      });
  }

  // ** Store a card */
  storeCard(card: any): void {
    this.apiService.storeCard(card)
      .then(response => {
        console.log('Card stored successfully', response);
        this.fetchStoredCards();
      })
      .catch(error => {
        console.error('Error storing card', error);
      })
  }

  // Next button click method
  nextRandomCard(): void {
    this.fetchRandomCard();
  }
}
