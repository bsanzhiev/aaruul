package firebaseapp

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebaseApp() error {
	opt := option.WithCredentialsFile("firebaseapp/tsurhai-2a88b-firebase-adminsdk-ylazx-69b66a7521.json")
	ctx := context.Background()
	// config := firebase.Config{
	// 	ProjectID: viper.GetString("FIREBASE_PROJECT_ID"),
	// }

	var err error
	FirebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}
	log.Println("The Firebase application has been initialized.")
	return nil
}
