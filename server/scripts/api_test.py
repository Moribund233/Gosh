#!/usr/bin/env python3
"""
Gosh Mall API 端到端测试脚本

用法:
  pip install requests
  python scripts/api_test.py

环境变量:
  BASE_URL (默认: http://localhost:8080/api/v1)
"""

import os
import sys
import time
import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:8080/api/v1")

def log(step, status, detail=""):
    icon = "✅" if status == "PASS" else "❌"
    print(f"{icon} [{status}] {step}  {detail}")

def check(label, ok, detail=""):
    if ok:
        log(label, "PASS", detail)
    else:
        log(label, "FAIL", detail)
        sys.exit(1)

def main():
    session = requests.Session()

    # === 1. Health Check ===
    r = session.get(f"{BASE_URL.replace('/api/v1', '')}/health")
    check("Health check", r.status_code == 200)

    # === 2. Register ===
    phone = f"138{int(time.time()) % 10000000000:010d}"[-11:]
    r = session.post(f"{BASE_URL}/user/register", json={
        "phone": phone,
        "password": "test1234",
        "nickname": "测试用户"
    })
    check("User register", r.status_code == 201, f"phone={phone}")
    token = r.json().get("data", {}).get("token")
    check("Token returned", bool(token))

    session.headers["Authorization"] = f"Bearer {token}"

    # === 3. Get Profile ===
    r = session.get(f"{BASE_URL}/user/profile")
    check("Get profile", r.status_code == 200)
    profile_data = r.json()["data"]
    nickname = profile_data.get("user", {}).get("nickname", "") or profile_data.get("nickname", "")
    check("Profile has nickname", nickname == "测试用户")

    # === 4. Create Address ===
    r = session.post(f"{BASE_URL}/addresses", json={
        "name": "张三", "phone": "13800138000",
        "province": "广东省", "city": "深圳市",
        "district": "南山区", "detail": "科技园南区",
        "is_default": True
    })
    check("Create address", r.status_code == 201)

    # === 5. Get Categories ===
    r = session.get(f"{BASE_URL}/categories")
    check("Get categories", r.status_code == 200)

    # === 6. Get Products (empty) ===
    r = session.get(f"{BASE_URL}/products")
    check("Get products (empty)", r.status_code == 200)
    check("Empty list returned", len(r.json()["data"]["list"]) == 0)

    # === 7. Search Products (empty) ===
    r = session.get(f"{BASE_URL}/products/search?keyword=test")
    check("Search products", r.status_code == 200)

    # === 8. Get Banners ===
    r = session.get(f"{BASE_URL}/banners")
    check("Get banners", r.status_code == 200)

    # === 9. Get Payment Methods ===
    r = session.get(f"{BASE_URL}/payment/methods")
    check("Payment methods", r.status_code == 200)
    check("Has 3 methods", len(r.json()["data"]) == 3)

    # === 10. Get Active Flash Sales ===
    r = session.get(f"{BASE_URL}/flash-sales")
    check("Flash sales", r.status_code == 200)

    # === 11. Points ===
    r = session.get(f"{BASE_URL}/points")
    check("Query points", r.status_code == 200)
    check("Initial points = 0", r.json()["data"]["points"] == 0)

    r = session.get(f"{BASE_URL}/points/logs")
    check("Empty point logs", r.status_code == 200)

    # === 12. Favorites (empty) ===
    r = session.get(f"{BASE_URL}/favorites")
    check("Empty favorites", r.status_code == 200)

    # === 13. Browse History (empty) ===
    r = session.get(f"{BASE_URL}/browse-history")
    check("Empty browse history", r.status_code == 200)

    # === 14. Cart (empty) ===
    r = session.get(f"{BASE_URL}/cart")
    check("Empty cart", r.status_code == 200)
    r = session.get(f"{BASE_URL}/cart/count")
    check("Cart count 0", r.status_code == 200)

    # === 15. Orders (empty) ===
    r = session.get(f"{BASE_URL}/orders")
    check("Empty orders", r.status_code == 200)

    # === 16. Available Coupons ===
    r = session.get(f"{BASE_URL}/coupons/available?amount=10000")
    check("Available coupons", r.status_code == 200)

    # === 17. Logout / Unauthorized ===
    session.headers.pop("Authorization")
    r = session.get(f"{BASE_URL}/user/profile")
    check("Unauthorized access rejected", r.status_code == 401)

    print(f"\n{'='*50}")
    print(f"✅ All smoke tests passed!")
    print(f"{'='*50}")

if __name__ == "__main__":
    main()
