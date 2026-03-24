package dto


//что клиент обязан прислать, чтобы создать комнату
type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity" validate:"required,min=1,max=1000"`
	Location    string `json:"location"`
}

//что клиент должен прислать чтобы обновить комнату
type UpdateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=120"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity" validate:"required,min=1,max=1000"`
	Location    string `json:"location"`
}

//то что возвращается клтенту
type RoomResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	Location    string `json:"location"`
	IsActive    bool   `json:"is_active"`
}