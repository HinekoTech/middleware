package common

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/FourWD/middleware/orm"

	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var FirebaseMessageClient *messaging.Client

func ConnectFirebaseNotification(key string) error {
	opt := option.WithCredentialsFile(key)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	FirebaseMessageClient, err = app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return err
	}

	return nil
}

// config struct message
var MessageConfig = struct {
	AndroidConfig *messaging.AndroidConfig
	APNSConfig    *messaging.APNSConfig
}{
	AndroidConfig: &messaging.AndroidConfig{
		Priority: "high",
	},
	APNSConfig: &messaging.APNSConfig{
		Headers: map[string]string{"apns-priority": "10"},
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Sound: "default",
			},
		},
	},
}

func SendMessageToUser(userToken string, title string, body string, data map[string]string) error {
	// Access title and body directly from the data map

	message := &messaging.Message{
		Data:  data,
		Token: userToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: MessageConfig.AndroidConfig,
		APNS:    MessageConfig.APNSConfig,
	}

	result, err := FirebaseMessageClient.Send(context.Background(), message)
	if err != nil {
		// Instead of log.Fatalf, return the error
		return fmt.Errorf("error sending message: %s, %v", result, err)
	}

	return nil
}

// func AddUserToSubscription(userToken string, topic string) error { // เอาคน (topic) เข้า กรุป auction
// 	fmt.Printf("UserToken: %s, Topic: %s\n", userToken, topic)
// 	_, err := FirebaseMessageClient.SubscribeToTopic(context.Background(), []string{userToken}, topic)
// 	if err != nil {
// 		log.Fatalf("error subscribing user to topic: %v\n", err)
// 		return err
// 	}

//		return nil
//	}
func AddUserToSubscription(topic string, userID string, userToken string) error {

	notificationTopic := orm.NotificationTopic{
		ID:   uuid.NewString(),
		Name: topic,
	}
	if err := Database.Debug().Create(&notificationTopic).Error; err != nil {
		log.Fatalf("failed to insert notification topic: %v\n", err)
	}
	if err := Database.Where("name = ?", topic).First(&notificationTopic).Error; err != nil {
		log.Fatalf("failed to insert notification topic: %v\n", err)
	}
	if notificationTopic.ID != "" {
		notificationTopicUser := orm.NotificationTopicUser{
			ID:                  uuid.NewString(),
			NotificationTopicID: notificationTopic.ID,
			UserID:              userID,
		}
		if err := Database.Debug().Create(&notificationTopicUser).Error; err != nil {
			log.Fatalf("failed to insert notification user topic user: %v\n", err)
		}
		_, err := FirebaseMessageClient.SubscribeToTopic(context.Background(), []string{userToken}, topic)
		if err != nil {
			log.Fatalf("error subscribing user to topic: %v\n", err)
			return err

		}
	}

	return nil
}

// func RemoveUserFromSubscription(topic string, userID string, userToken string) error { // เอาคน (topic) ออก กรุป auction
// 	fmt.Printf("UserToken: %s, Topic: %s\n", userToken, topic)

// 	_, err := FirebaseMessageClient.UnsubscribeFromTopic(context.Background(), []string{userToken}, topic)
// 	if err != nil {
// 		log.Fatalf("error unsubscribing user from topic: %v\n", err)
// 		return err
// 	}
// 	return nil
// }

func RemoveUserFromSubscription(topic string, userID string, userToken string) error { // เอาคน (topic) ออก กรุป auction
	fmt.Printf("UserToken: %s, userID: %s, Topic: %s\n", userToken, userID, topic)

	_, err := FirebaseMessageClient.UnsubscribeFromTopic(context.Background(), []string{userToken}, topic)
	if err != nil {
		log.Fatalf("error unsubscribing user from topic: %v\n", err)
		return err
	}
	var notificationTopicID string
	err = Database.Table("notification_topics").Select("id").Where("name = ?", topic).Scan(&notificationTopicID).Error
	if err != nil {
		log.Fatalf("error finding notification topic ID: %v\n", err)
		return err
	}

	if notificationTopicID == "" {
		log.Fatalf("notification topic ID not found for topic: %s\n", topic)
		return fmt.Errorf("notification topic ID not found for topic: %s", topic)
	}

	err = Database.Where("notification_topic_id = ? AND user_id = ?", notificationTopicID, userID).Unscoped().Debug().Delete(&orm.NotificationTopicUser{}).Error
	if err != nil {
		log.Fatalf("error removing user from topic in database: %v\n", err)
		return err
	}

	return nil
}

func SendMessageToSubscriber(topic string, title string, body string, data map[string]string) error {
	Print("data", title)
	message := &messaging.Message{
		// Title Body // R001 = ประกาศผล
		Data:  data,
		Topic: topic, // all_users
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	_, err := FirebaseMessageClient.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
		return err
	}

	return nil
}
