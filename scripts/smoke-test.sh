#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# Pick a free ephemeral port and an isolated workspace so concurrent smoke
# tests do not collide on port or SQLite files.
find_free_port() {
  python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1]); s.close()'
}
PORT=$(find_free_port)
WORKDIR=$(mktemp -d)
CONFIG="${WORKDIR}/config.yml"
trap 'rm -rf "$WORKDIR"; kill "$PID" 2>/dev/null || true' EXIT

# Build and initialize.
MUARA_BINARY="${MUARA_BINARY:-./bin/muara}"
if [[ "$MUARA_BINARY" == "./bin/muara" ]]; then
  go build -o bin/muara ./cmd/muara
fi
"$MUARA_BINARY" --config "$CONFIG" init

# Configure a random server port and point webhooks at the test receiver.
sed -i.bak "s|port: 9000|port: ${PORT}|" "$CONFIG"
sed -i.bak "s|^  url: \"\"|  url: \"http://127.0.0.1:${PORT}/_admin/webhook-receiver\"|" "$CONFIG"
rm -f "${CONFIG}.bak"

# Enable Stripe provider for the smoke test.
python3 - "$CONFIG" <<'PY'
import sys
path = sys.argv[1]
with open(path) as f:
    text = f.read()
text = text.replace('stripe:\n    enabled: false', 'stripe:\n    enabled: true')
text = text.replace('billplz:\n    enabled: false', 'billplz:\n    enabled: true')
text = text.replace('toyyibpay:\n    enabled: false', 'toyyibpay:\n    enabled: true')
text = text.replace('ipay88:\n    enabled: false', 'ipay88:\n    enabled: true')
with open(path, 'w') as f:
    f.write(text)
PY

# Start server in background.
"$MUARA_BINARY" --config "$CONFIG" start &
PID=$!
sleep 2

# Default config values.
MERCHANT_CODE="muara-merchant-code"
SECRET="muara-fawry-secret"
REF="test-ref-$(date +%s)"
USER_ID="user-123"
RETURN_URL="http://127.0.0.1:9999/callback"
ITEM_ID="prod_test_123"
PRICE="99.99"
QUANTITY="1"

# Compute Fawry-style SHA256 signature.
MSG="${MERCHANT_CODE}${REF}${USER_ID}${RETURN_URL}${ITEM_ID}${QUANTITY}${PRICE}${SECRET}"
SIGNATURE=$(printf '%s' "$MSG" | shasum -a 256 | awk '{print $1}')

EXPIRY=$(python3 -c 'import time; print(int((time.time()+600)*1000))')

# Charge request should succeed.
RESPONSE=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/fawry/charge" \
  -H "Content-Type: application/json" \
  -d "{\"merchantCode\":\"${MERCHANT_CODE}\",\"merchantRefNum\":\"${REF}\",\"customerEmail\":\"test@example.com\",\"customerName\":\"Test\",\"customerProfileId\":\"${USER_ID}\",\"paymentExpiry\":${EXPIRY},\"language\":\"ar-eg\",\"chargeItems\":[{\"itemId\":\"${ITEM_ID}\",\"price\":${PRICE},\"quantity\":${QUANTITY}}],\"returnUrl\":\"${RETURN_URL}\",\"signature\":\"${SIGNATURE}\"}")

echo "$RESPONSE" | grep -q '"status":"ok"'
echo "$RESPONSE" | grep -q '"reference":"'"${REF}"'"'

# Health endpoint should respond.
curl -fsS "http://127.0.0.1:${PORT}/healthz" | grep -q '"status":"ok"'

# Fetch the escape page to obtain the CSRF cookie and token.
ESCAPE_URL="http://127.0.0.1:${PORT}/_admin/fawry-escape?ref=${REF}&returnUrl=${RETURN_URL}&amount=${PRICE}"
COOKIES="${WORKDIR}/cookies.txt"
CSRF_TOKEN=$(curl -fsS -c "$COOKIES" "$ESCAPE_URL" | grep -o 'name="csrf_token" value="[^"]*"' | sed 's/.*value="\([^"]*\)".*/\1/' | head -n 1)

# Escape action should redirect. Give the async charge webhook a moment to
# finish writing to SQLite before we mutate the same transaction.
sleep 1
REDIRECT=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/_admin/fawry-escape" \
  -b "$COOKIES" \
  -d "ref=${REF}&returnUrl=${RETURN_URL}&status=PAID&csrf_token=${CSRF_TOKEN}" -i | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')

