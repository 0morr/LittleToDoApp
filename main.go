package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"html/template"
	"log"
)

// Task representa una tarea en nuestra lista
type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

// TaskStore almacena nuestras tareas en memoria
type TaskStore struct {
	tasks []Task
}

// NewTaskStore crea un nuevo almacén de tareas
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make([]Task, 0),
	}
}

func main() {
	app := fiber.New(fiber.Config{
		// Configuración para manejar 404 correctamente
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Si la ruta no existe, servimos el index.html
			return c.SendFile("index.html")
		},
	})
	
	store := NewTaskStore()

	// Middleware para manejar CORS
	app.Use(cors.New())

	// API endpoints
	api := app.Group("/api")

	// Obtener todas las tareas
	api.Get("/tasks", func(c *fiber.Ctx) error {
		return c.JSON(store.tasks)
	})

	// Crear una nueva tarea
	api.Post("/tasks", func(c *fiber.Ctx) error {
		var task Task
		if err := c.BodyParser(&task); err != nil {
			return err
		}

		task.ID = uuid.New().String()
		task.Completed = false
		store.tasks = append(store.tasks, task)

		// Renderizar el HTML para la nueva tarea
		html := template.Must(template.New("task").Parse(`
            <div class="flex items-center justify-between p-3 bg-gray-50 rounded"
                 hx-target="this"
                 hx-swap="outerHTML">
                <span>{{.Name}}</span>
                <div class="flex gap-2">
                    <button hx-put="/api/tasks/toggle/{{.ID}}"
                            class="text-sm px-2 py-1 rounded bg-green-100 text-green-700 hover:bg-green-200">
                        ✓
                    </button>
                    <button hx-delete="/api/tasks/{{.ID}}"
                            class="text-sm px-2 py-1 rounded bg-red-100 text-red-700 hover:bg-red-200">
                        ×
                    </button>
                </div>
            </div>
        `))

		return html.Execute(c, task)
	})

	// Marcar una tarea como completada/no completada
	api.Put("/tasks/toggle/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i := range store.tasks {
			if store.tasks[i].ID == id {
				store.tasks[i].Completed = !store.tasks[i].Completed
				
				// Renderizar el HTML actualizado
				html := template.Must(template.New("task").Parse(`
                    <div class="flex items-center justify-between p-3 bg-gray-50 rounded"
                         hx-target="this"
                         hx-swap="outerHTML">
                        <span class="{{if .Completed}}line-through{{end}}">{{.Name}}</span>
                        <div class="flex gap-2">
                            <button hx-put="/api/tasks/toggle/{{.ID}}"
                                    class="text-sm px-2 py-1 rounded bg-green-100 text-green-700 hover:bg-green-200">
                                ✓
                            </button>
                            <button hx-delete="/api/tasks/{{.ID}}"
                                    class="text-sm px-2 py-1 rounded bg-red-100 text-red-700 hover:bg-red-200">
                                ×
                            </button>
                        </div>
                    </div>
                `))

				return html.Execute(c, store.tasks[i])
			}
		}
		return c.Status(404).SendString("Task not found")
	})

	// Eliminar una tarea
	api.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i := range store.tasks {
			if store.tasks[i].ID == id {
				// Eliminar la tarea del slice
				store.tasks = append(store.tasks[:i], store.tasks[i+1:]...)
				return c.SendString("")
			}
		}
		return c.Status(404).SendString("Task not found")
	})

	// Ruta raíz - sirve el index.html
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	log.Fatal(app.Listen(":3000"))
}