package auction

import (
	"context"
	"fmt"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuctionAutoClose(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start MongoDB container: %v", err)
	}
	t.Cleanup(func() { mongoContainer.Terminate(ctx) })

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}
	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s/testauctions", host, port.Port())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	t.Cleanup(func() { client.Disconnect(ctx) })

	repo := NewAuctionRepository(client.Database("testauctions"))

	// Use a short duration so the test runs quickly
	t.Setenv("AUCTION_DURATION", "3s")

	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Test Laptop",
		"Electronics",
		"A high-performance laptop for testing the auto-close functionality",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatalf("failed to create auction entity: %v", internalErr)
	}

	if createErr := repo.CreateAuction(ctx, auctionEntity); createErr != nil {
		t.Fatalf("failed to persist auction: %v", createErr)
	}

	// Auction should start as Active
	found, findErr := repo.FindAuctionById(ctx, auctionEntity.Id)
	if findErr != nil {
		t.Fatalf("failed to find auction: %v", findErr)
	}
	if found.Status != auction_entity.Active {
		t.Errorf("expected initial status Active (%d), got %d", auction_entity.Active, found.Status)
	}

	// Wait for the goroutine to fire (duration + 2s buffer)
	time.Sleep(5 * time.Second)

	found, findErr = repo.FindAuctionById(ctx, auctionEntity.Id)
	if findErr != nil {
		t.Fatalf("failed to find auction after timeout: %v", findErr)
	}
	if found.Status != auction_entity.Completed {
		t.Errorf("expected status Completed (%d) after auto-close, got %d", auction_entity.Completed, found.Status)
	}
}
