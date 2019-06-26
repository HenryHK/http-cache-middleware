package cache

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"time"

	api "github.com/HenryHK/http-cache-middleware/api"
)

// Response is the cached response data structure.
type Response struct {
	// Value is the response value to be cached.
	Value []byte

	// Header is the response header to be cached.
	Header http.Header

	// Expiration is the cache expiration date.
	Expiration time.Time
}

// Client data structure for HTTP cache middleware.
type Client struct {
	adapter Adapter
	ttl     time.Duration
}

// ClientOption is used to set Client settings.
type ClientOption func(c *Client) error

// Adapter interface for HTTP cache middleware client.
type Adapter interface {
	// Get retrieves the cached response by a given key. It also
	// returns true or false, whether it exists or not.
	Get(key string) ([]byte, bool)

	// Set caches a response for a given key until an expiration date.
	Set(key string, response []byte, expiration time.Time)

	// Release frees cache for a given key.
	Release(key string)
}

// Middleware is the HTTP cache middleware handler.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get access key of request
		accessKey := r.Header.Get("autopilotapikey")
		if r.Method == "GET" {
			// get id/email from url and combine with access key as key
			id := strings.TrimPrefix(r.URL.Path, "/contact/")
			key := fmt.Sprintf("%s-%s", accessKey, id)
			// retrieve value from redis
			b, ok := c.adapter.Get(key)
			response := BytesToResponse(b)
			if ok {
				if response.Expiration.After(time.Now()) {
					c.adapter.Set(key, response.Bytes(), response.Expiration)

					w.Write(response.Value)
					return
				}
				// if the cache expired
				c.adapter.Release(key)
			}

			// redirect to api endpoint
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)
			result := rec.Result()

			value := rec.Body.Bytes()
			// fmt.Println(rec.Code)
			// fmt.Println(result.Header)
			// Naive method to test whether recorder had a successfult return
			// PS. httptest recorder always have status = 200 from result, need to retrieve from body
			if rec.Code == 200 {
				now := time.Now()

				response := Response{
					Value:      value,
					Header:     result.Header,
					Expiration: now.Add(c.ttl),
				}
				c.adapter.Set(key, response.Bytes(), response.Expiration)
			}

			w.Write(value)
			return
		}

		// POST handles add/update operation
		if r.Method == "POST" || r.Method == "PUT" {
			decoder := json.NewDecoder(r.Body)

			var contactRequest api.ContactRequest
			err := decoder.Decode(&contactRequest)
			if err != nil {
				panic(err)
			}
			// invalidate using email
			email := contactRequest.Contact.Email
			key := fmt.Sprintf("%s-%s", accessKey, email)
			c.adapter.Release(key)
			// invalidate using contact id
			id := contactRequest.Contact.ContactID
			key = fmt.Sprintf("%s-%s", accessKey, id)
			c.adapter.Release(key)
			// rebuild request, body consumed by previous reading
			rebuiltBody, _ := json.Marshal(contactRequest)
			r, _ = http.NewRequest("POST", r.URL.Path, bytes.NewBuffer(rebuiltBody))
			r.Header.Add("autopilotapikey", accessKey)
			r.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			value := rec.Body.Bytes()
			if rec.Code == 200 {
				// hacking way to get json from reponse
				re, _ := regexp.Compile("{[^}]+}")
				contactID := re.FindString(string(value))
				var addOrUpdateResponse api.AddOrUpdateResponse
				err = json.Unmarshal([]byte(contactID), &addOrUpdateResponse)
				if err != nil {
					log.Printf("%v", err)
				}
				key = fmt.Sprintf("%s-%s", accessKey, addOrUpdateResponse.ContactID)
				c.adapter.Release(key)
			}

			w.WriteHeader(rec.Code)
			w.Write(value)
			return
		}

		// DELETE handler, invalidate contact and pass to next handler
		if r.Method == "DELETE" {
			id := strings.TrimPrefix(r.URL.Path, "/contact/")
			key := fmt.Sprintf("%s-%s", accessKey, id)
			// retrieve value from redis
			c.adapter.Release(key)
		}

		next.ServeHTTP(w, r)
	})
}

// BytesToResponse converts bytes array into Response data structure.
func BytesToResponse(b []byte) Response {
	var r Response
	dec := gob.NewDecoder(bytes.NewReader(b))
	dec.Decode(&r)

	return r
}

// Bytes converts Response data structure into bytes array.
func (r Response) Bytes() []byte {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(&r)

	return b.Bytes()
}

// New initializes the cache HTTP middleware client with the given options.
func New(opts ...ClientOption) (*Client, error) {
	c := &Client{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.adapter == nil {
		return nil, errors.New("cache client adapter is not set")
	}
	if int64(c.ttl) < 1 {
		return nil, errors.New("cache client ttl is not set")
	}

	return c, nil
}

// ClientWithAdapter sets the adapter type for the HTTP cache middleware client.
func ClientWithAdapter(a Adapter) ClientOption {
	return func(c *Client) error {
		c.adapter = a
		return nil
	}
}

// ClientWithTTL sets how long each response is going to be cached.
func ClientWithTTL(ttl time.Duration) ClientOption {
	return func(c *Client) error {
		if int64(ttl) < 1 {
			return fmt.Errorf("cache client ttl %v is invalid", ttl)
		}

		c.ttl = ttl

		return nil
	}
}
