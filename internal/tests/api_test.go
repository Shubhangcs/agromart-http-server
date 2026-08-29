// Package tests runs the HTTP API end-to-end against a real Postgres database.
// Set TEST_DATABASE_NAME (and optionally TEST_DATABASE_HOST/USER/PASSWORD/PORT) to enable;
// the schema in that database is dropped and re-migrated on every run.
package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/shubhangcs/agromart-server/internal/app"
	"github.com/shubhangcs/agromart-server/internal/mailer"
	"github.com/shubhangcs/agromart-server/internal/routes"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	name := os.Getenv("TEST_DATABASE_NAME")
	if name == "" {
		fmt.Println("TEST_DATABASE_NAME not set: skipping API integration tests")
		os.Exit(0)
	}
	set := func(k, v string) {
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
	os.Setenv("DATABASE_NAME", name)
	os.Setenv("DATABASE_HOST", env("TEST_DATABASE_HOST", "localhost"))
	os.Setenv("DATABASE_PORT", env("TEST_DATABASE_PORT", "5432"))
	os.Setenv("DATABASE_USER", env("TEST_DATABASE_USER", "postgres"))
	os.Setenv("DATABASE_PASSWORD", env("TEST_DATABASE_PASSWORD", "postgres"))
	os.Setenv("DATABASE_SSL_MODE", "disable")
	set("JWT_SECRET_KEY", "test-secret")
	set("JWT_TOKEN_ISSUER", "test-issuer")
	set("ACCESS_KEY", "test")
	set("SECRET_KEY", "test")
	set("BUCKET_NAME", "test-bucket")
	os.Setenv("PUSH_DRY_RUN", "true")
	os.Setenv("ADMIN_BOOTSTRAP_SECRET", "bootstrap-secret")
	os.Setenv("AUTH_RATE_LIMIT_PER_MIN", "40")
	os.Unsetenv("RESEND_API_KEY")

	// fresh schema every run
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), name)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	if _, err = db.Exec(`DROP SCHEMA public CASCADE; CREATE SCHEMA public;`); err != nil {
		panic(fmt.Errorf("reset schema: %w", err))
	}
	db.Close()

	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	srv = httptest.NewServer(routes.SetupRoutes(application))
	code := m.Run()
	srv.Close()
	os.Exit(code)
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// ---- helpers -------------------------------------------------------------

type resp struct {
	Status int
	Body   map[string]any
}

func call(t *testing.T, method, path, token string, body any, headers ...string) resp {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, srv.URL+path, rdr)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer res.Body.Close()
	out := resp{Status: res.StatusCode, Body: map[string]any{}}
	_ = json.NewDecoder(res.Body).Decode(&out.Body)
	return out
}

func want(t *testing.T, r resp, status int, what string) resp {
	t.Helper()
	if r.Status != status {
		t.Fatalf("%s: want %d got %d body=%v", what, status, r.Status, r.Body)
	}
	return r
}

func str(m map[string]any, k string) string { s, _ := m[k].(string); return s }

func claims(token string) map[string]any {
	parts := strings.Split(token, ".")
	payload := parts[1] + strings.Repeat("=", (4-len(parts[1])%4)%4)
	var out map[string]any
	b := make([]byte, 0)
	for _, r := range payload { // base64url -> base64
		switch r {
		case '-':
			b = append(b, '+')
		case '_':
			b = append(b, '/')
		default:
			b = append(b, byte(r))
		}
	}
	dec := make([]byte, len(b))
	n, _ := base64Decode(dec, b)
	_ = json.Unmarshal(dec[:n], &out)
	return out
}

// ---- the scenario ----------------------------------------------------------

type actors struct {
	admin, seller, buyer           string // tokens
	sellerID, buyerID              string
	sellerBiz, buyerBiz            string
	category, subCategory, product string
	rfq                            string
}

var a actors

