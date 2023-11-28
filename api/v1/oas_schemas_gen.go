// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"
	"time"
)

func (s *ErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

type ApiKey struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *ApiKey) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *ApiKey) SetAPIKey(val string) {
	s.APIKey = val
}

// Ref: #/components/schemas/Error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *Error) GetCode() int64 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *Error) SetCode(val int64) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// ErrorStatusCode wraps Error with StatusCode.
type ErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

// Ref: #/components/schemas/GetOrganizationRequest
type GetOrganizationRequest struct {
	ID string `json:"id"`
}

// GetID returns the value of ID.
func (s *GetOrganizationRequest) GetID() string {
	return s.ID
}

// SetID sets the value of ID.
func (s *GetOrganizationRequest) SetID(val string) {
	s.ID = val
}

// Ref: #/components/schemas/GetOrganizationsByIdRequest
type GetOrganizationsByIdRequest struct {
	ID []string `json:"id"`
}

// GetID returns the value of ID.
func (s *GetOrganizationsByIdRequest) GetID() []string {
	return s.ID
}

// SetID sets the value of ID.
func (s *GetOrganizationsByIdRequest) SetID(val []string) {
	s.ID = val
}

// Ref: #/components/schemas/GetOrganizationsRequest
type GetOrganizationsRequest struct {
	Cursor string `json:"cursor"`
}

// GetCursor returns the value of Cursor.
func (s *GetOrganizationsRequest) GetCursor() string {
	return s.Cursor
}

// SetCursor sets the value of Cursor.
func (s *GetOrganizationsRequest) SetCursor(val string) {
	s.Cursor = val
}

// Ref: #/components/schemas/GetPeopleByIdRequest
type GetPeopleByIdRequest struct {
	ID []string `json:"id"`
}

// GetID returns the value of ID.
func (s *GetPeopleByIdRequest) GetID() []string {
	return s.ID
}

// SetID sets the value of ID.
func (s *GetPeopleByIdRequest) SetID(val []string) {
	s.ID = val
}

// Ref: #/components/schemas/GetPeopleRequest
type GetPeopleRequest struct {
	Cursor string `json:"cursor"`
}

// GetCursor returns the value of Cursor.
func (s *GetPeopleRequest) GetCursor() string {
	return s.Cursor
}

// SetCursor sets the value of Cursor.
func (s *GetPeopleRequest) SetCursor(val string) {
	s.Cursor = val
}

// Ref: #/components/schemas/GetPersonRequest
type GetPersonRequest struct {
	ID string `json:"id"`
}

// GetID returns the value of ID.
func (s *GetPersonRequest) GetID() string {
	return s.ID
}

// SetID sets the value of ID.
func (s *GetPersonRequest) SetID(val string) {
	s.ID = val
}

// NewOptBool returns new OptBool with value set to v.
func NewOptBool(v bool) OptBool {
	return OptBool{
		Value: v,
		Set:   true,
	}
}

// OptBool is optional bool.
type OptBool struct {
	Value bool
	Set   bool
}

