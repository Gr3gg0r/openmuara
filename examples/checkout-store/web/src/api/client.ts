import type { AppConfig, CheckoutRequest, CheckoutResponse, Payment, Product } from '../types';

const API_PREFIX = '/api';

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || `Request failed with ${res.status}`);
  }
  return (await res.json()) as T;
}

export async function fetchProduct(): Promise<Product> {
  const res = await fetch(`${API_PREFIX}/product`);
  return handleResponse<Product>(res);
}

export async function fetchConfig(): Promise<AppConfig> {
  const res = await fetch(`${API_PREFIX}/config`);
  return handleResponse<AppConfig>(res);
}

export async function createCheckout(payload: CheckoutRequest): Promise<CheckoutResponse> {
  const res = await fetch(`${API_PREFIX}/checkout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return handleResponse<CheckoutResponse>(res);
}

export async function fetchPayment(ref: string): Promise<Payment> {
  const res = await fetch(`${API_PREFIX}/payment/${encodeURIComponent(ref)}`);
  return handleResponse<Payment>(res);
}
