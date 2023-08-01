package models

type AvailableGPTModels struct {
	Object string  `json:"object"`
	Data   []Datum `json:"data"`
}
type Datum struct {
	ID         string       `json:"id"`
	Object     DatumObject  `json:"object"`
	Created    int64        `json:"created"`
	OwnedBy    OwnedBy      `json:"owned_by"`
	Permission []Permission `json:"permission"`
	Root       string       `json:"root"`
	Parent     interface{}  `json:"parent"`
}
type Permission struct {
	ID                 string           `json:"id"`
	Object             PermissionObject `json:"object"`
	Created            int64            `json:"created"`
	AllowCreateEngine  bool             `json:"allow_create_engine"`
	AllowSampling      bool             `json:"allow_sampling"`
	AllowLogprobs      bool             `json:"allow_logprobs"`
	AllowSearchIndices bool             `json:"allow_search_indices"`
	AllowView          bool             `json:"allow_view"`
	AllowFineTuning    bool             `json:"allow_fine_tuning"`
	Organization       Organization     `json:"organization"`
	Group              interface{}      `json:"group"`
	IsBlocking         bool             `json:"is_blocking"`
}
type DatumObject string

const (
	Model DatumObject = "model"
)

type OwnedBy string

const (
	Openai         OwnedBy = "openai"
	OpenaiDev      OwnedBy = "openai-dev"
	OpenaiInternal OwnedBy = "openai-internal"
)

type PermissionObject string

const (
	ModelPermission PermissionObject = "model_permission"
)

type Organization string

const (
	Empty Organization = "*"
)
