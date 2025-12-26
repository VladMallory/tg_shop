// product.go - Описание сущности товара.
// Этот файл - часть Domain Layer (слоя предметной области).
// Он ничего не знает о базе данных или телеграме. Это просто структура данных.
package domain

// ProductType - специальный тип для категории духов.
// Используем его вместо string, чтобы избежать опечаток (например, "femal" вместо "female").
type ProductType string

// Константы для типов духов. Это как Enum в других языках.
const (
	TypeFemale ProductType = "female" // Женские
	TypeMale   ProductType = "male"   // Мужские
	TypeUnisex ProductType = "unisex" // Унисекс
)

// Product - структура, описывающая одни духи.
// JSON теги нужны, если мы захотим превратить эту структуру в текст (например, для логов или API).
type Product struct {
	ID          int64       `json:"id"`          // Уникальный номер в базе данных
	Type        ProductType `json:"type"`        // Тип (муж/жен/уни)
	Name        string      `json:"name"`        // Название (например, "Chanel No. 5")
	Description string      `json:"description"` // Описание аромата
	Price       float64     `json:"price"`       // Цена (дробное число)
	ImageID     string      `json:"image_id"`    // ID файла картинки в Телеграме (мы не храним само фото, только ссылку)
}
