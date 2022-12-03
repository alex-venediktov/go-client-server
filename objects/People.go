package objects

import "github.com/google/uuid"

type People struct {
	Id   uuid.UUID
	Name string
	Age  int
}

func (people People) GetId() any {
	return people.Id
}

func (people People) SetId(id any) any {
	people.Id = id.(uuid.UUID)
	return people
}
