package domain
import(
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Name         string
	Role         string 
	CreatedAt    time.Time
}