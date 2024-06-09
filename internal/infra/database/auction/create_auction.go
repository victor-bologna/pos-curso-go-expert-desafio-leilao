package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection         *mongo.Collection
	AuctionStatusMutex *sync.Mutex
}

const AUCTION_INTERVAL = "AUCTION_INTERVAL"

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection:         database.Collection("auctions"),
		AuctionStatusMutex: &sync.Mutex{},
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}
	ar.checkIfAuctionStillOpen(ctx, auctionEntityMongo)
	return nil
}

func (ar *AuctionRepository) checkIfAuctionStillOpen(ctx context.Context, auctionEntityMongo *AuctionEntityMongo) {
	go func(auction AuctionEntityMongo) {
		for {
			time.Sleep(time.Second)
			if ar.isAuctionExpired(auction) {
				ar.AuctionStatusMutex.Lock()
				update := bson.M{
					"$set": bson.M{
						"status": auction_entity.Completed,
					},
				}
				logger.Info(fmt.Sprintf("Auction id %s closed. Updating status to %d", auction.Id, auction_entity.Completed))
				_, err := ar.Collection.UpdateByID(ctx, auction.Id, update)
				ar.AuctionStatusMutex.Unlock()
				if err != nil {
					logger.Error(fmt.Sprintf("Error while updating auction ID: %s", auction.Id), err)
					return
				}
				logger.Info(fmt.Sprintf("Auction id %s updated with status completed", auction.Id))
				return
			}
		}
	}(*auctionEntityMongo)
}

func (ar *AuctionRepository) isAuctionExpired(auctionEntity AuctionEntityMongo) bool {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	if auctionInterval == "" {
		auctionInterval = "20s"
	}
	fmt.Println("auctionInterval", auctionInterval)

	timestamp := time.Unix(auctionEntity.Timestamp, 0)
	fmt.Println("timestamp", timestamp)
	intervalDuration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		fmt.Println("Error parsing auction interval:", err)
		return false
	}

	timeUntilExpires := timestamp.Add(intervalDuration)
	currentTime := time.Now()
	fmt.Println("timeUntilExpires", timeUntilExpires)
	fmt.Println("currentTime", currentTime)

	if currentTime.After(timeUntilExpires) || currentTime.Equal(timeUntilExpires) {
		return true
	}
	return false
}
