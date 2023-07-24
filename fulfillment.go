package goshopify

import (
	"fmt"
	"time"
)

// FulfillmentService is an interface for interfacing with the fulfillment endpoints
// of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/2023-07/resources/fulfillment
type FulfillmentService interface {
	Create(FulfillmentRequest) (*Fulfillment, error)
}

// FulfillmentsService is an interface for other Shopify resources
// to interface with the fulfillment endpoints of the Shopify API.
// https://shopify.dev/docs/api/admin-rest/2023-07/resources/fulfillment
type FulfillmentsService interface {
	ListFulfillments(int64, interface{}) ([]Fulfillment, error)
	CountFulfillments(int64, interface{}) (int, error)
	GetFulfillment(int64, int64, interface{}) (*Fulfillment, error)
	CreateFulfillment(int64, Fulfillment) (*Fulfillment, error)
	UpdateFulfillment(int64, Fulfillment) (*Fulfillment, error)
	CompleteFulfillment(int64, int64) (*Fulfillment, error)
	TransitionFulfillment(int64, int64) (*Fulfillment, error)
	CancelFulfillment(int64, int64) (*Fulfillment, error)
}

// FulfillmentServiceOp handles communication with the fulfillment
// related methods of the Shopify API.
type FulfillmentServiceOp struct {
	client     *Client
	resource   string
	resourceID int64
}

// FulfillmentRequest represents a request to fulfill a FulfillmentOrder
type FulfillmentRequest struct {
	Fulfillment *FulfillmentRequestData `json:"fulfillment"`
}

type FulfillmentRequestData struct {
	LineItemsByFulfillmentOrder []LineItemsByFulfillmentOrder `json:"line_items_by_fulfillment_order"`
	Message                     string                        `json:"message,omitempty"`
	NotifyCustomer              bool                          `json:"notify_customer,omitempty"`
	OriginAddress               OriginAddress                 `json:"origin_address,omitempty"`
	TrackingInfo                TrackingInfo                  `json:"tracking_info,omitempty"`
}

type LineItemsByFulfillmentOrder struct {
	ID        int64                 `json:"fulfillment_order_id"`
	LineItems []FulfillmentLineItem `json:"fulfillment_order_line_items,omitempty"`
}

type FulfillmentLineItem struct {
	ID  int64 `json:"id"`
	Qty int   `json:"qty"`
}

type OriginAddress struct {
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	CountryCode  string `json:"country_code"`
	ProvinceCode string `json:"province_code,omitempty"`
	Zip          string `json:"zip,omitempty"`
}

type TrackingInfo struct {
	Company string `json:"company,omitempty"`
	Number  string `json:"number,omitempty"`
	URL     string `json:"url,omitempty"`
}

// Fulfillment represents a Shopify fulfillment.
type Fulfillment struct {
	ID              int64      `json:"id,omitempty"`
	OrderID         int64      `json:"order_id,omitempty"`
	LocationID      int64      `json:"location_id,omitempty"`
	Status          string     `json:"status,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	Service         string     `json:"service,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	TrackingCompany string     `json:"tracking_company,omitempty"`
	ShipmentStatus  string     `json:"shipment_status,omitempty"`
	TrackingNumber  string     `json:"tracking_number,omitempty"`
	TrackingNumbers []string   `json:"tracking_numbers,omitempty"`
	TrackingUrl     string     `json:"tracking_url,omitempty"`
	TrackingUrls    []string   `json:"tracking_urls,omitempty"`
	Receipt         Receipt    `json:"receipt,omitempty"`
	LineItems       []LineItem `json:"line_items,omitempty"`
	NotifyCustomer  bool       `json:"notify_customer"`
}

// Receipt represents a Shopify receipt.
type Receipt struct {
	TestCase      bool   `json:"testcase,omitempty"`
	Authorization string `json:"authorization,omitempty"`
}

// FulfillmentResource represents the result from the fulfillments/X.json endpoint
type FulfillmentResource struct {
	Fulfillment *Fulfillment `json:"fulfillment"`
}

// FulfillmentsResource represents the result from the fullfilments.json endpoint
type FulfillmentsResource struct {
	Fulfillments []Fulfillment `json:"fulfillments"`
}

// List fulfillments
func (s *FulfillmentServiceOp) List(options interface{}) ([]Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s.json", prefix)
	resource := new(FulfillmentsResource)
	err := s.client.Get(path, resource, options)
	return resource.Fulfillments, err
}

// Count fulfillments
func (s *FulfillmentServiceOp) Count(options interface{}) (int, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/count.json", prefix)
	return s.client.Count(path, options)
}

// Get individual fulfillment
func (s *FulfillmentServiceOp) Get(fulfillmentID int64, options interface{}) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d.json", prefix, fulfillmentID)
	resource := new(FulfillmentResource)
	err := s.client.Get(path, resource, options)
	return resource.Fulfillment, err
}

// Create a new fulfillment
func (s *FulfillmentServiceOp) Create(fulfillment FulfillmentRequestData) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s.json", prefix)
	wrappedData := FulfillmentRequest{Fulfillment: &fulfillment}
	resource := new(FulfillmentResource)
	err := s.client.Post(path, wrappedData, resource)
	return resource.Fulfillment, err
}

// Update an existing fulfillment
func (s *FulfillmentServiceOp) Update(fulfillment Fulfillment) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d.json", prefix, fulfillment.ID)
	wrappedData := FulfillmentResource{Fulfillment: &fulfillment}
	resource := new(FulfillmentResource)
	err := s.client.Put(path, wrappedData, resource)
	return resource.Fulfillment, err
}

// Complete an existing fulfillment
func (s *FulfillmentServiceOp) Complete(fulfillmentID int64) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d/complete.json", prefix, fulfillmentID)
	resource := new(FulfillmentResource)
	err := s.client.Post(path, nil, resource)
	return resource.Fulfillment, err
}

// Transition an existing fulfillment
func (s *FulfillmentServiceOp) Transition(fulfillmentID int64) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d/open.json", prefix, fulfillmentID)
	resource := new(FulfillmentResource)
	err := s.client.Post(path, nil, resource)
	return resource.Fulfillment, err
}

// Cancel an existing fulfillment
func (s *FulfillmentServiceOp) Cancel(fulfillmentID int64) (*Fulfillment, error) {
	prefix := FulfillmentPathPrefix(s.resource, s.resourceID)
	path := fmt.Sprintf("%s/%d/cancel.json", prefix, fulfillmentID)
	resource := new(FulfillmentResource)
	err := s.client.Post(path, nil, resource)
	return resource.Fulfillment, err
}
