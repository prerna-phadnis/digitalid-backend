package handlers

// RegisterRequest is the top-level DTO for the registration request body.
type RegisterRequest struct {
	PersonalInfo PersonalInfo  `json:"personal_info"  validate:"required,dive"`
	GovernmentID GovernmentID  `json:"government_id"  validate:"required,dive"`
	Contact      ContactInfo   `json:"contact_details" validate:"required,dive"`
	Travel       TravelDetails `json:"travel_details" validate:"required,dive"`
	Emergency    EmergencyInfo `json:"emergency_info" validate:"required,dive"`
	Consent      Consent       `json:"consent"        validate:"required,dive"`
}

// PersonalInfo holds basic personal details.
type PersonalInfo struct {
	FullName    string `json:"full_name"   validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required,datetime=2006-01-02"` // YYYY-MM-DD
	Gender      string `json:"gender"      validate:"required,oneof=Male Female Other"`
	Nationality string `json:"nationality" validate:"required"`
}

// GovernmentID represents an ID document reference.
type GovernmentID struct {
	IDType        string `json:"id_type"         validate:"required"`
	IDNumber      string `json:"id_number"       validate:"required"`
	IDDocumentURL string `json:"id_document_url" validate:"omitempty,url"`
}

// ContactInfo contains phone/email.
type ContactInfo struct {
	Mobile string `json:"mobile_number" validate:"required,e164"` // e164 style recommended: +<countrycode><number>
	Email  string `json:"email"         validate:"required,email"`
}

// TravelDetails contains trip-level fields.
type TravelDetails struct {
	TripItinerary    []ItineraryItem `json:"trip_itinerary" validate:"required,min=1,dive"`
	ArrivalDate      string          `json:"arrival_date"  validate:"required,datetime=2006-01-02"`
	DepartureDate    string          `json:"departure_date" validate:"required,datetime=2006-01-02"`
	BookingReference string          `json:"booking_reference" validate:"required"`
}

// ItineraryItem describes a single stop.
type ItineraryItem struct {
	City          string `json:"city"       validate:"required"`
	CheckIn       string `json:"check_in"   validate:"required,datetime=2006-01-02"`
	CheckOut      string `json:"check_out"  validate:"required,datetime=2006-01-02"`
	Accommodation string `json:"accommodation" validate:"omitempty"`
}

// EmergencyInfo contains emergency contacts and medical data.
type EmergencyInfo struct {
	Contacts          []EmergencyContact `json:"contacts" validate:"required,dive"`
	BloodGroup        string             `json:"blood_group" validate:"omitempty"`
	MedicalConditions string             `json:"medical_conditions" validate:"omitempty"`
}

type EmergencyContact struct {
	Name         string `json:"name"         validate:"required"`
	Relationship string `json:"relationship" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required,e164"`
}

// Consent booleans and acknowledgements.
type Consent struct {
	LiveTracking             bool `json:"live_tracking"`
	HealthDataSharing        bool `json:"health_data_sharing"`
	DataUsageAcknowledgement bool `json:"data_usage_acknowledgement"`
}
