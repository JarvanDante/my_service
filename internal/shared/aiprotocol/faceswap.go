// Package aiprotocol 是 my_service 与 my_ai_worker 的 Kafka JSON 契约。
// 只传对象引用，不传图片字节；输入输出都在 my-storage。
// 换脸与去衣共用 jobs/results topic，靠 JobMessage.Biz 分发。
package aiprotocol

const (
	TopicJobs    = "ai.faceswap.jobs"
	TopicResults = "ai.faceswap.results"

	BizFaceSwap = "faceswap"
	BizUndress  = "undress"

	MediaPhoto = "photo"
	MediaVideo = "video"

	StatusReady  = "ready"
	StatusFailed = "failed"

	SchemaVersion = 1

	DefaultBucket = "my-storage"
)

// ObjectRef MinIO 对象定位。
type ObjectRef struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func (r ObjectRef) OK() bool {
	return r.Bucket != "" && r.Key != ""
}

// OutputRef 输出目录前缀，最终 output_key = prefix + "result.bnc"。
type OutputRef struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

// JobMessage 换脸任务。
type JobMessage struct {
	SchemaVersion int       `json:"schema_version"`
	JobID         string    `json:"job_id"`
	Biz           string    `json:"biz"`
	MediaType     string    `json:"media_type"`
	Source        ObjectRef `json:"source"`
	Target        ObjectRef `json:"target"`
	Output        OutputRef `json:"output"`
	CreatedAt     string    `json:"created_at,omitempty"`
}

// ResultMessage 换脸结果。
type ResultMessage struct {
	SchemaVersion int    `json:"schema_version"`
	JobID         string `json:"job_id"`
	Biz           string `json:"biz"`
	MediaType     string `json:"media_type"`
	Status        string `json:"status"`
	OutputKey     string `json:"output_key,omitempty"`
	OutputURL     string `json:"output_url,omitempty"`
	Error         string `json:"error,omitempty"`
	FinishedAt    string `json:"finished_at,omitempty"`
}
