import { TestBed } from '@angular/core/testing';
import { ApiService } from './api.service';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';

const baseURL = 'http://localhost:3000';

describe('ApiService', () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ApiService],
    });
    service = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify(); // Ensure that there are no outstanding requests
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should fetch a random card', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };

    service.getRandomCard().then(card => {
      expect(card).toEqual(mockCard);
    });

    const req = httpMock.expectOne(`${baseURL}/random-card`);
    expect(req.request.method).toBe('GET');
    req.flush(mockCard); // Simulate a successful response
  });

  it('should store a card', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };
    const mockResponse = 'Card stored successfully';

    service.storeCard(mockCard).then(response => {
      expect(response).toBe(mockResponse);
    });

    const req = httpMock.expectOne(`${baseURL}/store`);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(mockCard); // Check if the correct card is sent
    req.flush(mockResponse); // Simulate a successful response
  });

  it('should handle error when fetching a random card', () => {
    service.getRandomCard().catch(error => {
      expect(error).toBeTruthy();
    });

    const req = httpMock.expectOne(`${baseURL}/random-card`);
    req.error(new ErrorEvent('Network error')); // Simulate a network error
  });

  it('should handle error when storing a card', () => {
    const mockCard = { id: '1', name: 'Pikachu', images: { large: 'url_to_image' } };

    service.storeCard(mockCard).catch(error => {
      expect(error).toBeTruthy();
    });

    const req = httpMock.expectOne(`${baseURL}/store`);
    req.error(new ErrorEvent('Network error')); // Simulate a network error
  });
});