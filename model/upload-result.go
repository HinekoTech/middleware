package model

type UploadResult struct {
	ID        string `json:"id"`
	Cdn       string `json:"cdn"`
	FileName  string `json:"file_name"`
	Extension string `json:"extension"`
	Path      string `json:"path"`
	FullPath  string `json:"full_path"`
}
