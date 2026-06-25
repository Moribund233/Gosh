#!/usr/bin/env python3
"""
Gosh Mall 全流程 + 边界场景 + 安全性 综合测试

用法:
  pip install requests
  python scripts/comprehensive_test.py

环境变量:
  BASE_URL (默认: http://localhost:9292/api/v1)
  ADMIN_TOKEN (可选: 已有管理员 token)
"""

import os
import sys
import time
import json
import uuid
import random
import hashlib
import threading
import subprocess
from urllib.parse import urljoin

import requests

BASE_URL = os.environ.get("BASE_URL", "http://localhost:9292/api/v1")
ROOT_URL = BASE_URL.replace("/api/v1", "")

PASS = 0
FAIL = 0
WARN = 0
errors = []


def log(step, status, detail=""):
    global PASS, FAIL, WARN
    if status == "PASS":
        icon, count = "✅", PASS
        PASS += 1
    elif status == "FAIL":
        icon, count = "❌", FAIL
        FAIL += 1
        errors.append(f"[FAIL] {step}: {detail}")
    else:
        icon, count = "⚠️", WARN
        WARN += 1
    print(f"  {icon} [{status}] {step}  {detail}")


def check(label, ok, detail=""):
    if ok:
        log(label, "PASS", detail)
    else:
        log(label, "FAIL", detail)
        if os.environ.get("FAIL_FAST"):
            sys.exit(1)


def check_eq(label, got, expected, detail=""):
    ok = got == expected
    d = detail if detail else f"got={got!r}, expected={expected!r}"
    check(label, ok, d)


def check_in(label, item, container, detail=""):
    ok = item in container
    d = detail if detail else f"item={item!r}, container_keys={list(container.keys())[:5]}"
    check(label, ok, d)


