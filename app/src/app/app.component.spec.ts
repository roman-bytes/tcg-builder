import { TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { ApiService } from './api.service';
import { CookieService } from 'ngx-cookie-service';
import { of, throwError } from 'rxjs';

describe('AppComponent', () => {
  let app: AppComponent;
  let apiService: jasmine.SpyObj<ApiService>;
  let cookieService: jasmine.SpyObj<CookieService>;

  beforeEach(async () => {
    const apiServiceSpy = jasmine.createSpyObj('ApiService', ['getRandomCard', 'storeCard']);
    const cookieServiceSpy = jasmine.createSpyObj('CookieService', ['get']);

    await TestBed.configureTestingModule({
      declarations: [AppComponent],
      providers: [
        { provide: ApiService, useValue: apiServiceSpy },
        { provide: CookieService, useValue: cookieServiceSpy }
      ],
    }).compileComponents();

    apiService = TestBed.inject(ApiService) as jasmine.SpyObj<ApiService>;
    cookieService = TestBed.inject(CookieService) as jasmine.SpyObj<CookieService>;
    const fixture = TestBed.createComponent(AppComponent);
    app = fixture.componentInstance;
  });

  it('should create the app', () => {
    expect(app).toBeTruthy();
  });

  it('should retrieve stored cards from cookie', () => {
    const mockStoredCards = JSON.stringify([{ id: '1', name: 'Pikachu', image: 'url_to_image' }]);
    cookieService.get.and.returnValue(mockStoredCards);

    app.getCookie();

    expect(app.storedCards).toEqual([{ id: '1', name: 'Pikachu', image: 'url_to_image' }]);
  });

  it('should handle error when storing a card', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };
    apiService.storeCard.and.returnValue(Promise.resolve('Could not store card'));

    app.storeCard(mockCard);

    expect(app.error).toBe('Error storing card');
  });

  it('should fetch a random card successfully', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };
    apiService.getRandomCard.and.returnValue(Promise.resolve(mockCard));

    app.fetchRandomCard();

    expect(app.randomCard).toEqual(mockCard);
  });

  it('should display error message when card already stored', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };
    apiService.storeCard.and.returnValue(Promise.resolve('Card already stored'));

    app.storeCard(mockCard);

    expect(app.error).toBe('Card already stored');
  });

  it('should render title', () => {
    const fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.querySelector('h1')?.textContent).toContain('Pokemon Pick Six of the day');
  });
});