if [[ "$REDIRECT" != "${RETURN_URL}?orderStatus=PAID&statusCode=200" ]]; then
  echo "Unexpected redirect: $REDIRECT"
  exit 1
fi

# Wait for async webhook delivery and verify it was received.
sleep 1
LIST=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$LIST" | grep -q '"ref":"'"${REF}"'"'
echo "$LIST" | grep -q '"status":"delivered"'

echo "Webhook delivery verified"

# === Fawry versioned webhook endpoints ===
# V1 webhook accepts a signed legacy payload.
V1_SIG=$(python3 - <<PY
import hmac, hashlib
secret = "muara-webhook-secret"
msg = "ref-smoke-v1PAID"
print(hmac.new(secret.encode(), msg.encode(), hashlib.sha256).hexdigest())
PY
)
curl -fsS -X POST "http://127.0.0.1:${PORT}/fawry/v1/webhook?token=muara-webhook-secret" \
  -H "Content-Type: application/json" \
  -d "{\"merchantRefNumber\":\"ref-smoke-v1\",\"orderStatus\":\"PAID\",\"messageSignature\":\"${V1_SIG}\"}" >/dev/null

# V2 webhook rejects a V1-style payload (missing V2 signature).
V2_STATUS=$(curl -o /dev/null -sS -w '%{http_code}' -X POST "http://127.0.0.1:${PORT}/fawry/v2/webhook?token=muara-webhook-secret" \
  -H "Content-Type: application/json" \
  -d '{"merchantRefNumber":"ref-smoke-v2","orderStatus":"PAID"}')
if [[ "$V2_STATUS" != "401" ]]; then
  echo "Expected V2 webhook to reject V1 payload, got HTTP $V2_STATUS"
  exit 1
fi

echo "Fawry v1/v2 webhook endpoints verified"

# === Stripe checkout (card default) + webhook flow ===
STRIPE_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/checkout/sessions" \
  -H "Content-Type: application/json" \
  -d '{"success_url":"http://127.0.0.1:9999/success","cancel_url":"http://127.0.0.1:9999/cancel","mode":"payment","line_items":[{"price_data":{"currency":"usd","unit_amount":999,"product_data":{"name":"Test Product"}},"quantity":1}]}')