func TestAPI(t *testing.T) {
	t.Run("01 bootstrap admin", func(t *testing.T) {
		body := map[string]any{"first_name": "Root", "email": "root@t.com", "phone": "9000000000", "password": "Passw0rd!"}
		want(t, call(t, "POST", "/admin/bootstrap", "", body), 403, "bootstrap without secret")
		want(t, call(t, "POST", "/admin/bootstrap", "", body, "X-Bootstrap-Secret", "bootstrap-secret"), 201, "bootstrap with secret")
		want(t, call(t, "POST", "/admin/bootstrap", "", body, "X-Bootstrap-Secret", "bootstrap-secret"), 403, "second bootstrap")
		want(t, call(t, "POST", "/admin/create", "", body), 401, "public admin create is gone")
		r := want(t, call(t, "POST", "/admin/login", "", map[string]any{"email": "root@t.com", "password": "Passw0rd!"}), 200, "admin login")
		a.admin = str(r.Body, "token")
		if claims(a.admin)["role"] != "admin" {
			t.Fatalf("admin token role = %v", claims(a.admin)["role"])
		}
	})

	t.Run("02 users sign up and log in", func(t *testing.T) {
		want(t, call(t, "POST", "/user/create", "", map[string]any{"first_name": "Seller", "email": "seller@t.com", "phone": "9000000001", "password": "Passw0rd!"}), 201, "create seller")
		want(t, call(t, "POST", "/user/create", "", map[string]any{"first_name": "Buyer", "email": "buyer@t.com", "phone": "9000000002", "password": "Passw0rd!"}), 201, "create buyer")
		if m, ok := mailer.LastCaptured("buyer@t.com"); !ok || !strings.Contains(m.Subject, "Welcome") {
			t.Fatalf("welcome email not sent: %+v %v", m, ok)
		}
		a.seller = str(want(t, call(t, "POST", "/user/login", "", map[string]any{"email": "seller@t.com", "password": "Passw0rd!"}), 200, "seller login").Body, "token")
		a.buyer = str(want(t, call(t, "POST", "/user/login", "", map[string]any{"email": "buyer@t.com", "password": "Passw0rd!"}), 200, "buyer login").Body, "token")
		a.sellerID = str(claims(a.seller), "user_id")
		a.buyerID = str(claims(a.buyer), "user_id")
		if claims(a.seller)["role"] != "user" {
			t.Fatalf("user token role = %v", claims(a.seller)["role"])
		}
	})

	t.Run("03 admin-only catalogue", func(t *testing.T) {
		cat := map[string]any{"name": "Cashew", "description": "Cashews"}
		want(t, call(t, "POST", "/category/create", a.seller, cat), 403, "user creates category")
		r := want(t, call(t, "POST", "/category/create", a.admin, cat), 201, "admin creates category")
		a.category = firstString(r.Body, "category_id", "id")
		r = want(t, call(t, "POST", "/category/sub/create", a.admin, map[string]any{"category_id": a.category, "name": "W320", "description": "grade"}), 201, "sub category")
		a.subCategory = firstString(r.Body, "sub_category_id", "id")
		want(t, call(t, "GET", "/category/get/all", a.buyer, nil), 200, "users can read categories")
	})

	t.Run("04 businesses and products", func(t *testing.T) {
		biz := func(uid, name, email, phone string) map[string]any {
			return map[string]any{"user_id": uid, "name": name, "email": email, "phone": phone, "address": "Surathkal", "city": "Mangalore", "state": "Karnataka", "pincode": "575011", "business_type": "Trader"}
		}
		a.sellerBiz = str(want(t, call(t, "POST", "/business/create", a.seller, biz(a.sellerID, "Canara Nuts", "cn@t.com", "9000000003")), 201, "seller business").Body, "business_id")
		a.buyerBiz = str(want(t, call(t, "POST", "/business/create", a.buyer, biz(a.buyerID, "Buyer Co", "bc@t.com", "9000000004")), 201, "buyer business").Body, "business_id")
		prod := map[string]any{"business_id": a.sellerBiz, "category_id": a.category, "sub_category_id": a.subCategory, "name": "Cashew W320", "description": "d", "quantity": 100, "unit": "kg", "price": 800, "moq": "10", "is_product_active": true}
		want(t, call(t, "POST", "/product/create", a.buyer, prod), 403, "buyer creates product under seller business")
		a.product = str(want(t, call(t, "POST", "/product/create", a.seller, prod), 201, "seller creates product").Body, "product_id")
		want(t, call(t, "PUT", "/business/social/update/"+a.sellerBiz, a.seller, map[string]any{"website": "https://b.com"}), 404, "social update before create")
		want(t, call(t, "POST", "/business/social/create", a.seller, map[string]any{"id": a.sellerBiz, "website": "https://a.com"}), 201, "social create")
		want(t, call(t, "PUT", "/business/social/update/"+a.sellerBiz, a.seller, map[string]any{"website": "https://b.com"}), 200, "social update without id in body")
		want(t, call(t, "PUT", "/business/update/"+a.sellerBiz, a.buyer, biz(a.sellerID, "Hijack", "cn@t.com", "9000000003")), 403, "buyer edits seller business")
	})

	t.Run("05 ownership of products", func(t *testing.T) {
		upd := map[string]any{"name": "Cashew W320 premium", "description": "d", "quantity": 100, "unit": "kg", "price": 850, "moq": "10"}
		want(t, call(t, "PUT", "/product/update/"+a.product, a.buyer, upd), 403, "buyer updates seller product")
		want(t, call(t, "PUT", "/product/update/"+a.product, a.seller, upd), 200, "seller updates own product")
		want(t, call(t, "PUT", "/product/update/"+a.product, a.admin, upd), 200, "admin updates product")
		want(t, call(t, "DELETE", "/product/delete/"+a.product, a.buyer, nil), 403, "buyer deletes seller product")
	})

	t.Run("06 ratings, follows and embedded counts", func(t *testing.T) {
		want(t, call(t, "POST", "/product/rate", a.buyer, map[string]any{"product_id": a.product, "user_id": a.sellerID, "rating": 4}), 200, "rate product (spoofed user_id)")
		want(t, call(t, "POST", "/business/rate", a.buyer, map[string]any{"business_id": a.sellerBiz, "user_id": a.buyerID, "rating": 5}), 200, "rate business")
		want(t, call(t, "POST", "/follower/follow", a.buyer, map[string]any{"business_id": a.sellerBiz, "user_id": a.buyerID}), 201, "follow")
		r := want(t, call(t, "GET", "/product/rate/get/"+a.product, a.buyer, nil), 200, "ratings list")
		ratings := r.Body["ratings"].([]any)
		if str(ratings[0].(map[string]any), "user_id") != a.buyerID {
			t.Fatalf("rating user_id was taken from the body, not the token")
		}
		r = want(t, call(t, "GET", "/product/get/all", a.buyer, nil), 200, "products")
		p := r.Body["products"].([]any)[0].(map[string]any)
		if p["average_rating"] != 4.0 || p["rating_count"] != 1.0 {
			t.Fatalf("product counts: %v %v", p["average_rating"], p["rating_count"])
		}
		r = want(t, call(t, "GET", "/business/get/"+a.sellerBiz, a.buyer, nil), 200, "business")
		d := r.Body["details"].(map[string]any)
		if d["average_rating"] != 5.0 || d["rating_count"] != 1.0 || d["followers_count"] != 1.0 {
			t.Fatalf("business counts: %v", d)
		}
	})

	t.Run("07 product filters on category endpoints", func(t *testing.T) {
		count := func(q string) int {
			r := want(t, call(t, "GET", "/product/get/category/"+a.category+"?"+q, a.buyer, nil), 200, q)
			if r.Body["products"] == nil {
				return 0
			}
			return len(r.Body["products"].([]any))
		}
		if count("q=xyz") != 0 || count("q=cashew") != 1 || count("city=mangalore") != 1 || count("city=delhi") != 0 {
			t.Fatalf("category filters not applied")
		}
	})

	t.Run("08 RFQs", func(t *testing.T) {
		rfq := map[string]any{"business_id": a.sellerBiz, "category_id": a.category, "sub_category_id": a.subCategory, "product_name": "Almonds", "quantity": 50, "unit": "kg", "price": 700, "is_rfq_active": true}
		a.rfq = str(want(t, call(t, "POST", "/rfq/create", a.seller, rfq), 201, "create rfq").Body, "rfq_id")
		r := want(t, call(t, "GET", "/rfq/get/one/"+a.rfq, a.buyer, nil), 200, "rfq by id")
		if str(r.Body["rfq"].(map[string]any), "category_name") != "Cashew" {
			t.Fatalf("rfq/get/one missing category: %v", r.Body)
		}
		want(t, call(t, "GET", "/rfq/get/one/00000000-0000-0000-0000-000000000000", a.buyer, nil), 404, "unknown rfq")
		want(t, call(t, "PUT", "/rfq/update/status/"+a.rfq, a.buyer, map[string]any{"is_rfq_active": false}), 403, "buyer deactivates seller rfq")
		r = want(t, call(t, "GET", "/rfq/get/"+a.sellerBiz, a.seller, nil), 200, "rfqs by business")
		if str(r.Body["rfqs"].([]any)[0].(map[string]any), "category_name") != "Cashew" {
			t.Fatalf("rfq by business missing category name")
		}
	})

	t.Run("09 leads and notifications", func(t *testing.T) {
		lead := map[string]any{"enquirer_business_id": a.buyerBiz, "enquire_to_id": a.sellerBiz, "product_id": a.product, "enquiry_message": "price for 500kg?", "order_quantity": 500}
		want(t, call(t, "POST", "/leads/create", a.seller, lead), 403, "seller enquires as buyer business")
		want(t, call(t, "POST", "/leads/create", a.buyer, lead), 201, "buyer enquires")
		want(t, call(t, "GET", "/leads/received/"+a.sellerBiz, a.buyer, nil), 403, "buyer reads seller leads")
		want(t, call(t, "GET", "/leads/received/"+a.sellerBiz, a.seller, nil), 200, "seller reads own leads")
		want(t, call(t, "POST", "/user/push-token", a.seller, map[string]any{"token": "garbage"}), 400, "bad push token")
		want(t, call(t, "POST", "/user/push-token", a.seller, map[string]any{"token": "ExponentPushToken[abc]", "platform": "android"}), 200, "push token")
	})

	t.Run("10 chat newest first", func(t *testing.T) {
		want(t, call(t, "POST", "/chat/send", a.buyer, map[string]any{"receiver_id": a.sellerID, "content": "first"}), 201, "send 1")
		want(t, call(t, "POST", "/chat/send", a.buyer, map[string]any{"receiver_id": a.sellerID, "content": "second"}), 201, "send 2")
		r := want(t, call(t, "GET", "/chat/history?with_user_id="+a.buyerID, a.seller, nil), 200, "history")
		msgs := r.Body["messages"].([]any)
		if str(msgs[0].(map[string]any), "content") != "second" {
			t.Fatalf("history not newest-first: %v", msgs)
		}
	})

	t.Run("11 banners and app config", func(t *testing.T) {
		want(t, call(t, "POST", "/banners/create", a.seller, map[string]any{"title": "x"}), 403, "user creates banner")
		id := str(want(t, call(t, "POST", "/banners/create", a.admin, map[string]any{"title": "Monsoon offers", "sort_order": 1}), 201, "admin creates banner").Body, "banner_id")
		if n := len(want(t, call(t, "GET", "/banners/get/active", a.buyer, nil), 200, "active").Body["banners"].([]any)); n != 0 {
			t.Fatalf("banner without image should not be active, got %d", n)
		}
		want(t, call(t, "PUT", "/banners/update/image/"+id, a.admin, nil), 200, "presign banner image")
		if n := len(want(t, call(t, "GET", "/banners/get/active", a.buyer, nil), 200, "active").Body["banners"].([]any)); n != 1 {
			t.Fatalf("expected 1 active banner, got %d", n)
		}
		want(t, call(t, "DELETE", "/banners/delete/"+id, a.admin, nil), 200, "delete banner")
		want(t, call(t, "DELETE", "/banners/delete/"+id, a.admin, nil), 404, "delete again")
		r := want(t, call(t, "GET", "/app/config", "", nil), 200, "app config is public")
		if str(r.Body, "min_app_version") == "" {
			t.Fatalf("app config missing min_app_version")
		}
	})

	t.Run("12 forgot / reset password", func(t *testing.T) {
		want(t, call(t, "POST", "/user/forgot-password", "", map[string]any{"email": "nobody@t.com"}), 200, "unknown email does not leak")
		want(t, call(t, "POST", "/user/forgot-password", "", map[string]any{"email": "buyer@t.com"}), 200, "forgot")
		m, ok := mailer.LastCaptured("buyer@t.com")
		code := regexp.MustCompile(`code is (\d{6})`).FindStringSubmatch(m.Text)
		if !ok || code == nil {
			t.Fatalf("reset email not captured: %+v", m)
		}
		want(t, call(t, "POST", "/user/reset-password", "", map[string]any{"email": "buyer@t.com", "code": "000000", "new_password": "NewPassw0rd!"}), 401, "wrong code")
		want(t, call(t, "POST", "/user/reset-password", "", map[string]any{"email": "buyer@t.com", "code": code[1], "new_password": "NewPassw0rd!"}), 200, "right code")
		want(t, call(t, "POST", "/user/reset-password", "", map[string]any{"email": "buyer@t.com", "code": code[1], "new_password": "Another1!"}), 401, "code is single use")
		want(t, call(t, "POST", "/user/login", "", map[string]any{"email": "buyer@t.com", "password": "Passw0rd!"}), 401, "old password")
		want(t, call(t, "POST", "/user/login", "", map[string]any{"email": "buyer@t.com", "password": "NewPassw0rd!"}), 200, "new password")
	})

	t.Run("13 self-or-admin on user routes", func(t *testing.T) {
		upd := map[string]any{"first_name": "Seller", "last_name": "One", "email": "seller@t.com", "phone": "9000000001"}
		want(t, call(t, "PUT", "/user/update/details/"+a.sellerID, a.buyer, upd), 403, "buyer edits seller")
		want(t, call(t, "PUT", "/user/update/details/"+a.sellerID, a.seller, upd), 200, "seller edits self")
		want(t, call(t, "PUT", "/user/update/details/"+a.sellerID, a.admin, upd), 200, "admin edits seller")
		want(t, call(t, "GET", "/user/get/all", a.buyer, nil), 403, "user lists users")
		want(t, call(t, "GET", "/user/get/all", a.admin, nil), 200, "admin lists users")
	})

	t.Run("14 auth rate limit", func(t *testing.T) {
		limited := false
		for i := 0; i < 60; i++ {
			if call(t, "POST", "/user/login", "", map[string]any{"email": "nobody@t.com", "password": "x"}).Status == 429 {
				limited = true
				break
			}
		}
		if !limited {
			t.Fatalf("login was never rate limited")
		}
	})
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if s := str(m, k); s != "" {
			return s
		}
	}
	return ""
}