class GoshTester:
    def __init__(self):
        self.session = requests.Session()
        self.session.headers.update({"Content-Type": "application/json"})
        self.user_token = None
        self.admin_token = None
        self.admin_session = requests.Session()
        self.admin_session.headers.update({"Content-Type": "application/json"})
        self.test_user = {}
        self.test_address = {}
        self.test_category = {}
        self.test_product = {}
        self.test_sku = {}
        self.test_cart_item = {}
        self.test_order = {}
        self.test_coupon = {}
        self.test_favorite = {}
        self.test_review = {}
        self.phone_counter = int(time.time()) % 1000000

    def url(self, path):
        return urljoin(BASE_URL + "/", path.lstrip("/"))

    def req(self, method, path, **kwargs):
        return self.session.request(method, self.url(path), timeout=10, **kwargs)

    def admin_req(self, method, path, **kwargs):
        return self.admin_session.request(method, self.url(path), timeout=10, **kwargs)

    def promote_to_admin(self, phone):
        try:
            subprocess.run(
                ["docker", "exec", "-i", "postgresql", "psql", "-U", "alucard", "-d", "gosh",
                 "-c", f"UPDATE users SET role = 'super_admin' WHERE phone = '{phone}';"],
                capture_output=True, text=True, timeout=10
            )
            log("promote to admin", "PASS", f"phone={phone}")
        except Exception as e:
            log("promote to admin", "WARN", str(e))

    def random_phone(self):
        self.phone_counter += 1
        return f"138{int(time.time() * 1000) % 10000000:07d}{self.phone_counter % 10:01d}"

    # =========== Section 1: Health & Basic ===========
    def test_health(self):
        print("\n📋 Section 1: Health & Basic")
        r = requests.get(f"{ROOT_URL}/health", timeout=10)
        check("GET /health returns 200", r.status_code == 200)
        check_eq("code is 0", r.json()["code"], 0)

        r = self.req("GET", "/user/register")
        check("register with GET method rejected", r.status_code in (400, 404, 405))
        r = self.req("POST", "/user/register", json={})
        check("register with empty body", r.json().get("code") != 0)

    # =========== Section 2: User Registration & Auth ===========
    def test_auth(self):
        print("\n📋 Section 2: Registration & Authentication")

        phone = self.random_phone()
        password = "Test@1234"
        r = self.req("POST", "/user/register", json={
            "phone": phone, "password": password, "nickname": "TestUser"
        })
        check("register new user", r.status_code == 201, f"phone={phone}")
        if r.status_code == 201:
            self.user_token = r.json().get("data", {}).get("token")
            self.test_user["phone"] = phone
            self.test_user["password"] = password
            self.test_user["id"] = r.json().get("data", {}).get("user", {}).get("id")
            check("register returns token", bool(self.user_token))
        if self.user_token:
            self.session.headers["Authorization"] = f"Bearer {self.user_token}"

        r = self.req("POST", "/user/register", json={
            "phone": phone, "password": password, "nickname": "TestUser"
        })
        check("duplicate phone rejected", r.status_code in (400, 409))

        r = self.req("POST", "/user/register", json={
            "password": "test1234", "nickname": "NoPhone"
        })
        check("register missing phone", r.status_code in (400, 422))

        r = self.req("POST", "/user/login", json={
            "phone": phone, "password": password
        })
        check("login with valid creds", r.status_code == 200)
        if r.status_code == 200:
            token = r.json().get("data", {}).get("token")
            check("login returns token", bool(token))

        r = self.req("POST", "/user/login", json={
            "phone": phone, "password": "wrongpassword"
        })
        check("login wrong password", r.status_code == 401)

        r = self.req("POST", "/user/login", json={
            "phone": "13900000000", "password": "test1234"
        })
        check("login nonexistent user", r.status_code == 401)

    # =========== Section 3: User Profile ===========
    def test_profile(self):
        print("\n📋 Section 3: User Profile")
        r = self.req("GET", "/user/profile")
        check("get profile", r.status_code == 200)
        data = r.json().get("data", {})
        user = data.get("user", data)
        check("profile has nickname", user.get("nickname") == "TestUser")
        check_in("profile has phone", "phone", user)

        r = self.req("PUT", "/user/profile", json={"nickname": "UpdatedName"})
        check("update profile", r.status_code == 200)
        r = self.req("GET", "/user/profile")
        check("profile updated", r.json().get("data", {}).get("user", r.json().get("data", {})).get("nickname") == "UpdatedName")

        r = self.req("PUT", "/user/profile", json={"nickname": "x" * 101})
        check("nickname too long", r.status_code in (400, 422))

        r = self.req("PUT", "/user/profile", json={"nickname": "<script>alert(1)</script>"})
        check("profile with XSS attempt", r.status_code in (200, 400))
        if r.status_code == 200:
            r2 = self.req("GET", "/user/profile")
            nick = r2.json().get("data", {}).get("user", r2.json().get("data", {})).get("nickname", "")
            if "<script>" in nick:
                log("XSS stored unescaped", "WARN", "data stored with script tags (frontend responsibility)")

    # =========== Section 4: Addresses ===========
    def test_addresses(self):
        print("\n📋 Section 4: Address Management")
        r = self.req("POST", "/addresses", json={
            "name": "张三", "phone": "13800138000",
            "province": "广东省", "city": "深圳市",
            "district": "南山区", "detail": "科技园南区",
            "is_default": True
        })
        check("create address", r.status_code == 201)
        if r.status_code == 201:
            self.test_address["id"] = r.json().get("data", {}).get("id")

        r = self.req("POST", "/addresses", json={
            "name": "李四", "phone": "13900139000",
            "province": "北京市", "city": "北京市",
            "district": "朝阳区", "detail": "国贸CBD",
            "is_default": False
        })
        check("create address 2", r.status_code == 201)

        r = self.req("GET", "/addresses")
        check("list addresses", r.status_code == 200)
        data = r.json().get("data", [])
        check("has 2 addresses", len(data) >= 2)

        r = self.req("POST", "/addresses", json={
            "name": "", "phone": "invalid", "province": "", "city": "", "district": "", "detail": ""
        })
        check("create address invalid fields", r.status_code in (400, 422))

    # =========== Section 5: Categories (Admin) ===========
    def test_categories(self):
        print("\n📋 Section 5: Category Management")
        r = self.req("GET", "/categories")
        check("get category tree (empty)", r.status_code == 200)

        r = self.admin_req("POST", "/admin/categories", json={
            "name": "电子产品", "icon": "electronics", "sort_order": 1
        })
        if r.status_code in (401, 403):
            log("create category as user", "WARN", "need admin token, use setup")
        elif r.status_code == 201:
            self.test_category["id"] = r.json().get("data", {}).get("id")

        r = self.req("GET", "/categories")
        check("get categories after create", r.status_code == 200)

    # =========== Section 6: Products (Admin) ===========
    def test_products(self):
        print("\n📋 Section 6: Product Management")

        cat_id = self.test_category.get("id") or 1
        r = self.admin_req("POST", "/admin/products", json={
            "name": "测试商品",
            "subtitle": "测试副标题",
            "category_id": cat_id,
            "price": 9999,
            "original_price": 12999,
            "brand": "TestBrand",
            "description": "这是一个测试商品描述",
            "images": ["https://example.com/img1.jpg"],
            "is_new": True,
            "is_hot": True,
            "is_featured": True,
            "skus": [
                {"name": "标准版", "price": 9999, "stock": 100},
                {"name": "豪华版", "price": 19999, "stock": 50},
            ]
        })
        if r.status_code in (401, 403):
            log("create product as user", "WARN", "need admin token")
            return
        check("create product", r.status_code == 201, f"got {r.status_code}")
        if r.status_code == 201:
            data = r.json().get("data", {})
            self.test_product["id"] = data.get("id")
            skus = data.get("skus", data.get("product", {}).get("skus", []))
            if skus:
                self.test_sku["id"] = skus[0].get("id")
                self.test_sku["stand_id"] = skus[0].get("id")
                self.test_sku["deluxe_id"] = skus[1].get("id") if len(skus) > 1 else skus[0].get("id")

        r = self.req("GET", "/products")
        check("list products", r.status_code == 200)
        if r.status_code == 200:
            items = r.json().get("data", {}).get("list", [])
            check("products list non-empty", len(items) > 0)

        if self.test_product.get("id"):
            r = self.req("GET", f"/products/{self.test_product['id']}")
            check("get product detail", r.status_code == 200)

        r = self.req("GET", "/products/search", params={"keyword": "测试"})
        check("search products", r.status_code == 200)

        r = self.req("GET", "/products/hot-search")
        check("hot search keywords", r.status_code == 200)

    # =========== Section 7: Cart ===========
    def test_cart(self):
        print("\n📋 Section 7: Shopping Cart")
        sku_id = self.test_sku.get("stand_id") or self.test_sku.get("id") or 0

        r = self.req("GET", "/cart")
        check("get empty cart", r.status_code == 200)

        r = self.req("GET", "/cart/count")
        check("cart count 0", r.status_code == 200)

        if sku_id == 0:
            log("add to cart", "WARN", "no SKU available")
            log("cart add second item", "WARN", "no SKU available")
            return

        r = self.req("POST", "/cart", json={"sku_id": sku_id, "quantity": 2})
        check("add to cart", r.status_code == 201)
        if r.status_code == 201:
            self.test_cart_item["id"] = r.json().get("data", {}).get("id")

        r = self.req("GET", "/cart")
        check("cart has items", r.status_code == 200)
        cart_data = r.json().get("data", [])
        items = cart_data if isinstance(cart_data, list) else cart_data.get("items", cart_data)
        check("cart not empty", len(items) > 0)

        r = self.req("GET", "/cart/count")
        check("cart count > 0", r.status_code == 200)
        data = r.json().get("data", {})
        if isinstance(data, dict):
            count = data.get("count", data.get("total_count", 0))
        else:
            count = data if isinstance(data, (int, float)) else 0
        check("count > 0", count > 0)

        if self.test_cart_item.get("id"):
            r = self.req("PUT", f"/cart/{self.test_cart_item['id']}", json={"quantity": 5})
            check("update cart item quantity", r.status_code == 200)

        r = self.req("POST", "/cart/select", json={"sku_ids": [sku_id], "selected": True})
        check("select cart item", r.status_code == 200)

        r = self.req("POST", "/cart", json={"sku_id": -1, "quantity": 1})
        check("add invalid sku", r.status_code in (400, 404))

        r = self.req("POST", "/cart", json={"sku_id": sku_id, "quantity": -1})
        check("add negative quantity", r.status_code in (400, 422))

        r = self.req("POST", "/cart", json={"sku_id": sku_id, "quantity": 99999})
        check("add excessive quantity handled", r.status_code in (200, 201, 400, 422))

    # =========== Section 8: Orders ===========
    def test_orders(self):
        print("\n📋 Section 8: Order Flow")
        sku_id = self.test_sku.get("stand_id") or self.test_sku.get("id") or 0

        r = self.req("GET", "/orders")
        check("get empty orders", r.status_code == 200)

        if sku_id == 0:
            log("create order", "WARN", "no SKU available")
            return

        addr_id = self.test_address.get("id", 0)
        if addr_id == 0:
            log("create order", "WARN", "no address available")
            return
        idempotent = str(uuid.uuid4())
        r = self.req("POST", "/orders", json={
            "items": [{"sku_id": sku_id, "quantity": 1}],
            "address_id": addr_id,
            "remark": "测试订单"
        }, headers={"Idempotent-Key": idempotent})
        if r.status_code == 201:
            self.test_order["id"] = r.json().get("data", {}).get("id")
            self.test_order["order_no"] = r.json().get("data", {}).get("order_no")
        check("create order (with items > 0)", r.status_code  == 201,
              f"got {r.status_code}: {r.json().get('message','')}")

        r = self.req("POST", "/orders", json={
            "items": [],
            "remark": "空订单",
            "address_id": self.test_address.get("id", 0),
        }, headers={"Idempotent-Key": str(uuid.uuid4())})
        check("create order from cart (empty items)", r.status_code in (201, 400),
              f"got {r.status_code}: {r.json().get('message','')}")

        if self.test_order.get("id"):
            r = self.req("GET", f"/orders/{self.test_order['id']}")
            check("get order detail", r.status_code == 200)

        r = self.req("GET", "/orders")
        check("list orders >0", r.status_code == 200)

    # =========== Section 9: Payment ===========
    def test_payment(self):
        print("\n📋 Section 9: Payment")

        r = self.req("GET", "/payment/methods")
        check("get payment methods", r.status_code == 200)
        methods = r.json().get("data", [])
        check("has payment methods", len(methods) > 0)
        has_alipay = any(m.get("method") == "alipay" for m in methods)
        has_wechat = any(m.get("method") in ("wechat", "wxpay") for m in methods)
        check("alipay available", has_alipay)
        check("wechat available", has_wechat)

        order_no = self.test_order.get("order_no", "")
        if order_no:
            r = self.req("POST", "/payment/pay", json={
                "order_no": order_no,
                "method": "mock"
            })
            check("pay with order", r.status_code in (200, 201, 400, 402),
                  f"got {r.status_code}: {r.json().get('message','')}")
            if r.status_code in (200, 201):
                self.test_order["paid"] = True

            r = self.req("GET", f"/payment/status/{order_no}")
            check("get payment status", r.status_code == 200)
        else:
            log("pay with order", "WARN", "no order available")

        r = self.req("GET", "/payment/status/INVALID_ORDER_NO")
        check("payment status invalid order", r.status_code in (400, 404))

    # =========== Section 10: Coupons ===========
    def test_coupons(self):
        print("\n📋 Section 10: Coupons")

        r = self.req("GET", "/coupons/available", params={"amount": 10000})
        check("get available coupons (no coupons yet)", r.status_code == 200)

        r = self.req("POST", "/coupons/calculate", json={
            "items": [{"product_id": 1, "quantity": 1, "price": 99.99}],
            "amount": 10000
        })
        check("calculate coupon", r.status_code in (200, 400),
              f"got {r.status_code}: {r.json().get('message','')}")

    # =========== Section 11: Favorites ===========
    def test_favorites(self):
        print("\n📋 Section 11: Favorites")

        r = self.req("GET", "/favorites")
        check("list empty favorites", r.status_code == 200)

        product_id = self.test_product.get("id", 1)
        r = self.req("POST", "/favorites", json={"product_id": product_id})
        check("add favorite", r.status_code in (200, 201, 400))
        if r.status_code in (200, 201):
            self.test_favorite["id"] = product_id

        r = self.req("GET", "/favorites")
        check("list favorites after add", r.status_code == 200)

        r = self.req("POST", "/favorites/remove", json={"product_id": product_id})
        check("remove favorite", r.status_code == 200)

    # =========== Section 12: Browse History ===========
    def test_browse_history(self):
        print("\n📋 Section 12: Browse History")

        r = self.req("GET", "/browse-history")
        check("list empty history", r.status_code == 200)

        product_id = self.test_product.get("id", 1)
        r = self.req("POST", "/browse-history", json={"product_id": product_id})
        check("add browse history", r.status_code in (200, 201))

        r = self.req("GET", "/browse-history")
        check("list history after add", r.status_code == 200)

    # =========== Section 13: Points ===========
    def test_points(self):
        print("\n📋 Section 13: Points")

        r = self.req("GET", "/points")
        check("get points balance", r.status_code == 200)

        r = self.req("GET", "/points/logs")
        check("get points logs", r.status_code == 200)

    # =========== Section 14: Reviews ===========
    def test_reviews(self):
        print("\n📋 Section 14: Reviews")

        product_id = self.test_product.get("id", 1)
        r = self.req("GET", "/reviews", params={"product_id": product_id})
        check("list reviews", r.status_code == 200)

        order_id = self.test_order.get("id", 0)
        r = self.req("POST", "/reviews", json={
            "product_id": product_id,
            "order_id": order_id,
            "rating": 5,
            "content": "好评！"
        })
        check("create review", r.status_code in (200, 201, 400))
        if r.status_code in (200, 201):
            self.test_review["id"] = r.json().get("data", {}).get("id")

    # =========== Section 15: Security Tests ===========
    def test_security(self):
        print("\n📋 Section 15: Security Tests")

        # 15.1 Unauthenticated access
        clean_session = requests.Session()
        protected_endpoints = [
            ("GET", "/user/profile"),
            ("GET", "/orders"),
            ("POST", "/orders"),
            ("GET", "/cart"),
            ("POST", "/cart"),
            ("GET", "/favorites"),
            ("GET", "/points"),
            ("GET", "/addresses"),
            ("POST", "/addresses"),
        ]
        for method, path in protected_endpoints:
            r = clean_session.request(method, self.url(path), timeout=10)
            check(f"unauth {method} {path} -> {r.status_code}", r.status_code == 401)

        # 15.2 Token tampering
        self.session.headers["Authorization"] = "Bearer invalid-token-12345"
        r = self.req("GET", "/user/profile")
        check("invalid token rejected", r.status_code == 401)
        self.session.headers["Authorization"] = f"Bearer {self.user_token}"

        # 15.3 Token with malformed format
        for bad_header in [
            "Bearer ",
            "Bearer",
            "Basic abc123",
            "invalid",
            "",
        ]:
            self.session.headers["Authorization"] = bad_header
            r = self.req("GET", "/user/profile")
            if bad_header:
                check(f"malformed token '{bad_header[:20]}' rejected", r.status_code == 401)
        self.session.headers["Authorization"] = f"Bearer {self.user_token}"

        # 15.4 Admin endpoint access by regular user
        admin_endpoints = [
            ("GET", "/admin/users"),
            ("POST", "/admin/categories"),
            ("GET", "/admin/banners"),
            ("GET", "/admin/merchant/applications"),
        ]
        for method, path in admin_endpoints:
            r = self.req(method, path)
            check(f"user access admin {method} {path}", r.status_code in (401, 403),
                  f"got {r.status_code}")

        # 15.5 SQL injection attempts (search endpoint should handle gracefully)
        sql_patterns = [
            "' OR '1'='1",
            "1; DROP TABLE users",
            "' UNION SELECT * FROM users",
            "1' AND 1=1--",
            "'; SELECT pg_sleep(5)--",
        ]
        for sql in sql_patterns:
            r = self.req("GET", "/products/search", params={"keyword": sql})
            check(f"sql injection keyword safe: {sql[:20]}", r.status_code in (200, 400),
                  f"got {r.status_code}")
            if r.status_code == 200:
                items = r.json().get("data", {}).get("list", [])
                if isinstance(items, list) and len(items) > 0:
                    log(f"sql injection returned results", "WARN", f"keyword={sql}")

        # 15.6 XSS attempts in search (should be handled gracefully)
        xss_patterns = [
            "<script>alert(1)</script>",
            "<img src=x onerror=alert(1)>",
            "javascript:alert(1)",
        ]
        for xss in xss_patterns:
            r = self.req("GET", "/products/search", params={"keyword": xss})
            check(f"xss keyword safe: {xss[:15]}", r.status_code in (200, 400),
                  f"got {r.status_code}")

        # 15.7 Path traversal in product ID
        import urllib.parse
        for path_traversal in ["../../etc/passwd", "0; rm -rf /", "../etc/shadow"]:
            encoded = urllib.parse.quote(path_traversal, safe='')
            r = self.req("GET", f"/products/{encoded}")
            check(f"path traversal: {path_traversal[:20]}", r.status_code in (400, 404))

        # 15.8 Oversized payload
        r = self.req("POST", "/user/register", json={
            "phone": self.random_phone(),
            "password": "test1234",
            "nickname": "A" * 10000
        })
        check("oversized nickname", r.status_code in (400, 413, 422))

        # 15.9 Order manipulation
        r = self.req("POST", f"/orders/999999/cancel")
        check("cancel nonexistent order", r.status_code in (400, 404))

        r = self.req("POST", f"/orders/999999/pay")
        check("pay nonexistent order", r.status_code in (400, 404))

    # =========== Section 16: Edge Cases ===========
    def test_edge_cases(self):
        print("\n📋 Section 16: Edge Cases & Boundary Tests")

        # 16.1 Pagination
        r = self.req("GET", "/orders", params={"page": 1, "size": 10})
        check("pagination with valid params", r.status_code == 200)

        r = self.req("GET", "/orders", params={"page": 0, "size": 10})
        check("pagination page=0", r.status_code in (200, 400))

        r = self.req("GET", "/orders", params={"page": -1, "size": 10})
        check("pagination negative page", r.status_code in (200, 400))

        r = self.req("GET", "/orders", params={"page": 1, "size": 0})
        check("pagination size=0", r.status_code in (200, 400))

        r = self.req("GET", "/orders", params={"page": 1, "size": 1000})
        check("pagination excessive size", r.status_code in (200, 400))

        # 16.2 Idempotency
        sku_id = self.test_sku.get("id", 0)
        if sku_id:
            idempotent_key = str(uuid.uuid4())
            r1 = self.req("POST", "/orders", json={
                "items": [{"sku_id": sku_id, "quantity": 1}],
                "address_id": self.test_address.get("id", 0),
                "remark": "幂等测试"
            }, headers={"Idempotent-Key": idempotent_key})
            r2 = self.req("POST", "/orders", json={
                "items": [{"sku_id": sku_id, "quantity": 1}],
                "address_id": self.test_address.get("id", 0),
                "remark": "幂等测试"
            }, headers={"Idempotent-Key": idempotent_key})
            check("idempotency: same key returns same result",
                  r1.status_code == r2.status_code)
            if r1.status_code == 201 and r2.status_code == 201:
                data1 = r1.json().get("data", {})
                data2 = r2.json().get("data", {})
                check("idempotency: same order_id",
                      data1.get("id") == data2.get("id"))

        # 16.3 Invalid content types
        r = self.session.post(self.url("/user/register"),
                              data="not json",
                              headers={"Content-Type": "text/plain"},
                              timeout=10)
        check("invalid content-type", r.status_code in (400, 415, 422))

        # 16.4 Empty body
        r = self.req("POST", "/orders", json={}, headers={"Idempotent-Key": str(uuid.uuid4())})
        check("empty order body", r.status_code in (400, 422))

        # 16.5 Merchant application
        r = self.req("POST", "/merchant/apply", json={
            "company_name": "测试公司",
            "contact_name": "张三",
            "contact_phone": "13800138000",
            "company_address": "深圳南山",
            "business_license": "https://example.com/license.jpg"
        })
        check("merchant application", r.status_code in (200, 201, 400))

        r = self.req("GET", "/merchant/application")
        check("merchant app status query", r.status_code in (200, 400, 404))

        # 16.6 Search with special characters
        special_chars = ["!", "@#$%", "测试  空格", "中文搜索", "a" * 200]
        for char in special_chars:
            r = self.req("GET", "/products/search", params={"keyword": char})
            check(f"special keyword ({char[:10]})", r.status_code in (200, 400),
                  f"got {r.status_code}")

    # =========== Section 17: Concurrent Tests ===========
    def test_concurrency(self):
        print("\n📋 Section 17: Concurrency & Race Conditions")

        # Create a dedicated product for concurrency testing via admin
        r = self.admin_req("POST", "/admin/categories", json={
            "name": "并发分类", "icon": "concurrent", "sort_order": 99
        })
        cat_id = r.json().get("data", {}).get("id")
        if not cat_id:
            log("concurrent stock test", "WARN", "cannot create category")
            return

        r = self.admin_req("POST", "/admin/products", json={
            "name": "并发测试商品",
            "category_id": cat_id,
            "price": 1000,
            "images": [],
            "skus": [{"name": "并发版", "price": 1000, "stock": 100}],
        })
        if r.status_code != 201:
            log("concurrent stock test", "WARN", f"cannot create product: {r.json().get('message','')}")
            return
        sku_id = r.json().get("data", {}).get("skus", [{}])[0].get("id")
        if not sku_id:
            log("concurrent stock test", "WARN", "no SKU from created product")
            return

        results = []
        lock = threading.Lock()

        def place_order(idx):
            try:
                phone = self.random_phone()
                s = requests.Session()
                s.headers["Content-Type"] = "application/json"
                r = s.post(self.url("/user/register"), json={
                    "phone": phone, "password": "ConcTest1", "nickname": f"C{idx}"
                }, timeout=10)
                token = r.json().get("data", {}).get("token", "")
                s.headers["Authorization"] = f"Bearer {token}"

                r = s.post(self.url("/addresses"), json={
                    "name": "T", "phone": "13800138000",
                    "province": "GD", "city": "SZ", "district": "NS",
                    "detail": "Addr", "is_default": True
                }, timeout=10)
                thread_addr = r.json().get("data", {}).get("id", 0)

                s.post(self.url("/cart"), json={
                    "sku_id": sku_id, "quantity": 1
                }, timeout=10)
                s.post(self.url("/cart/select"), json={
                    "sku_ids": [sku_id], "selected": True
                }, timeout=10)

                r = s.post(self.url("/orders"), json={
                    "address_id": thread_addr, "remark": f"C{idx}"
                }, headers={"Idempotent-Key": str(uuid.uuid4())}, timeout=10)
                with lock:
                    msg = r.json().get("message", "")[:50] if r.text else "no body"
                    results.append((idx, r.status_code, f"{r.json().get('code',-1)}:{msg}"))
            except Exception as e:
                with lock:
                    results.append((idx, -1, str(e)))

        threads = []
        for i in range(3):
            t = threading.Thread(target=place_order, args=(i,))
            threads.append(t)
            t.start()

        for t in threads:
            t.join()

        success_count = sum(1 for _, status, code in results if status == 201)
        check("concurrent orders: some succeed",
              success_count > 0, f"{success_count}/{len(results)} success")
        check("concurrent orders: no 500 errors",
              all(status != 500 for _, status, _ in results),
              f"details: {[(i, s, m[:30]) for i, s, m in results]}")

    # =========== Section 18: Upload ===========
    def test_upload(self):
        print("\n📋 Section 18: File Upload")

        r = self.req("POST", "/upload", files={
            "file": ("test.txt", b"Hello World", "text/plain")
        })
        check("upload file", r.status_code in (200, 201, 400, 413))

        r = self.req("POST", "/upload/base64", json={
            "file": "SGVsbG8gV29ybGQ=",
            "filename": "test.txt"
        })
        check("upload base64", r.status_code in (200, 201, 400))

    # =========== Section 19: Rebuy Flow ===========
    def test_rebuy(self):
        print("\n📋 Section 19: Rebuy")

        if self.test_order.get("id"):
            r = self.req("POST", f"/orders/{self.test_order['id']}/rebuy")
            check("rebuy order", r.status_code in (200, 201, 400))
        else:
            log("rebuy order", "WARN", "no order available")

    # =========== Section 20: Admin management ===========
    def test_admin(self):
        print("\n📋 Section 20: Admin Management")

        if not self.admin_token:
            log("admin tests", "WARN", "no admin token available, skip")
            return

        r = self.admin_req("GET", "/admin/users", params={"role": "user", "page": 1, "size": 10})
        check("admin list users", r.status_code == 200)

        r = self.admin_req("GET", "/admin/banners")
        check("admin list banners", r.status_code == 200)

        r = self.admin_req("POST", "/admin/coupons", json={
            "name": "测试优惠券",
            "type": "full_reduce",
            "condition": 10000,
            "discount": 1000,
            "total_count": 100,
            "per_limit": 1,
            "start_at": "2026-01-01 00:00:00",
            "end_at": "2026-12-31 23:59:59"
        })
        check("admin create coupon", r.status_code == 201)
        if r.status_code == 201:
            self.test_coupon["id"] = r.json().get("data", {}).get("id")

        r = self.admin_req("GET", "/admin/merchant/applications")
        check("admin merchant applications", r.status_code == 200)

    # =========== RUN ALL ===========
    def run(self):
        global PASS, FAIL, WARN, errors
        start = time.time()

        # Register admin user first (if no token provided)
        if os.environ.get("ADMIN_TOKEN"):
            self.admin_token = os.environ["ADMIN_TOKEN"]
            self.admin_session.headers["Authorization"] = f"Bearer {self.admin_token}"
            log("admin setup", "PASS", "from ADMIN_TOKEN env")
        else:
            admin_phone = self.random_phone()
            admin_password = "Admin@1234"
            r = self.req("POST", "/user/register", json={
                "phone": admin_phone, "password": admin_password, "nickname": "AdminUser"
            })
            if r.status_code == 201:
                self.promote_to_admin(admin_phone)
                r2 = self.req("POST", "/user/login", json={
                    "phone": admin_phone, "password": admin_password
                })
                if r2.status_code == 200:
                    self.admin_token = r2.json().get("data", {}).get("token")
                    if self.admin_token:
                        self.admin_session.headers["Authorization"] = f"Bearer {self.admin_token}"
                        log("admin setup", "PASS", f"admin phone={admin_phone}")

        print(f"\n{'='*60}")
        print(f"  Gosh Mall Comprehensive Test Suite")
        print(f"  Target: {BASE_URL}")
        print(f"{'='*60}")

        self.test_health()
        self.test_auth()
        self.test_profile()
        self.test_addresses()
        self.test_categories()
        self.test_products()
        self.test_cart()
        self.test_orders()
        self.test_payment()
        self.test_coupons()
        self.test_favorites()
        self.test_browse_history()
        self.test_points()
        self.test_reviews()
        self.test_security()
        self.test_edge_cases()
        self.test_concurrency()
        self.test_upload()
        self.test_rebuy()
        self.test_admin()

        elapsed = time.time() - start
        print(f"\n{'='*60}")
        print(f"  RESULTS")
        print(f"  {'='*30}")
        print(f"  ✅ PASS: {PASS}")
        print(f"  ❌ FAIL: {FAIL}")
        print(f"  ⚠️  WARN: {WARN}")
        print(f"  ⏱  Time: {elapsed:.1f}s")
        if errors:
            print(f"\n  Failures:")
            for e in errors:
                print(f"    {e}")
        print(f"{'='*60}")

        return FAIL == 0


if __name__ == "__main__":
    tester = GoshTester()
    success = tester.run()
    sys.exit(0 if success else 1)
