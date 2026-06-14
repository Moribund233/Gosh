import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 5 },   // Ramp up to 5 VUs
    { duration: '30s', target: 5 },   // Stay at 5 VUs
    { duration: '10s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:9292/api/v1';

export default function () {
  const base = BASE_URL.replace('/api/v1', '');
  const healthRes = http.get(base + '/health');
  check(healthRes, { 'health check ok': (r) => r.status === 200 });

  const bannerRes = http.get(`${BASE_URL}/banners`);
  check(bannerRes, { 'banners ok': (r) => r.status === 200 });

  const categoriesRes = http.get(`${BASE_URL}/categories`);
  check(categoriesRes, { 'categories ok': (r) => r.status === 200 });

  const productsRes = http.get(`${BASE_URL}/products`);
  check(productsRes, { 'products ok': (r) => r.status === 200 });

  const paymentMethodsRes = http.get(`${BASE_URL}/payment/methods`);
  check(paymentMethodsRes, { 'payment methods ok': (r) => r.status === 200 });

  // Register + login
  const phone = `138${String(Date.now()).slice(-8)}`;
  const registerRes = http.post(`${BASE_URL}/user/register`, JSON.stringify({
    phone, password: 'test1234', nickname: 'k6_user',
  }), { headers: { 'Content-Type': 'application/json' } });
  check(registerRes, { 'register ok': (r) => r.status === 201 });

  const token = registerRes.json().data?.token;

  if (token) {
    const authHeaders = {
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    };

    const profileRes = http.get(`${BASE_URL}/user/profile`, authHeaders);
    check(profileRes, { 'profile ok': (r) => r.status === 200 });

    const pointsRes = http.get(`${BASE_URL}/points`, authHeaders);
    check(pointsRes, { 'points ok': (r) => r.status === 200 });

    const cartRes = http.get(`${BASE_URL}/cart`, authHeaders);
    check(cartRes, { 'cart ok': (r) => r.status === 200 });

    const ordersRes = http.get(`${BASE_URL}/orders`, authHeaders);
    check(ordersRes, { 'orders ok': (r) => r.status === 200 });

    const favoritesRes = http.get(`${BASE_URL}/favorites`, authHeaders);
    check(favoritesRes, { 'favorites ok': (r) => r.status === 200 });

    // Create address
    const addrRes = http.post(`${BASE_URL}/addresses`, JSON.stringify({
      name: 'k6', phone: '13800138000',
      province: '广东', city: '深圳', district: '南山', detail: '测试地址',
      is_default: true,
    }), authHeaders);
    check(addrRes, { 'address created': (r) => r.status === 201 });
  }

  sleep(1);
}
