package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var client *mongo.Client
var collection *mongo.Collection

func connectToMongoDB() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://kunalkhare21:KUNALKHARE21@cluster0.lqma2.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	// Check the connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}
	fmt.Println("Connected to MongoDB successfully!")
	collection = client.Database("urlshortener").Collection("urls") // Create a collection named "urls"
}

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB := URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	// Insert into MongoDB
	_, err := collection.InsertOne(context.TODO(), urlDB)
	if err != nil {
		fmt.Println("Error inserting URL into MongoDB:", err)
		return ""
	}
	fmt.Println("Stored URL:", urlDB)
	return shortURL
}

func getURL(id string) (URL, error) {
	var url URL
	err := collection.FindOne(context.TODO(), map[string]string{"id": id}).Decode(&url)
	if err != nil {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}
	shortURL := createURL(data.URL)
	response := struct {
		ShortURL_ string `json:"short_url"`
	}{ShortURL_: shortURL}
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Request is invalid", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {

	// Connect to MongoDB
	connectToMongoDB()
	// http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "test.html") // path to your HTML file
	// })
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	fmt.Println("Starting server at port no. 611 ...")
	err := http.ListenAndServe(":611", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
