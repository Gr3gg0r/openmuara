export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  imageUrl: string;
}

export interface Payment {
  ref: string;
  method: 'fawry' | 'stripe';
  status: 'pending' | 'paid' | 'canceled';
  email: string;
  name: string;
  amount: number;
  currency: string;
  productId: string;
  createdAt: number;
}

export interface CheckoutRequest {
  method: 'fawry' | 'stripe';
  email: string;
  name: string;
}

export interface CheckoutResponse {
  ok: boolean;
  ref: string;
  redirectUrl: string;
}
