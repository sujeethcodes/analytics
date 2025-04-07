package main
import(
	"github.com/labstack/echo/v4"
	"analytics/connectors"
	"analytics/repository"
	"log"
	"os"
)


func init() {
	connectors.LoadEnv()
	connectors.LoadLogger()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	mongoInstance := mongodb.MongoConnect()
	if mongoInstance != nil {
		log.Println("Database Connected Successfully")
	}

	e := echo.New()
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}