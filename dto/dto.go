package dto

type DTO interface {
	IsValid() error
	ToObject(
		data []byte,
		obj any,
	) error
}
