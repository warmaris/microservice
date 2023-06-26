package exibillia

// UpdateRequest is just DTO. It has limited set of fields against of main entity.
// If we don't allow user to update some fields of entity after creation, we use simpler DTO to call update method.
type UpdateRequest struct {
	ID          uint64
	Description string
	Tags        []string
}
