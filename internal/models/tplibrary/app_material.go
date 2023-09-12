package tplibrary

type IdName struct {
	Id   int64
	Name string
}

type MaterialRequest struct {
	Id                 int64
	Title              string
	SubjectId          int64
	DifficultcyLevelId int64
	MaterialTypeId     int64
	FileLink           *string
	File               []byte
}

type MaterialResponse struct {
	Id               int64
	Title            string
	Size             *int64
	MaterialType     IdName
	DifficultcyLevel IdName
	Subject          IdName
	FileLink         string
}

type FileInfo struct {
	Key  string
	Size int64
	Url  string
}
