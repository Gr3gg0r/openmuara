export interface ProviderVersionDetails {
  base_url?: string;
  sample_route?: string;
}

export interface ProviderInfo {
  name: string;
  display_name?: string;
  description?: string;
  enabled?: boolean;
  active?: boolean;
  category?: string;
  real_providers?: string[];
  sample_method?: string;
  sample_route?: string;
  docs_path?: string;
  version?: string;
  versions?: string[];
  version_details?: Record<string, ProviderVersionDetails>;
  env_vars?: string[];
  base_url?: string;
  webhook_target_url?: string;
  is_recommended_for_first_time?: boolean;
}

export interface ProvidersResponse {
  active?: string;
  enabled?: string[];
  available?: string[];
  providers?: Record<string, ProviderInfo>;
}

export interface ProviderDetailResponse extends ProviderInfo {}

export interface OnboardingResponse {
  server_ready?: boolean;
  providers_enabled?: boolean;
  first_transaction?: boolean;
  first_webhook_received?: boolean;
  webhooks_enabled?: boolean;
  next_step?: {
    hint?: string;
    method?: string;
    route?: string;
  };
}

export interface LedgerEvent {
  id?: string;
  time?: string;
  type: 'transaction' | 'webhook' | string;
  provider?: string;
  reference: string;
  status?: string;
  summary?: string;
}

export interface LedgerResponse {
  results?: LedgerEvent[];
}

export interface Transaction {
  reference: string;
  provider?: string;
  amount?: number;
  currency?: string;
  status?: string;
  trace_id?: string;
  customerRef?: string;
  createdAt?: string;
  updatedAt?: string;
  items?: unknown[];
}

export interface TransactionDetailResponse {
  transaction?: Transaction;
}

export interface TransactionsResponse {
  results?: Transaction[];
}

export interface WebhookAttemptEvent {
  time?: string;
  status?: string;
  error?: string;
}

export interface WebhookAttempt {
  ref: string;
  provider?: string;
  provider_name?: string;
  url?: string;
  status?: string;
  attempts?: number;
  last_error?: string;
  trace_id?: string;
  signature_valid?: boolean;
  headers?: Record<string, string>;
  payload?: string | unknown;
  attempt_events?: WebhookAttemptEvent[];
  createdAt?: string;
  updatedAt?: string;
}

export interface WebhooksResponse {
  results?: WebhookAttempt[];
}

export interface WebhookDetailResponse {
  webhook?: WebhookAttempt;
}

export interface FailedWebhooksResponse {
  results?: WebhookAttempt[];
}

export interface WebhookConfigResponse {
  url?: string;
  max_retries?: number;
  targets?: Record<string, string>;
  events?: Record<string, string[]>;
}

export interface ConfigResponse {
  server?: {
    host?: string;
    port?: number;
    admin_port?: number;
  };
  providers?: Record<string, { enabled?: boolean; config?: Record<string, unknown> }>;
  webhook?: WebhookConfigResponse;
}
