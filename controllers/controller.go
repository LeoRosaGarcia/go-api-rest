package controllers

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/LeoRosaGarcia/go-api-rest/config"
	"github.com/LeoRosaGarcia/go-api-rest/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllTodoContents(c *fiber.Ctx) error {

	TodoContentCollection := config.MI.DB.Collection("TodoContents")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var TodoContents []models.TodoContent

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"TodoTitle": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"TodoContent": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := TodoContentCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := TodoContentCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "TodoContents Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var TodoContent models.TodoContent
		cursor.Decode(&TodoContent)
		TodoContents = append(TodoContents, TodoContent)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      TodoContents,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetTodoContent(c *fiber.Ctx) error {

	TodoContentCollection := config.MI.DB.Collection("TodoContents")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	var TodoContent models.TodoContent
	objId, err := primitive.ObjectIDFromHex(c.Params("_id"))
	findResult := TodoContentCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&TodoContent)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    TodoContent,
		"success": true,
	})
}

func AddTodoContent(c *fiber.Ctx) error {
	TodoContentCollection := config.MI.DB.Collection("TodoContents")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	TodoContent := new(models.TodoContent)

	if err := c.BodyParser(TodoContent); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := TodoContentCollection.InsertOne(ctx, TodoContent)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "TodoContent inserted successfully",
	})
}

func UpdateTodoContent(c *fiber.Ctx) error {

	TodoContentCollection := config.MI.DB.Collection("TodoContents")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	TodoContent := new(models.TodoContent)

	if err := c.BodyParser(TodoContent); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("_id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": TodoContent,
	}
	_, err = TodoContentCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "TodoContent updated successfully",
	})
}

func DeleteTodoContent(c *fiber.Ctx) error {

	TodoContentCollection := config.MI.DB.Collection("TodoContents")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	objId, err := primitive.ObjectIDFromHex(c.Params("_id"))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent not found",
			"error":   err,
		})
	}
	_, err = TodoContentCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "TodoContent failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "TodoContent deleted successfully",
	})
}