SESSION_ID=$(echo "$STRIPE_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')
echo "$STRIPE_RESP" | grep -q '"status":"open"'
echo "$STRIPE_RESP" | grep -q '"payment_method_types":\["card"\]'

# Simulate customer completing checkout.
curl -fsS -X POST "http://127.0.0.1:${PORT}/_admin/stripe/success?session_id=${SESSION_ID}" \
  -b "$COOKIES" \
  -H "X-CSRF-Token: ${CSRF_TOKEN}" | grep -q '"status":"complete"'

# Wait for async delivery and verify the webhook was recorded.
sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${SESSION_ID}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

# Verify Stripe-Signature using the stored attempt.
STRIPE_WEBHOOKS="${WORKDIR}/stripe_webhooks.json"
echo "$WEBHOOKS" > "$STRIPE_WEBHOOKS"
python3 - "$SESSION_ID" "$STRIPE_WEBHOOKS" <<'PY'
import json, hmac, hashlib, sys, base64
session_id, path = sys.argv[1], sys.argv[2]
data = json.load(open(path))
attempts = data.get("results", [])
attempt = next((a for a in attempts if a.get("ref") == session_id), None)
if not attempt:
    sys.exit("stripe webhook attempt not found")
sig = attempt["headers"].get("Stripe-Signature", "")
parts = dict(p.split("=", 1) for p in sig.split(",") if "=" in p)
t = parts.get("t")
v1 = parts.get("v1")
if not t or not v1:
    sys.exit("stripe signature header malformed")
payload = base64.b64decode(attempt["payload"])
secret = b"whsec_muara"
expected = hmac.new(secret, f"{t}.".encode() + payload, hashlib.sha256).hexdigest()
if not hmac.compare_digest(expected, v1):
    print(f"expected {expected}, got {v1}")
    sys.exit("stripe signature mismatch")
print("Stripe signature verified")
PY

echo "Stripe checkout card-default flow verified"

# === Stripe FPX checkout flow ===
FPX_SUCCESS="http://127.0.0.1:9999/fpx/success"
FPX_CANCEL="http://127.0.0.1:9999/fpx/cancel"
FPX_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/checkout/sessions" \
  -H "Content-Type: application/json" \
  -d "{\"success_url\":\"${FPX_SUCCESS}\",\"cancel_url\":\"${FPX_CANCEL}\",\"mode\":\"payment\",\"payment_method_types\":[\"fpx\"],\"line_items\":[{\"price_data\":{\"currency\":\"myr\",\"unit_amount\":5000,\"product_data\":{\"name\":\"FPX Product\"}},\"quantity\":1}]}")
echo "$FPX_RESP" | grep -q '"status":"open"'
FPX_SESSION=$(echo "$FPX_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')

# Fetch checkout page.
curl -fsS "http://127.0.0.1:${PORT}/v1/checkout/sessions/${FPX_SESSION}/pay" | grep -q 'Maybank2U'

# Confirm payment.
FPX_REDIRECT=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/checkout/sessions/${FPX_SESSION}/pay" \
  -b "$COOKIES" \
  -H "X-CSRF-Token: ${CSRF_TOKEN}" \
  -d "action=confirm&bank=maybank2u&csrf_token=${CSRF_TOKEN}" -i | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')

if [[ "$FPX_REDIRECT" != "${FPX_SUCCESS}" ]]; then
  echo "Unexpected FPX redirect: $FPX_REDIRECT"
  exit 1
fi

echo "Stripe FPX checkout flow verified"

# === Stripe card checkout flow ===
CARD_SUCCESS="http://127.0.0.1:9999/card/success"
CARD_CANCEL="http://127.0.0.1:9999/card/cancel"
CARD_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/checkout/sessions" \
  -H "Content-Type: application/json" \
  -d "{\"success_url\":\"${CARD_SUCCESS}\",\"cancel_url\":\"${CARD_CANCEL}\",\"mode\":\"payment\",\"payment_method_types\":[\"card\"],\"line_items\":[{\"price_data\":{\"currency\":\"usd\",\"unit_amount\":2500,\"product_data\":{\"name\":\"Card Product\"}},\"quantity\":1}]}")
echo "$CARD_RESP" | grep -q '"status":"open"'
CARD_SESSION=$(echo "$CARD_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')

# Fetch checkout page.
curl -fsS "http://127.0.0.1:${PORT}/v1/checkout/sessions/${CARD_SESSION}/pay" | grep -q 'Card number'

# Confirm payment.
CARD_REDIRECT=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/checkout/sessions/${CARD_SESSION}/pay" \
  -b "$COOKIES" \
  -H "X-CSRF-Token: ${CSRF_TOKEN}" \
  -d "action=confirm&card_number=4242424242424242&expiry=12/30&cvc=123&csrf_token=${CSRF_TOKEN}" -i | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')

if [[ "$CARD_REDIRECT" != "${CARD_SUCCESS}" ]]; then
  echo "Unexpected card redirect: $CARD_REDIRECT"
  exit 1
fi

echo "Stripe card checkout flow verified"

# === Stripe PaymentIntents FPX flow ===
PI_FPX_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/payment_intents" \
  -H "Content-Type: application/json" \
  -d '{"amount":5000,"currency":"myr","payment_method_types":["fpx"]}')
echo "$PI_FPX_RESP" | grep -q '"status":"requires_confirmation"'
PI_FPX_ID=$(echo "$PI_FPX_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')

PI_FPX_CONFIRM=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/payment_intents/${PI_FPX_ID}/confirm" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"pm_fpx_maybank"}')
echo "$PI_FPX_CONFIRM" | grep -q '"status":"requires_action"'
PI_FPX_AUTH_URL=$(echo "$PI_FPX_CONFIRM" | python3 -c 'import sys,json; print(json.load(sys.stdin)["next_action"]["redirect_to_url"]["url"])')

# Visit local auth page and confirm via admin endpoint.
curl -fsS "${PI_FPX_AUTH_URL}" | grep -q 'Maybank2U'
curl -fsS -X POST "http://127.0.0.1:${PORT}/_admin/stripe/payment_intent/${PI_FPX_ID}" \
  -b "$COOKIES" \
  -H "X-CSRF-Token: ${CSRF_TOKEN}" \
  -d "action=confirm&csrf_token=${CSRF_TOKEN}" | grep -q '"status":"succeeded"'

sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${PI_FPX_ID}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

WEBHOOK_FILE="${WORKDIR}/payment_intent_webhooks.json"
echo "$WEBHOOKS" > "$WEBHOOK_FILE"
python3 - "$PI_FPX_ID" "$WEBHOOK_FILE" <<'PY'
import json, hmac, hashlib, sys, base64
ref, path = sys.argv[1], sys.argv[2]
data = json.load(open(path))
attempts = data.get("results", [])
attempt = next((a for a in attempts if a.get("ref") == ref), None)
if not attempt:
    sys.exit("payment intent webhook attempt not found")
sig = attempt["headers"].get("Stripe-Signature", "")
parts = dict(p.split("=", 1) for p in sig.split(",") if "=" in p)
t = parts.get("t")
v1 = parts.get("v1")
if not t or not v1:
    sys.exit("stripe signature header malformed")
payload = base64.b64decode(attempt["payload"])
secret = b"whsec_muara"
expected = hmac.new(secret, f"{t}.".encode() + payload, hashlib.sha256).hexdigest()
if not hmac.compare_digest(expected, v1):
    print(f"expected {expected}, got {v1}")
    sys.exit("stripe signature mismatch")
print("PaymentIntent FPX signature verified")
PY

echo "Stripe PaymentIntents FPX flow verified"

# === Stripe PaymentIntents card flow ===
PI_CARD_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/payment_intents" \
  -H "Content-Type: application/json" \
  -d '{"amount":2500,"currency":"usd","payment_method_types":["card"]}')
echo "$PI_CARD_RESP" | grep -q '"status":"requires_confirmation"'
PI_CARD_ID=$(echo "$PI_CARD_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["id"])')

curl -fsS -X POST "http://127.0.0.1:${PORT}/v1/payment_intents/${PI_CARD_ID}/confirm" \
  -H "Content-Type: application/json" \
  -d '{"payment_method":"pm_card_visa"}' | grep -q '"status":"succeeded"'

sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${PI_CARD_ID}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

echo "Stripe PaymentIntents card flow verified"

# Invalid signature should return 400.
BAD_RESPONSE=$(curl -sS -X POST "http://127.0.0.1:${PORT}/fawry/charge" \
  -H "Content-Type: application/json" \
  -d "{\"merchantCode\":\"${MERCHANT_CODE}\",\"merchantRefNum\":\"${REF}-bad\",\"customerProfileId\":\"${USER_ID}\",\"returnUrl\":\"${RETURN_URL}\",\"chargeItems\":[{\"itemId\":\"${ITEM_ID}\",\"price\":${PRICE},\"quantity\":${QUANTITY}}],\"signature\":\"invalid\"}")

echo "$BAD_RESPONSE" | grep -q '"code":"OPENMUARA_INVALID_SIGNATURE"'

# === Billplz collection + bill + pay flow ===
BILLPLZ_API_KEY="muara-billplz-api-key"
BILLPLZ_CALLBACK="http://127.0.0.1:${PORT}/_admin/webhook-receiver"
BILLPLZ_REDIRECT="http://127.0.0.1:9999/billplz/return"

COLL_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/api/v3/collections" \
  -u "${BILLPLZ_API_KEY}:" \
  -H "Content-Type: application/json" \
  -d '{"title":"Smoke Collection"}')
COLL_ID=$(echo "$COLL_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["collection"]["id"])')

BILL_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/api/v3/bills" \
  -u "${BILLPLZ_API_KEY}:" \
  -H "Content-Type: application/json" \
  -d "{\"collection_id\":\"${COLL_ID}\",\"email\":\"test@example.com\",\"name\":\"Test\",\"amount\":1250,\"callback_url\":\"${BILLPLZ_CALLBACK}\",\"description\":\"Smoke bill\",\"redirect_url\":\"${BILLPLZ_REDIRECT}\"}")
BILL_ID=$(echo "$BILL_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["bill"]["id"])')
echo "$BILL_RESP" | grep -q '"state":"due"'

curl -fsS -c "$COOKIES" -b "$COOKIES" "http://127.0.0.1:${PORT}/_admin/billplz/pay/${BILL_ID}" >/dev/null
BILLPLZ_REDIRECT_LOCATION=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/_admin/billplz/pay/${BILL_ID}" \
  -b "$COOKIES" \
  -d "outcome=pay&method=fpx&csrf_token=${CSRF_TOKEN}" -i | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')

BILLPLZ_STATE=$(python3 -c 'import sys, urllib.parse; print(urllib.parse.parse_qs(urllib.parse.urlparse(sys.argv[1]).query).get("billplz[state]", [""])[0])' "$BILLPLZ_REDIRECT_LOCATION")
BILLPLZ_HAS_SIG=$(python3 -c 'import sys, urllib.parse; print("x_signature" in urllib.parse.parse_qs(urllib.parse.urlparse(sys.argv[1]).query))' "$BILLPLZ_REDIRECT_LOCATION")
if [[ "$BILLPLZ_HAS_SIG" != "True" ]]; then
  echo "Billplz redirect missing x_signature: $BILLPLZ_REDIRECT_LOCATION"
  exit 1
fi
if [[ "$BILLPLZ_STATE" != "paid" ]]; then
  echo "Billplz redirect state not paid: $BILLPLZ_REDIRECT_LOCATION"
  exit 1
fi

sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${BILL_ID}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

echo "Billplz flow verified"

# === ToyyibPay category + bill + pay flow ===
TOYYIB_SECRET="muara-toyyibpay-secret"
TOYYIB_REF="toyyib-ref-$(date +%s)"
TOYYIB_RETURN="http://127.0.0.1:9999/toyyib/return"
TOYYIB_CALLBACK="http://127.0.0.1:${PORT}/_admin/webhook-receiver"

CAT_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/index.php/api/createCategory" \
  --data-urlencode "userSecretKey=${TOYYIB_SECRET}" \
  --data-urlencode "categoryName=Smoke Category")
CAT_CODE=$(echo "$CAT_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["data"]["categoryCode"])')

BILL_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/index.php/api/createBill" \
  --data-urlencode "userSecretKey=${TOYYIB_SECRET}" \
  --data-urlencode "categoryCode=${CAT_CODE}" \
  --data-urlencode "billName=Smoke Bill" \
  --data-urlencode "billAmount=1250" \
  --data-urlencode "billReturnUrl=${TOYYIB_RETURN}" \
  --data-urlencode "billCallbackUrl=${TOYYIB_CALLBACK}" \
  --data-urlencode "billExternalReferenceNo=${TOYYIB_REF}" \
  --data-urlencode "billTo=Test" \
  --data-urlencode "billEmail=test@example.com" \
  --data-urlencode "billPhone=0123456789")
BILL_CODE=$(echo "$BILL_RESP" | python3 -c 'import sys,json; print(json.load(sys.stdin)["bill"]["billCode"])')

curl -fsS -c "$COOKIES" -b "$COOKIES" "http://127.0.0.1:${PORT}/_admin/toyyibpay/pay/${BILL_CODE}" >/dev/null
TOYYIB_REDIRECT=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/_admin/toyyibpay/pay/${BILL_CODE}" \
  -b "$COOKIES" \
  -d "status=1&payment_method=fpx&csrf_token=${CSRF_TOKEN}" -i | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')

if ! echo "$TOYYIB_REDIRECT" | grep -q '^/toyyibpay/return'; then
  echo "Unexpected ToyyibPay first redirect: $TOYYIB_REDIRECT"
  exit 1
fi
if ! echo "$TOYYIB_REDIRECT" | grep -q 'status_id=1'; then
  echo "ToyyibPay first redirect missing status_id=1: $TOYYIB_REDIRECT"
  exit 1
fi

TOYYIB_FINAL_REDIRECT=$(curl -fsS -i "http://127.0.0.1:${PORT}${TOYYIB_REDIRECT}" | grep -i '^Location:' | awk '{print $2}' | tr -d '\r')
if [[ "$TOYYIB_FINAL_REDIRECT" != "${TOYYIB_RETURN}"* ]]; then
  echo "Unexpected ToyyibPay final redirect: $TOYYIB_FINAL_REDIRECT"
  exit 1
fi
if ! echo "$TOYYIB_FINAL_REDIRECT" | grep -q 'status_id=1'; then
  echo "ToyyibPay final redirect missing status_id=1: $TOYYIB_FINAL_REDIRECT"
  exit 1
fi

sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${TOYYIB_REF}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

echo "ToyyibPay flow verified"

# === iPay88 entry + pay + backend + requery flow ===
IPAY88_KEY="muara-ipay88-key"
IPAY88_CODE="muara-ipay88-merchant"
IPAY88_REF="ipay88-ref-$(date +%s)"
IPAY88_AMOUNT="12.50"
IPAY88_CURRENCY="MYR"
IPAY88_PAYMENTID="2"
IPAY88_RESPONSE_URL="https://example.com/response"
IPAY88_BACKEND_URL="https://example.com/backend"

AMOUNT_STRIPPED=$(echo "$IPAY88_AMOUNT" | tr -d '.,')
ENTRY_SIG=$(printf '%s' "${IPAY88_KEY}${IPAY88_CODE}${IPAY88_REF}${AMOUNT_STRIPPED}${IPAY88_CURRENCY}" | shasum -a 256 | awk '{print $1}')

curl -fsS -X POST "http://127.0.0.1:${PORT}/ePayment/entry.asp" \
  --data-urlencode "MerchantCode=${IPAY88_CODE}" \
  --data-urlencode "PaymentId=${IPAY88_PAYMENTID}" \
  --data-urlencode "RefNo=${IPAY88_REF}" \
  --data-urlencode "Amount=${IPAY88_AMOUNT}" \
  --data-urlencode "Currency=${IPAY88_CURRENCY}" \
  --data-urlencode "ProdDesc=Smoke test" \
  --data-urlencode "UserName=Test" \
  --data-urlencode "UserEmail=test@example.com" \
  --data-urlencode "UserContact=0123456789" \
  --data-urlencode "Signature=${ENTRY_SIG}" \
  --data-urlencode "SignatureType=SHA256" \
  --data-urlencode "ResponseURL=${IPAY88_RESPONSE_URL}" \
  --data-urlencode "BackendURL=${IPAY88_BACKEND_URL}" >/dev/null

curl -fsS -c "$COOKIES" -b "$COOKIES" "http://127.0.0.1:${PORT}/_admin/ipay88/pay/${IPAY88_REF}" >/dev/null
curl -sS -X POST "http://127.0.0.1:${PORT}/_admin/ipay88/pay/${IPAY88_REF}" \
  -b "$COOKIES" \
  -d "outcome=pay&payment_method=${IPAY88_PAYMENTID}&csrf_token=${CSRF_TOKEN}" >/dev/null

BACKEND_SIG=$(printf '%s' "${IPAY88_KEY}${IPAY88_CODE}${IPAY88_PAYMENTID}${IPAY88_REF}${AMOUNT_STRIPPED}${IPAY88_CURRENCY}1" | shasum -a 256 | awk '{print $1}')
BACKEND_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/ipay88/backend" \
  --data-urlencode "MerchantCode=${IPAY88_CODE}" \
  --data-urlencode "PaymentId=${IPAY88_PAYMENTID}" \
  --data-urlencode "RefNo=${IPAY88_REF}" \
  --data-urlencode "Amount=${IPAY88_AMOUNT}" \
  --data-urlencode "Currency=${IPAY88_CURRENCY}" \
  --data-urlencode "Status=1" \
  --data-urlencode "Signature=${BACKEND_SIG}" \
  --data-urlencode "SignatureType=SHA256")

if [[ "$BACKEND_RESP" != "RECEIVEOK" ]]; then
  echo "Unexpected iPay88 backend response: $BACKEND_RESP"
  exit 1
fi

sleep 1
WEBHOOKS=$(curl -fsS "http://127.0.0.1:${PORT}/_admin/webhooks")
echo "$WEBHOOKS" | grep -q "\"ref\":\"${IPAY88_REF}\""
echo "$WEBHOOKS" | grep -q '"status":"delivered"'

REQUERY_RESP=$(curl -fsS -X POST "http://127.0.0.1:${PORT}/ePayment/enquiry.asp" \
  --data-urlencode "MerchantCode=${IPAY88_CODE}" \
  --data-urlencode "RefNo=${IPAY88_REF}" \
  --data-urlencode "Amount=${IPAY88_AMOUNT}")

if [[ "$REQUERY_RESP" != "00" ]]; then
  echo "Unexpected iPay88 requery response: $REQUERY_RESP"
  exit 1
fi

echo "iPay88 flow verified"

echo "Smoke test passed"
