package users

import (
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nomadphone/lib/models"
	"github.com/nomadphone/telegram-bot/database"
	"github.com/nomadphone/telegram-bot/phonenumbers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserFromTwillioPhone(twillioPhone string) models.User {
	client, ctx, cancel := database.GetClient()
	defer client.Disconnect(ctx)
	defer cancel()
	var result models.User
	p := phonenumbers.NumbersOnly(twillioPhone)
	c := client.Database("nomadphone").Collection("users")
	filter := bson.M{"twilliophone": p}
	err := c.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		panic(err)
	}
	return result
}

func ManageTelegramUser(telegramUser *telegram.User, chatID int64) models.User {
	client, ctx, cancel := database.GetClient()
	defer client.Disconnect(ctx)
	defer cancel()
	c := client.Database("nomadphone").Collection("users")
	filter := bson.M{"telegramusername": telegramUser.UserName}
	// Try to find telegram username on the database
	var user models.User
	err := c.FindOne(ctx, filter).Decode(&user)
	if err != nil && err != mongo.ErrNoDocuments {
		panic(err)
	}

	// If not found, create a new user
	if err == mongo.ErrNoDocuments {
		log.Printf("No documents found")
		// Register a twillio phone number for this lovely person
		phoneNumber := phonenumbers.NewPhoneNumberProvider().ProvideNewNumber()

		// Create a new user
		newUser := models.User{
			TelegramUsername: telegramUser.UserName,
			Name:             telegramUser.FirstName + " " + telegramUser.LastName,
			TwillioPhone:     phonenumbers.NumbersOnly(phoneNumber),
			TelegramChatID:   chatID,
		}
		_, err := c.InsertOne(ctx, newUser)
		if err != nil {
			panic(err)
		}
		return newUser
	}

	// If found just update
	result, err := c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"telegramchatid": chatID}})
	if err != nil {
		panic(err)
	}
	log.Printf("Updated %v documents!", result.ModifiedCount)
	user.TelegramChatID = chatID
	return user
}
