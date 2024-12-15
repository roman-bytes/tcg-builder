import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  //TODO: update for production. ENV?
  private baseURL = 'http://localhost:3000';

  constructor() { }

  //** Get a random card */
  getRandomCard(): Promise<any> {
    return fetch(`${this.baseURL}/random-card`)
      .then(response => {
        if (!response.ok) {
          throw new Error('Could not fetch random card');
        }
        return response.json();
      })
  }

  //** Store a card */
  storeCard(card: any): Promise<any> {
    return fetch(`${this.baseURL}/store`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(card)
    }).then(response => {
      if (!response.ok) {
        throw new Error('Could not store card');
      }
      return response;
    });
  }

  //** Get stored cards */
  getStoredCards(): Promise<any> {
    return fetch(`${this.baseURL}/stored`)
      .then(response => {
        if (!response.ok) {
          throw new Error('Could not fetch stored cards');
        }
        return response.json();
      });
  }
}
