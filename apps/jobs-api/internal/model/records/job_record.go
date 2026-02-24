package records

type JobRecord struct {
	JobID     string
	Status    string
	PDFKey    string
	CreatedAt int64
	UpdatedAt int64
}