// IsSet returns true if OptBool was set.
func (o OptBool) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptBool) Reset() {
	var v bool
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptBool) SetTo(v bool) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptBool) Get() (v bool, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptBool) Or(d bool) bool {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptPersonSettings returns new OptPersonSettings with value set to v.
func NewOptPersonSettings(v PersonSettings) OptPersonSettings {
	return OptPersonSettings{
		Value: v,
		Set:   true,
	}
}

// OptPersonSettings is optional PersonSettings.
type OptPersonSettings struct {
	Value PersonSettings
	Set   bool
}

// IsSet returns true if OptPersonSettings was set.
func (o OptPersonSettings) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptPersonSettings) Reset() {
	var v PersonSettings
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptPersonSettings) SetTo(v PersonSettings) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptPersonSettings) Get() (v PersonSettings, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptPersonSettings) Or(d PersonSettings) PersonSettings {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/Organization
type Organization struct {
	ID          OptString            `json:"id"`
	DateCreated OptDateTime          `json:"date_created"`
	DateUpdated OptDateTime          `json:"date_updated"`
	Type        OptString            `json:"type"`
	Acronym     OptString            `json:"acronym"`
	NameDut     OptString            `json:"name_dut"`
	NameEng     OptString            `json:"name_eng"`
	Parent      []OrganizationParent `json:"parent"`
	Identifier  []PropertyValue      `json:"identifier"`
}

// GetID returns the value of ID.
func (s *Organization) GetID() OptString {
	return s.ID
}

// GetDateCreated returns the value of DateCreated.
func (s *Organization) GetDateCreated() OptDateTime {
	return s.DateCreated
}

// GetDateUpdated returns the value of DateUpdated.
func (s *Organization) GetDateUpdated() OptDateTime {
	return s.DateUpdated
}

// GetType returns the value of Type.
func (s *Organization) GetType() OptString {
	return s.Type
}

// GetAcronym returns the value of Acronym.
func (s *Organization) GetAcronym() OptString {
	return s.Acronym
}

// GetNameDut returns the value of NameDut.
func (s *Organization) GetNameDut() OptString {
	return s.NameDut
}

// GetNameEng returns the value of NameEng.
func (s *Organization) GetNameEng() OptString {
	return s.NameEng
}

// GetParent returns the value of Parent.
func (s *Organization) GetParent() []OrganizationParent {
	return s.Parent
}

// GetIdentifier returns the value of Identifier.
func (s *Organization) GetIdentifier() []PropertyValue {
	return s.Identifier
}

// SetID sets the value of ID.
func (s *Organization) SetID(val OptString) {
	s.ID = val
}

// SetDateCreated sets the value of DateCreated.
func (s *Organization) SetDateCreated(val OptDateTime) {
	s.DateCreated = val
}

// SetDateUpdated sets the value of DateUpdated.
func (s *Organization) SetDateUpdated(val OptDateTime) {
	s.DateUpdated = val
}

// SetType sets the value of Type.
func (s *Organization) SetType(val OptString) {
	s.Type = val
}

// SetAcronym sets the value of Acronym.
func (s *Organization) SetAcronym(val OptString) {
	s.Acronym = val
}

// SetNameDut sets the value of NameDut.
func (s *Organization) SetNameDut(val OptString) {
	s.NameDut = val
}

// SetNameEng sets the value of NameEng.
func (s *Organization) SetNameEng(val OptString) {
	s.NameEng = val
}

// SetParent sets the value of Parent.
func (s *Organization) SetParent(val []OrganizationParent) {
	s.Parent = val
}

// SetIdentifier sets the value of Identifier.
func (s *Organization) SetIdentifier(val []PropertyValue) {
	s.Identifier = val
}

// Ref: #/components/schemas/OrganizationListResponse
type OrganizationListResponse struct {
	Cursor OptString      `json:"cursor"`
	Data   []Organization `json:"data"`
}

// GetCursor returns the value of Cursor.
func (s *OrganizationListResponse) GetCursor() OptString {
	return s.Cursor
}

// GetData returns the value of Data.
func (s *OrganizationListResponse) GetData() []Organization {
	return s.Data
}

// SetCursor sets the value of Cursor.
func (s *OrganizationListResponse) SetCursor(val OptString) {
	s.Cursor = val
}

// SetData sets the value of Data.
func (s *OrganizationListResponse) SetData(val []Organization) {
	s.Data = val
}

// Ref: #/components/schemas/OrganizationMember
type OrganizationMember struct {
	ID          string      `json:"id"`
	DateCreated OptDateTime `json:"date_created"`
	DateUpdated OptDateTime `json:"date_updated"`
}

// GetID returns the value of ID.
func (s *OrganizationMember) GetID() string {
	return s.ID
}

// GetDateCreated returns the value of DateCreated.
func (s *OrganizationMember) GetDateCreated() OptDateTime {
	return s.DateCreated
}

// GetDateUpdated returns the value of DateUpdated.
func (s *OrganizationMember) GetDateUpdated() OptDateTime {
	return s.DateUpdated
}

// SetID sets the value of ID.
func (s *OrganizationMember) SetID(val string) {
	s.ID = val
}

// SetDateCreated sets the value of DateCreated.
func (s *OrganizationMember) SetDateCreated(val OptDateTime) {
	s.DateCreated = val
}

// SetDateUpdated sets the value of DateUpdated.
func (s *OrganizationMember) SetDateUpdated(val OptDateTime) {
	s.DateUpdated = val
}

// Ref: #/components/schemas/OrganizationParent
type OrganizationParent struct {
	ID          string      `json:"id"`
	DateCreated OptDateTime `json:"date_created"`
	DateUpdated OptDateTime `json:"date_updated"`
	From        time.Time   `json:"from"`
	Until       OptDateTime `json:"until"`
}

// GetID returns the value of ID.
func (s *OrganizationParent) GetID() string {
	return s.ID
}

// GetDateCreated returns the value of DateCreated.
func (s *OrganizationParent) GetDateCreated() OptDateTime {
	return s.DateCreated
}

// GetDateUpdated returns the value of DateUpdated.
func (s *OrganizationParent) GetDateUpdated() OptDateTime {
	return s.DateUpdated
}

// GetFrom returns the value of From.
func (s *OrganizationParent) GetFrom() time.Time {
	return s.From
}

// GetUntil returns the value of Until.
func (s *OrganizationParent) GetUntil() OptDateTime {
	return s.Until
}

// SetID sets the value of ID.
func (s *OrganizationParent) SetID(val string) {
	s.ID = val
}

// SetDateCreated sets the value of DateCreated.
func (s *OrganizationParent) SetDateCreated(val OptDateTime) {
	s.DateCreated = val
}

// SetDateUpdated sets the value of DateUpdated.
func (s *OrganizationParent) SetDateUpdated(val OptDateTime) {
	s.DateUpdated = val
}

// SetFrom sets the value of From.
func (s *OrganizationParent) SetFrom(val time.Time) {
	s.From = val
}

// SetUntil sets the value of Until.
func (s *OrganizationParent) SetUntil(val OptDateTime) {
	s.Until = val
}

// Ref: #/components/schemas/Person
type Person struct {
	ID                  OptString            `json:"id"`
	Active              OptBool              `json:"active"`
	DateCreated         OptDateTime          `json:"date_created"`
	DateUpdated         OptDateTime          `json:"date_updated"`
	Name                OptString            `json:"name"`
	GivenName           OptString            `json:"given_name"`
	FamilyName          OptString            `json:"family_name"`
	Email               OptString            `json:"email"`
	Token               []PropertyValue      `json:"token"`
	PreferredGivenName  OptString            `json:"preferred_given_name"`
	PreferredFamilyName OptString            `json:"preferred_family_name"`
	BirthDate           OptString            `json:"birth_date"`
	HonorificPrefix     OptString            `json:"honorific_prefix"`
	Identifier          []PropertyValue      `json:"identifier"`
	Organization        []OrganizationMember `json:"organization"`
	JobCategory         []string             `json:"job_category"`
	Role                []string             `json:"role"`
	Settings            OptPersonSettings    `json:"settings"`
	ObjectClass         []string             `json:"object_class"`
	ExpirationDate      OptString            `json:"expiration_date"`
}

// GetID returns the value of ID.
func (s *Person) GetID() OptString {
	return s.ID
}

// GetActive returns the value of Active.
func (s *Person) GetActive() OptBool {
	return s.Active
}

// GetDateCreated returns the value of DateCreated.
func (s *Person) GetDateCreated() OptDateTime {
	return s.DateCreated
}

// GetDateUpdated returns the value of DateUpdated.
func (s *Person) GetDateUpdated() OptDateTime {
	return s.DateUpdated
}

// GetName returns the value of Name.
func (s *Person) GetName() OptString {
	return s.Name
}

// GetGivenName returns the value of GivenName.
func (s *Person) GetGivenName() OptString {
	return s.GivenName
}

// GetFamilyName returns the value of FamilyName.
func (s *Person) GetFamilyName() OptString {
	return s.FamilyName
}

// GetEmail returns the value of Email.
func (s *Person) GetEmail() OptString {
	return s.Email
}

// GetToken returns the value of Token.
func (s *Person) GetToken() []PropertyValue {
	return s.Token
}

// GetPreferredGivenName returns the value of PreferredGivenName.
func (s *Person) GetPreferredGivenName() OptString {
	return s.PreferredGivenName
}

// GetPreferredFamilyName returns the value of PreferredFamilyName.
func (s *Person) GetPreferredFamilyName() OptString {
	return s.PreferredFamilyName
}

// GetBirthDate returns the value of BirthDate.
func (s *Person) GetBirthDate() OptString {
	return s.BirthDate
}

// GetHonorificPrefix returns the value of HonorificPrefix.
func (s *Person) GetHonorificPrefix() OptString {
	return s.HonorificPrefix
}

// GetIdentifier returns the value of Identifier.
func (s *Person) GetIdentifier() []PropertyValue {
	return s.Identifier
}

// GetOrganization returns the value of Organization.
func (s *Person) GetOrganization() []OrganizationMember {
	return s.Organization
}

// GetJobCategory returns the value of JobCategory.
func (s *Person) GetJobCategory() []string {
	return s.JobCategory
}

// GetRole returns the value of Role.
func (s *Person) GetRole() []string {
	return s.Role
}

// GetSettings returns the value of Settings.
func (s *Person) GetSettings() OptPersonSettings {
	return s.Settings
}

// GetObjectClass returns the value of ObjectClass.
func (s *Person) GetObjectClass() []string {
	return s.ObjectClass
}

// GetExpirationDate returns the value of ExpirationDate.
func (s *Person) GetExpirationDate() OptString {
	return s.ExpirationDate
}

// SetID sets the value of ID.
func (s *Person) SetID(val OptString) {
	s.ID = val
}

// SetActive sets the value of Active.
func (s *Person) SetActive(val OptBool) {
	s.Active = val
}

// SetDateCreated sets the value of DateCreated.
func (s *Person) SetDateCreated(val OptDateTime) {
	s.DateCreated = val
}

// SetDateUpdated sets the value of DateUpdated.
func (s *Person) SetDateUpdated(val OptDateTime) {
	s.DateUpdated = val
}

// SetName sets the value of Name.
func (s *Person) SetName(val OptString) {
	s.Name = val
}

// SetGivenName sets the value of GivenName.
func (s *Person) SetGivenName(val OptString) {
	s.GivenName = val
}

// SetFamilyName sets the value of FamilyName.
func (s *Person) SetFamilyName(val OptString) {
	s.FamilyName = val
}

// SetEmail sets the value of Email.
func (s *Person) SetEmail(val OptString) {
	s.Email = val
}

// SetToken sets the value of Token.
func (s *Person) SetToken(val []PropertyValue) {
	s.Token = val
}

// SetPreferredGivenName sets the value of PreferredGivenName.
func (s *Person) SetPreferredGivenName(val OptString) {
	s.PreferredGivenName = val
}

// SetPreferredFamilyName sets the value of PreferredFamilyName.
func (s *Person) SetPreferredFamilyName(val OptString) {
	s.PreferredFamilyName = val
}

// SetBirthDate sets the value of BirthDate.
func (s *Person) SetBirthDate(val OptString) {
	s.BirthDate = val
}

// SetHonorificPrefix sets the value of HonorificPrefix.
func (s *Person) SetHonorificPrefix(val OptString) {
	s.HonorificPrefix = val
}

// SetIdentifier sets the value of Identifier.
func (s *Person) SetIdentifier(val []PropertyValue) {
	s.Identifier = val
}

// SetOrganization sets the value of Organization.
func (s *Person) SetOrganization(val []OrganizationMember) {
	s.Organization = val
}

// SetJobCategory sets the value of JobCategory.
func (s *Person) SetJobCategory(val []string) {
	s.JobCategory = val
}

// SetRole sets the value of Role.
func (s *Person) SetRole(val []string) {
	s.Role = val
}

// SetSettings sets the value of Settings.
func (s *Person) SetSettings(val OptPersonSettings) {
	s.Settings = val
}

// SetObjectClass sets the value of ObjectClass.
func (s *Person) SetObjectClass(val []string) {
	s.ObjectClass = val
}

// SetExpirationDate sets the value of ExpirationDate.
func (s *Person) SetExpirationDate(val OptString) {
	s.ExpirationDate = val
}

// Ref: #/components/schemas/PersonListResponse
type PersonListResponse struct {
	Cursor OptString `json:"cursor"`
	Data   []Person  `json:"data"`
}

// GetCursor returns the value of Cursor.
func (s *PersonListResponse) GetCursor() OptString {
	return s.Cursor
}

// GetData returns the value of Data.
func (s *PersonListResponse) GetData() []Person {
	return s.Data
}

// SetCursor sets the value of Cursor.
func (s *PersonListResponse) SetCursor(val OptString) {
	s.Cursor = val
}

// SetData sets the value of Data.
func (s *PersonListResponse) SetData(val []Person) {
	s.Data = val
}

type PersonSettings map[string]string

func (s *PersonSettings) init() PersonSettings {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

// Ref: #/components/schemas/PropertyValue
type PropertyValue struct {
	Type       string `json:"type"`
	PropertyID string `json:"property_id"`
	Value      string `json:"value"`
}

// GetType returns the value of Type.
func (s *PropertyValue) GetType() string {
	return s.Type
}

// GetPropertyID returns the value of PropertyID.
func (s *PropertyValue) GetPropertyID() string {
	return s.PropertyID
}

// GetValue returns the value of Value.
func (s *PropertyValue) GetValue() string {
	return s.Value
}

// SetType sets the value of Type.
func (s *PropertyValue) SetType(val string) {
	s.Type = val
}

// SetPropertyID sets the value of PropertyID.
func (s *PropertyValue) SetPropertyID(val string) {
	s.PropertyID = val
}

// SetValue sets the value of Value.
func (s *PropertyValue) SetValue(val string) {
	s.Value = val
}

// Ref: #/components/schemas/SetPersonOrcidRequest
type SetPersonOrcidRequest struct {
	ID    string `json:"id"`
	Orcid string `json:"orcid"`
}

// GetID returns the value of ID.
func (s *SetPersonOrcidRequest) GetID() string {
	return s.ID
}

// GetOrcid returns the value of Orcid.
func (s *SetPersonOrcidRequest) GetOrcid() string {
	return s.Orcid
}

// SetID sets the value of ID.
func (s *SetPersonOrcidRequest) SetID(val string) {
	s.ID = val
}

// SetOrcid sets the value of Orcid.
func (s *SetPersonOrcidRequest) SetOrcid(val string) {
	s.Orcid = val
}

// Ref: #/components/schemas/SetPersonOrcidTokenRequest
type SetPersonOrcidTokenRequest struct {
	ID         string `json:"id"`
	OrcidToken string `json:"orcid_token"`
}

// GetID returns the value of ID.
func (s *SetPersonOrcidTokenRequest) GetID() string {
	return s.ID
}

// GetOrcidToken returns the value of OrcidToken.
func (s *SetPersonOrcidTokenRequest) GetOrcidToken() string {
	return s.OrcidToken
}

// SetID sets the value of ID.
func (s *SetPersonOrcidTokenRequest) SetID(val string) {
	s.ID = val
}

// SetOrcidToken sets the value of OrcidToken.
func (s *SetPersonOrcidTokenRequest) SetOrcidToken(val string) {
	s.OrcidToken = val
}

// Ref: #/components/schemas/SetPersonRoleRequest
type SetPersonRoleRequest struct {
	ID   string   `json:"id"`
	Role []string `json:"role"`
}

// GetID returns the value of ID.
func (s *SetPersonRoleRequest) GetID() string {
	return s.ID
}

// GetRole returns the value of Role.
func (s *SetPersonRoleRequest) GetRole() []string {
	return s.Role
}

// SetID sets the value of ID.
func (s *SetPersonRoleRequest) SetID(val string) {
	s.ID = val
}

// SetRole sets the value of Role.
func (s *SetPersonRoleRequest) SetRole(val []string) {
	s.Role = val
}

// Ref: #/components/schemas/SetPersonSettingsRequest
type SetPersonSettingsRequest struct {
	ID       string                           `json:"id"`
	Settings SetPersonSettingsRequestSettings `json:"settings"`
}

// GetID returns the value of ID.
func (s *SetPersonSettingsRequest) GetID() string {
	return s.ID
}

// GetSettings returns the value of Settings.
func (s *SetPersonSettingsRequest) GetSettings() SetPersonSettingsRequestSettings {
	return s.Settings
}

// SetID sets the value of ID.
func (s *SetPersonSettingsRequest) SetID(val string) {
	s.ID = val
}

// SetSettings sets the value of Settings.
func (s *SetPersonSettingsRequest) SetSettings(val SetPersonSettingsRequestSettings) {
	s.Settings = val
}

type SetPersonSettingsRequestSettings map[string]string

func (s *SetPersonSettingsRequestSettings) init() SetPersonSettingsRequestSettings {
	m := *s
	if m == nil {
		m = map[string]string{}
		*s = m
	}
	return m
}

// Ref: #/components/schemas/SuggestOrganizationsRequest
type SuggestOrganizationsRequest struct {
	Query string `json:"query"`
}

// GetQuery returns the value of Query.
func (s *SuggestOrganizationsRequest) GetQuery() string {
	return s.Query
}

// SetQuery sets the value of Query.
func (s *SuggestOrganizationsRequest) SetQuery(val string) {
	s.Query = val
}

// Ref: #/components/schemas/SuggestPeopleRequest
type SuggestPeopleRequest struct {
	Query string `json:"query"`
}

// GetQuery returns the value of Query.
func (s *SuggestPeopleRequest) GetQuery() string {
	return s.Query
}

// SetQuery sets the value of Query.
func (s *SuggestPeopleRequest) SetQuery(val string) {
	s.Query = val
}
