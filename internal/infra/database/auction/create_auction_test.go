package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateAuction(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "2s")
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("given a valid auction when create auction should return auction status completed", func(mt *mtest.T) {
		ctx := context.Background()
		auctionRepo := NewAuctionRepository(mt.DB)
		auction := auction_entity.Auction{
			Id:          "a9c60b9e-6eec-4222-bf8e-47e5a0103712",
			ProductName: "Car",
			Category:    "Car",
			Description: "Car auction",
			Condition:   0,
			Status:      0,
			Timestamp:   time.Now(),
		}

		expectedAuction := bson.D{
			{Key: "ok", Value: 1},
			{Key: "value", Value: bson.D{
				{Key: "_id", Value: auction.Id},
				{Key: "product_name", Value: auction.ProductName},
				{Key: "category", Value: auction.Category},
				{Key: "description", Value: auction.Description},
				{Key: "condition", Value: auction.Condition},
				{Key: "status", Value: auction.Status},
				{Key: "timestamp", Value: auction.Timestamp},
			}},
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse(), expectedAuction)

		err := auctionRepo.CreateAuction(ctx, &auction)
		assert.Nil(t, err)

		createAuction := mtest.CreateCursorResponse(1, "test.auction", mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: auction.Id},
				{Key: "product_name", Value: auction.ProductName},
				{Key: "category", Value: auction.Category},
				{Key: "description", Value: auction.Description},
				{Key: "condition", Value: auction.Condition},
				{Key: "status", Value: auction_entity.AuctionStatus(1)},
				{Key: "timestamp", Value: auction.Timestamp.Unix()},
			})

		mt.AddMockResponses(createAuction)

		time.Sleep(time.Second * 3)

		auctionResp, err2 := auctionRepo.FindAuctionById(ctx, auction.Id)
		assert.Nil(t, err2)
		assert.NotNil(t, auctionResp)
		assert.NotNil(t, auctionResp.Id)
		assert.Equal(t, auction_entity.AuctionStatus(1), auctionResp.Status)
	})
}
