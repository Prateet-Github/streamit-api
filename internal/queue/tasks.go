package queue

const (
	TypeProcessVideo = "video:process"
)

type VideoTask struct {
	VideoID string `json:"videoId"`
	S3Key   string `json:"s3Key"`
}
