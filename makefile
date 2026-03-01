PRODUCT_APP := product
PRODUCT_CMD_DIR := ./cmd/$(PRODUCT_APP)
BIN_DIR := ./bin
PRODUCT_BIN := $(BIN_DIR)/$(PRODUCT_APP)

PRODUCT_APP_PORT := 8080
PRODUCT_APP_DATABASE_URL := "postgres://postgres:NataZf0192274@localhost:5432/dalty?sslmode=disable"

export ADDR = localhost:$(PRODUCT_APP_PORT)
export DATABASE = $(PRODUCT_APP_DATABASE_URL)

GO := go

$(BIN_DIR):
	mkdir -p $(BIN_DIR)


build: $(BIN_DIR)
	$(GO) build -o $(PRODUCT_BIN) $(PRODUCT_CMD_DIR)

run:
	$(GO) run $(PRODUCT_CMD_DIR)

clean:
	rm $(BIN_DIR)