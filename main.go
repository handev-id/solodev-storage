package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"storage/utils"
)

func main() {
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	store, err := utils.NewS3Storage(cfg)
	if err != nil {
		log.Fatalf("s3 init error: %v", err)
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/", utils.BearerAuth(cfg.UploadSecretKey), func(c *fiber.Ctx) error {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "multipart field 'file' is required"})
		}

		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to open uploaded file"})
		}
		defer file.Close()

		key, err := utils.ResolveUploadKey(c, fileHeader.Filename)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if err := store.UploadObject(c.Context(), key, file, contentType); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to upload file"})
		}

		publicURL := utils.BuildPublicURL(cfg.PublicBaseURL, key)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":    "file uploaded",
			"key":        key,
			"public_url": publicURL,
		})
	})

	app.Get("/*", func(c *fiber.Ctx) error {
		key := utils.NormalizeKey(c.Params("*"))
		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "path parameter 'key' is required"})
		}

		obj, err := store.GetObject(c.Context(), key)
		if err != nil {
			if store.IsNotFound(err) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch file"})
		}
		defer obj.Body.Close()

		if obj.ContentType != "" {
			c.Set("Content-Type", obj.ContentType)
		} else {
			c.Set("Content-Type", "application/octet-stream")
		}

		if obj.ETag != "" {
			c.Set("ETag", obj.ETag)
		}

		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", utils.SafeFilename(key)))

		if obj.ContentLength > 0 {
			return c.SendStream(obj.Body, int(obj.ContentLength))
		}

		return c.SendStream(obj.Body)
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("running at %s", addr)
	log.Fatal(app.Listen(addr))
}
